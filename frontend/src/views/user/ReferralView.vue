<template>
  <AppLayout>
    <div class="mx-auto max-w-3xl space-y-6">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <template v-else>
        <!-- Invite Code Card -->
        <div class="card overflow-hidden">
          <div class="bg-gradient-to-br from-primary-500 to-primary-600 px-6 py-8 text-center">
            <div class="mb-4 inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 backdrop-blur-sm">
              <Icon name="share" size="xl" class="text-white" />
            </div>
            <p class="text-sm font-medium text-primary-100">{{ t('referral.yourInviteCode') }}</p>
            <p class="mt-2 text-3xl font-bold tracking-wider text-white">{{ inviteCode }}</p>
            <div class="mt-4 flex items-center justify-center gap-3">
              <button @click="copyCode" class="inline-flex items-center gap-2 rounded-lg bg-white/20 px-4 py-2 text-sm font-medium text-white backdrop-blur-sm transition hover:bg-white/30">
                <Icon name="clipboard" size="sm" />
                {{ t('referral.copyCode') }}
              </button>
              <button @click="copyLink" class="inline-flex items-center gap-2 rounded-lg bg-white/20 px-4 py-2 text-sm font-medium text-white backdrop-blur-sm transition hover:bg-white/30">
                <Icon name="link" size="sm" />
                {{ t('referral.copyLink') }}
              </button>
            </div>
            <p v-if="copied" class="mt-2 text-sm text-primary-100">{{ t('referral.copied') }}</p>
          </div>
        </div>

        <!-- Stats Cards -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div class="card p-5 text-center">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('referral.totalInvitees') }}</p>
            <p class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total_invitees }}</p>
          </div>
          <div class="card p-5 text-center">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('referral.rewarded') }}</p>
            <p class="mt-1 text-2xl font-bold text-green-600 dark:text-green-400">{{ stats.rewarded_count }}</p>
          </div>
          <div class="card p-5 text-center">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('referral.totalReward') }}</p>
            <p class="mt-1 text-2xl font-bold text-primary-600 dark:text-primary-400">${{ stats.total_reward_amount?.toFixed(2) || '0.00' }}</p>
          </div>
        </div>

        <!-- Invitees Table -->
        <div class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('referral.inviteeList') }}</h2>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-gray-100 dark:border-dark-700">
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.inviteeEmail') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.status') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.rewardAmount') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.invitedAt') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                <tr v-if="invitees.length === 0">
                  <td colspan="4" class="px-6 py-8 text-center text-gray-500 dark:text-dark-400">
                    {{ t('referral.noInvitees') }}
                  </td>
                </tr>
                <tr v-for="item in invitees" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-dark-800">
                  <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">{{ item.invitee_email }}</td>
                  <td class="px-6 py-4">
                    <span :class="item.reward_status === 'rewarded'
                      ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                      : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'"
                      class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium">
                      {{ item.reward_status === 'rewarded' ? t('referral.statusRewarded') : t('referral.statusPending') }}
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
            <button @click="page > 1 && loadInvitees(page - 1)" :disabled="page <= 1" class="btn btn-secondary btn-sm">{{ t('common.previous') }}</button>
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ page }} / {{ totalPages }}</span>
            <button @click="page < totalPages && loadInvitees(page + 1)" :disabled="page >= totalPages" class="btn btn-secondary btn-sm">{{ t('common.next') }}</button>
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
import { useAppStore } from '@/stores'
import { getInviteCode, listInvitees, getStats, type ReferralInvitee, type ReferralStats } from '@/api/referral'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const inviteCode = ref('')
const copied = ref(false)
const page = ref(1)
const totalPages = ref(1)
const invitees = ref<ReferralInvitee[]>([])
const stats = ref<ReferralStats>({ total_invitees: 0, rewarded_count: 0, pending_count: 0, total_reward_amount: 0 })

onMounted(async () => {
  try {
    const [codeRes, statsRes] = await Promise.all([getInviteCode(), getStats()])
    inviteCode.value = codeRes.invite_code
    stats.value = statsRes
    await loadInvitees(1)
  } catch (err) {
    console.error('Failed to load referral data:', err)
    appStore.showError(t('referral.loadError'))
  } finally {
    loading.value = false
  }
})

async function loadInvitees(p: number) {
  try {
    const res = await listInvitees(p, 20)
    invitees.value = res.items || []
    page.value = res.page
    totalPages.value = res.pages || Math.ceil((res.total || 0) / (res.page_size || 20))
  } catch (err) {
    console.error('Failed to load invitees:', err)
  }
}

function copyCode() {
  navigator.clipboard.writeText(inviteCode.value)
  copied.value = true
  appStore.showSuccess(t('referral.copied'))
  setTimeout(() => { copied.value = false }, 2000)
}

function copyLink() {
  const link = `${window.location.origin}/register?invite=${inviteCode.value}`
  navigator.clipboard.writeText(link)
  copied.value = true
  appStore.showSuccess(t('referral.copied'))
  setTimeout(() => { copied.value = false }, 2000)
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}
</script>
