import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'
import type { OwnedSubSite, OwnedSubSiteLedgerEntry, UpdateOwnedSubSiteRequest } from '../subsite'

export interface SubSiteAdminDashboardStats {
  user_count: number
  active_users: number
  requests: number
  total_cost: number
  revenue_fen: number
  pool_balance_fen: number
  total_topup_fen: number
  total_consumed_fen: number
  range_start: string
}

export interface SubSiteAdminUser {
  id: number
  email: string
  username?: string
  role: string
  status: string
  balance: number
  bind_source?: string
  created_at: string
  bound_at: string
}

export interface SubSiteAdminOrder {
  id: number
  order_no: string
  user_id: number
  user_email?: string
  plan_key?: string
  order_type: string
  amount_fen: number
  status: string
  pay_method?: string
  paid_at?: string | null
  created_at: string
}

export interface SubSiteAdminUsage {
  id: number
  user_id: number
  user_email?: string
  model: string
  input_tokens: number
  output_tokens: number
  cache_creation_tokens: number
  cache_read_tokens: number
  total_cost: number
  actual_cost: number
  created_at: string
}

const base = (siteID: number) => `/subsite/admin/${siteID}`

export async function getDashboard(siteID: number): Promise<SubSiteAdminDashboardStats> {
  const { data } = await apiClient.get<SubSiteAdminDashboardStats>(`${base(siteID)}/dashboard`)
  return data
}

export async function getSite(siteID: number): Promise<OwnedSubSite> {
  const { data } = await apiClient.get<OwnedSubSite>(`${base(siteID)}/site`)
  return data
}

export async function updateSite(siteID: number, payload: UpdateOwnedSubSiteRequest): Promise<OwnedSubSite> {
  const { data } = await apiClient.put<OwnedSubSite>(`${base(siteID)}/site`, payload)
  return data
}

export async function listUsers(
  siteID: number,
  page = 1,
  pageSize = 20,
  search?: string,
): Promise<PaginatedResponse<SubSiteAdminUser>> {
  const { data } = await apiClient.get<PaginatedResponse<SubSiteAdminUser>>(`${base(siteID)}/users`, {
    params: { page, page_size: pageSize, search: search || undefined },
  })
  return data
}

export async function offlineTopup(
  siteID: number,
  payload: { user_id: number; amount_fen: number; note?: string },
): Promise<OwnedSubSite> {
  const { data } = await apiClient.post<OwnedSubSite>(`${base(siteID)}/users/offline-topup`, payload)
  return data
}

export async function listOrders(
  siteID: number,
  page = 1,
  pageSize = 20,
  status?: string,
): Promise<PaginatedResponse<SubSiteAdminOrder>> {
  const { data } = await apiClient.get<PaginatedResponse<SubSiteAdminOrder>>(`${base(siteID)}/orders`, {
    params: { page, page_size: pageSize, status: status || undefined },
  })
  return data
}

export async function listUsage(
  siteID: number,
  page = 1,
  pageSize = 20,
  model?: string,
): Promise<PaginatedResponse<SubSiteAdminUsage>> {
  const { data } = await apiClient.get<PaginatedResponse<SubSiteAdminUsage>>(`${base(siteID)}/usage`, {
    params: { page, page_size: pageSize, model: model || undefined },
  })
  return data
}

export async function listLedger(
  siteID: number,
  page = 1,
  pageSize = 20,
  txType?: string,
): Promise<PaginatedResponse<OwnedSubSiteLedgerEntry>> {
  const { data } = await apiClient.get<PaginatedResponse<OwnedSubSiteLedgerEntry>>(`${base(siteID)}/ledger`, {
    params: { page, page_size: pageSize, tx_type: txType || undefined },
  })
  return data
}

export const subSiteAdminAPI = {
  getDashboard,
  getSite,
  updateSite,
  listUsers,
  offlineTopup,
  listOrders,
  listUsage,
  listLedger,
}

export default subSiteAdminAPI
