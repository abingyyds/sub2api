import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'

export interface AdminPaymentOrder {
  id: number
  order_no: string
  user_id: number
  plan_key: string
  group_id: number
  amount_fen: number
  validity_days: number
  order_type: string
  balance_amount: number
  status: string
  pay_method: string
  wechat_transaction_id: string | null
  paid_at: string | null
  expired_at: string
  created_at: string
}

export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: { status?: string; order_type?: string },
  options?: { signal?: AbortSignal }
): Promise<PaginatedResponse<AdminPaymentOrder>> {
  const { data } = await apiClient.get<PaginatedResponse<AdminPaymentOrder>>('/admin/orders', {
    params: {
      page,
      page_size: pageSize,
      status: filters?.status || undefined,
      order_type: filters?.order_type || undefined
    },
    signal: options?.signal
  })
  return data
}

export const ordersAPI = { list }
export default ordersAPI
