package repository

import (
	"context"
	"database/sql"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/predicate"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type userUsageReportRepository struct {
	client *dbent.Client
	sql    *sql.DB
}

// NewUserUsageReportRepository creates a new UserUsageReportRepository
func NewUserUsageReportRepository(client *dbent.Client, sqlDB *sql.DB) service.UserUsageReportRepository {
	return &userUsageReportRepository{client: client, sql: sqlDB}
}

// GetUserReportConfig retrieves user's usage report configuration
func (r *userUsageReportRepository) GetUserReportConfig(ctx context.Context, userID int64) (*service.UserUsageReportConfig, error) {
	user, err := r.client.User.Query().
		Where(dbuser.IDEQ(userID)).
		Select(
			dbuser.FieldUsageReportEnabled,
			dbuser.FieldUsageReportSchedule,
			dbuser.FieldUsageReportTimezone,
		).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrUserNotFound, nil)
	}

	return &service.UserUsageReportConfig{
		Enabled:  user.UsageReportEnabled,
		Schedule: user.UsageReportSchedule,
		Timezone: user.UsageReportTimezone,
	}, nil
}

// UpdateUserReportConfig updates user's usage report configuration
func (r *userUsageReportRepository) UpdateUserReportConfig(ctx context.Context, userID int64, config *service.UserUsageReportConfig) error {
	_, err := r.client.User.UpdateOneID(userID).
		SetUsageReportEnabled(config.Enabled).
		SetUsageReportSchedule(config.Schedule).
		SetUsageReportTimezone(config.Timezone).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
	}
	return nil
}

// GetUsersForUsageReport retrieves users eligible for usage reports based on scope
func (r *userUsageReportRepository) GetUsersForUsageReport(ctx context.Context, scope string, now time.Time) ([]service.User, error) {
	query := r.client.User.Query().
		Where(
			dbuser.DeletedAtIsNil(),
			dbuser.StatusEQ(service.StatusActive),
			// Email must be valid (not empty and not ending with .invalid)
			dbuser.EmailNEQ(""),
			predicate.User(dbuser.Not(dbuser.EmailContains(".invalid"))),
		)

	switch scope {
	case service.UsageReportScopeAll:
		// All users with valid email - no additional filters

	case service.UsageReportScopeActiveToday:
		// Get users who had activity YESTERDAY (since report sends yesterday's data)
		// This ensures users who used the service yesterday will receive their report
		yesterday := now.AddDate(0, 0, -1)
		activeUserIDs, err := r.GetActiveUserIDs(ctx, yesterday)
		if err != nil {
			return nil, err
		}
		if len(activeUserIDs) == 0 {
			return []service.User{}, nil
		}
		query = query.Where(dbuser.IDIn(activeUserIDs...))

	case service.UsageReportScopeOptedIn:
		// Only users who opted in
		query = query.Where(dbuser.UsageReportEnabledEQ(true))

	default:
		// Default to opted_in
		query = query.Where(dbuser.UsageReportEnabledEQ(true))
	}

	users, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]service.User, 0, len(users))
	for _, u := range users {
		result = append(result, service.User{
			ID:                  u.ID,
			Email:               u.Email,
			Username:            u.Username,
			UsageReportEnabled:  u.UsageReportEnabled,
			UsageReportSchedule: u.UsageReportSchedule,
			UsageReportTimezone: u.UsageReportTimezone,
		})
	}
	return result, nil
}

// GetActiveUserIDs returns user IDs who had activity on the given date
func (r *userUsageReportRepository) GetActiveUserIDs(ctx context.Context, date time.Time) ([]int64, error) {
	// Calculate start and end of day in the configured server timezone
	// This ensures consistency with how usage_logs timestamps are stored and queried
	startOfDay := timezone.StartOfDay(date)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Query distinct user_ids from usage_logs for the date
	rows, err := r.sql.QueryContext(ctx, `
		SELECT DISTINCT user_id
		FROM usage_logs
		WHERE created_at >= $1 AND created_at < $2
	`, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, rows.Err()
}
