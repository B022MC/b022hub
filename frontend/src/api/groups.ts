/**
 * User Groups API endpoints (non-admin)
 * Handles group-related operations for regular users
 */

import { apiClient } from './client'
import type { Group } from '@/types'

export interface AvailableModelsGroup {
  group: Group
  models: string[]
  model_count: number
}

export interface AvailableModelsResponse {
  groups: AvailableModelsGroup[]
}

/**
 * Get available groups that the current user can bind to API keys
 * This returns groups based on user's permissions:
 * - Standard groups: public (non-exclusive) or explicitly allowed
 * - Subscription groups: user has active subscription
 * @returns List of available groups
 */
export async function getAvailable(): Promise<Group[]> {
  const { data } = await apiClient.get<Group[]>('/groups/available')
  return data
}

/**
 * Get current user's custom group rate multipliers
 * @returns Map of group_id to custom rate_multiplier
 */
export async function getUserGroupRates(): Promise<Record<number, number>> {
  const { data } = await apiClient.get<Record<number, number> | null>('/groups/rates')
  return data || {}
}

/**
 * Get available models for all groups accessible to the current user
 */
export async function getModels(): Promise<AvailableModelsResponse> {
  const { data } = await apiClient.get<AvailableModelsResponse>('/groups/models')
  return data
}

export const userGroupsAPI = {
  getAvailable,
  getUserGroupRates,
  getModels
}

export default userGroupsAPI
