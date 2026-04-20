import { apiClient } from './client'
import type { Group } from '@/types'

export interface GroupModels {
  group: Group
  models: string[]
}

export interface ModelPlazaPricingMetrics {
  input_per_million: number
  output_per_million: number
  cache_write_per_million: number
  cache_read_per_million: number
  input_per_million_above_200k: number | null
  output_per_million_above_200k: number | null
  cache_write_per_million_above_200k: number | null
  cache_read_per_million_above_200k: number | null
}

export interface ModelPlazaPricingGroup {
  id: number
  name: string
  platform: Group['platform']
  rate_multiplier: number
  display_rate_multiplier: number | null
  display_price: string
  display_discount: string
  description: string | null
  subscription_type: Group['subscription_type']
  has_explicit_models: boolean
  models_count: number
}

export interface ModelPlazaPricingItem {
  model: string
  aliases: string[]
  platform: Group['platform']
  provider: string
  mode: string
  supports_prompt_caching: boolean
  max_input_tokens: number
  max_output_tokens: number
  official: ModelPlazaPricingMetrics
  group_prices: Record<number, ModelPlazaPricingMetrics>
}

export interface ModelPlazaPricingTable {
  groups: ModelPlazaPricingGroup[]
  items: ModelPlazaPricingItem[]
}

const MODEL_PLAZA_CACHE_TTL_MS = 5 * 60 * 1000
const MODEL_PLAZA_SESSION_CACHE_KEY = 'sub2api:model-plaza'
const MODEL_PLAZA_PRICING_SESSION_CACHE_KEY = 'sub2api:model-plaza:pricing-table'

let cachedGroupModels: GroupModels[] | null = null
let cachedPricingTable: ModelPlazaPricingTable | null = null
let cacheExpiresAt = 0
let inFlightRequest: Promise<GroupModels[]> | null = null
let pricingInFlightRequest: Promise<ModelPlazaPricingTable> | null = null

interface ModelPlazaCachePayload<T> {
  expiresAt: number
  data: T
}

const normalizeGroupModels = (value: unknown): GroupModels[] => {
  if (!Array.isArray(value)) {
    return []
  }

  return value
    .map((item) => {
      if (!item || typeof item !== 'object' || !('group' in item)) {
        return null
      }

      const record = item as GroupModels & { models?: unknown }
      return {
        group: record.group,
        models: Array.isArray(record.models) ? record.models.filter((model): model is string => typeof model === 'string') : []
      }
    })
    .filter((item): item is GroupModels => item !== null)
}

const normalizePricingTable = (value: unknown): ModelPlazaPricingTable => {
  if (!value || typeof value !== 'object') {
    return { groups: [], items: [] }
  }

  const record = value as {
    groups?: unknown
    items?: unknown
  }

  return {
    groups: Array.isArray(record.groups) ? record.groups as ModelPlazaPricingGroup[] : [],
    items: Array.isArray(record.items)
      ? record.items.map((item) => {
          const pricingItem = (item ?? {}) as ModelPlazaPricingItem & {
            aliases?: unknown
            group_prices?: unknown
          }

          return {
            ...pricingItem,
            aliases: Array.isArray(pricingItem.aliases)
              ? pricingItem.aliases.filter((alias): alias is string => typeof alias === 'string')
              : [],
            group_prices:
              pricingItem.group_prices && typeof pricingItem.group_prices === 'object'
                ? pricingItem.group_prices as Record<number, ModelPlazaPricingMetrics>
                : {}
          }
        })
      : []
  }
}

const readSessionCache = <T>(key: string): ModelPlazaCachePayload<T> | null => {
  if (typeof window === 'undefined') {
    return null
  }

  try {
    const raw = window.sessionStorage.getItem(key)
    if (!raw) {
      return null
    }

    const payload = JSON.parse(raw) as ModelPlazaCachePayload<T>
    if (payload.expiresAt <= Date.now()) {
      window.sessionStorage.removeItem(key)
      return null
    }

    return payload
  } catch {
    return null
  }
}

const writeSessionCache = <T>(key: string, data: T, expiresAt: number) => {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.sessionStorage.setItem(
      key,
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
    const cachedPayload = readSessionCache<GroupModels[]>(MODEL_PLAZA_SESSION_CACHE_KEY)
    if (cachedPayload && Array.isArray(cachedPayload.data)) {
      cachedGroupModels = normalizeGroupModels(cachedPayload.data)
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
      cachedGroupModels = normalizeGroupModels(data)
      cacheExpiresAt = Date.now() + MODEL_PLAZA_CACHE_TTL_MS
      writeSessionCache(MODEL_PLAZA_SESSION_CACHE_KEY, cachedGroupModels, cacheExpiresAt)
      return cachedGroupModels
    })
    .finally(() => {
      inFlightRequest = null
    })

  return inFlightRequest
}

export async function getModelPlazaPricingTable(force = false): Promise<ModelPlazaPricingTable> {
  const now = Date.now()

  if (!force && cachedPricingTable && cacheExpiresAt > now) {
    return cachedPricingTable
  }

  if (!force && !cachedPricingTable) {
    const cachedPayload = readSessionCache<ModelPlazaPricingTable>(MODEL_PLAZA_PRICING_SESSION_CACHE_KEY)
    if (cachedPayload && cachedPayload.data) {
      cachedPricingTable = normalizePricingTable(cachedPayload.data)
      cacheExpiresAt = cachedPayload.expiresAt
      return cachedPricingTable
    }
  }

  if (!force && pricingInFlightRequest) {
    return pricingInFlightRequest
  }

  pricingInFlightRequest = apiClient
    .get<ModelPlazaPricingTable>('/model-plaza/pricing-table')
    .then(({ data }) => {
      cachedPricingTable = normalizePricingTable(data)
      cacheExpiresAt = Date.now() + MODEL_PLAZA_CACHE_TTL_MS
      writeSessionCache(MODEL_PLAZA_PRICING_SESSION_CACHE_KEY, cachedPricingTable, cacheExpiresAt)
      return cachedPricingTable
    })
    .finally(() => {
      pricingInFlightRequest = null
    })

  return pricingInFlightRequest
}

export default { getModelPlaza, getModelPlazaPricingTable }
