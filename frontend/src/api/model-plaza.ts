import { apiClient } from './client'
import type { Group } from '@/types'

export interface GroupModels {
  group: Group
  models: string[]
}

export async function getModelPlaza(): Promise<GroupModels[]> {
  const { data } = await apiClient.get<GroupModels[]>('/model-plaza')
  return data
}

export default { getModelPlaza }
