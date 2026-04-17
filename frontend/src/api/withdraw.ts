import { apiClient } from './client'
import type { PaginatedResponse } from '@/types'

export interface WithdrawRequest {
  id: number
  user_id: number
  user_email?: string
  amount: number
  alipay_name: string
  alipay_account: string
  alipay_qr_image?: string
  source_type: 'agent_commission' | 'sub_site_profit'
  source_sub_site_id?: number
  status: 'pending' | 'approved' | 'rejected' | 'paid' | 'cancelled'
  review_note?: string
  reviewed_at?: string
  reviewer_id?: number
  paid_at?: string
  paid_by?: number
  created_at: string
  updated_at: string
}

export interface ApplyWithdrawRequest {
  amount: number
  alipay_name: string
  alipay_account: string
  alipay_qr_image?: string
  source_type: 'agent_commission' | 'sub_site_profit'
  source_sub_site_id?: number
}

export async function applyWithdraw(payload: ApplyWithdrawRequest): Promise<WithdrawRequest> {
  const { data } = await apiClient.post<WithdrawRequest>('/withdraw/apply', payload)
  return data
}

export async function cancelWithdraw(id: number): Promise<WithdrawRequest> {
  const { data } = await apiClient.post<WithdrawRequest>(`/withdraw/${id}/cancel`)
  return data
}

export async function listWithdraws(
  page: number = 1,
  pageSize: number = 20,
  filters?: { source_type?: string; status?: string }
): Promise<PaginatedResponse<WithdrawRequest>> {
  const { data } = await apiClient.get<PaginatedResponse<WithdrawRequest>>('/withdraw/list', {
    params: {
      page,
      page_size: pageSize,
      source_type: filters?.source_type || undefined,
      status: filters?.status || undefined
    }
  })
  return data
}

export const withdrawAPI = {
  applyWithdraw,
  cancelWithdraw,
  listWithdraws,
}

export default withdrawAPI
