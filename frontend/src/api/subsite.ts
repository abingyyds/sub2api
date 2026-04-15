import { apiClient } from './client'
import type { CreateSubSiteActivationInput, SubSiteGroupPriceOverride, SubSiteRechargePriceOverride } from './payment'

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
  default_custom_config?: string
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
  enable_topup: boolean
  allow_sub_site: boolean
  sub_site_price_fen: number
  subscription_expired_at?: string | null
  user_count?: number
  child_site_count?: number
  entry_url?: string
  group_price_overrides?: SubSiteGroupPriceOverride[]
  recharge_price_overrides?: SubSiteRechargePriceOverride[]
  created_at: string
  updated_at: string
}

export type UpdateOwnedSubSiteRequest = CreateSubSiteActivationInput

export async function getOpenInfo(): Promise<SubSiteOpenInfo> {
  const { data } = await apiClient.get<SubSiteOpenInfo>('/subsite/open-info')
  return data
}

export async function listOwnedSites(): Promise<OwnedSubSite[]> {
  const { data } = await apiClient.get<OwnedSubSite[]>('/subsite/owned')
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

export const subSiteAPI = {
  getOpenInfo,
  listOwnedSites,
  getOwnedSite,
  updateOwnedSite,
}

export default subSiteAPI
