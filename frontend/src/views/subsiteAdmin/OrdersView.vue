<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-4">
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站订单</h1>
        <div class="flex items-center gap-2">
          <select v-model="statusFilter" class="input w-36" @change="handleFilter">
            <option value="">全部状态</option>
            <option value="pending">待支付</option>
            <option value="paid">已支付</option>
            <option value="cancelled">已取消</option>
            <option value="refunded">已退款</option>
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
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">订单号</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">用户</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">类型</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">套餐</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">金额</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">状态</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">支付方式</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">支付时间</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">创建时间</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-for="o in items" :key="o.id">
                <td class="px-4 py-3 font-mono text-xs text-gray-700 dark:text-dark-200">{{ o.order_no }}</td>
                <td class="px-4 py-3 text-xs">{{ o.user_email || `#${o.user_id}` }}</td>
                <td class="px-4 py-3 text-xs">{{ o.order_type }}</td>
                <td class="px-4 py-3 text-xs">{{ o.plan_key || '-' }}</td>
                <td class="px-4 py-3 text-right text-sm font-medium">¥{{ fenToYuan(o.amount_fen) }}</td>
                <td class="px-4 py-3 text-xs">
                  <span :class="['badge', statusBadgeClass(o.status)]">{{ o.status }}</span>
                </td>
                <td class="px-4 py-3 text-xs">{{ o.pay_method || '-' }}</td>
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ formatDate(o.paid_at) }}</td>
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ formatDate(o.created_at) }}</td>
              </tr>
              <tr v-if="!loading && items.length === 0">
                <td colspan="9" class="px-4 py-10 text-center text-xs text-gray-400">暂无订单</td>
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
import { subSiteAdminAPI, type SubSiteAdminOrder } from '@/api/subsiteAdmin'
import { useAppStore } from '@/stores'

const route = useRoute()
const appStore = useAppStore()

const siteId = computed(() => Number(route.params.siteId))

const loading = ref(false)
const statusFilter = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const items = ref<SubSiteAdminOrder[]>([])

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

async function loadData() {
  if (!siteId.value) return
  loading.value = true
  try {
    const res = await subSiteAdminAPI.listOrders(siteId.value, page.value, pageSize.value, statusFilter.value || undefined)
    items.value = res.items || []
    total.value = res.total || 0
  } catch (error: any) {
    appStore.showError(error?.message || '加载订单失败')
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

function statusBadgeClass(status: string): string {
  switch (status) {
    case 'paid':
      return 'badge-green'
    case 'pending':
      return 'badge-yellow'
    case 'cancelled':
    case 'refunded':
      return 'badge-gray'
    default:
      return 'badge-gray'
  }
}

function fenToYuan(value: number): string {
  const num = Number(value || 0) / 100
  return num.toFixed(num % 1 === 0 ? 0 : 2)
}

function formatDate(value?: string | null): string {
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
  statusFilter.value = ''
  loadData()
})
</script>
