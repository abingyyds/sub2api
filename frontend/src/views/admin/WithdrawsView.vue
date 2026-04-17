<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-6">
      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">提现管理</h1>
          <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
            审核和处理代理佣金、分站利润的提现申请
          </p>
        </div>
      </section>

      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="grid gap-3 md:grid-cols-[180px_180px_auto]">
          <select v-model="filters.source_type" class="input" @change="loadData(1)">
            <option value="">全部来源</option>
            <option value="agent_commission">代理佣金</option>
            <option value="sub_site_profit">分站利润</option>
          </select>
          <select v-model="filters.status" class="input" @change="loadData(1)">
            <option value="">全部状态</option>
            <option value="pending">待审核</option>
            <option value="approved">已批准</option>
            <option value="rejected">已拒绝</option>
            <option value="paid">已打款</option>
            <option value="cancelled">已取消</option>
          </select>
          <button class="btn btn-secondary" @click="loadData(1)">刷新</button>
        </div>
      </section>

      <section class="overflow-hidden rounded-3xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-dark-700">
            <thead class="bg-gray-50 dark:bg-dark-800/60">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">ID / 用户</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">来源</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">金额</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">支付宝信息</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">状态</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">时间</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-if="loading">
                <td colspan="7" class="px-4 py-10 text-center text-sm text-gray-500 dark:text-dark-400">加载中...</td>
              </tr>
              <tr v-else-if="items.length === 0">
                <td colspan="7" class="px-4 py-10 text-center text-sm text-gray-500 dark:text-dark-400">暂无提现申请</td>
              </tr>
              <tr v-for="item in items" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50">
                <td class="px-4 py-4 align-top">
                  <div class="font-medium text-gray-900 dark:text-white">ID: {{ item.id }}</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">用户: {{ item.user_id }}</div>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium"
                    :class="item.source_type === 'agent_commission'
                      ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
                      : 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300'">
                    {{ item.source_type === 'agent_commission' ? '代理佣金' : '分站利润' }}
                  </span>
                  <div v-if="item.source_sub_site_id" class="mt-1 text-xs text-gray-500 dark:text-dark-400">
                    分站ID: {{ item.source_sub_site_id }}
                  </div>
                </td>
                <td class="px-4 py-4 align-top">
                  <div class="font-medium text-gray-900 dark:text-white">￥{{ fenToYuan(item.amount) }}</div>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <div>{{ item.alipay_name }}</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ item.alipay_account }}</div>
                  <div v-if="item.alipay_qr_image" class="mt-1">
                    <a :href="item.alipay_qr_image" target="_blank" class="text-xs text-primary-600 hover:underline dark:text-primary-400">查看收款码</a>
                  </div>
                </td>
                <td class="px-4 py-4 align-top">
                  <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium"
                    :class="getStatusClass(item.status)">
                    {{ getStatusText(item.status) }}
                  </span>
                  <div v-if="item.review_note" class="mt-1 text-xs text-gray-500 dark:text-dark-400">
                    备注: {{ item.review_note }}
                  </div>
                </td>
                <td class="px-4 py-4 align-top text-xs text-gray-500 dark:text-dark-400">
                  <div>申请: {{ formatDate(item.created_at) }}</div>
                  <div v-if="item.reviewed_at" class="mt-1">审核: {{ formatDate(item.reviewed_at) }}</div>
                  <div v-if="item.paid_at" class="mt-1">打款: {{ formatDate(item.paid_at) }}</div>
                </td>
                <td class="px-4 py-4 align-top text-right">
                  <div class="flex justify-end gap-2">
                    <button v-if="item.status === 'pending'" class="btn btn-secondary btn-sm" @click="openReview(item, true)">批准</button>
                    <button v-if="item.status === 'pending'" class="btn btn-secondary btn-sm text-red-600 dark:text-red-300" @click="openReview(item, false)">拒绝</button>
                    <button v-if="item.status === 'approved'" class="btn btn-primary btn-sm" @click="handlePay(item)">确认打款</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-if="pages > 1" class="flex items-center justify-between border-t border-gray-100 px-4 py-3 text-sm dark:border-dark-800">
          <span class="text-gray-500 dark:text-dark-400">第 {{ page }} / {{ pages }} 页</span>
          <div class="flex gap-2">
            <button class="btn btn-secondary btn-sm" :disabled="page <= 1" @click="loadData(page - 1)">上一页</button>
            <button class="btn btn-secondary btn-sm" :disabled="page >= pages" @click="loadData(page + 1)">下一页</button>
          </div>
        </div>
      </section>
    </div>

    <BaseDialog :show="showReviewDialog" :title="reviewApprove ? '批准提现' : '拒绝提现'" @close="closeReview">
      <div class="space-y-4">
        <div v-if="reviewItem">
          <p class="text-sm text-gray-600 dark:text-dark-300">
            用户 {{ reviewItem.user_id }} 申请提现 <span class="font-medium">￥{{ fenToYuan(reviewItem.amount) }}</span>
          </p>
          <p class="mt-2 text-sm text-gray-600 dark:text-dark-300">
            支付宝: {{ reviewItem.alipay_name }} ({{ reviewItem.alipay_account }})
          </p>
        </div>
        <div>
          <label class="input-label">审核备注</label>
          <textarea v-model="reviewNote" class="input mt-1 w-full" rows="3" placeholder="可选"></textarea>
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="closeReview">取消</button>
        <button class="btn" :class="reviewApprove ? 'btn-primary' : 'btn-secondary text-red-600'" :disabled="reviewing" @click="handleReview">
          {{ reviewing ? '处理中...' : (reviewApprove ? '确认批准' : '确认拒绝') }}
        </button>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminWithdrawsAPI } from '@/api/admin/withdraws'
import type { WithdrawRequest } from '@/api/withdraw'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'

const loading = ref(false)
const items = ref<WithdrawRequest[]>([])
const page = ref(1)
const pages = ref(1)
const pageSize = 20

const filters = ref({
  source_type: '',
  status: ''
})

const showReviewDialog = ref(false)
const reviewItem = ref<WithdrawRequest | null>(null)
const reviewApprove = ref(true)
const reviewNote = ref('')
const reviewing = ref(false)

function fenToYuan(fen: number): string {
  return (fen / 100).toFixed(2)
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getStatusClass(status: string): string {
  switch (status) {
    case 'pending':
      return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
    case 'approved':
      return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
    case 'rejected':
      return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
    case 'paid':
      return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
    case 'cancelled':
      return 'bg-gray-200 text-gray-700 dark:bg-dark-700 dark:text-dark-300'
    default:
      return 'bg-gray-200 text-gray-700 dark:bg-dark-700 dark:text-dark-300'
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

async function loadData(p: number = 1) {
  loading.value = true
  try {
    const res = await adminWithdrawsAPI.listWithdraws(p, pageSize, {
      source_type: filters.value.source_type || undefined,
      status: filters.value.status || undefined
    })
    items.value = res.items || []
    page.value = res.page
    pages.value = res.pages
  } catch (err: any) {
    console.error('Failed to load withdraws:', err)
    alert(err.response?.data?.error || '加载失败')
  } finally {
    loading.value = false
  }
}

function openReview(item: WithdrawRequest, approve: boolean) {
  reviewItem.value = item
  reviewApprove.value = approve
  reviewNote.value = ''
  showReviewDialog.value = true
}

function closeReview() {
  showReviewDialog.value = false
  reviewItem.value = null
  reviewNote.value = ''
}

async function handleReview() {
  if (!reviewItem.value) return
  reviewing.value = true
  try {
    await adminWithdrawsAPI.reviewWithdraw(reviewItem.value.id, {
      approve: reviewApprove.value,
      review_note: reviewNote.value || undefined
    })
    closeReview()
    loadData(page.value)
  } catch (err: any) {
    console.error('Failed to review:', err)
    alert(err.response?.data?.error || '操作失败')
  } finally {
    reviewing.value = false
  }
}

async function handlePay(item: WithdrawRequest) {
  if (!confirm(`确认已向 ${item.alipay_name} (${item.alipay_account}) 打款 ￥${fenToYuan(item.amount)}？`)) return
  try {
    await adminWithdrawsAPI.payWithdraw(item.id)
    loadData(page.value)
  } catch (err: any) {
    console.error('Failed to mark as paid:', err)
    alert(err.response?.data?.error || '操作失败')
  }
}

onMounted(() => {
  loadData(1)
})
</script>
