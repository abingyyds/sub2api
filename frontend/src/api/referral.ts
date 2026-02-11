import { apiClient } from './client'
import type { PaginatedResponse } from '@/types'

export interface ReferralInvitee {
  id: number
  invitee_id: number
  invitee_email: string
  reward_status: 'pending' | 'rewarded'
  reward_amount: number
  rewarded_at: string | null
  created_at: string
}

export interface ReferralStats {
  total_invitees: number
  rewarded_count: number
  pending_count: number
  total_reward_amount: number
}

export async function getInviteCode(): Promise<{ invite_code: string }> {
  const { data } = await apiClient.get<{ invite_code: string }>('/referral/code')
  return data
}

export async function listInvitees(
  page: number = 1,
  pageSize: number = 20
): Promise<PaginatedResponse<ReferralInvitee>> {
  const { data } = await apiClient.get<PaginatedResponse<ReferralInvitee>>('/referral/invitees', {
    params: { page, page_size: pageSize }
  })
  return data
}

export async function getStats(): Promise<ReferralStats> {
  const { data } = await apiClient.get<ReferralStats>('/referral/stats')
  return data
}

export default { getInviteCode, listInvitees, getStats }
