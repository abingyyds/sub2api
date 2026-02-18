import { apiClient } from '../client'
import type { Announcement, BasePaginationResponse } from '@/types'

export async function list(page = 1, pageSize = 20): Promise<BasePaginationResponse<Announcement>> {
  const { data } = await apiClient.get<BasePaginationResponse<Announcement>>('/admin/announcements', { params: { page, page_size: pageSize } })
  return data
}

export async function create(req: { title: string; content?: string; status?: string; priority?: number }): Promise<Announcement> {
  const { data } = await apiClient.post<Announcement>('/admin/announcements', req)
  return data
}

export async function update(id: number, req: Partial<{ title: string; content: string; status: string; priority: number }>): Promise<Announcement> {
  const { data } = await apiClient.put<Announcement>(`/admin/announcements/${id}`, req)
  return data
}

export async function remove(id: number): Promise<void> {
  await apiClient.delete(`/admin/announcements/${id}`)
}

export const announcementsAPI = { list, create, update, remove }
export default announcementsAPI
