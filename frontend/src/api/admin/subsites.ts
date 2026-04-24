import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'

export interface AdminSubSite {
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
  created_at: string
  updated_at: string
}

export interface SaveSubSiteRequest {
  owner_user_id: number
  parent_sub_site_id?: number
  name: string
  slug: string
  custom_domain?: string
  status: 'pending' | 'active' | 'disabled'
  site_logo?: string
  site_favicon?: string
  site_subtitle?: string
  announcement?: string
  contact_info?: string
  doc_url?: string
  home_content?: string
  theme_template?: string
  registration_mode?: 'open' | 'invite' | 'closed'
  enable_topup?: boolean
  allow_sub_site?: boolean
  sub_site_price_fen?: number
  consume_rate_multiplier?: number
  subscription_expired_at?: string | null
  allow_online_topup?: boolean
  allow_offline_topup?: boolean
}

export interface SubSiteLedgerEntry {
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

export interface PlatformSubSiteConfig {
  entry_enabled: boolean
  enabled: boolean
  activation_price_fen: number
  validity_days: number
  max_level: number
  default_theme_template: string
  theme_templates: Array<{ key: string; label: string; description: string }>
}

export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: { search?: string; status?: string }
): Promise<PaginatedResponse<AdminSubSite>> {
  const { data } = await apiClient.get<PaginatedResponse<AdminSubSite>>('/admin/subsites', {
    params: {
      page,
      page_size: pageSize,
      search: filters?.search,
      status: filters?.status || undefined
    }
  })
  return data
}

export async function create(payload: SaveSubSiteRequest): Promise<AdminSubSite> {
  const { data } = await apiClient.post<AdminSubSite>('/admin/subsites', payload)
  return data
}

export async function update(id: number, payload: SaveSubSiteRequest): Promise<AdminSubSite> {
  const { data } = await apiClient.put<AdminSubSite>(`/admin/subsites/${id}`, payload)
  return data
}

export async function remove(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/subsites/${id}`)
  return data
}

export async function getPlatformConfig(): Promise<PlatformSubSiteConfig> {
  const { data } = await apiClient.get<PlatformSubSiteConfig>('/admin/subsites/platform-config')
  return data
}

export async function updatePlatformConfig(payload: PlatformSubSiteConfig): Promise<PlatformSubSiteConfig> {
  const { data } = await apiClient.put<PlatformSubSiteConfig>('/admin/subsites/platform-config', payload)
  return data
}

export async function topupPool(
  id: number,
  amountFen: number,
  note: string = ''
): Promise<AdminSubSite> {
  const { data } = await apiClient.post<AdminSubSite>(`/admin/subsites/${id}/topup`, {
    amount_fen: amountFen,
    note
  })
  return data
}

export async function listLedger(
  id: number,
  page: number = 1,
  pageSize: number = 20,
  txType?: string
): Promise<PaginatedResponse<SubSiteLedgerEntry>> {
  const { data } = await apiClient.get<PaginatedResponse<SubSiteLedgerEntry>>(
    `/admin/subsites/${id}/ledger`,
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

export async function setMode(id: number, mode: 'pool' | 'rate'): Promise<AdminSubSite> {
  const { data } = await apiClient.put<AdminSubSite>(`/admin/subsites/${id}/mode`, { mode })
  return data
}

export async function reviewHomeContent(
  id: number,
  approve: boolean,
  reviewNote: string = ''
): Promise<AdminSubSite> {
  const { data } = await apiClient.post<AdminSubSite>(`/admin/subsites/${id}/home-content/review`, {
    approve,
    review_note: reviewNote
  })
  return data
}

export const subSitesAPI = {
  list,
  create,
  update,
  remove,
  getPlatformConfig,
  updatePlatformConfig,
  topupPool,
  listLedger,
  setMode,
  reviewHomeContent,
}

export default subSitesAPI
