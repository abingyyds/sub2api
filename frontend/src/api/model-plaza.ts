import { apiClient } from './client'
import type { Group } from '@/types'

export interface GroupModels {
  group: Group
  models: string[]
}

const MODEL_PLAZA_CACHE_TTL_MS = 5 * 60 * 1000
const MODEL_PLAZA_SESSION_CACHE_KEY = 'sub2api:model-plaza'

let cachedGroupModels: GroupModels[] | null = null
let cacheExpiresAt = 0
let inFlightRequest: Promise<GroupModels[]> | null = null

interface ModelPlazaCachePayload {
  expiresAt: number
  data: GroupModels[]
}

const readSessionCache = (): ModelPlazaCachePayload | null => {
  if (typeof window === 'undefined') {
    return null
  }

  try {
    const raw = window.sessionStorage.getItem(MODEL_PLAZA_SESSION_CACHE_KEY)
    if (!raw) {
      return null
    }

    const payload = JSON.parse(raw) as ModelPlazaCachePayload
    if (!Array.isArray(payload.data) || payload.expiresAt <= Date.now()) {
      window.sessionStorage.removeItem(MODEL_PLAZA_SESSION_CACHE_KEY)
      return null
    }

    return payload
  } catch {
    return null
  }
}

const writeSessionCache = (data: GroupModels[], expiresAt: number) => {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.sessionStorage.setItem(
      MODEL_PLAZA_SESSION_CACHE_KEY,
      JSON.stringify({ data, expiresAt })
    )
  } catch {
    // Ignore storage failures and continue with the in-memory cache.
  }
}

export async function getModelPlaza(force = false): Promise<GroupModels[]> {
  const now = Date.now()

  if (!force && cachedGroupModels && cacheExpiresAt > now) {
    return cachedGroupModels
  }

  if (!force && !cachedGroupModels) {
    const cachedPayload = readSessionCache()
    if (cachedPayload) {
      cachedGroupModels = cachedPayload.data
      cacheExpiresAt = cachedPayload.expiresAt
      return cachedGroupModels
    }
  }

  if (!force && inFlightRequest) {
    return inFlightRequest
  }

  inFlightRequest = apiClient
    .get<GroupModels[]>('/model-plaza')
    .then(({ data }) => {
      cachedGroupModels = data
      cacheExpiresAt = Date.now() + MODEL_PLAZA_CACHE_TTL_MS
      writeSessionCache(data, cacheExpiresAt)
      return data
    })
    .finally(() => {
      inFlightRequest = null
    })

  return inFlightRequest
}

export default { getModelPlaza }
