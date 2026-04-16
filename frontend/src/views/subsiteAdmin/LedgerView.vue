<template>
  <AppLayout>
    <div class="mx-auto max-w-5xl space-y-4">
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站池流水</h1>
        <div class="flex items-center gap-2">
          <select v-model="txTypeFilter" class="input w-40" @change="handleFilter">
            <option value="">全部类型</option>
            <option value="topup">充值</option>
            <option value="consume">消耗</option>
            <option value="offline_topup">线下加余额</option>
            <option value="refund">退款</option>
          </select>
          <button class="btn btn-secondary" :disabled="loading" @click="loadData">
            {{ loading ? '...' : '刷新' }}
          </button>
        </div>
      </div>

      <div class="rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 text-sm dark:divide-dark-700">
            <thead class="bg-gray-50 dark:bg-dark-800/60">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">时间</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">类型</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">变动</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">变动后</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">关联</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">备注</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-for="row in items" :key="row.id">
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ formatDate(row.created_at) }}</td>
                <td class="px-4 py-3 text-xs">{{ row.tx_type }}</td>
                <td class="px-4 py-3 text-right text-xs" :class="row.delta_fen >= 0 ? 'text-emerald-600' : 'text-red-600'">
                  {{ row.delta_fen >= 0 ? '+' : '' }}¥{{ fenToYuan(row.delta_fen) }}
                </td>
                <td class="px-4 py-3 text-right text-xs">¥{{ fenToYuan(row.balance_after_fen) }}</td>
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">
                  <span v-if="row.related_user_id">user #{{ row.related_user_id }}</span>
                  <span v-else-if="row.related_order_id">order #{{ row.related_order_id }}</span>
                  <span v-else-if="row.related_usage_log_id">usage #{{ row.related_usage_log_id }}</span>
                  <span v-else>-</span>
                </td>
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ row.note || '-' }}</td>
              </tr>
              <tr v-if="!loading && items.length === 0">
                <td colspan="6" class="px-4 py-10 text-center text-xs text-gray-400">暂无流水</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-if="total > 0" class="flex items-center justify-end gap-2 border-t border-gray-100 px-4 py-3 dark:border-dark-800">
          <button class="btn btn-secondary btn-sm" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
          <span class="text-xs text-gray-500 dark:text-dark-400">第 {{ page }} / {{ totalPages }} 页（共 {{ total }}）</span>
          <button class="btn btn-secondary btn-sm" :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { subSiteAdminAPI } from '@/api/subsiteAdmin'
import type { OwnedSubSiteLedgerEntry } from '@/api/subsite'
import { useAppStore } from '@/stores'

const route = useRoute()
const appStore = useAppStore()

const siteId = computed(() => Number(route.params.siteId))

const loading = ref(false)
const txTypeFilter = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const items = ref<OwnedSubSiteLedgerEntry[]>([])

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

async function loadData() {
  if (!siteId.value) return
  loading.value = true
  try {
    const res = await subSiteAdminAPI.listLedger(siteId.value, page.value, pageSize.value, txTypeFilter.value || undefined)
    items.value = res.items || []
    total.value = res.total || 0
  } catch (error: any) {
    appStore.showError(error?.message || '加载流水失败')
  } finally {
    loading.value = false
  }
}

function changePage(newPage: number) {
  if (newPage < 1 || newPage > totalPages.value) return
  page.value = newPage
  loadData()
}

function handleFilter() {
  page.value = 1
  loadData()
}

function fenToYuan(value: number): string {
  const num = Number(value || 0) / 100
  return num.toFixed(num % 1 === 0 ? 0 : 2)
}

function formatDate(value: string): string {
  if (!value) return '-'
  try {
    return new Date(value).toLocaleString('zh-CN', { hour12: false })
  } catch {
    return value
  }
}

onMounted(loadData)
watch(siteId, () => {
  page.value = 1
  txTypeFilter.value = ''
  loadData()
})
</script>
