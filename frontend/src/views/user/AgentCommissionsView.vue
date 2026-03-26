<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-5xl space-y-6">
        <!-- Header -->
        <div class="flex items-center gap-3">
          <router-link to="/agent" class="text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200">
            <Icon name="arrowLeft" size="md" />
          </router-link>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('agent.commissions') }}</h1>
        </div>

        <!-- Filter -->
        <div class="flex gap-2">
          <button v-for="s in ['all', 'pending', 'settled']" :key="s"
            @click="filterStatus = s; loadData(1)"
            :class="filterStatus === s ? 'bg-primary-600 text-white' : 'bg-white text-gray-700 dark:bg-dark-800 dark:text-dark-300'"
            class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium transition hover:shadow-sm dark:border-dark-600">
            {{ s === 'all' ? t('agent.all') : s === 'pending' ? t('agent.pendingStatus') : t('agent.settledStatus') }}
          </button>
        </div>

        <!-- Table -->
        <SlideIn direction="up" :delay="100">
          <GlowCard glow-color="rgb(217, 119, 87)">
            <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft overflow-hidden">
              <div class="overflow-x-auto">
                <table class="w-full">
                  <thead>
                    <tr class="border-b border-gray-100 dark:border-dark-700">
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.userEmail') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.sourceType') }}</th>
                      <th class="px-6 py-3 text-right text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.sourceAmount') }}</th>
                      <th class="px-6 py-3 text-right text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.commissionRate') }}</th>
                      <th class="px-6 py-3 text-right text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.commissionAmount') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.orderStatus') }}</th>
                      <th class="px-6 py-3 text-left text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('agent.createdAt') }}</th>
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
                        {{ t('agent.noCommissions') }}
                      </td>
                    </tr>
                    <tr v-for="c in items" :key="c.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50 transition-colors">
                      <td class="px-6 py-4 text-sm text-gray-900 dark:text-white">{{ c.user_email }}</td>
                      <td class="px-6 py-4 text-sm text-gray-600 dark:text-dark-300">{{ c.source_type }}</td>
                      <td class="px-6 py-4 text-sm text-right font-medium text-gray-900 dark:text-white">${{ c.source_amount.toFixed(2) }}</td>
                      <td class="px-6 py-4 text-sm text-right text-gray-600 dark:text-dark-300">{{ (c.commission_rate * 100).toFixed(1) }}%</td>
                      <td class="px-6 py-4 text-sm text-right font-bold text-primary-600 dark:text-primary-400">${{ c.commission_amount.toFixed(2) }}</td>
                      <td class="px-6 py-4">
                        <span :class="c.status === 'settled' ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400' : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400'"
                          class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
                          {{ c.status === 'settled' ? t('agent.settledStatus') : t('agent.pendingStatus') }}
                        </span>
                      </td>
                      <td class="px-6 py-4 text-sm text-gray-500 dark:text-dark-400">{{ formatDate(c.created_at) }}</td>
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
import { agentAPI, type AgentCommission } from '@/api/agent'
import { FadeIn, SlideIn, GlowCard } from '@/components/animations'

const { t } = useI18n()

const loading = ref(true)
const page = ref(1)
const totalPages = ref(1)
const filterStatus = ref('all')
const items = ref<AgentCommission[]>([])

onMounted(() => loadData(1))

async function loadData(p: number) {
  loading.value = true
  try {
    const filters = filterStatus.value !== 'all' ? { status: filterStatus.value } : undefined
    const res = await agentAPI.listCommissions(p, 20, filters)
    items.value = res.items || []
    page.value = res.page
    totalPages.value = res.pages || Math.ceil((res.total || 0) / (res.page_size || 20))
  } catch (err) {
    console.error('Failed to load commissions:', err)
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}
</script>
