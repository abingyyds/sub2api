import { apiClient } from './client'
import type {
  OrgDashboard,
  OrgMember,
  CreateOrgMemberRequest,
  UpdateOrgMemberRequest,
  OrgProject,
  CreateOrgProjectRequest,
  UpdateOrgProjectRequest,
  OrgAuditLog,
  AuditLogFilters,
  AuditConfig,
  BasePaginationResponse
} from '@/types'

// Dashboard
export async function getDashboard(): Promise<OrgDashboard> {
  const { data } = await apiClient.get<OrgDashboard>('/org/dashboard')
  return data
}

// Members
export async function listMembers(
  page = 1,
  pageSize = 20,
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<OrgMember>> {
  const { data } = await apiClient.get<BasePaginationResponse<OrgMember>>('/org/members', {
    params: { page, page_size: pageSize },
    signal: options?.signal
  })
  return data
}

export async function getMemberById(id: number): Promise<OrgMember> {
  const { data } = await apiClient.get<OrgMember>(`/org/members/${id}`)
  return data
}

export async function createMember(req: CreateOrgMemberRequest): Promise<OrgMember> {
  const { data } = await apiClient.post<OrgMember>('/org/members', req)
  return data
}

export async function updateMember(id: number, req: UpdateOrgMemberRequest): Promise<OrgMember> {
  const { data } = await apiClient.put<OrgMember>(`/org/members/${id}`, req)
  return data
}

export async function removeMember(id: number): Promise<void> {
  await apiClient.delete(`/org/members/${id}`)
}

export async function suspendMember(id: number): Promise<OrgMember> {
  const { data } = await apiClient.post<OrgMember>(`/org/members/${id}/suspend`)
  return data
}

// Projects
export async function listProjects(
  page = 1,
  pageSize = 20,
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<OrgProject>> {
  const { data } = await apiClient.get<BasePaginationResponse<OrgProject>>('/org/projects', {
    params: { page, page_size: pageSize },
    signal: options?.signal
  })
  return data
}

export async function getProjectById(id: number): Promise<OrgProject> {
  const { data } = await apiClient.get<OrgProject>(`/org/projects/${id}`)
  return data
}

export async function createProject(req: CreateOrgProjectRequest): Promise<OrgProject> {
  const { data } = await apiClient.post<OrgProject>('/org/projects', req)
  return data
}

export async function updateProject(id: number, req: UpdateOrgProjectRequest): Promise<OrgProject> {
  const { data } = await apiClient.put<OrgProject>(`/org/projects/${id}`, req)
  return data
}

export async function removeProject(id: number): Promise<void> {
  await apiClient.delete(`/org/projects/${id}`)
}

// Audit Logs
export async function listAuditLogs(
  page = 1,
  pageSize = 20,
  filters?: AuditLogFilters,
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<OrgAuditLog>> {
  const params: Record<string, unknown> = { page, page_size: pageSize }
  if (filters) {
    if (filters.member_id) params.member_id = filters.member_id
    if (filters.project_id) params.project_id = filters.project_id
    if (filters.action) params.action = filters.action
    if (filters.model) params.model = filters.model
    if (filters.flagged !== undefined) params.flagged = filters.flagged
    if (filters.start_date) params.start_date = filters.start_date
    if (filters.end_date) params.end_date = filters.end_date
  }
  const { data } = await apiClient.get<BasePaginationResponse<OrgAuditLog>>('/org/audit-logs', {
    params,
    signal: options?.signal
  })
  return data
}

export async function getAuditLogById(id: number): Promise<OrgAuditLog> {
  const { data } = await apiClient.get<OrgAuditLog>(`/org/audit-logs/${id}`)
  return data
}

export async function getFlaggedCount(): Promise<{ count: number }> {
  const { data } = await apiClient.get<{ count: number }>('/org/audit-logs/flagged-count')
  return data
}

// Audit Config
export async function getAuditConfig(): Promise<AuditConfig> {
  const { data } = await apiClient.get<AuditConfig>('/org/audit-config')
  return data
}

export async function updateAuditConfig(auditMode: string): Promise<AuditConfig> {
  const { data } = await apiClient.put<AuditConfig>('/org/audit-config', { audit_mode: auditMode })
  return data
}

export const orgAPI = {
  getDashboard,
  listMembers,
  getMemberById,
  createMember,
  updateMember,
  removeMember,
  suspendMember,
  listProjects,
  getProjectById,
  createProject,
  updateProject,
  removeProject,
  listAuditLogs,
  getAuditLogById,
  getFlaggedCount,
  getAuditConfig,
  updateAuditConfig
}
export default orgAPI
