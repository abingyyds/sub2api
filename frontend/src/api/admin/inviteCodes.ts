import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'

export interface AdminInviteCode {
  id: number
  code: string
  source_name: string
  created_by: number
  used_count: number
  max_uses: number | null
  enabled: boolean
  notes: string
  created_at: string
  updated_at: string
}

export interface CreateInviteCodeRequest {
  source_name: string
  max_uses?: number | null
  notes?: string
}

export interface UpdateInviteCodeRequest {
  source_name?: string
  max_uses?: number | null
  enabled?: boolean
  notes?: string
}

export async function list(
  page: number = 1,
  pageSize: number = 20
): Promise<PaginatedResponse<AdminInviteCode>> {
  const { data } = await apiClient.get<PaginatedResponse<AdminInviteCode>>('/admin/invite-codes', {
    params: { page, page_size: pageSize }
  })
  return data
}

export async function create(req: CreateInviteCodeRequest): Promise<AdminInviteCode> {
  const { data } = await apiClient.post<AdminInviteCode>('/admin/invite-codes', req)
  return data
}

export async function update(id: number, req: UpdateInviteCodeRequest): Promise<AdminInviteCode> {
  const { data } = await apiClient.put<AdminInviteCode>(`/admin/invite-codes/${id}`, req)
  return data
}

export async function remove(id: number): Promise<void> {
  await apiClient.delete(`/admin/invite-codes/${id}`)
}

export const inviteCodesAPI = { list, create, update, remove }
export default inviteCodesAPI
