import { apiClient } from './client'
import type { PaginatedResponse } from '@/types'

export interface AgentStatus {
  is_agent: boolean
  agent_status: string
  agent_commission_rate: number
  agent_note: string
  agent_approved_at: string | null
}

export interface AgentDashboardStats {
  total_sub_users: number
  total_recharge: number
  total_consumed: number
  total_commission: number
  pending_commission: number
  settled_commission: number
  direct_commission: number
  differential_commission: number
}

export interface AgentSubUser {
  user_id: number
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
  order_no: string
  user_id: number
  user_email: string
  plan_key: string
  amount_fen: number
  balance_amount: number
  order_type: string
  status: string
  paid_at: string | null
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
  invite_url: string
}

export async function getStatus(): Promise<AgentStatus> {
  const { data } = await apiClient.get<AgentStatus>('/agent/status')
  return data
}

export interface AgentApplyRequest {
  contact: string
  social?: string
  promotion?: string
}

export async function apply(data: AgentApplyRequest): Promise<void> {
  await apiClient.post('/agent/apply', data)
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
  pageSize: number = 20,
  filters?: { user_id?: number }
): Promise<PaginatedResponse<AgentFinancialLog>> {
  const { data } = await apiClient.get<PaginatedResponse<AgentFinancialLog>>('/agent/financial-logs', {
    params: {
      page,
      page_size: pageSize,
      user_id: filters?.user_id || undefined
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

export async function setSubUserRate(subUserId: number, rate: number): Promise<void> {
  await apiClient.put(`/agent/sub-users/${subUserId}/rate`, { commission_rate: rate })
}

export const agentAPI = {
  getStatus,
  apply,
  getDashboard,
  getLink,
  listSubUsers,
  listFinancialLogs,
  listCommissions,
  setSubUserRate
}

export default agentAPI
