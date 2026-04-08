<template>
  <AppLayout>
    <FadeIn>
      <div class="mx-auto max-w-2xl space-y-8">
        <!-- Title -->
        <SlideIn direction="up" :delay="100">
          <div class="text-center">
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('recharge.title') }}</h1>
            <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('recharge.subtitle') }}</p>
          </div>
        </SlideIn>

        <SlideIn direction="up" :delay="200">
          <GlowCard glow-color="rgb(59, 130, 246)">
            <div class="rounded-2xl border-2 border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900 shadow-soft p-6 space-y-6">
              <!-- Quick Select Amounts -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-3">
                  {{ t('recharge.quickSelect') }}
                </label>
                <div class="grid grid-cols-3 gap-3">
                  <button
                    v-for="preset in presets"
                    :key="preset"
                    class="rounded-lg border-2 px-4 py-3 text-center font-medium transition-all hover:border-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20"
                    :class="selectedAmount === preset
                      ? 'border-primary-500 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-400'
                      : 'border-gray-200 text-gray-700 dark:border-dark-600 dark:text-dark-300'"
                    @click="selectPreset(preset)"
                  >
                    ¥{{ preset }}
                  </button>
                </div>
              </div>

              <!-- Custom Input -->
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-2">
                  {{ t('recharge.customAmount') }}
                </label>
                <div class="relative">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-dark-400 font-medium">¥</span>
                  <input
                    v-model="customInput"
                    type="number"
                    min="1"
                    step="1"
                    :placeholder="t('recharge.inputPlaceholder')"
                    class="w-full rounded-lg border border-gray-300 bg-white py-3 pl-8 pr-4 text-lg font-medium text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-800 dark:text-white dark:placeholder-dark-500 transition-all duration-300"
                    @input="onCustomInput"
                  />
                </div>
              </div>

              <!-- Amount Summary -->
              <div v-if="finalAmount > 0" class="rounded-lg bg-gray-50 dark:bg-dark-800 p-4">
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-600 dark:text-dark-400">{{ t('recharge.rechargeAmount') }}</span>
                  <span class="text-lg font-bold text-primary-600">¥{{ finalAmount.toFixed(2) }}</span>
                </div>
                <div class="mt-1 flex items-center justify-between">
                  <span class="text-sm text-gray-600 dark:text-dark-400">{{ t('recharge.payAmount') }}</span>
                  <span class="text-2xl font-bold text-gray-900 dark:text-white">¥{{ finalAmount.toFixed(2) }}</span>
                </div>
              </div>

              <!-- Submit Button -->
              <MagneticButton>
                <button
                  class="w-full rounded-lg bg-primary-600 py-3 text-base font-medium text-white transition-all duration-300 hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  :disabled="finalAmount <= 0 || creatingOrder"
                  @click="handleRecharge"
                >
                  <span v-if="creatingOrder" class="flex items-center justify-center gap-2">
                    <div class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
                    {{ t('recharge.creating') }}
                  </span>
                  <span v-else>
                    {{ finalAmount > 0 ? t('recharge.payNow', { amount: finalAmount.toFixed(2) }) : t('recharge.enterAmount') }}
                  </span>
                </button>
              </MagneticButton>
            </div>
          </GlowCard>
        </SlideIn>
      </div>
    </FadeIn>

    <!-- Payment Method Selection Modal -->
    <Teleport to="body">
      <div
        v-if="showPayMethodModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="showPayMethodModal = false"
      >
        <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-2xl dark:bg-dark-700 mx-4">
          <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-4">选择支付方式</h3>

          <div class="space-y-3 mb-6">
            <div
              v-for="method in availablePayMethods"
              :key="method"
              class="flex items-center gap-3 rounded-lg border-2 p-4 cursor-pointer transition-all"
              :class="selectedPayMethod === method
                ? (payMethodMeta(method).color === 'green' ? 'border-green-500 bg-green-50 dark:bg-green-900/20 shadow-sm' : 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 shadow-sm')
                : (payMethodMeta(method).color === 'green' ? 'border-gray-200 dark:border-dark-600 hover:border-green-300' : 'border-gray-200 dark:border-dark-600 hover:border-blue-300')"
              @click="selectedPayMethod = method"
            >
              <input type="radio" :checked="selectedPayMethod === method" class="h-5 w-5" :class="payMethodMeta(method).color === 'green' ? 'text-green-600' : 'text-blue-600'" />
              <svg class="h-8 w-8" viewBox="0 0 24 24" fill="none" v-html="payMethodMeta(method).icon"></svg>
              <span class="font-bold text-gray-900 dark:text-white">{{ payMethodMeta(method).label }}</span>
            </div>
          </div>

          <div class="flex gap-3">
            <button
              class="btn flex-1 border border-gray-300 dark:border-dark-500"
              @click="showPayMethodModal = false"
            >
              取消
            </button>
            <button
              class="btn btn-primary flex-1"
              @click="confirmPayMethod"
            >
              确认支付
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Payment QR Code Modal -->
    <Teleport to="body">
      <div
        v-if="showPaymentModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="cancelPayment"
      >
        <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-2xl dark:bg-dark-700 mx-4">
          <div class="mb-6 text-center">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('recharge.payment.scanToPay') }}</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
              {{ t('recharge.payment.amount') }}: <span class="font-bold text-primary-600">¥{{ currentOrderAmount }}</span>
            </p>
          </div>

          <div class="flex justify-center mb-6">
            <div v-if="qrLoading" class="flex h-48 w-48 items-center justify-center">
              <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
            </div>
            <canvas v-else ref="qrCanvas" class="rounded-lg"></canvas>
          </div>

          <div class="mb-4 text-center">
            <div v-if="paymentStatus === 'pending'" class="flex items-center justify-center gap-2 text-sm text-gray-500 dark:text-dark-400">
              <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
              {{ t('recharge.payment.waiting') }}
            </div>
            <div v-else-if="paymentStatus === 'paid'" class="flex items-center justify-center gap-2 text-sm text-green-600">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ t('recharge.payment.success') }}
            </div>
            <div v-else-if="paymentStatus === 'closed'" class="text-sm text-red-500">
              {{ t('recharge.payment.expired') }}
            </div>
          </div>

          <div v-if="paymentStatus === 'pending' && countdown > 0" class="mb-4 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('recharge.payment.expiresIn', { minutes: Math.floor(countdown / 60), seconds: countdown % 60 }) }}
          </div>

          <div class="flex gap-3">
            <button
              v-if="paymentStatus !== 'paid'"
              class="btn flex-1 border border-gray-300 dark:border-dark-500"
              @click="cancelPayment"
            >
              {{ t('recharge.payment.cancel') }}
            </button>
            <button
              v-if="paymentStatus === 'paid'"
              class="btn btn-primary flex-1"
              @click="goToDashboard"
            >
              {{ t('recharge.payment.viewDashboard') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { FadeIn, SlideIn, GlowCard, MagneticButton } from '@/components/animations'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import { paymentAPI, type PayMethod } from '@/api/payment'
import { useAppStore, useAuthStore } from '@/stores'

const { t } = useI18n()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

const presets = [10, 50, 100, 200, 500, 1000]
const selectedAmount = ref<number | null>(null)
const customInput = ref<string>('')
const creatingOrder = ref(false)
const availablePayMethods = ref<PayMethod[]>([])
const selectedPayMethod = ref<PayMethod>('wechat')
const showPayMethodModal = ref(false)
const pendingPaymentAction = ref<(() => Promise<void>) | null>(null)

const showPaymentModal = ref(false)
const qrLoading = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)
const paymentStatus = ref<'pending' | 'paid' | 'closed'>('pending')
const currentOrderNo = ref('')
const currentOrderAmount = ref('')
const countdown = ref(0)

let pollTimer: ReturnType<typeof setInterval> | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null

const finalAmount = computed(() => {
  if (customInput.value !== '') {
    const val = parseFloat(customInput.value)
    return isNaN(val) || val <= 0 ? 0 : val
  }
  return selectedAmount.value ?? 0
})

function selectPreset(amount: number) {
  selectedAmount.value = amount
  customInput.value = ''
}

function onCustomInput() {
  selectedAmount.value = null
}

function clearTimers() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  if (countdownTimer) { clearInterval(countdownTimer); countdownTimer = null }
}

onMounted(async () => {
  try {
    const methods = await paymentAPI.getPayMethods()
    availablePayMethods.value = methods || []
    if (methods.length > 0) {
      selectedPayMethod.value = methods[0]
    }
  } catch {
    // fallback: empty means no methods available
  }
})

onUnmounted(() => clearTimers())

// === Payment method helpers ===
const WECHAT_ICON = '<rect width="24" height="24" rx="4" fill="#07C160"/><path d="M16.7 10.5c-.2 0-.4 0-.6.1.1-.4.1-.8.1-1.2 0-2.8-2.7-5-5.9-5C7 4.4 4.4 6.6 4.4 9.4c0 1.5.8 2.9 2.1 3.9l-.5 1.6 1.9-1c.6.2 1.2.3 1.8.3h.2c-.1-.3-.1-.7-.1-1 0-2.4 2.3-4.4 5.1-4.4l.2-.1c-.1-.1-.2-.2-.4-.2zm-3.3-2.2c.4 0 .7.3.7.7s-.3.7-.7.7-.7-.3-.7-.7.3-.7.7-.7zm-4.2 1.4c-.4 0-.7-.3-.7-.7s.3-.7.7-.7.7.3.7.7-.3.7-.7.7zm10.4 4.5c0-2.2-2.2-4-4.8-4s-4.8 1.8-4.8 4 2.2 4 4.8 4c.5 0 1-.1 1.5-.2l1.5.8-.4-1.3c1.2-.8 2.2-2 2.2-3.3zm-6.4-.6c-.3 0-.6-.3-.6-.6s.3-.6.6-.6.6.3.6.6-.3.6-.6.6zm3.2 0c-.3 0-.6-.3-.6-.6s.3-.6.6-.6.6.3.6.6-.3.6-.6.6z" fill="white"/>'
const ALIPAY_ICON = '<rect width="24" height="24" rx="4" fill="#1677FF"/><path d="M17.5 14.2c-1.2-.5-2.3-1-3.2-1.4.4-.8.7-1.7.9-2.6h-2.5v-1h3V8.4h-3V6.5h-1.4v1.9h-3v.8h3v1h-2.5v.8h4.6c-.2.7-.5 1.3-.8 1.9-1.3-.5-2.7-.8-3.8-.8-2 0-3.3.9-3.3 2.3 0 1.5 1.3 2.3 3.2 2.3 1.5 0 2.9-.6 4-1.5.9.5 1.9 1 3 1.5l.8-1.2zM9.2 16.3c-1.3 0-1.9-.5-1.9-1.2 0-.8.7-1.3 1.9-1.3.9 0 1.9.2 2.9.7-.9.8-1.9 1.8-2.9 1.8z" fill="white"/>'

const PAY_METHOD_META: Record<string, { label: string; icon: string; color: 'green' | 'blue' }> = {
  wechat: { label: '微信支付', icon: WECHAT_ICON, color: 'green' },
  alipay: { label: '支付宝', icon: ALIPAY_ICON, color: 'blue' },
  epay_alipay: { label: 'Epay-支付宝', icon: ALIPAY_ICON, color: 'blue' },
  epay_wxpay: { label: 'Epay-微信', icon: WECHAT_ICON, color: 'green' },
}

function payMethodMeta(method: string) {
  return PAY_METHOD_META[method] || { label: method, icon: '', color: 'blue' as const }
}

function showPayMethodOrDirect(action: () => Promise<void>) {
  pendingPaymentAction.value = action
  if (availablePayMethods.value.length === 0) {
    pendingPaymentAction.value = null
    alert(t('recharge.payment.noMethodsAvailable'))
    return
  }
  if (availablePayMethods.value.length <= 1) {
    confirmPayMethod()
  } else {
    showPayMethodModal.value = true
  }
}

async function confirmPayMethod() {
  showPayMethodModal.value = false
  if (pendingPaymentAction.value) {
    await pendingPaymentAction.value()
    pendingPaymentAction.value = null
  }
}

async function handleRecharge() {
  if (finalAmount.value <= 0 || creatingOrder.value) return

  showPayMethodOrDirect(async () => {
    creatingOrder.value = true
    try {
      const order = await paymentAPI.createRechargeOrder(finalAmount.value, undefined, selectedPayMethod.value)
      currentOrderNo.value = order.order_no
      currentOrderAmount.value = (order.amount_fen / 100).toFixed(order.amount_fen % 100 === 0 ? 0 : 2)
      paymentStatus.value = 'pending'
      showPaymentModal.value = true

      const expiresAt = new Date(order.expired_at).getTime()
      countdown.value = Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))

      if (order.code_url) {
        qrLoading.value = false
        await nextTick()
        if (qrCanvas.value) {
          await QRCode.toCanvas(qrCanvas.value, order.code_url, {
            width: 192,
            margin: 2,
            color: { dark: '#000000', light: '#ffffff' },
          })
        }
      }

      startPolling()

      countdownTimer = setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) {
          paymentStatus.value = 'closed'
          clearTimers()
        }
      }, 1000)
    } catch (err: any) {
      const msg = err?.message || err?.response?.data?.message || t('recharge.payment.createFailed')
      alert(msg)
    } finally {
      creatingOrder.value = false
    }
  })
}

function startPolling() {
  pollTimer = setInterval(async () => {
    try {
      const order = await paymentAPI.queryOrder(currentOrderNo.value)
      if (order.status === 'paid') {
        paymentStatus.value = 'paid'
        clearTimers()
        await syncPaymentSuccessState()
      } else if (order.status === 'closed') {
        paymentStatus.value = 'closed'
        clearTimers()
      }
    } catch {
      // ignore poll errors
    }
  }, 3000)
}

async function syncPaymentSuccessState() {
  try {
    await authStore.refreshUser()
  } catch (error) {
    console.error('Failed to refresh recharge success state:', error)
    appStore.showError('充值成功，但页面余额刷新失败，请手动刷新页面查看最新状态')
  }
}

function cancelPayment() {
  showPaymentModal.value = false
  clearTimers()
}

function goToDashboard() {
  showPaymentModal.value = false
  clearTimers()
  router.push('/dashboard')
}
</script>
