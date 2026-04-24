import { apiClient } from './client'
import type { CreateSubSiteActivationInput } from './payment'
import type { PaginatedResponse } from '@/types'

export interface SubSiteThemeTemplateOption {
  key: string
  label: string
  description: string
}

export interface SubSiteOpenInfo {
  enabled: boolean
  scope: 'platform' | 'subsite'
  parent_sub_site_id?: number
  parent_sub_site_name?: string
  level: number
  max_level: number
  price_fen: number
  validity_days: number
  currency: string
  allow_custom_domain: boolean
  default_theme_template: string
  theme_templates: SubSiteThemeTemplateOption[]
}

export interface OwnedSubSite {
  id: number
  owner_user_id: number
  owner_email?: string
  parent_sub_site_id?: number
  parent_sub_site_name?: string
  level: number
  name: string
  slug: string
  custom_domain?: string
  status: 'pending' | 'active' | 'disabled'
  mode: 'pool' | 'rate'
  site_logo?: string
  site_favicon?: string
  site_subtitle?: string
  announcement?: string
  contact_info?: string
  doc_url?: string
  home_content?: string
  pending_home_content?: string
  home_content_review_status?: 'none' | 'pending' | 'approved' | 'rejected'
  home_content_review_note?: string
  home_content_submitted_at?: string
  home_content_reviewed_at?: string
  home_content_reviewed_by?: number
  theme_template?: string
  registration_mode?: 'open' | 'invite' | 'closed'
  enable_topup: boolean
  allow_sub_site: boolean
  sub_site_price_fen: number
  consume_rate_multiplier?: number
  subscription_expired_at?: string | null
  user_count?: number
  child_site_count?: number
  entry_url?: string
  balance_fen?: number
  total_topup_fen?: number
  total_consumed_fen?: number
  total_withdrawn_fen?: number
  allow_online_topup?: boolean
  allow_offline_topup?: boolean
  owner_payment_config?: OwnerPaymentConfig | null
  created_at: string
  updated_at: string
}

export interface OwnerPaymentConfig {
  wechat?: { enabled: boolean; app_id: string; mch_id: string; apiv3_key: string; mch_serial_no: string; public_key_id: string; public_key: string; private_key: string; notify_url: string }
  alipay?: { enabled: boolean; app_id: string; private_key: string; public_key: string; notify_url: string; is_production: boolean }
  epay?: { enabled: boolean; gateway: string; pid: string; pkey: string; notify_url: string }
}

export interface OwnedSubSiteLedgerEntry {
  id: number
  sub_site_id: number
  tx_type: string
  delta_fen: number
  balance_after_fen: number
  related_user_id?: number
  related_usage_log_id?: number
  related_order_id?: number
  operator_id?: number
  note?: string
  created_at: string
}

export type UpdateOwnedSubSiteRequest = CreateSubSiteActivationInput

export async function getOpenInfo(): Promise<SubSiteOpenInfo> {
  const { data } = await apiClient.get<SubSiteOpenInfo>('/subsite/open-info')
  return {
    ...data,
    theme_templates: Array.isArray(data?.theme_templates) ? data.theme_templates : [],
  }
}

export async function listOwnedSites(): Promise<OwnedSubSite[]> {
  const { data } = await apiClient.get<OwnedSubSite[] | null>('/subsite/owned')
  if (!Array.isArray(data)) {
    return []
  }
  return data
}

export async function getOwnedSite(id: number): Promise<OwnedSubSite> {
  const { data } = await apiClient.get<OwnedSubSite>(`/subsite/owned/${id}`)
  return data
}

export async function updateOwnedSite(id: number, payload: UpdateOwnedSubSiteRequest): Promise<OwnedSubSite> {
  const { data } = await apiClient.put<OwnedSubSite>(`/subsite/owned/${id}`, payload)
  return data
}

export async function offlineTopupUser(
  siteID: number,
  payload: { user_id: number; amount_fen: number; note?: string }
): Promise<OwnedSubSite> {
  const { data } = await apiClient.post<OwnedSubSite>(`/subsite/owned/${siteID}/offline-topup`, payload)
  return data
}

export async function listOwnedLedger(
  siteID: number,
  page: number = 1,
  pageSize: number = 20,
  txType?: string
): Promise<PaginatedResponse<OwnedSubSiteLedgerEntry>> {
  const { data } = await apiClient.get<PaginatedResponse<OwnedSubSiteLedgerEntry>>(
    `/subsite/owned/${siteID}/ledger`,
    {
      params: {
        page,
        page_size: pageSize,
        tx_type: txType || undefined
      }
    }
  )
  return data
}

export async function createPoolTopupOrder(payload: {
  site_id: number
  amount_fen: number
  pay_method?: string
}): Promise<{
  order_no: string
  code_url?: string
  amount_fen: number
  expired_at: string
}> {
  const { data } = await apiClient.post('/payment/subsite-topup', payload)
  return data as any
}

export const subSiteAPI = {
  getOpenInfo,
  listOwnedSites,
  getOwnedSite,
  updateOwnedSite,
  offlineTopupUser,
  listOwnedLedger,
  createPoolTopupOrder,
}

export default subSiteAPI
