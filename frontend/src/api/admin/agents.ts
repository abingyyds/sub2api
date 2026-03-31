import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'

export interface AdminAgent {
  id: number
  email: string
  username: string
  is_agent: boolean
  agent_status: string
  agent_commission_rate: number
  agent_note: string
  agent_approved_at: string | null
  sub_user_count: number
  total_commission: number
  pending_commission: number
  created_at: string
}

export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: { status?: string; search?: string },
  options?: { signal?: AbortSignal }
): Promise<PaginatedResponse<AdminAgent>> {
  const { data } = await apiClient.get<PaginatedResponse<AdminAgent>>('/admin/agents', {
    params: {
      page,
      page_size: pageSize,
      status: filters?.status || undefined,
      search: filters?.search || undefined
    },
    signal: options?.signal
  })
  return data
}

export async function approve(id: number): Promise<void> {
  await apiClient.post(`/admin/agents/${id}/approve`)
}

export async function reject(id: number, reason?: string): Promise<void> {
  await apiClient.post(`/admin/agents/${id}/reject`, { reason })
}

export async function update(id: number, data: { commission_rate: number }): Promise<void> {
  await apiClient.put(`/admin/agents/${id}`, data)
}

export async function settle(id: number): Promise<{ settled_count: number; settled_amount: number }> {
  const { data } = await apiClient.post<{ settled_count: number; settled_amount: number }>(
    `/admin/agents/${id}/settle`
  )
  return data
}

export async function updateParent(userId: number, parentId: number): Promise<void> {
  await apiClient.put(`/admin/agents/${userId}/parent`, { parent_id: parentId })
}

export const agentsAPI = { list, approve, reject, update, settle, updateParent }
export default agentsAPI
