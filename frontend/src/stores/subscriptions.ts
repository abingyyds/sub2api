/**
 * Subscription Store
 * Global state management for user subscriptions with caching and deduplication
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import subscriptionsAPI from '@/api/subscriptions'
import type { UserSubscription } from '@/types'

// Cache TTL: 60 seconds
const CACHE_TTL_MS = 60_000

// Request generation counters to invalidate stale in-flight responses
let activeRequestGeneration = 0
let allRequestGeneration = 0

export const useSubscriptionStore = defineStore('subscriptions', () => {
  // State
  const activeSubscriptions = ref<UserSubscription[]>([])
  const allSubscriptions = ref<UserSubscription[]>([])
  const loading = ref(false)
  const allLoading = ref(false)
  const loaded = ref(false)
  const allLoaded = ref(false)
  const lastFetchedAt = ref<number | null>(null)
  const allLastFetchedAt = ref<number | null>(null)

  // In-flight request deduplication
  let activePromise: Promise<UserSubscription[]> | null = null
  let allPromise: Promise<UserSubscription[]> | null = null

  // Auto-refresh interval
  let pollerInterval: ReturnType<typeof setInterval> | null = null

  // Computed
  const hasActiveSubscriptions = computed(() => activeSubscriptions.value.length > 0)

  /**
   * Fetch all subscriptions with caching and deduplication
   * @param force - Force refresh even if cache is valid
   */
  async function fetchMySubscriptions(force = false): Promise<UserSubscription[]> {
    const now = Date.now()

    if (
      !force &&
      allLoaded.value &&
      allLastFetchedAt.value &&
      now - allLastFetchedAt.value < CACHE_TTL_MS
    ) {
      return allSubscriptions.value
    }

    if (allPromise && !force) {
      return allPromise
    }

    const currentGeneration = ++allRequestGeneration

    allLoading.value = true
    const requestPromise = subscriptionsAPI
      .getMySubscriptions()
      .then((data) => {
        if (currentGeneration === allRequestGeneration) {
          allSubscriptions.value = data
          allLoaded.value = true
          allLastFetchedAt.value = Date.now()
          activeSubscriptions.value = data.filter((subscription) => subscription.status === 'active')
          loaded.value = true
          lastFetchedAt.value = allLastFetchedAt.value
        }
        return data
      })
      .catch((error) => {
        console.error('Failed to fetch subscriptions:', error)
        throw error
      })
      .finally(() => {
        if (allPromise === requestPromise) {
          allLoading.value = false
          allPromise = null
        }
      })

    allPromise = requestPromise

    return allPromise
  }

  /**
   * Fetch active subscriptions with caching and deduplication
   * @param force - Force refresh even if cache is valid
   */
  async function fetchActiveSubscriptions(force = false): Promise<UserSubscription[]> {
    const now = Date.now()

    // Return cached data if valid
    if (
      !force &&
      loaded.value &&
      lastFetchedAt.value &&
      now - lastFetchedAt.value < CACHE_TTL_MS
    ) {
      return activeSubscriptions.value
    }

    // Return in-flight request if exists (deduplication)
    if (activePromise && !force) {
      return activePromise
    }

    const currentGeneration = ++activeRequestGeneration

    // Start new request
    loading.value = true
    const requestPromise = subscriptionsAPI
      .getActiveSubscriptions()
      .then((data) => {
        if (currentGeneration === activeRequestGeneration) {
          activeSubscriptions.value = data
          loaded.value = true
          lastFetchedAt.value = Date.now()
          if (allLoaded.value) {
            allLastFetchedAt.value = null
          }
        }
        return data
      })
      .catch((error) => {
        console.error('Failed to fetch active subscriptions:', error)
        throw error
      })
      .finally(() => {
        if (activePromise === requestPromise) {
          loading.value = false
          activePromise = null
        }
      })

    activePromise = requestPromise

    return activePromise
  }

  /**
   * Start auto-refresh polling (every 10 minutes)
   */
  function startPolling() {
    if (pollerInterval) return

    pollerInterval = setInterval(() => {
      fetchActiveSubscriptions(true).catch((error) => {
        console.error('Subscription polling failed:', error)
      })
    }, 10 * 60 * 1000) // 10 minutes
  }

  /**
   * Stop auto-refresh polling
   */
  function stopPolling() {
    if (pollerInterval) {
      clearInterval(pollerInterval)
      pollerInterval = null
    }
  }

  /**
   * Clear all subscription data and stop polling
   */
  function clear() {
    activeRequestGeneration++
    allRequestGeneration++
    activePromise = null
    allPromise = null
    activeSubscriptions.value = []
    allSubscriptions.value = []
    loaded.value = false
    allLoaded.value = false
    lastFetchedAt.value = null
    allLastFetchedAt.value = null
    loading.value = false
    allLoading.value = false
    stopPolling()
  }

  /**
   * Invalidate cache (force next fetch to reload)
   */
  function invalidateCache() {
    lastFetchedAt.value = null
    allLastFetchedAt.value = null
  }

  return {
    // State
    activeSubscriptions,
    allSubscriptions,
    loading,
    allLoading,
    loaded,
    allLoaded,
    hasActiveSubscriptions,

    // Actions
    fetchMySubscriptions,
    fetchActiveSubscriptions,
    startPolling,
    stopPolling,
    clear,
    invalidateCache
  }
})
