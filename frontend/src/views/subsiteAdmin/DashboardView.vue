<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-6">
      <div class="flex flex-col gap-2">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ site?.name || '分站概览' }}
        </h1>
        <p class="text-sm text-gray-500 dark:text-dark-400">
          统计范围：{{ stats ? formatRangeStart(stats.range_start) : '近 30 天' }} 至今
        </p>
      </div>

      <div v-if="loading" class="flex justify-center py-16">
        <LoadingSpinner />
      </div>

      <template v-else-if="stats">
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">用户总数</p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.user_count }}</p>
            <p class="text-xs text-gray-500 dark:text-dark-400">
              近 30 天活跃：{{ stats.active_users }}
            </p>
          </div>
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">请求数（近 30 天）</p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.requests }}</p>
            <p class="text-xs text-gray-500 dark:text-dark-400">
              总费用：${{ (stats.total_cost || 0).toFixed(4) }}
            </p>
          </div>
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">支付收入（近 30 天）</p>
            <p class="mt-2 text-2xl font-bold text-emerald-600 dark:text-emerald-400">
              ¥{{ fenToYuan(stats.revenue_fen) }}
            </p>
          </div>
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">分站池余额</p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              ¥{{ fenToYuan(stats.pool_balance_fen) }}
            </p>
          </div>
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">累计充值</p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              ¥{{ fenToYuan(stats.total_topup_fen) }}
            </p>
          </div>
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">累计消耗</p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              ¥{{ fenToYuan(stats.total_consumed_fen) }}
            </p>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { subSiteAdminAPI, type SubSiteAdminDashboardStats } from '@/api/subsiteAdmin'
import type { OwnedSubSite } from '@/api/subsite'
import { useAppStore } from '@/stores'

const route = useRoute()
const appStore = useAppStore()

const loading = ref(true)
const stats = ref<SubSiteAdminDashboardStats | null>(null)
const site = ref<OwnedSubSite | null>(null)

const siteId = computed(() => Number(route.params.siteId))

async function loadData() {
  if (!siteId.value) return
  loading.value = true
  try {
    const [statsRes, siteRes] = await Promise.all([
      subSiteAdminAPI.getDashboard(siteId.value),
      subSiteAdminAPI.getSite(siteId.value),
    ])
    stats.value = statsRes
    site.value = siteRes
  } catch (error: any) {
    appStore.showError(error?.message || '加载分站概览失败')
  } finally {
    loading.value = false
  }
}

function fenToYuan(value: number): string {
  const num = Number(value || 0) / 100
  return num.toFixed(num % 1 === 0 ? 0 : 2)
}

function formatRangeStart(value: string): string {
  if (!value) return ''
  try {
    return new Date(value).toLocaleDateString('zh-CN')
  } catch {
    return value
  }
}

onMounted(loadData)
watch(siteId, loadData)
</script>
