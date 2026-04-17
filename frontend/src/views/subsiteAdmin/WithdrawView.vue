<template>
  <AppLayout>
    <div class="mx-auto max-w-5xl space-y-6">
      <div class="flex items-center gap-3">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站利润提现</h1>
        <span v-if="site?.mode" class="inline-flex rounded-full px-3 py-1 text-xs font-medium"
          :class="site.mode === 'rate'
            ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300'
            : 'bg-gray-200 text-gray-700 dark:bg-dark-700 dark:text-dark-300'">
          {{ site.mode === 'rate' ? '倍率分成模式' : '资金池模式' }}
        </span>
      </div>

      <div v-if="site?.mode === 'pool'" class="rounded-2xl border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800 dark:border-amber-800/50 dark:bg-amber-900/20 dark:text-amber-200">
        当前分站为资金池模式，盈利通过线下售卖用户余额实现，不支持利润提现。
      </div>

      <template v-else-if="site">
        <div class="grid gap-4 md:grid-cols-3">
          <div class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-xs font-medium text-gray-500 dark:text-dark-400">可提现利润</p>
            <p class="mt-2 text-2xl font-bold text-emerald-600 dark:text-emerald-400">
              ￥{{ fenToYuan(site.balance_fen || 0) }}
            </p>
          </div>
          <div class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-xs font-medium text-gray-500 dark:text-dark-400">累计利润</p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              ￥{{ fenToYuan(site.total_topup_fen || 0) }}
            </p>
          </div>
          <div class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <p class="text-xs font-medium text-gray-500 dark:text-dark-400">累计已提现</p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              ￥{{ fenToYuan(site.total_withdrawn_fen || 0) }}
            </p>
          </div>
        </div>

        <section class="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">申请提现</h2>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">申请提交后需管理员审核，审核通过并打款后到账。</p>
          <div class="mt-4 grid gap-4 md:grid-cols-2">
            <div>
              <label class="input-label">提现金额（元）</label>
              <input v-model.number="form.amountYuan" type="number" min="1" step="1" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">支付宝姓名</label>
              <input v-model="form.alipay_name" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">支付宝账号</label>
              <input v-model="form.alipay_account" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">收款码图片 URL（可选）</label>
              <input v-model="form.alipay_qr_image" class="input mt-1 w-full" />
            </div>
          </div>
          <div class="mt-4 flex justify-end">
            <button class="btn btn-primary" :disabled="submitting || !canSubmit" @click="handleApply">
              {{ submitting ? '提交中...' : '提交申请' }}
            </button>
          </div>
        </section>

        <section class="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-800">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">提现历史</h2>
          </div>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 text-sm dark:divide-dark-700">
              <thead class="bg-gray-50 dark:bg-dark-800/60">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-semibold uppercase">ID</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold uppercase">金额</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold uppercase">收款信息</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold uppercase">状态</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold uppercase">时间</th>
                  <th class="px-4 py-3 text-right text-xs font-semibold uppercase">操作</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
                <tr v-if="listLoading">
                  <td colspan="6" class="px-4 py-10 text-center text-gray-400">加载中...</td>
                </tr>
                <tr v-else-if="withdraws.length === 0">
                  <td colspan="6" class="px-4 py-10 text-center text-gray-400">暂无提现记录</td>
                </tr>
                <tr v-for="w in withdraws" :key="w.id">
                  <td class="px-4 py-3 text-xs">{{ w.id }}</td>
                  <td class="px-4 py-3 font-medium">￥{{ fenToYuan(w.amount) }}</td>
                  <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">
                    <div>{{ w.alipay_name }}</div>
                    <div>{{ w.alipay_account }}</div>
                  </td>
                  <td class="px-4 py-3">
                    <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium" :class="getStatusClass(w.status)">
                      {{ getStatusText(w.status) }}
                    </span>
                    <div v-if="w.review_note" class="mt-1 text-xs text-gray-500">{{ w.review_note }}</div>
                  </td>
                  <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400">{{ formatDate(w.created_at) }}</td>
                  <td class="px-4 py-3 text-right">
                    <button v-if="w.status === 'pending'" class="btn btn-secondary btn-sm text-red-600 dark:text-red-300" @click="handleCancel(w)">撤销</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { subSiteAdminAPI } from '@/api/subsiteAdmin'
import type { OwnedSubSite } from '@/api/subsite'
import { withdrawAPI, type WithdrawRequest } from '@/api/withdraw'
import { useAppStore } from '@/stores'

const route = useRoute()
const appStore = useAppStore()

const siteId = computed(() => Number(route.params.siteId))

const site = ref<OwnedSubSite | null>(null)
const withdraws = ref<WithdrawRequest[]>([])
const listLoading = ref(false)
const submitting = ref(false)

const form = reactive({
  amountYuan: 0,
  alipay_name: '',
  alipay_account: '',
  alipay_qr_image: ''
})

const canSubmit = computed(() =>
  form.amountYuan > 0 &&
  form.alipay_name.trim() !== '' &&
  form.alipay_account.trim() !== '' &&
  site.value &&
  Math.round(form.amountYuan * 100) <= (site.value.balance_fen || 0)
)

function fenToYuan(fen: number): string {
  return (fen / 100).toFixed(2)
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', { hour12: false })
}

function getStatusClass(status: string): string {
  switch (status) {
    case 'pending': return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
    case 'approved': return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
    case 'rejected': return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
    case 'paid': return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
    case 'cancelled': return 'bg-gray-200 text-gray-700 dark:bg-dark-700 dark:text-dark-300'
    default: return 'bg-gray-200 text-gray-700 dark:bg-dark-700 dark:text-dark-300'
  }
}

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    pending: '待审核',
    approved: '已批准',
    rejected: '已拒绝',
    paid: '已打款',
    cancelled: '已取消'
  }
  return map[status] || status
}

async function loadSite() {
  if (!siteId.value) return
  try {
    site.value = await subSiteAdminAPI.getSite(siteId.value)
  } catch (e: any) {
    appStore.showError(e?.message || '加载分站失败')
  }
}

async function loadWithdraws() {
  listLoading.value = true
  try {
    const res = await withdrawAPI.listWithdraws(1, 50, {
      source_type: 'sub_site_profit'
    })
    withdraws.value = (res.items || []).filter((w) => w.source_sub_site_id === siteId.value)
  } catch (e: any) {
    appStore.showError(e?.message || '加载提现历史失败')
  } finally {
    listLoading.value = false
  }
}

async function handleApply() {
  if (!canSubmit.value || !siteId.value) return
  submitting.value = true
  try {
    await withdrawAPI.applyWithdraw({
      amount: Math.round(form.amountYuan * 100),
      alipay_name: form.alipay_name.trim(),
      alipay_account: form.alipay_account.trim(),
      alipay_qr_image: form.alipay_qr_image.trim() || undefined,
      source_type: 'sub_site_profit',
      source_sub_site_id: siteId.value
    })
    appStore.showSuccess('提现申请已提交')
    form.amountYuan = 0
    form.alipay_qr_image = ''
    await Promise.all([loadSite(), loadWithdraws()])
  } catch (e: any) {
    appStore.showError(e?.response?.data?.error || e?.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

async function handleCancel(w: WithdrawRequest) {
  if (!confirm('确定撤销这笔提现申请？撤销后金额将退回利润余额。')) return
  try {
    await withdrawAPI.cancelWithdraw(w.id)
    appStore.showSuccess('已撤销')
    await Promise.all([loadSite(), loadWithdraws()])
  } catch (e: any) {
    appStore.showError(e?.response?.data?.error || '撤销失败')
  }
}

onMounted(async () => {
  await loadSite()
  if (site.value?.mode === 'rate') {
    await loadWithdraws()
  }
})

watch(siteId, async () => {
  await loadSite()
  if (site.value?.mode === 'rate') {
    await loadWithdraws()
  }
})
</script>
