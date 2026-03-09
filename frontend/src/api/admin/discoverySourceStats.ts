import { apiClient } from '../client'

export interface DiscoverySourceStat {
  source: string
  count: number
}

export interface DiscoverySourceStatsResponse {
  day7: DiscoverySourceStat[]
  day1: DiscoverySourceStat[]
  total: number
}

export async function getStats(): Promise<DiscoverySourceStatsResponse> {
  const { data } = await apiClient.get<DiscoverySourceStatsResponse>('/admin/discovery-source-stats')
  return data
}

export const discoverySourceStatsAPI = { getStats }
export default discoverySourceStatsAPI
