<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-5xl space-y-6">
        <!-- Header -->
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <router-link to="/agent" class="text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200">
              <Icon name="arrowLeft" size="md" />
            </router-link>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('agent.subUsers.title') }}</h1>
          </div>
        </div>

        <!-- Table -->
        <SlideIn direction="up" :delay="100">
          <GlowCard glow-color="rgb(59, 130, 246)">
            <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft overflow-hidden">
              <div class="overflow-x-auto">
                <table class="w-full">
                  <thead>
                    <tr class="border-b border-gray-100 dark:border-dark-700">
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.email') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.username') }}</th>
                      <th class="px-6 py-3 text-right text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.commissionRate') }}</th>
                      <th class="px-6 py-3 text-right text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.totalRecharge') }}</th>
                      <th class="px-6 py-3 text-right text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.totalConsumed') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.registeredAt') }}</th>
                      <th class="px-6 py-3 text-center text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.actions') }}</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                    <tr v-if="loading">
                      <td colspan="7" class="px-6 py-12 text-center">
                        <div class="inline-block h-6 w-6 animate-spin rounded-full border-b-2 border-primary-600"></div>
                      </td>
                    </tr>
                    <tr v-else-if="items.length === 0">
                      <td colspan="7" class="px-6 py-12 text-center text-gray-500 dark:text-dark-400">
                        {{ t('agent.noSubUsers') }}
                      </td>
                    </tr>
                    <tr v-for="user in items" :key="user.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50 transition-colors">
                      <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">
                        {{ user.email }}
                        <span v-if="user.is_agent" class="ml-1 inline-flex items-center rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700 dark:bg-blue-900/30 dark:text-blue-300">{{ t('agent.isSubAgent') }}</span>
                      </td>
                      <td class="px-6 py-4 text-sm text-gray-600 dark:text-dark-300">{{ user.username || '-' }}</td>
                      <td class="px-6 py-4 text-sm text-right font-medium text-gray-900 dark:text-white">
                        <span v-if="user.commission_rate != null">{{ (user.commission_rate * 100).toFixed(1) }}%</span>
                        <span v-else class="text-gray-400 dark:text-dark-500">{{ t('agent.rateDefault') }}</span>
                      </td>
                      <td class="px-6 py-4 text-sm text-right font-medium text-green-600 dark:text-green-400">${{ user.total_recharge.toFixed(2) }}</td>
                      <td class="px-6 py-4 text-sm text-right font-medium text-red-600 dark:text-red-400">${{ user.total_consumed.toFixed(2) }}</td>
                      <td class="px-6 py-4 text-sm text-gray-500 dark:text-dark-400">{{ formatDate(user.created_at) }}</td>
                      <td class="px-6 py-4 text-center">
                        <button @click="openSetRate(user)" class="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 transition hover:bg-gray-50 dark:border-dark-600 dark:text-dark-300 dark:hover:bg-dark-800">
                          {{ t('agent.setRate') }}
                        </button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <!-- Pagination -->
              <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-gray-100 px-6 py-3 dark:border-dark-700">
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('agent.page') }} {{ page }} / {{ totalPages }}</p>
                <div class="flex gap-2">
                  <button @click="loadData(page - 1)" :disabled="page <= 1" class="rounded-lg border border-gray-200 px-3 py-1.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 disabled:opacity-50 dark:border-dark-600 dark:text-dark-300 dark:hover:bg-dark-800">
                    {{ t('agent.prev') }}
                  </button>
                  <button @click="loadData(page + 1)" :disabled="page >= totalPages" class="rounded-lg border border-gray-200 px-3 py-1.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 disabled:opacity-50 dark:border-dark-600 dark:text-dark-300 dark:hover:bg-dark-800">
                    {{ t('agent.next') }}
                  </button>
                </div>
              </div>
            </div>
          </GlowCard>
        </SlideIn>
      </div>
    </FadeIn>

    <!-- Set Rate Modal -->
    <Teleport to="body">
      <div v-if="showRateModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="fixed inset-0 bg-black/50" @click="showRateModal = false"></div>
        <div class="relative w-full max-w-md rounded-xl bg-white p-6 shadow-xl dark:bg-dark-800">
          <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">{{ t('agent.setRateTitle') }}</h3>
          <div class="space-y-4">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ editingUser?.email }}</p>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1">{{ t('agent.commissionRate') }} (%)</label>
              <input v-model.number="editRate" type="number" min="0" :max="100" step="0.1" class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm dark:border-dark-600 dark:bg-dark-700 dark:text-white" />
            </div>
          </div>
          <div class="mt-6 flex justify-end gap-3">
            <button @click="showRateModal = false" class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50 dark:border-dark-600 dark:text-dark-300 dark:hover:bg-dark-800">{{ t('common.cancel') }}</button>
            <button @click="saveRate" :disabled="saving" class="rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-primary-700 disabled:opacity-50">{{ saving ? t('common.saving') : t('common.save') }}</button>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { agentAPI, type AgentSubUser } from '@/api/agent'
import { FadeIn, SlideIn, GlowCard } from '@/components/animations'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const saving = ref(false)
const page = ref(1)
const totalPages = ref(1)
const items = ref<AgentSubUser[]>([])

// Set rate modal
const showRateModal = ref(false)
const editRate = ref(0)
const editingUser = ref<AgentSubUser | null>(null)

onMounted(() => loadData(1))

async function loadData(p: number) {
  loading.value = true
  try {
    const res = await agentAPI.listSubUsers(p, 20)
    items.value = res.items || []
    page.value = res.page
    totalPages.value = res.pages || Math.ceil((res.total || 0) / (res.page_size || 20))
  } catch (err) {
    console.error('Failed to load sub users:', err)
  } finally {
    loading.value = false
  }
}

function openSetRate(user: AgentSubUser) {
  editingUser.value = user
  editRate.value = user.commission_rate != null ? user.commission_rate * 100 : 0
  showRateModal.value = true
}

async function saveRate() {
  if (!editingUser.value) return
  saving.value = true
  try {
    await agentAPI.setSubUserRate(editingUser.value.id, editRate.value / 100)
    appStore.showSuccess(t('agent.rateSaved'))
    showRateModal.value = false
    loadData(page.value)
  } catch {
    appStore.showError(t('agent.rateError'))
  } finally {
    saving.value = false
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}
</script>
