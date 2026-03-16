<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-5xl space-y-6">
        <!-- Title -->
        <SlideIn direction="up" :delay="100">
          <div class="text-center">
            <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white tracking-tight">{{ t('userSubscriptions.title') || '我的订阅' }}</h1>
            <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('userSubscriptions.subtitle') || '管理您的活跃订阅与用量' }}</p>
          </div>
        </SlideIn>

        <!-- Loading State -->
        <div v-if="loading" class="flex justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
        </div>

        <!-- Empty State -->
        <SlideIn v-else-if="subscriptions.length === 0" direction="up" :delay="200">
          <GlowCard glow-color="rgb(156, 163, 175)">
            <div class="rounded-3xl border-2 border-gray-200 bg-white p-16 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
              <div class="mx-auto mb-5 flex h-20 w-20 items-center justify-center rounded-2xl bg-gray-100 dark:bg-dark-700">
                <Icon name="creditCard" size="xl" class="text-gray-400" />
              </div>
              <h3 class="mb-2 text-xl font-bold text-gray-900 dark:text-white">
                {{ t('userSubscriptions.noActiveSubscriptions') }}
              </h3>
              <p class="text-gray-500 dark:text-dark-400 mb-6">
                {{ t('userSubscriptions.noActiveSubscriptionsDesc') }}
              </p>
              <MagneticButton>
                <router-link to="/pricing" class="inline-flex items-center gap-2 rounded-xl bg-primary-600 px-6 py-3 text-sm font-bold text-white hover:bg-primary-700 transition shadow-md">
                  <Icon name="sparkles" size="sm" />
                  立即订阅
                </router-link>
              </MagneticButton>
            </div>
          </GlowCard>
        </SlideIn>

        <!-- Subscriptions Grid -->
        <StaggerContainer v-else :stagger-delay="150" :delay="200">
          <div class="grid gap-6 lg:grid-cols-2">
            <GlowCard
              v-for="subscription in subscriptions"
              :key="subscription.id"
              glow-color="rgb(168, 85, 247)"
            >
              <div class="overflow-hidden rounded-2xl border-2 bg-white dark:bg-dark-900 shadow-soft"
                :class="subscription.status === 'active'
                  ? 'border-purple-200 dark:border-purple-800/30'
                  : 'border-gray-200 dark:border-dark-700'"
              >
                <!-- Header -->
                <div class="flex items-center justify-between border-b border-gray-100 p-5 dark:border-dark-700">
                  <div class="flex items-center gap-4">
                    <div class="flex h-12 w-12 items-center justify-center rounded-xl"
                      :class="subscription.status === 'active'
                        ? 'bg-purple-100 dark:bg-purple-900/30'
                        : 'bg-gray-100 dark:bg-dark-700'"
                    >
                      <Icon name="creditCard" size="md"
                        :class="subscription.status === 'active'
                          ? 'text-purple-600 dark:text-purple-400'
                          : 'text-gray-400 dark:text-dark-500'"
                      />
                    </div>
                    <div>
                      <h3 class="text-lg font-bold text-gray-900 dark:text-white">
                        {{ subscription.group?.name || `Group #${subscription.group_id}` }}
                      </h3>
                      <p class="text-xs text-gray-500 dark:text-dark-400">
                        {{ subscription.group?.description || '' }}
                      </p>
                    </div>
                  </div>
                  <span
                    :class="[
                      'rounded-full px-3 py-1 text-xs font-bold',
                      subscription.status === 'active'
                        ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                        : subscription.status === 'expired'
                          ? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
                          : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
                    ]"
                  >
                    {{ t(`userSubscriptions.status.${subscription.status}`) }}
                  </span>
                </div>

                <!-- Usage Progress -->
                <div class="space-y-4 p-5">
                  <!-- Expiration Info -->
                  <div v-if="subscription.expires_at" class="flex items-center justify-between rounded-xl bg-gray-50 dark:bg-dark-800 px-4 py-3">
                    <span class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ t('userSubscriptions.expires') }}</span>
                    <span class="text-sm font-bold" :class="getExpirationClass(subscription.expires_at)">
                      {{ formatExpirationDate(subscription.expires_at) }}
                    </span>
                  </div>
                  <div v-else class="flex items-center justify-between rounded-xl bg-gray-50 dark:bg-dark-800 px-4 py-3">
                    <span class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ t('userSubscriptions.expires') }}</span>
                    <span class="text-sm font-bold text-gray-700 dark:text-gray-300">{{ t('userSubscriptions.noExpiration') }}</span>
                  </div>

                  <!-- Daily Usage -->
                  <div v-if="subscription.group?.daily_limit_usd" class="space-y-2">
                    <div class="flex items-center justify-between">
                      <span class="text-sm font-bold text-gray-700 dark:text-gray-300">
                        {{ t('userSubscriptions.daily') }}
                      </span>
                      <span class="text-sm font-medium text-gray-500 dark:text-dark-400">
                        ${{ (subscription.daily_usage_usd || 0).toFixed(2) }} / ${{ subscription.group.daily_limit_usd.toFixed(2) }}
                      </span>
                    </div>
                    <div class="relative h-2.5 overflow-hidden rounded-full bg-gray-200 dark:bg-dark-600">
                      <div
                        class="absolute inset-y-0 left-0 rounded-full transition-all duration-500"
                        :class="getProgressBarClass(subscription.daily_usage_usd, subscription.group.daily_limit_usd)"
                        :style="{ width: getProgressWidth(subscription.daily_usage_usd, subscription.group.daily_limit_usd) }"
                      ></div>
                    </div>
                    <p v-if="subscription.daily_window_start" class="text-xs text-gray-500 dark:text-dark-400">
                      {{ t('userSubscriptions.resetIn', { time: formatResetTime(subscription.daily_window_start, 24) }) }}
                    </p>
                  </div>

                  <!-- Weekly Usage -->
                  <div v-if="subscription.group?.weekly_limit_usd" class="space-y-2">
                    <div class="flex items-center justify-between">
                      <span class="text-sm font-bold text-gray-700 dark:text-gray-300">
                        {{ t('userSubscriptions.weekly') }}
                      </span>
                      <span class="text-sm font-medium text-gray-500 dark:text-dark-400">
                        ${{ (subscription.weekly_usage_usd || 0).toFixed(2) }} / ${{ subscription.group.weekly_limit_usd.toFixed(2) }}
                      </span>
                    </div>
                    <div class="relative h-2.5 overflow-hidden rounded-full bg-gray-200 dark:bg-dark-600">
                      <div
                        class="absolute inset-y-0 left-0 rounded-full transition-all duration-500"
                        :class="getProgressBarClass(subscription.weekly_usage_usd, subscription.group.weekly_limit_usd)"
                        :style="{ width: getProgressWidth(subscription.weekly_usage_usd, subscription.group.weekly_limit_usd) }"
                      ></div>
                    </div>
                    <p v-if="subscription.weekly_window_start" class="text-xs text-gray-500 dark:text-dark-400">
                      {{ t('userSubscriptions.resetIn', { time: formatResetTime(subscription.weekly_window_start, 168) }) }}
                    </p>
                  </div>

                  <!-- Monthly Usage -->
                  <div v-if="subscription.group?.monthly_limit_usd" class="space-y-2">
                    <div class="flex items-center justify-between">
                      <span class="text-sm font-bold text-gray-700 dark:text-gray-300">
                        {{ t('userSubscriptions.monthly') }}
                      </span>
                      <span class="text-sm font-medium text-gray-500 dark:text-dark-400">
                        ${{ (subscription.monthly_usage_usd || 0).toFixed(2) }} / ${{ subscription.group.monthly_limit_usd.toFixed(2) }}
                      </span>
                    </div>
                    <div class="relative h-2.5 overflow-hidden rounded-full bg-gray-200 dark:bg-dark-600">
                      <div
                        class="absolute inset-y-0 left-0 rounded-full transition-all duration-500"
                        :class="getProgressBarClass(subscription.monthly_usage_usd, subscription.group.monthly_limit_usd)"
                        :style="{ width: getProgressWidth(subscription.monthly_usage_usd, subscription.group.monthly_limit_usd) }"
                      ></div>
                    </div>
                    <p v-if="subscription.monthly_window_start" class="text-xs text-gray-500 dark:text-dark-400">
                      {{ t('userSubscriptions.resetIn', { time: formatResetTime(subscription.monthly_window_start, 720) }) }}
                    </p>
                  </div>

                  <!-- No limits configured - Unlimited badge -->
                  <div
                    v-if="!subscription.group?.daily_limit_usd && !subscription.group?.weekly_limit_usd && !subscription.group?.monthly_limit_usd"
                    class="flex items-center justify-center rounded-2xl bg-gradient-to-r from-emerald-50 to-teal-50 py-8 dark:from-emerald-900/20 dark:to-teal-900/20 border border-emerald-100 dark:border-emerald-800/30"
                  >
                    <div class="flex items-center gap-4">
                      <span class="text-5xl text-emerald-600 dark:text-emerald-400">∞</span>
                      <div>
                        <p class="text-base font-bold text-emerald-700 dark:text-emerald-300">
                          {{ t('userSubscriptions.unlimited') }}
                        </p>
                        <p class="text-xs text-emerald-600/70 dark:text-emerald-400/70">
                          {{ t('userSubscriptions.unlimitedDesc') }}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </GlowCard>
          </div>
        </StaggerContainer>
      </div>
    </FadeIn>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import subscriptionsAPI from '@/api/subscriptions'
import type { UserSubscription } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateOnly } from '@/utils/format'
import { FadeIn, SlideIn, StaggerContainer, GlowCard, MagneticButton } from '@/components/animations'

const { t } = useI18n()
const appStore = useAppStore()

const subscriptions = ref<UserSubscription[]>([])
const loading = ref(true)

async function loadSubscriptions() {
  try {
    loading.value = true
    subscriptions.value = await subscriptionsAPI.getMySubscriptions()
  } catch (error) {
    console.error('Failed to load subscriptions:', error)
    appStore.showError(t('userSubscriptions.failedToLoad'))
  } finally {
    loading.value = false
  }
}

function getProgressWidth(used: number | undefined, limit: number | null | undefined): string {
  if (!limit || limit === 0) return '0%'
  const percentage = Math.min(((used || 0) / limit) * 100, 100)
  return `${percentage}%`
}

function getProgressBarClass(used: number | undefined, limit: number | null | undefined): string {
  if (!limit || limit === 0) return 'bg-gray-400'
  const percentage = ((used || 0) / limit) * 100
  if (percentage >= 90) return 'bg-red-500'
  if (percentage >= 70) return 'bg-orange-500'
  return 'bg-green-500'
}

function formatExpirationDate(expiresAt: string): string {
  const now = new Date()
  const expires = new Date(expiresAt)
  const diff = expires.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))

  if (days < 0) {
    return t('userSubscriptions.status.expired')
  }

  const dateStr = formatDateOnly(expires)

  if (days === 0) {
    return `${dateStr} (Today)`
  }
  if (days === 1) {
    return `${dateStr} (Tomorrow)`
  }

  return t('userSubscriptions.daysRemaining', { days }) + ` (${dateStr})`
}

function getExpirationClass(expiresAt: string): string {
  const now = new Date()
  const expires = new Date(expiresAt)
  const diff = expires.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))

  if (days <= 0) return 'text-red-600 dark:text-red-400'
  if (days <= 3) return 'text-red-600 dark:text-red-400'
  if (days <= 7) return 'text-orange-600 dark:text-orange-400'
  return 'text-gray-700 dark:text-gray-300'
}

function formatResetTime(windowStart: string | null, windowHours: number): string {
  if (!windowStart) return t('userSubscriptions.windowNotActive')

  const start = new Date(windowStart)
  const end = new Date(start.getTime() + windowHours * 60 * 60 * 1000)
  const now = new Date()
  const diff = end.getTime() - now.getTime()

  if (diff <= 0) return t('userSubscriptions.windowNotActive')

  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))

  if (hours > 24) {
    const days = Math.floor(hours / 24)
    const remainingHours = hours % 24
    return `${days}d ${remainingHours}h`
  }

  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }

  return `${minutes}m`
}

onMounted(() => {
  loadSubscriptions()
})
</script>
