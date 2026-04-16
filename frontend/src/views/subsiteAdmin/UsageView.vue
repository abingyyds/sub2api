<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-4">
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站用量</h1>
        <div class="flex items-center gap-2">
          <input
            v-model="modelFilter"
            type="text"
            class="input w-48"
            placeholder="按模型过滤（可选）"
            @keyup.enter="handleFilter"
          />
          <button class="btn btn-secondary" :disabled="loading" @click="handleFilter">
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
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">用户</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">模型</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">输入</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">输出</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">缓存写</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">缓存读</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">原始费用</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">实际费用</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-for="row in items" :key="row.id">
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ formatDate(row.created_at) }}</td>
                <td class="px-4 py-3 text-xs">{{ row.user_email || `#${row.user_id}` }}</td>
                <td class="px-4 py-3 text-xs font-mono">{{ row.model }}</td>
                <td class="px-4 py-3 text-right text-xs">{{ row.input_tokens }}</td>
                <td class="px-4 py-3 text-right text-xs">{{ row.output_tokens }}</td>
                <td class="px-4 py-3 text-right text-xs">{{ row.cache_creation_tokens }}</td>
                <td class="px-4 py-3 text-right text-xs">{{ row.cache_read_tokens }}</td>
                <td class="px-4 py-3 text-right text-xs text-gray-500 dark:text-dark-400">
                  ${{ (row.total_cost || 0).toFixed(6) }}
                </td>
                <td class="px-4 py-3 text-right text-sm font-medium">${{ (row.actual_cost || 0).toFixed(6) }}</td>
              </tr>
              <tr v-if="!loading && items.length === 0">
                <td colspan="9" class="px-4 py-10 text-center text-xs text-gray-400">暂无用量</td>
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
import { subSiteAdminAPI, type SubSiteAdminUsage } from '@/api/subsiteAdmin'
import { useAppStore } from '@/stores'

const route = useRoute()
const appStore = useAppStore()

const siteId = computed(() => Number(route.params.siteId))

const loading = ref(false)
const modelFilter = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const items = ref<SubSiteAdminUsage[]>([])

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

async function loadData() {
  if (!siteId.value) return
  loading.value = true
  try {
    const res = await subSiteAdminAPI.listUsage(siteId.value, page.value, pageSize.value, modelFilter.value || undefined)
    items.value = res.items || []
    total.value = res.total || 0
  } catch (error: any) {
    appStore.showError(error?.message || '加载用量失败')
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
  modelFilter.value = ''
  loadData()
})
</script>
