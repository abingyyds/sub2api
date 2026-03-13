/**
 * Payment API endpoints
 * Handles WeChat Pay payment operations
 */

import { apiClient } from './client'

export interface PaymentPlan {
  key: string
  name: string
  description?: string
  features?: string[]
  amount_fen: number
  group_id: number
  validity_days: number
  type?: 'subscription' | 'balance'
  balance_amount?: number
}

export interface RechargePlan {
  key: string
  name: string
  description?: string
  pay_amount_fen: number
  balance_amount: number
  popular?: boolean
}

export interface RechargeInfo {
  min_amount: number
  plans: RechargePlan[]
}

export interface CreateOrderResponse {
  order_no: string
  code_url: string | null
  amount_fen: number
  discount_amount: number
  expired_at: string
}

export interface PaymentOrder {
  order_no: string
  plan_key: string
  amount_fen: number
  promo_code: string
  discount_amount: number
  status: 'pending' | 'paid' | 'closed' | 'refunded'
  pay_method: string
  code_url?: string | null
  paid_at: string | null
  expired_at: string
  created_at: string
}

/**
 * Get available payment plans
 */
export async function getPlans(): Promise<PaymentPlan[]> {
  const { data } = await apiClient.get<PaymentPlan[]>('/payment/plans')
  return data
}

/**
 * Get recharge info (plans + min amount)
 */
export async function getRechargeInfo(): Promise<RechargeInfo> {
  const { data } = await apiClient.get<RechargeInfo>('/payment/recharge-info')
  return data
}

/**
 * Create a payment order
 */
export async function createOrder(planKey: string, promoCode?: string, payMethod?: 'wechat' | 'alipay'): Promise<CreateOrderResponse> {
  const { data } = await apiClient.post<CreateOrderResponse>('/payment/orders', {
    plan_key: planKey,
    promo_code: promoCode || '',
    pay_method: payMethod || 'wechat',
  })
  return data
}

/**
 * Create a balance recharge order with custom amount
 */
export async function createRechargeOrder(amount: number, promoCode?: string, payMethod?: 'wechat' | 'alipay'): Promise<CreateOrderResponse> {
  const { data } = await apiClient.post<CreateOrderResponse>('/payment/recharge', {
    amount,
    promo_code: promoCode || '',
    pay_method: payMethod || 'wechat',
  })
  return data
}

/**
 * Query order status
 */
export async function queryOrder(orderNo: string): Promise<PaymentOrder> {
  const { data } = await apiClient.get<PaymentOrder>(`/payment/orders/${orderNo}`)
  return data
}

/**
 * List user's payment orders
 */
export async function listOrders(params?: {
  page?: number
  page_size?: number
}): Promise<{
  items: PaymentOrder[]
  total: number
  page: number
  page_size: number
  pages: number
}> {
  const { data } = await apiClient.get('/payment/orders', { params })
  return data
}

export const paymentAPI = {
  getPlans,
  getRechargeInfo,
  createOrder,
  createRechargeOrder,
  queryOrder,
  listOrders,
}

export default paymentAPI
