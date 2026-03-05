import { apiClient } from '../client'
import type {
  Organization,
  CreateOrganizationRequest,
  UpdateOrganizationRequest,
  UpdateOrgBalanceRequest,
  BasePaginationResponse
} from '@/types'

export async function list(
  page = 1,
  pageSize = 20,
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<Organization>> {
  const { data } = await apiClient.get<BasePaginationResponse<Organization>>(
    '/admin/organizations',
    { params: { page, page_size: pageSize }, signal: options?.signal }
  )
  return data
}

export async function getById(id: number): Promise<Organization> {
  const { data } = await apiClient.get<Organization>(`/admin/organizations/${id}`)
  return data
}

export async function create(req: CreateOrganizationRequest): Promise<Organization> {
  const { data } = await apiClient.post<Organization>('/admin/organizations', req)
  return data
}

export async function update(id: number, req: UpdateOrganizationRequest): Promise<Organization> {
  const { data } = await apiClient.put<Organization>(`/admin/organizations/${id}`, req)
  return data
}

export async function remove(id: number): Promise<void> {
  await apiClient.delete(`/admin/organizations/${id}`)
}

export async function updateBalance(
  id: number,
  req: UpdateOrgBalanceRequest
): Promise<Organization> {
  const { data } = await apiClient.post<Organization>(`/admin/organizations/${id}/balance`, req)
  return data
}

export const organizationsAPI = { list, getById, create, update, remove, updateBalance }
export default organizationsAPI
