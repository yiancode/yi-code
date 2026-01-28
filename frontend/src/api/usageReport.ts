/**
 * User Usage Report API endpoints
 * Handles user-level usage report email configuration
 */

import { apiClient } from './client'

export interface UsageReportConfig {
  enabled: boolean
  schedule: string // HH:MM format
  timezone: string
  global_enabled: boolean
  email_bound: boolean
}

export interface UpdateUsageReportConfigRequest {
  enabled?: boolean
  schedule?: string
  timezone?: string
}

/**
 * Get current user's usage report configuration
 * @returns Usage report config including enabled state and schedule
 */
export async function getUsageReportConfig(): Promise<UsageReportConfig> {
  const { data } = await apiClient.get<UsageReportConfig>('/usage-report/config')
  return data
}

/**
 * Update user's usage report configuration
 * @param request - Fields to update
 * @returns Updated configuration
 */
export async function updateUsageReportConfig(request: UpdateUsageReportConfigRequest): Promise<UsageReportConfig> {
  const { data } = await apiClient.put<UsageReportConfig>('/usage-report/config', request)
  return data
}

/**
 * Send a test usage report email to the user
 * @returns Success response
 */
export async function sendTestReport(): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>('/usage-report/test')
  return data
}

export const usageReportAPI = {
  getConfig: getUsageReportConfig,
  updateConfig: updateUsageReportConfig,
  sendTestReport
}

export default usageReportAPI
