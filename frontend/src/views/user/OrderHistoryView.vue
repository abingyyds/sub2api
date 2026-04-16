<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-5xl space-y-6">
        <!-- Header -->
        <SlideIn direction="up" :delay="100">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('orderHistory.title') }}</h1>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('orderHistory.subtitle') }}</p>
            </div>
            <div class="flex items-center gap-3">
              <button
                v-if="!invoiceMode"
                class="rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
                @click="invoiceMode = true"
              >
                {{ t('orderHistory.invoice') }}
              </button>
              <button
                v-if="invoiceMode"
                class="rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
                @click="invoiceMode = false; selectedOrders.clear()"
              >
                {{ t('common.cancel') }}
              </button>
              <button
                v-if="invoiceMode"
                class="rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-primary-700 disabled:opacity-50"
                :disabled="selectedOrders.size === 0"
                @click="showInvoiceModal = true"
              >
                {{ t('orderHistory.submitInvoice') }} ({{ selectedOrders.size }})
              </button>
            </div>
          </div>
        </SlideIn>

        <!-- Loading -->
        <div v-if="loading" class="flex justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
        </div>

        <!-- Error -->
        <div v-else-if="error" class="rounded-2xl border border-red-200 bg-red-50 p-6 text-center dark:border-red-800/30 dark:bg-red-900/10">
          <p class="text-sm text-red-600 dark:text-red-400">{{ t('orderHistory.loadError') }}</p>
        </div>

        <SlideIn v-if="!loading && !error && showSuccessBanner" direction="up" :delay="150">
          <div class="rounded-2xl border border-green-200 bg-green-50 px-5 py-4 text-sm text-green-700 dark:border-green-800/30 dark:bg-green-900/10 dark:text-green-300">
            <div class="flex items-center justify-between gap-4">
              <span>{{ t('orderHistory.paymentSuccessArrived') }}</span>
              <button
                class="rounded-lg border border-green-300 px-3 py-1.5 text-xs font-medium text-green-700 transition hover:bg-green-100 dark:border-green-700 dark:text-green-300 dark:hover:bg-green-900/30"
                @click="fetchOrders(currentPage)"
              >
                {{ t('orderHistory.refreshOrders') }}
              </button>
            </div>
          </div>
        </SlideIn>

        <!-- Empty state -->
        <SlideIn v-if="!loading && !error && orders.length === 0" direction="up" :delay="200">
          <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft">
            <div class="flex flex-col items-center justify-center py-16 px-6 text-center">
              <div class="mb-6 inline-flex h-20 w-20 items-center justify-center rounded-2xl bg-primary-100 dark:bg-primary-900/30">
                <Icon name="clock" size="xl" class="text-primary-500 dark:text-primary-400" />
              </div>
              <h2 class="text-xl font-bold text-gray-900 dark:text-white">{{ t('orderHistory.noOrders') }}</h2>
              <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ t('orderHistory.noOrdersDesc') }}</p>
            </div>
          </div>
        </SlideIn>

        <!-- Order list -->
        <SlideIn v-else-if="!loading && !error" direction="up" :delay="200">
          <div class="rounded-2xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft overflow-hidden">
            <div class="overflow-x-auto">
              <table class="w-full text-left text-sm">
                <thead>
                  <tr class="border-b border-gray-100 bg-gray-50 dark:border-dark-700 dark:bg-dark-800">
                    <th v-if="invoiceMode" class="px-4 py-3"></th>
                    <th class="px-4 py-3 font-medium text-gray-500 dark:text-dark-400">{{ t('orderHistory.orderNo') }}</th>
                    <th class="px-4 py-3 font-medium text-gray-500 dark:text-dark-400">{{ t('orderHistory.type') }}</th>
                    <th class="px-4 py-3 font-medium text-gray-500 dark:text-dark-400">{{ t('orderHistory.amount') }}</th>
                    <th class="px-4 py-3 font-medium text-gray-500 dark:text-dark-400">{{ t('orderHistory.status') }}</th>
                    <th class="px-4 py-3 font-medium text-gray-500 dark:text-dark-400">{{ t('orderHistory.invoiceStatus') }}</th>
                    <th class="px-4 py-3 font-medium text-gray-500 dark:text-dark-400">{{ t('orderHistory.payMethod') }}</th>
                    <th class="px-4 py-3 font-medium text-gray-500 dark:text-dark-400">{{ t('orderHistory.createdAt') }}</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                  <tr
                    v-for="order in orders"
                    :key="order.order_no"
                    class="transition hover:bg-gray-50 dark:hover:bg-dark-800/50"
                    :class="isHighlightedOrder(order.order_no) ? 'bg-green-50/80 dark:bg-green-900/10' : ''"
                  >
                    <!-- Checkbox for invoice mode -->
                    <td v-if="invoiceMode" class="px-4 py-3">
                      <input
                        type="checkbox"
                        class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-500 dark:bg-dark-700"
                        :checked="selectedOrders.has(order.order_no)"
                        :disabled="!canRequestInvoice(order)"
                        @change="toggleOrder(order)"
                      />
                    </td>
                    <!-- Order No -->
                    <td class="px-4 py-3">
                      <div class="flex items-center gap-1.5">
                        <span class="font-mono text-xs text-gray-700 dark:text-dark-300">{{ order.order_no.slice(0, 16) }}...</span>
                        <span
                          v-if="order.invoice_requested_at"
                          class="inline-flex items-center rounded-full px-2 py-0.5 text-[10px] font-medium"
                          :class="invoiceBadgeClass(order)"
                        >
                          {{ invoiceStatusText(order) }}
                        </span>
                        <button
                          class="text-gray-400 hover:text-gray-600 dark:hover:text-dark-200"
                          :title="t('orderHistory.copyOrderNo')"
                          @click="copyToClipboard(order.order_no)"
                        >
                          <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                          </svg>
                        </button>
                        <span
                          v-if="isHighlightedOrder(order.order_no)"
                          class="inline-flex items-center rounded-full bg-green-100 px-2 py-0.5 text-[10px] font-medium text-green-700 dark:bg-green-900/30 dark:text-green-300"
                        >
                          {{ t('orderHistory.justPaid') }}
                        </span>
                      </div>
                    </td>
                    <!-- Type -->
                    <td class="px-4 py-3">
                      <span
                        class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                        :class="getOrderType(order) === 'recharge'
                          ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
                          : 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400'"
                      >
                        {{ getOrderType(order) === 'recharge' ? t('orderHistory.typeRecharge') : t('orderHistory.typeSubscription') }}
                      </span>
                    </td>
                    <!-- Amount -->
                    <td class="px-4 py-3 font-medium text-gray-900 dark:text-white">
                      ¥{{ (order.amount_fen / 100).toFixed(2) }}
                      <span v-if="order.discount_amount > 0" class="ml-1 text-xs text-green-600 dark:text-green-400">
                        -¥{{ (order.discount_amount / 100).toFixed(2) }}
                      </span>
                    </td>
                    <!-- Status -->
                    <td class="px-4 py-3">
                      <span
                        class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                        :class="statusClass(order.status)"
                      >
                        {{ statusText(order.status) }}
                      </span>
                    </td>
                    <td class="px-4 py-3">
                      <span
                        class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                        :class="invoiceBadgeClass(order)"
                      >
                        {{ invoiceStatusText(order) }}
                      </span>
                    </td>
                    <!-- Pay method -->
                    <td class="px-4 py-3 text-gray-500 dark:text-dark-400 text-xs">
                      {{ payMethodText(order.pay_method) }}
                    </td>
                    <!-- Created at -->
                    <td class="px-4 py-3 text-xs text-gray-500 dark:text-dark-400 whitespace-nowrap">
                      {{ formatDate(order.created_at) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-gray-100 px-4 py-3 dark:border-dark-700">
              <span class="text-sm text-gray-500 dark:text-dark-400">
                {{ t('common.page') }} {{ currentPage }} / {{ totalPages }}
              </span>
              <div class="flex gap-2">
                <button
                  class="rounded-lg border border-gray-200 px-3 py-1.5 text-sm text-gray-600 transition hover:bg-gray-50 disabled:opacity-40 dark:border-dark-600 dark:text-dark-300 dark:hover:bg-dark-700"
                  :disabled="currentPage <= 1"
                  @click="goPage(currentPage - 1)"
                >
                  {{ t('common.previous') }}
                </button>
                <button
                  class="rounded-lg border border-gray-200 px-3 py-1.5 text-sm text-gray-600 transition hover:bg-gray-50 disabled:opacity-40 dark:border-dark-600 dark:text-dark-300 dark:hover:bg-dark-700"
                  :disabled="currentPage >= totalPages"
                  @click="goPage(currentPage + 1)"
                >
                  {{ t('common.next') }}
                </button>
              </div>
            </div>
          </div>
        </SlideIn>

        <!-- Invoice mode summary bar -->
        <div
          v-if="invoiceMode && selectedOrders.size > 0"
          class="fixed bottom-6 left-1/2 -translate-x-1/2 flex items-center gap-4 rounded-2xl border border-gray-200 bg-white px-6 py-3 shadow-lg dark:border-dark-600 dark:bg-dark-800 z-50"
        >
          <span class="text-sm text-gray-600 dark:text-dark-300">
            {{ t('orderHistory.selectedCount', { count: selectedOrders.size }) }}
          </span>
          <span class="text-sm font-bold text-gray-900 dark:text-white">
            ¥{{ selectedTotalAmount.toFixed(2) }}
          </span>
          <button
            class="rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-primary-700"
            @click="showInvoiceModal = true"
          >
            {{ t('orderHistory.submitInvoice') }}
          </button>
        </div>
      </div>
    </FadeIn>

    <!-- Invoice Modal -->
    <Teleport to="body">
      <transition name="fade">
        <div
          v-if="showInvoiceModal"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
          @click.self="showInvoiceModal = false"
        >
          <div class="mx-4 w-full max-w-md rounded-2xl bg-white p-6 shadow-xl dark:bg-dark-900 space-y-5">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('orderHistory.invoiceTitle') }}</h3>

            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1">{{ t('orderHistory.companyName') }}</label>
                <input
                  v-model="invoiceForm.companyName"
                  type="text"
                  class="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm dark:border-dark-600 dark:bg-dark-800 dark:text-white focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20"
                  :placeholder="t('orderHistory.companyNamePlaceholder')"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1">{{ t('orderHistory.taxId') }}</label>
                <input
                  v-model="invoiceForm.taxId"
                  type="text"
                  class="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm dark:border-dark-600 dark:bg-dark-800 dark:text-white focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20"
                  :placeholder="t('orderHistory.taxIdPlaceholder')"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1">{{ t('orderHistory.email') }}</label>
                <input
                  v-model="invoiceForm.email"
                  type="email"
                  class="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm dark:border-dark-600 dark:bg-dark-800 dark:text-white focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20"
                  :placeholder="t('orderHistory.emailPlaceholder')"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1">{{ t('orderHistory.remark') }}</label>
                <textarea
                  v-model="invoiceForm.remark"
                  rows="2"
                  class="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm dark:border-dark-600 dark:bg-dark-800 dark:text-white focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 resize-none"
                  :placeholder="t('orderHistory.remarkPlaceholder')"
                ></textarea>
              </div>

              <!-- Selected summary -->
              <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800">
                <div class="flex items-center justify-between text-sm">
                  <span class="text-gray-500 dark:text-dark-400">{{ t('orderHistory.selectedAmount') }}</span>
                  <span class="font-bold text-gray-900 dark:text-white">¥{{ selectedTotalAmount.toFixed(2) }}</span>
                </div>
                <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">
                  {{ t('orderHistory.selectedCount', { count: selectedOrders.size }) }}
                </p>
              </div>
            </div>

            <div class="flex gap-3 pt-2">
              <button
                class="flex-1 rounded-lg border border-gray-200 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 dark:border-dark-600 dark:text-dark-300 dark:hover:bg-dark-700"
                @click="showInvoiceModal = false"
              >
                {{ t('common.cancel') }}
              </button>
              <button
                class="flex-1 rounded-lg bg-primary-600 py-2.5 text-sm font-medium text-white transition hover:bg-primary-700 disabled:opacity-50"
                :disabled="!invoiceForm.companyName || !invoiceForm.taxId || !invoiceForm.email || submittingInvoice"
                @click="submitInvoice"
              >
                {{ submittingInvoice ? t('common.loading') : t('orderHistory.submitInvoice') }}
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { FadeIn, SlideIn } from '@/components/animations'
import { useAppStore } from '@/stores/app'
import { paymentAPI, type PaymentOrder } from '@/api/payment'

const { t } = useI18n()
const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

const loading = ref(true)
const error = ref(false)
const orders = ref<PaymentOrder[]>([])
const currentPage = ref(1)
const totalPages = ref(1)
const pageSize = 20
const highlightedOrderNo = ref('')
const showSuccessBanner = ref(false)

// Invoice mode
const invoiceMode = ref(false)
const selectedOrders = reactive(new Set<string>())
const showInvoiceModal = ref(false)
const submittingInvoice = ref(false)
const invoiceForm = reactive({
  companyName: '',
  taxId: '',
  email: '',
  remark: '',
})

const selectedTotalAmount = computed(() => {
  let total = 0
  for (const orderNo of selectedOrders) {
    const order = orders.value.find(o => o.order_no === orderNo)
    if (order) {
      total += order.amount_fen / 100
    }
  }
  return total
})

function toggleOrder(order: PaymentOrder) {
  if (!canRequestInvoice(order)) return
  if (selectedOrders.has(order.order_no)) {
    selectedOrders.delete(order.order_no)
  } else {
    selectedOrders.add(order.order_no)
  }
}

function isLegacySubscriptionFallback(order: PaymentOrder): boolean {
  return order.order_no.startsWith('legacy-sub-') || order.pay_method.startsWith('legacy_subscription:')
}

function canRequestInvoice(order: PaymentOrder) {
  return !isLegacySubscriptionFallback(order) && order.status === 'paid' && !order.invoice_requested_at
}

function getOrderType(order: PaymentOrder): 'recharge' | 'subscription' {
  // If plan_key contains "recharge" or is empty, treat as recharge
  if (!order.plan_key || order.plan_key.includes('recharge') || order.plan_key.startsWith('custom')) {
    return 'recharge'
  }
  return 'subscription'
}

function statusClass(status: string): string {
  switch (status) {
    case 'paid': return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
    case 'pending': return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
    case 'closed': return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-400'
    case 'refunded': return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
    default: return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-400'
  }
}

function statusText(status: string): string {
  switch (status) {
    case 'paid': return t('orderHistory.statusPaid')
    case 'pending': return t('orderHistory.statusPending')
    case 'closed': return t('orderHistory.statusClosed')
    case 'refunded': return t('orderHistory.statusRefunded')
    default: return status
  }
}

function payMethodText(method: string): string {
  if (method.startsWith('legacy_subscription:')) {
    return t('orderHistory.legacySubscriptionSource', { group: method.replace('legacy_subscription:', '') })
  }
  switch (method) {
    case 'wechat': return t('orderHistory.payMethodWechat')
    case 'alipay': return t('orderHistory.payMethodAlipay')
    case 'epay_alipay': return 'Epay-支付宝'
    case 'epay_wxpay': return 'Epay-微信'
    default: return method
  }
}

function invoiceStatusText(order: PaymentOrder): string {
  if (isLegacySubscriptionFallback(order)) return t('orderHistory.invoiceUnsupported')
  if (order.invoice_processed_at) return t('orderHistory.invoiceProcessed')
  if (order.invoice_requested_at) return t('orderHistory.invoiceRequested')
  return t('orderHistory.invoiceNotRequested')
}

function invoiceBadgeClass(order: PaymentOrder): string {
  if (isLegacySubscriptionFallback(order)) return 'bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-dark-400'
  if (order.invoice_processed_at) return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400'
  if (order.invoice_requested_at) return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
  return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-400'
}

function resetInvoiceForm() {
  invoiceForm.companyName = ''
  invoiceForm.taxId = ''
  invoiceForm.email = ''
  invoiceForm.remark = ''
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function isHighlightedOrder(orderNo: string): boolean {
  return highlightedOrderNo.value !== '' && highlightedOrderNo.value === orderNo
}

async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text)
  } catch {
    // fallback
  }
}

async function fetchOrders(page = 1) {
  loading.value = true
  error.value = false
  try {
    const result = await paymentAPI.listOrders({ page, page_size: pageSize })
    orders.value = result.items || []
    const selectable = new Set(orders.value.filter(canRequestInvoice).map(order => order.order_no))
    for (const orderNo of Array.from(selectedOrders)) {
      if (!selectable.has(orderNo)) {
        selectedOrders.delete(orderNo)
      }
    }
    currentPage.value = result.page
    totalPages.value = result.pages
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

function goPage(page: number) {
  if (page < 1 || page > totalPages.value) return
  fetchOrders(page)
}

async function submitInvoice() {
  if (selectedOrders.size === 0) return
  submittingInvoice.value = true
  try {
    await paymentAPI.submitInvoiceRequest({
      order_nos: Array.from(selectedOrders),
      company_name: invoiceForm.companyName.trim(),
      tax_id: invoiceForm.taxId.trim(),
      email: invoiceForm.email.trim(),
      remark: invoiceForm.remark.trim(),
    })
    showInvoiceModal.value = false
    invoiceMode.value = false
    selectedOrders.clear()
    resetInvoiceForm()
    await fetchOrders(currentPage.value)
    appStore.showSuccess(t('orderHistory.invoiceSubmitted'))
  } catch (err: any) {
    appStore.showError(err?.message || t('orderHistory.invoiceError'))
  } finally {
    submittingInvoice.value = false
  }
}

onMounted(() => {
  const highlight = typeof route.query.highlight === 'string' ? route.query.highlight : ''
  highlightedOrderNo.value = highlight
  showSuccessBanner.value = route.query.success === '1'

  fetchOrders().finally(() => {
    if (showSuccessBanner.value) {
      router.replace({
        query: {
          ...route.query,
          success: undefined,
        }
      })
    }
  })
})
</script>
