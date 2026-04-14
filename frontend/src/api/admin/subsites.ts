import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'

export interface AdminSubSite {
  id: number
  owner_user_id: number
  owner_email?: string
  name: string
  slug: string
  custom_domain?: string
  status: 'active' | 'disabled'
  site_logo?: string
  site_favicon?: string
  site_subtitle?: string
  announcement?: string
  contact_info?: string
  doc_url?: string
  home_content?: string
  theme_config?: string
  user_count?: number
  entry_url?: string
  created_at: string
  updated_at: string
}

export interface SaveSubSiteRequest {
  owner_user_id: number
  name: string
  slug: string
  custom_domain?: string
  status: 'active' | 'disabled'
  site_logo?: string
  site_favicon?: string
  site_subtitle?: string
  announcement?: string
  contact_info?: string
  doc_url?: string
  home_content?: string
  theme_config?: string
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

export const subSitesAPI = {
  list,
  create,
  update,
  remove
}

export default subSitesAPI
