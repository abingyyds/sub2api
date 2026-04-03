/**
 * Admin Users API endpoints
 * Handles user management for administrators
 */

import { apiClient } from '../client'
import type { AdminUser, UpdateUserRequest, PaginatedResponse } from '@/types'

/**
 * List all users with pagination
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters (status, role, search, attributes)
 * @param options - Optional request options (signal)
 * @returns Paginated list of users
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    status?: 'active' | 'disabled'
    role?: 'admin' | 'user'
    search?: string
    attributes?: Record<number, string>  // attributeId -> value
    created_after?: string  // YYYY-MM-DD
    created_before?: string // YYYY-MM-DD
  },
  options?: {
    signal?: AbortSignal
  }
): Promise<PaginatedResponse<AdminUser>> {
  // Build params with attribute filters in attr[id]=value format
  const params: Record<string, any> = {
    page,
    page_size: pageSize,
    status: filters?.status,
    role: filters?.role,
    search: filters?.search,
    created_after: filters?.created_after,
    created_before: filters?.created_before
  }

  // Add attribute filters as attr[id]=value
  if (filters?.attributes) {
    for (const [attrId, value] of Object.entries(filters.attributes)) {
      if (value) {
        params[`attr[${attrId}]`] = value
      }
    }
  }
  const { data } = await apiClient.get<PaginatedResponse<AdminUser>>('/admin/users', {
    params,
    signal: options?.signal
  })
  return data
}

/**
 * Get user by ID
 * @param id - User ID
 * @returns User details
 */
export async function getById(id: number): Promise<AdminUser> {
  const { data } = await apiClient.get<AdminUser>(`/admin/users/${id}`)
  return data
}

/**
 * Create new user
 * @param userData - User data (email, password, etc.)
 * @returns Created user
 */
export async function create(userData: {
  email: string
  password: string
  balance?: number
  concurrency?: number
  allowed_groups?: number[] | null
}): Promise<AdminUser> {
  const { data } = await apiClient.post<AdminUser>('/admin/users', userData)
  return data
}

/**
 * Update user
 * @param id - User ID
 * @param updates - Fields to update
 * @returns Updated user
 */
export async function update(id: number, updates: UpdateUserRequest): Promise<AdminUser> {
  const { data } = await apiClient.put<AdminUser>(`/admin/users/${id}`, updates)
  return data
}

/**
 * Delete user
 * @param id - User ID
 * @returns Success confirmation
 */
export async function deleteUser(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/users/${id}`)
  return data
}

/**
 * Update user balance
 * @param id - User ID
 * @param balance - New balance
 * @param operation - Operation type ('set', 'add', 'subtract')
 * @param notes - Optional notes for the balance adjustment
 * @returns Updated user
 */
export async function updateBalance(
  id: number,
  balance: number,
  operation: 'set' | 'add' | 'subtract' = 'set',
  notes?: string
): Promise<AdminUser> {
  const { data } = await apiClient.post<AdminUser>(`/admin/users/${id}/balance`, {
    balance,
    operation,
    notes: notes || ''
  })
  return data
}

/**
 * Update user concurrency
 * @param id - User ID
 * @param concurrency - New concurrency limit
 * @returns Updated user
 */
export async function updateConcurrency(id: number, concurrency: number): Promise<AdminUser> {
  return update(id, { concurrency })
}

/**
 * Toggle user status
 * @param id - User ID
 * @param status - New status
 * @returns Updated user
 */
export async function toggleStatus(id: number, status: 'active' | 'disabled'): Promise<AdminUser> {
  return update(id, { status })
}

/**
 * Get user's API keys
 * @param id - User ID
 * @returns List of user's API keys
 */
export async function getUserApiKeys(id: number): Promise<PaginatedResponse<any>> {
  const { data } = await apiClient.get<PaginatedResponse<any>>(`/admin/users/${id}/api-keys`)
  return data
}

/**
 * Get user's usage statistics
 * @param id - User ID
 * @param period - Time period
 * @returns User usage statistics
 */
export async function getUserUsageStats(
  id: number,
  period: string = 'month'
): Promise<{
  total_requests: number
  total_cost: number
  total_tokens: number
}> {
  const { data } = await apiClient.get<{
    total_requests: number
    total_cost: number
    total_tokens: number
  }>(`/admin/users/${id}/usage`, {
    params: { period }
  })
  return data
}

/**
 * Get list of agents (users with is_agent=true)
 * @returns List of agent users
 */
export async function getAgents(): Promise<AdminUser[]> {
  const { data } = await apiClient.get<PaginatedResponse<AdminUser>>('/admin/users', {
    params: {
      page: 1,
      page_size: 1000,
      // Note: Backend needs to support filtering by is_agent, for now we get all and filter client-side
    }
  })
  console.log('All users:', data.items.length)
  console.log('Sample user:', data.items[0])
  // Filter agents client-side (backend should ideally support this filter)
  const agents = data.items.filter((u: any) => u.is_agent === true)
  console.log('Filtered agents:', agents.length, agents)
  return agents
}

/**
 * Get user's current inviter/referrer
 * @param id - User ID
 * @returns Inviter user or null
 */
export async function getUserInviter(id: number): Promise<AdminUser | null> {
  // TODO: Backend needs to implement /admin/users/:id/inviter endpoint
  // For now, we'll need to get this from the referral API or user details
  // Returning null as a temporary workaround
  return null
}

export const usersAPI = {
  list,
  getById,
  create,
  update,
  delete: deleteUser,
  updateBalance,
  updateConcurrency,
  toggleStatus,
  getUserApiKeys,
  getUserUsageStats,
  getAgents,
  getUserInviter
}

export default usersAPI
