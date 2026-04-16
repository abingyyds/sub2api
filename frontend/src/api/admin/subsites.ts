import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'
import type { SubSiteGroupPriceOverride, SubSiteRechargePriceOverride } from '../payment'

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
  theme_config?: string
  custom_config?: string
  registration_mode?: 'open' | 'invite' | 'closed'
  enable_topup?: boolean
  allow_sub_site?: boolean
  sub_site_price_fen?: number
  subscription_expired_at?: string | null
  group_price_overrides?: SubSiteGroupPriceOverride[]
  recharge_price_overrides?: SubSiteRechargePriceOverride[]
}

export interface PlatformSubSiteConfig {
  entry_enabled: boolean
  enabled: boolean
  activation_price_fen: number
  validity_days: number
  default_theme_template: string
  default_custom_config?: string
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

export const subSitesAPI = {
  list,
  create,
  update,
  remove,
  getPlatformConfig,
  updatePlatformConfig,
}

export default subSitesAPI
