<template>
  <AppLayout>
    <FadeIn>
      <div class="space-y-6">
        <div v-if="loading" class="flex items-center justify-center py-12"><LoadingSpinner /></div>
        <template v-else>
          <!-- Welcome Section -->
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ greeting }}，{{ userName }}
            </h1>
            <p class="mt-1 text-gray-500 dark:text-gray-400">{{ t('dashboard.welcomeBack') }}</p>
          </div>

          <!-- Subscription & Balance Cards -->
          <SlideIn direction="up" :delay="200">
            <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
              <!-- Subscription Card -->
              <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft p-6">
                <div class="flex items-center justify-between">
                  <div>
                    <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('dashboard.subscriptionQuota') }}</p>
                    <p class="mt-1 text-lg font-bold text-gray-900 dark:text-white">
                      {{ subscriptionName }}
                    </p>
                  </div>
                  <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-primary-100 dark:bg-primary-900/30">
                    <svg class="h-5 w-5 text-primary-600 dark:text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456z" />
                    </svg>
                  </div>
                </div>
                <p class="mt-3 text-sm text-gray-500 dark:text-gray-400">{{ t('dashboard.upgradePlan') }}</p>
                <router-link
                  to="/pricing"
                  class="mt-4 inline-flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
                >
                  {{ t('dashboard.viewPlans') }}
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                  </svg>
                </router-link>
              </div>

              <!-- Balance Card -->
              <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft p-6">
                <div class="flex items-center justify-between">
                  <div>
                    <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('dashboard.flexibleBalance') }}</p>
                    <p class="mt-1 text-2xl font-bold text-emerald-600 dark:text-emerald-400">
                      ${{ (user?.balance || 0).toFixed(2) }}
                    </p>
                  </div>
                  <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-100 dark:bg-emerald-900/30">
                    <svg class="h-5 w-5 text-emerald-600 dark:text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
                    </svg>
                  </div>
                </div>
                <p class="mt-3 text-sm text-gray-500 dark:text-gray-400">{{ t('dashboard.payAsYouGo') }}</p>
                <router-link
                  to="/pricing?tab=recharge"
                  class="mt-4 inline-flex items-center gap-2 text-sm font-medium text-emerald-600 transition-colors hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-emerald-300"
                >
                  {{ t('dashboard.buyBalance') }} →
                </router-link>
              </div>
            </div>
          </SlideIn>

          <!-- Recent Announcements -->
          <SlideIn direction="up" :delay="300">
            <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft">
              <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
                <h2 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('dashboard.recentAnnouncements') }}</h2>
                <router-link
                  to="/changelog"
                  class="flex items-center gap-1 text-sm font-medium text-primary-600 transition-colors hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                >
                  {{ t('dashboard.viewAll') }} →
                </router-link>
              </div>
              <div class="p-4">
                <div v-if="announcements.length === 0" class="py-6 text-center text-sm text-gray-500 dark:text-gray-400">
                  {{ t('dashboard.noAnnouncements') }}
                </div>
                <div v-else class="space-y-1">
                  <router-link
                    v-for="item in announcements.slice(0, 4)"
                    :key="item.id"
                    to="/changelog"
                    class="flex items-center justify-between rounded-xl px-4 py-3 transition-colors hover:bg-gray-50 dark:hover:bg-dark-800"
                  >
                    <span class="text-sm text-gray-700 dark:text-gray-300">{{ item.title }}</span>
                    <span class="ml-4 flex-shrink-0 text-xs text-gray-400 dark:text-gray-500">
                      {{ formatDate(item.published_at || item.created_at) }}
                    </span>
                  </router-link>
                </div>
              </div>
            </div>
          </SlideIn>

          <!-- Quick Actions -->
          <SlideIn direction="up" :delay="400">
            <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft">
              <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
                <h2 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('dashboard.quickActions') }}</h2>
              </div>
              <div class="grid grid-cols-2 gap-4 p-4 lg:grid-cols-4">
                <button
                  @click="router.push('/pricing')"
                  class="group flex flex-col items-center gap-3 rounded-xl bg-gray-50 p-5 transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800"
                >
                  <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-primary-100 transition-transform group-hover:scale-105 dark:bg-primary-900/30">
                    <svg class="h-6 w-6 text-primary-600 dark:text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.75 10.5V6a3.75 3.75 0 10-7.5 0v4.5m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007zM8.625 10.5a.375.375 0 11-.75 0 .375.375 0 01.75 0zm7.5 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                    </svg>
                  </div>
                  <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('nav.purchase') }}</span>
                </button>

                <button
                  @click="router.push('/keys')"
                  class="group flex flex-col items-center gap-3 rounded-xl bg-gray-50 p-5 transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800"
                >
                  <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-100 transition-transform group-hover:scale-105 dark:bg-blue-900/30">
                    <svg class="h-6 w-6 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.75 5.25a3 3 0 013 3m3 0a6 6 0 01-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1121.75 8.25z" />
                    </svg>
                  </div>
                  <span class="text-sm font-medium text-gray-700 dark:text-gray-300">API Key</span>
                </button>

                <button
                  @click="router.push('/usage')"
                  class="group flex flex-col items-center gap-3 rounded-xl bg-gray-50 p-5 transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800"
                >
                  <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-emerald-100 transition-transform group-hover:scale-105 dark:bg-emerald-900/30">
                    <svg class="h-6 w-6 text-emerald-600 dark:text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
                    </svg>
                  </div>
                  <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('nav.usage') }}</span>
                </button>

                <button
                  @click="router.push('/redeem')"
                  class="group flex flex-col items-center gap-3 rounded-xl bg-gray-50 p-5 transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800"
                >
                  <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-amber-100 transition-transform group-hover:scale-105 dark:bg-amber-900/30">
                    <svg class="h-6 w-6 text-amber-600 dark:text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 11.25v8.25a1.5 1.5 0 01-1.5 1.5H5.25a1.5 1.5 0 01-1.5-1.5v-8.25M12 4.875A2.625 2.625 0 109.375 7.5H12m0-2.625V7.5m0-2.625A2.625 2.625 0 1114.625 7.5H12m0 0V21m-8.625-9.75h18c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125h-18c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z" />
                    </svg>
                  </div>
                  <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('dashboard.useRedeemCode') }}</span>
                </button>
              </div>
            </div>
          </SlideIn>
        </template>
      </div>
    </FadeIn>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useSubscriptionStore } from '@/stores/subscriptions'
import { authAPI } from '@/api/auth'
import type { Announcement } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { FadeIn, SlideIn } from '@/components/animations'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const subscriptionStore = useSubscriptionStore()

const user = computed(() => authStore.user)
const userName = computed(() => user.value?.username || user.value?.email?.split('@')[0] || '')
const loading = ref(false)
const announcements = ref<Announcement[]>([])

// Time-based greeting
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return t('dashboard.greeting.morning')
  if (hour < 18) return t('dashboard.greeting.afternoon')
  return t('dashboard.greeting.evening')
})

// Subscription display
const subscriptionName = computed(() => {
  if (subscriptionStore.hasActiveSubscriptions) {
    const sub = subscriptionStore.activeSubscriptions[0]
    return sub.group?.name || t('dashboard.subscriptionQuota')
  }
  return t('dashboard.flexibleUser')
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getFullYear()}/${String(date.getMonth() + 1).padStart(2, '0')}/${String(date.getDate()).padStart(2, '0')}`
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      authStore.refreshUser(),
      subscriptionStore.fetchActiveSubscriptions(),
      authAPI.getActiveAnnouncements().then(data => { announcements.value = data }).catch(() => {})
    ])
  } catch (error) {
    console.error('Failed to load dashboard:', error)
  } finally {
    loading.value = false
  }
})
</script>
