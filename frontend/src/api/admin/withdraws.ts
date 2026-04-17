import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'
import type { WithdrawRequest } from '../withdraw'

export async function listWithdraws(
  page: number = 1,
  pageSize: number = 20,
  filters?: { source_type?: string; status?: string }
): Promise<PaginatedResponse<WithdrawRequest>> {
  const { data } = await apiClient.get<PaginatedResponse<WithdrawRequest>>('/admin/withdraws', {
    params: {
      page,
      page_size: pageSize,
      source_type: filters?.source_type || undefined,
      status: filters?.status || undefined
    }
  })
  return data
}

export async function reviewWithdraw(
  id: number,
  payload: { approve: boolean; review_note?: string }
): Promise<WithdrawRequest> {
  const { data } = await apiClient.post<WithdrawRequest>(`/admin/withdraws/${id}/review`, payload)
  return data
}

export async function payWithdraw(id: number): Promise<WithdrawRequest> {
  const { data } = await apiClient.post<WithdrawRequest>(`/admin/withdraws/${id}/pay`)
  return data
}

export const adminWithdrawsAPI = {
  listWithdraws,
  reviewWithdraw,
  payWithdraw,
}

export default adminWithdrawsAPI
