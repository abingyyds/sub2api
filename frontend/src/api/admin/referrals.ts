import { apiClient } from '../client'
import type { PaginatedResponse } from '@/types'

export interface AdminReferral {
  id: number
  inviter_id: number
  inviter_email: string
  invitee_id: number
  invitee_email: string
  reward_status: 'pending' | 'rewarded'
  reward_amount: number
  rewarded_at: string | null
  created_at: string
}

export interface AdminReferralStats {
  total_referrals: number
  rewarded_count: number
  pending_count: number
  total_reward_amount: number
}

export async function list(
  page: number = 1,
  pageSize: number = 20,
  search?: string,
  options?: { signal?: AbortSignal }
): Promise<PaginatedResponse<AdminReferral>> {
  const { data } = await apiClient.get<PaginatedResponse<AdminReferral>>('/admin/referrals', {
    params: { page, page_size: pageSize, search: search || undefined },
    signal: options?.signal
  })
  return data
}

export async function getStats(): Promise<AdminReferralStats> {
  const { data } = await apiClient.get<AdminReferralStats>('/admin/referrals/stats')
  return data
}

export const referralsAPI = { list, getStats }
export default referralsAPI
