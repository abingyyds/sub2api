<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-6">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <template v-else>
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-4">
          <div class="card p-5 text-center">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.referrals.totalReferrals') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total_referrals }}</p>
          </div>
          <div class="card p-5 text-center">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.referrals.rewardedCount') }}</p>
            <p class="mt-1 text-2xl font-bold text-green-600 dark:text-green-400">{{ stats.rewarded_count }}</p>
          </div>
          <div class="card p-5 text-center">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.referrals.pendingCount') }}</p>
            <p class="mt-1 text-2xl font-bold text-yellow-600 dark:text-yellow-400">{{ stats.pending_count }}</p>
          </div>
          <div class="card p-5 text-center">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.referrals.totalRewardAmount') }}</p>
            <p class="mt-1 text-2xl font-bold text-primary-600 dark:text-primary-400">${{ stats.total_reward_amount?.toFixed(2) || '0.00' }}</p>
          </div>
        </div>

        <!-- Referrals Table -->
        <div class="card">
          <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.referrals.allRecords') }}</h2>
            <div class="relative w-64">
              <input
                v-model="search"
                type="text"
                :placeholder="t('admin.referrals.searchPlaceholder')"
                class="input pl-10 text-sm"
                @input="handleSearch"
              />
              <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                <Icon name="search" size="sm" class="text-gray-400" />
              </div>
            </div>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-gray-100 dark:border-dark-700">
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.referrals.inviter') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.referrals.invitee') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.referrals.status') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.referrals.rewardAmount') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.referrals.createdAt') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                <tr v-if="referrals.length === 0">
                  <td colspan="5" class="px-6 py-8 text-center text-gray-500 dark:text-dark-400">
                    {{ t('admin.referrals.noRecords') }}
                  </td>
                </tr>
                <tr v-for="item in referrals" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-dark-800">
                  <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">{{ item.inviter_email }}</td>
                  <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">{{ item.invitee_email }}</td>
                  <td class="px-6 py-4">
                    <span :class="item.reward_status === 'rewarded'
                      ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                      : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'"
                      class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium">
                      {{ item.reward_status === 'rewarded' ? t('admin.referrals.statusRewarded') : t('admin.referrals.statusPending') }}
                    </span>
                  </td>
                  <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">${{ item.reward_amount?.toFixed(2) || '0.00' }}</td>
                  <td class="px-6 py-4 text-sm text-gray-500 dark:text-dark-400">{{ formatDate(item.created_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-gray-100 px-6 py-3 dark:border-dark-700">
            <button @click="page > 1 && loadReferrals(page - 1)" :disabled="page <= 1" class="btn btn-secondary btn-sm">{{ t('common.previous') }}</button>
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ page }} / {{ totalPages }}</span>
            <button @click="page < totalPages && loadReferrals(page + 1)" :disabled="page >= totalPages" class="btn btn-secondary btn-sm">{{ t('common.next') }}</button>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { list, getStats, type AdminReferral, type AdminReferralStats } from '@/api/admin/referrals'

const { t } = useI18n()

const loading = ref(true)
const search = ref('')
const page = ref(1)
const totalPages = ref(1)
const referrals = ref<AdminReferral[]>([])
const stats = ref<AdminReferralStats>({ total_referrals: 0, rewarded_count: 0, pending_count: 0, total_reward_amount: 0 })
let searchTimeout: ReturnType<typeof setTimeout> | null = null

onMounted(async () => {
  try {
    const [statsRes] = await Promise.all([getStats()])
    stats.value = statsRes
    await loadReferrals(1)
  } catch (err) {
    console.error('Failed to load referral data:', err)
  } finally {
    loading.value = false
  }
})

async function loadReferrals(p: number) {
  try {
    const res = await list(p, 20, search.value || undefined)
    referrals.value = res.items || []
    page.value = res.page
    totalPages.value = res.pages || Math.ceil((res.total || 0) / (res.page_size || 20))
  } catch (err) {
    console.error('Failed to load referrals:', err)
  }
}

function handleSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => { loadReferrals(1) }, 300)
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}
</script>
