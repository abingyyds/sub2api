<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-3xl space-y-6">
        <!-- Loading -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
        </div>

        <template v-else>
          <!-- Invite Code Card -->
          <SlideIn direction="up" :delay="100">
            <GlowCard glow-color="rgb(217, 119, 87)">
              <div class="overflow-hidden rounded-3xl border-2 border-primary-100 dark:border-primary-800/30 shadow-soft">
                <div class="bg-gradient-to-br from-primary-500 to-primary-600 px-6 py-10 text-center relative overflow-hidden">
                  <div class="absolute top-0 right-0 w-48 h-48 bg-white/5 rounded-full -translate-y-1/2 translate-x-1/2 pointer-events-none"></div>
                  <div class="mb-4 inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 backdrop-blur-sm">
                    <Icon name="gift" size="xl" class="text-white" />
                  </div>
                  <p class="text-sm font-medium text-primary-100">{{ t('referral.yourInviteCode') }}</p>
                  <p class="mt-2 text-3xl font-extrabold tracking-wider text-white">{{ inviteCode }}</p>
                  <div class="mt-5 flex items-center justify-center gap-3">
                    <MagneticButton>
                      <button @click="copyCode" class="inline-flex items-center gap-2 rounded-xl bg-white/20 px-5 py-2.5 text-sm font-bold text-white backdrop-blur-sm transition hover:bg-white/30 border border-white/20">
                        <Icon name="clipboard" size="sm" />
                        {{ t('referral.copyCode') }}
                      </button>
                    </MagneticButton>
                    <MagneticButton>
                      <button @click="copyLink" class="inline-flex items-center gap-2 rounded-xl bg-white/20 px-5 py-2.5 text-sm font-bold text-white backdrop-blur-sm transition hover:bg-white/30 border border-white/20">
                        <Icon name="link" size="sm" />
                        {{ t('referral.copyLink') }}
                      </button>
                    </MagneticButton>
                  </div>
                  <p v-if="copied" class="mt-3 text-sm font-medium text-white/80">{{ t('referral.copied') }}</p>
                </div>
              </div>
            </GlowCard>
          </SlideIn>

          <!-- Stats Cards -->
          <StaggerContainer :stagger-delay="100" :delay="200">
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
              <GlowCard glow-color="rgb(59, 130, 246)">
                <div class="rounded-2xl border-2 border-gray-200 bg-white p-6 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                  <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 dark:bg-blue-900/20">
                    <Icon name="users" size="md" class="text-blue-600 dark:text-blue-400" />
                  </div>
                  <p class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ t('referral.totalInvitees') }}</p>
                  <p class="mt-1 text-3xl font-extrabold text-gray-900 dark:text-white">
                    <AnimatedNumber :value="stats.total_invitees" />
                  </p>
                </div>
              </GlowCard>
              <GlowCard glow-color="rgb(34, 197, 94)">
                <div class="rounded-2xl border-2 border-gray-200 bg-white p-6 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                  <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-green-50 dark:bg-green-900/20">
                    <Icon name="checkCircle" size="md" class="text-green-600 dark:text-green-400" />
                  </div>
                  <p class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ t('referral.rewarded') }}</p>
                  <p class="mt-1 text-3xl font-extrabold text-green-600 dark:text-green-400">
                    <AnimatedNumber :value="stats.rewarded_count" />
                  </p>
                </div>
              </GlowCard>
              <GlowCard glow-color="rgb(217, 119, 87)">
                <div class="rounded-2xl border-2 border-gray-200 bg-white p-6 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                  <div class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-primary-50 dark:bg-primary-900/20">
                    <Icon name="dollar" size="md" class="text-primary-600 dark:text-primary-400" />
                  </div>
                  <p class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ t('referral.totalReward') }}</p>
                  <p class="mt-1 text-3xl font-extrabold text-primary-600 dark:text-primary-400">
                    $<AnimatedNumber :value="stats.total_reward_amount || 0" :decimals="2" />
                  </p>
                </div>
              </GlowCard>
            </div>
          </StaggerContainer>

          <!-- Invitees Table -->
          <SlideIn direction="up" :delay="400">
            <GlowCard glow-color="rgb(99, 102, 241)">
              <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft overflow-hidden">
                <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
                  <h2 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('referral.inviteeList') }}</h2>
                </div>
                <div class="overflow-x-auto">
                  <table class="w-full">
                    <thead>
                      <tr class="border-b border-gray-100 dark:border-dark-700">
                        <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.inviteeEmail') }}</th>
                        <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.status') }}</th>
                        <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.rewardAmount') }}</th>
                        <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('referral.invitedAt') }}</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                      <tr v-if="invitees.length === 0">
                        <td colspan="4" class="px-6 py-12 text-center text-gray-500 dark:text-dark-400">
                          <div class="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-gray-100 dark:bg-dark-800">
                            <Icon name="users" size="md" class="text-gray-400 dark:text-dark-500" />
                          </div>
                          {{ t('referral.noInvitees') }}
                        </td>
                      </tr>
                      <tr v-for="item in invitees" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-dark-800 transition">
                        <td class="px-6 py-4 text-sm font-medium text-gray-900 dark:text-white">{{ item.invitee_email }}</td>
                        <td class="px-6 py-4">
                          <span :class="item.reward_status === 'rewarded'
                            ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                            : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'"
                            class="inline-flex rounded-full px-3 py-1 text-xs font-bold">
                            {{ item.reward_status === 'rewarded' ? t('referral.statusRewarded') : t('referral.statusPending') }}
                          </span>
                        </td>
                        <td class="px-6 py-4 text-sm font-bold text-gray-900 dark:text-white">${{ item.reward_amount?.toFixed(2) || '0.00' }}</td>
                        <td class="px-6 py-4 text-sm text-gray-500 dark:text-dark-400">{{ formatDate(item.created_at) }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <!-- Pagination -->
                <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-gray-100 px-6 py-3 dark:border-dark-700">
                  <button @click="page > 1 && loadInvitees(page - 1)" :disabled="page <= 1" class="rounded-lg border border-gray-200 dark:border-dark-600 px-4 py-2 text-sm font-bold text-gray-700 dark:text-dark-300 hover:bg-gray-50 dark:hover:bg-dark-800 transition disabled:opacity-50">{{ t('common.previous') }}</button>
                  <span class="text-sm font-medium text-gray-500 dark:text-dark-400">{{ page }} / {{ totalPages }}</span>
                  <button @click="page < totalPages && loadInvitees(page + 1)" :disabled="page >= totalPages" class="rounded-lg border border-gray-200 dark:border-dark-600 px-4 py-2 text-sm font-bold text-gray-700 dark:text-dark-300 hover:bg-gray-50 dark:hover:bg-dark-800 transition disabled:opacity-50">{{ t('common.next') }}</button>
                </div>
              </div>
            </GlowCard>
          </SlideIn>
        </template>
      </div>
    </FadeIn>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores'
import { getInviteCode, listInvitees, getStats, type ReferralInvitee, type ReferralStats } from '@/api/referral'
import { FadeIn, SlideIn, StaggerContainer, GlowCard, MagneticButton, AnimatedNumber } from '@/components/animations'

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

<style scoped>
.shadow-soft {
  box-shadow: 0 4px 24px -4px rgba(0, 0, 0, 0.08);
}
</style>
