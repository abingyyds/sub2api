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
  alipay_trade_no: string | null
  epay_trade_no: string | null
  invoice_company_name: string
  invoice_tax_id: string
  invoice_email: string
  invoice_remark: string
  invoice_requested_at: string | null
  invoice_processed_at: string | null
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

export async function markInvoiceProcessed(id: number): Promise<{ success: boolean }> {
  const { data } = await apiClient.post<{ success: boolean }>(`/admin/orders/${id}/invoice/processed`)
  return data
}

export async function repair(id: number): Promise<{ success: boolean }> {
  const { data } = await apiClient.post<{ success: boolean }>(`/admin/orders/${id}/repair`)
  return data
}

export const ordersAPI = { list, markInvoiceProcessed, repair }
export default ordersAPI
