import { apiClient } from './client'
import type { PaginatedResponse } from '@/types'

export interface AgentProfile {
  user_id: number
  real_name: string
  id_card_no: string
  phone: string
  identity_status: string
  identity_submitted_at: string | null
  contract_status: string
  contract_version: string
  contract_signed_at: string | null
  contract_ip?: string
  activation_order_id: number | null
  activation_fee_paid_at: string | null
  frozen_balance: number
  withdrawable_balance: number
  total_withdrawn: number
  is_frozen: boolean
  frozen_reason: string
}

export interface AgentWalletSummary {
  site_balance: number
  frozen_balance: number
  withdrawable_balance: number
  total_withdrawn: number
}

export interface AgentWithdrawWindow {
  weekday: number
  start_hour: number
  end_hour: number
  label: string
}

export interface AgentStatus {
  enabled: boolean
  is_agent: boolean
  agent_status: string
  commission_rate: number
  invite_code?: string
  can_apply: boolean
  activation_fee: number
  profile: AgentProfile
  wallet: AgentWalletSummary
  withdraw_freeze_days: number
  withdraw_window: AgentWithdrawWindow
}

export interface AgentDashboardStats {
  total_sub_users: number
  total_recharge: number
  total_consumed: number
  total_commission: number
  pending_commission: number
  settled_commission: number
  today_new_users: number
  today_recharge: number
  site_balance: number
  frozen_balance: number
  withdrawable_balance: number
  total_withdrawn: number
}

export interface AgentSubUser {
  id: number
  email: string
  username: string
  total_recharge: number
  total_consumed: number
  created_at: string
  commission_rate: number | null
  is_agent: boolean
}

export interface AgentFinancialLog {
  id: number
  user_id: number
  user_email: string
  type: string
  amount: number
  detail: string
  created_at: string
}

export interface AgentCommission {
  id: number
  agent_id: number
  user_id: number
  user_email: string
  order_id: number
  order_no: string
  source_type: string
  source_amount: number
  commission_rate: number
  commission_amount: number
  status: string
  settled_at: string | null
  created_at: string
}

export interface AgentInviteLink {
  invite_code: string
  invite_url?: string
}

export interface AgentProfileRequest {
  real_name: string
  id_card_no: string
  phone: string
  contract_accepted: boolean
}

export async function getStatus(): Promise<AgentStatus> {
  const { data } = await apiClient.get<AgentStatus>('/agent/status')
  return data
}

export async function saveProfile(payload: AgentProfileRequest): Promise<void> {
  await apiClient.post('/agent/profile', payload)
}

export async function apply(): Promise<void> {
  await apiClient.post('/agent/apply', {})
}

export async function getDashboard(): Promise<AgentDashboardStats> {
  const { data } = await apiClient.get<AgentDashboardStats>('/agent/dashboard')
  return data
}

export async function getLink(): Promise<AgentInviteLink> {
  const { data } = await apiClient.get<AgentInviteLink>('/agent/link')
  return data
}

export async function listSubUsers(
  page: number = 1,
  pageSize: number = 20
): Promise<PaginatedResponse<AgentSubUser>> {
  const { data } = await apiClient.get<PaginatedResponse<AgentSubUser>>('/agent/sub-users', {
    params: { page, page_size: pageSize }
  })
  return data
}

export async function listFinancialLogs(
  page: number = 1,
  pageSize: number = 20
): Promise<PaginatedResponse<AgentFinancialLog>> {
  const { data } = await apiClient.get<PaginatedResponse<AgentFinancialLog>>('/agent/financial-logs', {
    params: {
      page,
      page_size: pageSize,
    }
  })
  return data
}

export async function listCommissions(
  page: number = 1,
  pageSize: number = 20,
  filters?: { status?: string }
): Promise<PaginatedResponse<AgentCommission>> {
  const { data } = await apiClient.get<PaginatedResponse<AgentCommission>>('/agent/commissions', {
    params: {
      page,
      page_size: pageSize,
      status: filters?.status || undefined
    }
  })
  return data
}

export const agentAPI = {
  getStatus,
  saveProfile,
  apply,
  getDashboard,
  getLink,
  listSubUsers,
  listFinancialLogs,
  listCommissions
}

export default agentAPI
