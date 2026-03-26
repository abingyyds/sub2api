<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-5xl space-y-6">
        <!-- Header -->
        <div class="flex items-center gap-3">
          <router-link to="/agent" class="text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200">
            <Icon name="arrowLeft" size="md" />
          </router-link>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('agent.financialLogs.title') }}</h1>
        </div>

        <!-- Table -->
        <SlideIn direction="up" :delay="100">
          <GlowCard glow-color="rgb(34, 197, 94)">
            <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft overflow-hidden">
              <div class="overflow-x-auto">
                <table class="w-full">
                  <thead>
                    <tr class="border-b border-gray-100 dark:border-dark-700">
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.orderNo') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.userEmail') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.planKey') }}</th>
                      <th class="px-6 py-3 text-right text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.amount') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.orderStatus') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.paidAt') }}</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                    <tr v-if="loading">
                      <td colspan="6" class="px-6 py-12 text-center">
                        <div class="inline-block h-6 w-6 animate-spin rounded-full border-b-2 border-primary-600"></div>
                      </td>
                    </tr>
                    <tr v-else-if="items.length === 0">
                      <td colspan="6" class="px-6 py-12 text-center text-gray-500 dark:text-dark-400">
                        {{ t('agent.noFinancialLogs') }}
                      </td>
                    </tr>
                    <tr v-for="log in items" :key="log.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50 transition-colors">
                      <td class="px-6 py-4 text-sm font-mono text-gray-900 dark:text-white">{{ log.order_no }}</td>
                      <td class="px-6 py-4 text-sm text-gray-600 dark:text-dark-300">{{ log.user_email }}</td>
                      <td class="px-6 py-4 text-sm text-gray-600 dark:text-dark-300">{{ log.plan_key || '-' }}</td>
                      <td class="px-6 py-4 text-sm text-right font-medium text-green-600 dark:text-green-400">
                        {{ log.balance_amount > 0 ? `$${log.balance_amount.toFixed(2)}` : `¥${(log.amount_fen / 100).toFixed(2)}` }}
                      </td>
                      <td class="px-6 py-4">
                        <span :class="statusClass(log.status)" class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                          {{ log.status }}
                        </span>
                      </td>
                      <td class="px-6 py-4 text-sm text-gray-500 dark:text-dark-400">{{ formatDate(log.paid_at || log.created_at) }}</td>
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
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { agentAPI, type AgentFinancialLog } from '@/api/agent'
import { FadeIn, SlideIn, GlowCard } from '@/components/animations'

const { t } = useI18n()

const loading = ref(true)
const page = ref(1)
const totalPages = ref(1)
const items = ref<AgentFinancialLog[]>([])

onMounted(() => loadData(1))

async function loadData(p: number) {
  loading.value = true
  try {
    const res = await agentAPI.listFinancialLogs(p, 20)
    items.value = res.items || []
    page.value = res.page
    totalPages.value = res.pages || Math.ceil((res.total || 0) / (res.page_size || 20))
  } catch (err) {
    console.error('Failed to load financial logs:', err)
  } finally {
    loading.value = false
  }
}

function statusClass(status: string): string {
  switch (status) {
    case 'paid': return 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
    case 'pending': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400'
    case 'expired': return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}
</script>
