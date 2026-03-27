<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-5xl space-y-6">
        <!-- Loading -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
        </div>

        <template v-else>
          <!-- Not Agent: Apply Section -->
          <template v-if="!status.is_agent">
            <SlideIn direction="up" :delay="100">
              <GlowCard glow-color="rgb(217, 119, 87)">
                <div class="rounded-3xl border-2 border-primary-100 dark:border-primary-800/30 shadow-soft overflow-hidden">
                  <div class="bg-gradient-to-br from-primary-500 to-primary-600 px-6 py-10 text-center relative overflow-hidden">
                    <div class="absolute top-0 right-0 w-48 h-48 bg-white/5 rounded-full -translate-y-1/2 translate-x-1/2 pointer-events-none"></div>
                    <div class="mb-4 inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 backdrop-blur-sm">
                      <Icon name="users" size="xl" class="text-white" />
                    </div>
                    <h2 class="text-2xl font-extrabold text-white">{{ t('agent.becomeAgent') }}</h2>
                    <p class="mt-2 text-sm text-primary-100">{{ t('agent.becomeAgentDesc') }}</p>
                    <div class="mt-6 space-y-3 max-w-md mx-auto text-left">
                      <div>
                        <label class="block text-xs font-medium text-primary-100 mb-1">{{ t('agent.applyContact') }} *</label>
                        <input
                          v-model="applyForm.contact"
                          :placeholder="t('agent.applyContactPlaceholder')"
                          class="w-full rounded-xl border border-white/20 bg-white/10 px-4 py-2.5 text-sm text-white placeholder-white/50 backdrop-blur-sm focus:border-white/40 focus:outline-none"
                        />
                      </div>
                      <div>
                        <label class="block text-xs font-medium text-primary-100 mb-1">{{ t('agent.applySocial') }}</label>
                        <input
                          v-model="applyForm.social"
                          :placeholder="t('agent.applySocialPlaceholder')"
                          class="w-full rounded-xl border border-white/20 bg-white/10 px-4 py-2.5 text-sm text-white placeholder-white/50 backdrop-blur-sm focus:border-white/40 focus:outline-none"
                        />
                      </div>
                      <div>
                        <label class="block text-xs font-medium text-primary-100 mb-1">{{ t('agent.applyPromotion') }}</label>
                        <textarea
                          v-model="applyForm.promotion"
                          :placeholder="t('agent.applyPromotionPlaceholder')"
                          rows="2"
                          class="w-full rounded-xl border border-white/20 bg-white/10 px-4 py-2.5 text-sm text-white placeholder-white/50 backdrop-blur-sm focus:border-white/40 focus:outline-none"
                        ></textarea>
                      </div>
                    </div>
                    <div class="mt-4">
                      <MagneticButton>
                        <button @click="handleApply" :disabled="applying || !applyForm.contact.trim()" class="inline-flex items-center gap-2 rounded-xl bg-white/20 px-6 py-2.5 text-sm font-bold text-white backdrop-blur-sm transition hover:bg-white/30 border border-white/20 disabled:opacity-50">
                          <Icon name="arrowRight" size="sm" />
                          {{ applying ? t('agent.applying') : t('agent.applyNow') }}
                        </button>
                      </MagneticButton>
                    </div>
                  </div>
                </div>
              </GlowCard>
            </SlideIn>
          </template>

          <!-- Pending Status -->
          <template v-else-if="status.agent_status === 'pending'">
            <SlideIn direction="up" :delay="100">
              <div class="rounded-2xl border-2 border-yellow-200 bg-yellow-50 p-8 text-center dark:border-yellow-800/30 dark:bg-yellow-900/10 shadow-soft">
                <div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-yellow-100 dark:bg-yellow-900/30">
                  <Icon name="clock" size="lg" class="text-yellow-600 dark:text-yellow-400" />
                </div>
                <h2 class="text-xl font-bold text-yellow-800 dark:text-yellow-300">{{ t('agent.pendingTitle') }}</h2>
                <p class="mt-2 text-sm text-yellow-600 dark:text-yellow-400">{{ t('agent.pendingDesc') }}</p>
              </div>
            </SlideIn>
          </template>

          <!-- Rejected Status -->
          <template v-else-if="status.agent_status === 'rejected'">
            <SlideIn direction="up" :delay="100">
              <div class="rounded-2xl border-2 border-red-200 bg-red-50 p-8 text-center dark:border-red-800/30 dark:bg-red-900/10 shadow-soft">
                <div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-red-100 dark:bg-red-900/30">
                  <Icon name="xCircle" size="lg" class="text-red-600 dark:text-red-400" />
                </div>
                <h2 class="text-xl font-bold text-red-800 dark:text-red-300">{{ t('agent.rejectedTitle') }}</h2>
                <p class="mt-2 text-sm text-red-600 dark:text-red-400">{{ t('agent.rejectedDesc') }}</p>
              </div>
            </SlideIn>
          </template>

          <!-- Approved Agent Dashboard -->
          <template v-else-if="status.agent_status === 'approved'">
            <!-- Invite Link Card -->
            <SlideIn direction="up" :delay="100">
              <GlowCard glow-color="rgb(217, 119, 87)">
                <div class="overflow-hidden rounded-3xl border-2 border-primary-100 dark:border-primary-800/30 shadow-soft">
                  <div class="bg-gradient-to-br from-primary-500 to-primary-600 px-6 py-8 text-center relative overflow-hidden">
                    <div class="absolute top-0 right-0 w-48 h-48 bg-white/5 rounded-full -translate-y-1/2 translate-x-1/2 pointer-events-none"></div>
                    <p class="text-sm font-medium text-primary-100">{{ t('agent.yourInviteLink') }}</p>
                    <p class="mt-2 text-lg font-bold text-white break-all">{{ inviteLink }}</p>
                    <div class="mt-4 flex items-center justify-center gap-3">
                      <MagneticButton>
                        <button @click="copyLink" class="inline-flex items-center gap-2 rounded-xl bg-white/20 px-5 py-2.5 text-sm font-bold text-white backdrop-blur-sm transition hover:bg-white/30 border border-white/20">
                          <Icon name="link" size="sm" />
                          {{ t('agent.copyLink') }}
                        </button>
                      </MagneticButton>
                    </div>
                  </div>
                </div>
              </GlowCard>
            </SlideIn>

            <!-- Stats Cards -->
            <StaggerContainer :stagger-delay="100" :delay="200">
              <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
                <GlowCard glow-color="rgb(59, 130, 246)">
                  <div class="rounded-2xl border-2 border-gray-200 bg-white p-5 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                    <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('agent.totalSubUsers') }}</p>
                    <p class="mt-1 text-2xl font-extrabold text-gray-900 dark:text-white">
                      <AnimatedNumber :value="dashboard.total_sub_users" />
                    </p>
                  </div>
                </GlowCard>
                <GlowCard glow-color="rgb(34, 197, 94)">
                  <div class="rounded-2xl border-2 border-gray-200 bg-white p-5 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                    <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('agent.totalRecharge') }}</p>
                    <p class="mt-1 text-2xl font-extrabold text-green-600 dark:text-green-400">
                      $<AnimatedNumber :value="dashboard.total_recharge" :decimals="2" />
                    </p>
                  </div>
                </GlowCard>
                <GlowCard glow-color="rgb(239, 68, 68)">
                  <div class="rounded-2xl border-2 border-gray-200 bg-white p-5 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                    <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('agent.totalConsumed') }}</p>
                    <p class="mt-1 text-2xl font-extrabold text-red-600 dark:text-red-400">
                      $<AnimatedNumber :value="dashboard.total_consumed" :decimals="2" />
                    </p>
                  </div>
                </GlowCard>
                <GlowCard glow-color="rgb(217, 119, 87)">
                  <div class="rounded-2xl border-2 border-gray-200 bg-white p-5 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                    <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('agent.totalCommission') }}</p>
                    <p class="mt-1 text-2xl font-extrabold text-primary-600 dark:text-primary-400">
                      $<AnimatedNumber :value="dashboard.total_commission" :decimals="2" />
                    </p>
                  </div>
                </GlowCard>
                <GlowCard glow-color="rgb(249, 115, 22)">
                  <div class="rounded-2xl border-2 border-gray-200 bg-white p-5 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                    <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('agent.pendingCommission') }}</p>
                    <p class="mt-1 text-2xl font-extrabold text-orange-600 dark:text-orange-400">
                      $<AnimatedNumber :value="dashboard.pending_commission" :decimals="2" />
                    </p>
                  </div>
                </GlowCard>
                <GlowCard glow-color="rgb(99, 102, 241)">
                  <div class="rounded-2xl border-2 border-gray-200 bg-white p-5 text-center dark:border-dark-700 dark:bg-dark-900 shadow-soft">
                    <p class="text-xs font-medium text-gray-500 dark:text-dark-400">{{ t('agent.settledCommission') }}</p>
                    <p class="mt-1 text-2xl font-extrabold text-indigo-600 dark:text-indigo-400">
                      $<AnimatedNumber :value="dashboard.settled_commission" :decimals="2" />
                    </p>
                  </div>
                </GlowCard>
              </div>
            </StaggerContainer>

            <!-- Quick Nav -->
            <SlideIn direction="up" :delay="400">
              <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
                <router-link to="/agent/sub-users" class="group rounded-2xl border-2 border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-900 shadow-soft transition hover:border-primary-300 dark:hover:border-primary-700 hover:shadow-md">
                  <div class="flex items-center gap-4">
                    <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 dark:bg-blue-900/20">
                      <Icon name="users" size="md" class="text-blue-600 dark:text-blue-400" />
                    </div>
                    <div>
                      <p class="font-bold text-gray-900 dark:text-white">{{ t('agent.subUsers.label') }}</p>
                      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('agent.subUsers.desc') }}</p>
                    </div>
                  </div>
                </router-link>
                <router-link to="/agent/financial-logs" class="group rounded-2xl border-2 border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-900 shadow-soft transition hover:border-primary-300 dark:hover:border-primary-700 hover:shadow-md">
                  <div class="flex items-center gap-4">
                    <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-green-50 dark:bg-green-900/20">
                      <Icon name="dollar" size="md" class="text-green-600 dark:text-green-400" />
                    </div>
                    <div>
                      <p class="font-bold text-gray-900 dark:text-white">{{ t('agent.financialLogs.label') }}</p>
                      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('agent.financialLogs.desc') }}</p>
                    </div>
                  </div>
                </router-link>
                <router-link to="/agent/commissions" class="group rounded-2xl border-2 border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-900 shadow-soft transition hover:border-primary-300 dark:hover:border-primary-700 hover:shadow-md">
                  <div class="flex items-center gap-4">
                    <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-primary-50 dark:bg-primary-900/20">
                      <Icon name="document" size="md" class="text-primary-600 dark:text-primary-400" />
                    </div>
                    <div>
                      <p class="font-bold text-gray-900 dark:text-white">{{ t('agent.commissions.label') }}</p>
                      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('agent.commissions.desc') }}</p>
                    </div>
                  </div>
                </router-link>
              </div>
            </SlideIn>
          </template>
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
import { agentAPI, type AgentStatus, type AgentDashboardStats } from '@/api/agent'
import { FadeIn, SlideIn, StaggerContainer, GlowCard, MagneticButton, AnimatedNumber } from '@/components/animations'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const applying = ref(false)
const applyForm = ref({
  contact: '',
  social: '',
  promotion: ''
})
const inviteLink = ref('')
const status = ref<AgentStatus>({
  is_agent: false,
  agent_status: '',
  agent_commission_rate: 0,
  agent_note: '',
  agent_approved_at: null
})
const dashboard = ref<AgentDashboardStats>({
  total_sub_users: 0,
  total_recharge: 0,
  total_consumed: 0,
  total_commission: 0,
  pending_commission: 0,
  settled_commission: 0
})

onMounted(async () => {
  try {
    status.value = await agentAPI.getStatus()
    if (status.value.agent_status === 'approved') {
      const [dashData, linkData] = await Promise.all([
        agentAPI.getDashboard(),
        agentAPI.getLink()
      ])
      dashboard.value = dashData
      inviteLink.value = linkData.invite_url || `${window.location.origin}/register?invite=${linkData.invite_code}`
    }
  } catch (err) {
    console.error('Failed to load agent data:', err)
  } finally {
    loading.value = false
  }
})

async function handleApply() {
  applying.value = true
  try {
    await agentAPI.apply({
      contact: applyForm.value.contact,
      social: applyForm.value.social || undefined,
      promotion: applyForm.value.promotion || undefined
    })
    status.value.is_agent = true
    status.value.agent_status = 'pending'
    appStore.showSuccess(t('agent.applySuccess'))
  } catch (err) {
    console.error('Apply failed:', err)
    appStore.showError(t('agent.applyError'))
  } finally {
    applying.value = false
  }
}

function copyLink() {
  navigator.clipboard.writeText(inviteLink.value)
  appStore.showSuccess(t('agent.linkCopied'))
}
</script>
