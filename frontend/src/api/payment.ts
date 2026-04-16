/**
 * Payment API endpoints
 * Handles WeChat Pay payment operations
 */

import { apiClient } from './client'

const CACHE_TTL_MS = 60_000

interface CacheEntry<T> {
  value: T | null
  fetchedAt: number
  promise: Promise<T> | null
}

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
  is_newcomer?: boolean
  max_purchases?: number
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
  invoice_company_name: string
  invoice_tax_id: string
  invoice_email: string
  invoice_remark: string
  invoice_requested_at: string | null
  invoice_processed_at: string | null
}

export type PayMethod = 'wechat' | 'alipay' | 'epay_alipay' | 'epay_wxpay'

export interface CreateSubSiteActivationInput {
  name: string
  slug: string
  custom_domain?: string
  site_logo?: string
  site_favicon?: string
  site_subtitle?: string
  announcement?: string
  contact_info?: string
  doc_url?: string
  home_content?: string
  theme_template?: string
  theme_config?: string
  custom_config?: string
  registration_mode?: 'open' | 'invite' | 'closed'
  enable_topup?: boolean
  allow_sub_site?: boolean
  sub_site_price_fen?: number
  consume_rate_multiplier?: number
}

const plansCache: CacheEntry<PaymentPlan[]> = {
  value: null,
  fetchedAt: 0,
  promise: null
}

const rechargeInfoCache: CacheEntry<RechargeInfo> = {
  value: null,
  fetchedAt: 0,
  promise: null
}

const payMethodsCache: CacheEntry<PayMethod[]> = {
  value: null,
  fetchedAt: 0,
  promise: null
}

const newcomerStatusCache: CacheEntry<{ eligible: boolean }> = {
  value: null,
  fetchedAt: 0,
  promise: null
}

function hasFreshCache<T>(entry: CacheEntry<T>): entry is CacheEntry<T> & { value: T } {
  return entry.value !== null && Date.now() - entry.fetchedAt < CACHE_TTL_MS
}

function withCache<T>(entry: CacheEntry<T>, loader: () => Promise<T>): Promise<T> {
  if (hasFreshCache(entry)) {
    return Promise.resolve(entry.value)
  }

  if (entry.promise) {
    return entry.promise
  }

  const request = loader()
    .then((data) => {
      entry.value = data
      entry.fetchedAt = Date.now()
      return data
    })
    .finally(() => {
      if (entry.promise === request) {
        entry.promise = null
      }
    })

  entry.promise = request

  return request
}

/**
 * Get available payment methods
 */
export async function getPayMethods(): Promise<PayMethod[]> {
  return withCache(payMethodsCache, async () => {
    const { data } = await apiClient.get<{ methods: PayMethod[] }>('/payment/methods')
    return data.methods
  })
}

/**
 * Get available payment plans
 */
export async function getPlans(): Promise<PaymentPlan[]> {
  return withCache(plansCache, async () => {
    const { data } = await apiClient.get<PaymentPlan[]>('/payment/plans')
    return data
  })
}

/**
 * Get recharge info (plans + min amount)
 */
export async function getRechargeInfo(): Promise<RechargeInfo> {
  return withCache(rechargeInfoCache, async () => {
    const { data } = await apiClient.get<RechargeInfo>('/payment/recharge-info')
    return data
  })
}

/**
 * Create a payment order
 */
export async function createOrder(planKey: string, promoCode?: string, payMethod?: PayMethod): Promise<CreateOrderResponse> {
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
export async function createRechargeOrder(amount: number, promoCode?: string, payMethod?: PayMethod, planKey?: string): Promise<CreateOrderResponse> {
  const { data } = await apiClient.post<CreateOrderResponse>('/payment/recharge', {
    amount,
    promo_code: promoCode || '',
    pay_method: payMethod || 'wechat',
    plan_key: planKey || '',
  })
  return data
}

/**
 * Create an agent activation fee order.
 */
export async function createAgentActivationOrder(payMethod?: PayMethod): Promise<CreateOrderResponse> {
  const { data } = await apiClient.post<CreateOrderResponse>('/payment/agent-activation', {
    pay_method: payMethod || 'wechat',
  })
  return data
}

export async function createSubSiteActivationOrder(
  activationInput: CreateSubSiteActivationInput,
  payMethod?: PayMethod
): Promise<CreateOrderResponse> {
  const { data } = await apiClient.post<CreateOrderResponse>('/payment/subsite-activation', {
    pay_method: payMethod || 'wechat',
    activation_input: activationInput,
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

export async function submitInvoiceRequest(payload: {
  order_nos: string[]
  company_name: string
  tax_id: string
  email: string
  remark?: string
}): Promise<{ success: boolean }> {
  const { data } = await apiClient.post<{ success: boolean }>('/payment/invoice-requests', payload)
  return data
}

/**
 * Check if user is eligible for newcomer plans
 */
export async function getNewcomerStatus(): Promise<{ eligible: boolean }> {
  return withCache(newcomerStatusCache, async () => {
    const { data } = await apiClient.get<{ eligible: boolean }>('/payment/newcomer-status')
    return data
  })
}

export const paymentAPI = {
  getPlans,
  getRechargeInfo,
  getPayMethods,
  createOrder,
  createRechargeOrder,
  createAgentActivationOrder,
  createSubSiteActivationOrder,
  queryOrder,
  listOrders,
  submitInvoiceRequest,
  getNewcomerStatus,
}

export default paymentAPI
