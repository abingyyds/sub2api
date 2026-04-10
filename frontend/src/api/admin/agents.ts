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
  real_name: string
  phone: string
  identity_status: string
  contract_status: string
  activation_fee_paid_at: string | null
  is_frozen: boolean
  frozen_reason: string
  frozen_balance: number
  withdrawable_balance: number
  total_withdrawn: number
  sub_user_count: number
  total_commission: number
  pending_commission: number
  created_at: string
}

export interface AdminAgentDetail extends AdminAgent {
  id_card_no: string
  identity_submitted_at: string | null
  contract_version: string
  contract_signed_at: string | null
  contract_ip: string
  contract_signature_data?: string
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

export async function getDetail(id: number): Promise<AdminAgentDetail> {
  const { data } = await apiClient.get<AdminAgentDetail>(`/admin/agents/${id}`)
  return data
}

export async function reject(id: number, reason?: string): Promise<void> {
  await apiClient.post(`/admin/agents/${id}/reject`, { reason })
}

export async function update(id: number, data: { commission_rate: number }): Promise<void> {
  await apiClient.put(`/admin/agents/${id}`, data)
}

export async function setFrozen(id: number, frozen: boolean, reason?: string): Promise<void> {
  await apiClient.post(`/admin/agents/${id}/freeze`, { frozen, reason: reason || '' })
}

export const agentsAPI = { list, getDetail, approve, reject, update, setFrozen }
export default agentsAPI
