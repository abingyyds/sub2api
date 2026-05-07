<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-4">
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站用户</h1>
        <div class="flex items-center gap-2">
          <input
            v-model="search"
            type="text"
            class="input w-56"
            placeholder="搜索邮箱/用户名"
            @keyup.enter="handleSearch"
          />
          <button class="btn btn-secondary" :disabled="loading" @click="handleSearch">
            {{ loading ? '...' : '搜索' }}
          </button>
        </div>
      </div>

      <div class="rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 text-sm dark:divide-dark-700">
            <thead class="bg-gray-50 dark:bg-dark-800/60">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">ID</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">邮箱</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">用户名</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">角色</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">状态</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">余额</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase">绑定时间</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-for="u in items" :key="u.id">
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ u.id }}</td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-white">{{ u.email }}</td>
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ u.username || '-' }}</td>
                <td class="px-4 py-3 text-xs">{{ u.role }}</td>
                <td class="px-4 py-3 text-xs">
                  <span :class="['badge', u.status === 'active' ? 'badge-green' : 'badge-gray']">{{ u.status }}</span>
                </td>
                <td class="px-4 py-3 text-right text-sm">¥{{ yuanAmount(u.balance) }}</td>
                <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ formatDate(u.bound_at) }}</td>
                <td class="px-4 py-3 text-right">
                  <button class="btn btn-secondary btn-sm" @click="openTopup(u)">加余额</button>
                </td>
              </tr>
              <tr v-if="!loading && items.length === 0">
                <td colspan="8" class="px-4 py-10 text-center text-xs text-gray-400">暂无用户</td>
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

    <BaseDialog :show="showTopup" title="线下给用户加余额" @close="showTopup = false">
      <div v-if="topupTarget" class="space-y-4">
        <div class="rounded-2xl border border-gray-200 px-4 py-3 text-sm dark:border-dark-700">
          <div class="font-medium text-gray-900 dark:text-white">{{ topupTarget.email }}</div>
          <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">当前余额：¥{{ yuanAmount(topupTarget.balance) }}</div>
          <div class="text-xs text-gray-500 dark:text-dark-400">说明：给用户加余额的同时会从分站池按等额扣除。</div>
        </div>
        <div>
          <label class="input-label">加余额金额（元）</label>
          <input v-model.number="topupAmountYuan" type="number" min="0" step="1" class="input mt-1 w-full" />
        </div>
        <div>
          <label class="input-label">备注（可选）</label>
          <input v-model="topupNote" class="input mt-1 w-full" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="showTopup = false">取消</button>
          <button class="btn btn-primary" :disabled="topupSubmitting" @click="handleTopupSubmit">
            {{ topupSubmitting ? '提交中...' : '确认加余额' }}
          </button>
        </div>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { subSiteAdminAPI, type SubSiteAdminUser } from '@/api/subsiteAdmin'
import { useAppStore } from '@/stores'

const route = useRoute()
const appStore = useAppStore()

const siteId = computed(() => Number(route.params.siteId))

const loading = ref(false)
const search = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const items = ref<SubSiteAdminUser[]>([])

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

async function loadData() {
  if (!siteId.value) return
  loading.value = true
  try {
    const res = await subSiteAdminAPI.listUsers(siteId.value, page.value, pageSize.value, search.value || undefined)
    items.value = res.items || []
    total.value = res.total || 0
  } catch (error: any) {
    appStore.showError(error?.message || '加载用户失败')
  } finally {
    loading.value = false
  }
}

function changePage(newPage: number) {
  if (newPage < 1 || newPage > totalPages.value) return
  page.value = newPage
  loadData()
}

function handleSearch() {
  page.value = 1
  loadData()
}

const showTopup = ref(false)
const topupTarget = ref<SubSiteAdminUser | null>(null)
const topupAmountYuan = ref<number>(0)
const topupNote = ref('')
const topupSubmitting = ref(false)

function openTopup(user: SubSiteAdminUser) {
  topupTarget.value = user
  topupAmountYuan.value = 0
  topupNote.value = ''
  showTopup.value = true
}

async function handleTopupSubmit() {
  if (!topupTarget.value || !siteId.value) return
  const amountFen = Math.round((topupAmountYuan.value || 0) * 100)
  if (amountFen <= 0) {
    appStore.showError('请输入有效金额')
    return
  }
  topupSubmitting.value = true
  try {
    await subSiteAdminAPI.offlineTopup(siteId.value, {
      user_id: topupTarget.value.id,
      amount_fen: amountFen,
      note: topupNote.value,
    })
    appStore.showSuccess('加余额成功')
    showTopup.value = false
    await loadData()
  } catch (error: any) {
    appStore.showError(error?.message || '加余额失败')
  } finally {
    topupSubmitting.value = false
  }
}

function yuanAmount(value: number): string {
  const num = Number(value || 0)
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
  search.value = ''
  loadData()
})
</script>
