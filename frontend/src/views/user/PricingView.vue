<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-6">
      <!-- Title -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('pricing.title') }}</h1>
        <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('pricing.subtitle') }}</p>
      </div>

      <!-- Tabs -->
      <div class="flex justify-center">
        <div class="inline-flex rounded-lg border border-gray-200 bg-gray-100 p-1 dark:border-dark-600 dark:bg-dark-800">
          <button
            class="rounded-md px-5 py-2 text-sm font-medium transition-all"
            :class="activeTab === 'subscription'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'"
            @click="activeTab = 'subscription'"
          >
            {{ t('pricing.tabSubscription') }}
          </button>
          <button
            class="rounded-md px-5 py-2 text-sm font-medium transition-all"
            :class="activeTab === 'recharge'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-dark-400 dark:hover:text-dark-200'"
            @click="activeTab = 'recharge'"
          >
            {{ t('pricing.tabRecharge') }}
          </button>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>

      <!-- Promo Code Input (shared across tabs) -->
      <div v-if="!loading" class="mx-auto max-w-lg">
        <div class="flex items-center gap-3">
          <div class="relative flex-1">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
              <Icon name="gift" size="md" :class="promoValidation.valid ? 'text-green-500' : 'text-gray-400 dark:text-dark-500'" />
            </div>
            <input
              v-model="promoCode"
              type="text"
              class="input pl-10 pr-10"
              :class="{
                'border-green-500 focus:border-green-500 focus:ring-green-500': promoValidation.valid,
                'border-red-500 focus:border-red-500 focus:ring-red-500': promoValidation.invalid
              }"
              :placeholder="t('pricing.promoCodePlaceholder')"
              @input="handlePromoCodeInput"
            />
            <div v-if="promoValidating" class="absolute inset-y-0 right-0 flex items-center pr-3">
              <svg class="h-4 w-4 animate-spin text-gray-400" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </div>
            <div v-else-if="promoValidation.valid" class="absolute inset-y-0 right-0 flex items-center pr-3">
              <Icon name="checkCircle" size="md" class="text-green-500" />
            </div>
            <div v-else-if="promoValidation.invalid" class="absolute inset-y-0 right-0 flex items-center pr-3">
              <Icon name="exclamationCircle" size="md" class="text-red-500" />
            </div>
          </div>
        </div>
        <!-- Promo code validation result -->
        <transition name="fade">
          <div v-if="promoValidation.valid" class="mt-2 flex items-center gap-2 rounded-lg bg-green-50 px-3 py-2 dark:bg-green-900/20">
            <Icon name="gift" size="sm" class="text-green-600 dark:text-green-400" />
            <span class="text-sm text-green-700 dark:text-green-400">
              {{ formatDiscountText() }}
            </span>
          </div>
          <p v-else-if="promoValidation.invalid" class="mt-1 text-sm text-red-500">
            {{ promoValidation.message }}
          </p>
        </transition>
      </div>

      <!-- ==================== Subscription Tab ==================== -->
      <template v-if="!loading && activeTab === 'subscription'">
        <div v-if="decoratedPlans.length > 0" class="grid gap-6 md:grid-cols-1 lg:grid-cols-3">
          <div
            v-for="dp in decoratedPlans"
            :key="dp.plan.key"
            class="relative flex flex-col overflow-hidden rounded-2xl border transition-all hover:shadow-xl hover:-translate-y-1 cursor-pointer"
            :class="dp.featured
              ? 'border-primary-500 ring-2 ring-primary-500 shadow-lg scale-[1.02]'
              : 'border-gray-200 dark:border-dark-600 shadow-sm'"
            @click="handleBuy(dp.plan)"
          >
            <!-- Top badge ribbon -->
            <div
              v-if="dp.badge"
              class="w-full py-1.5 text-center text-xs font-bold tracking-wide"
              :class="dp.featured
                ? 'bg-primary-600 text-white'
                : 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-300'"
            >
              {{ dp.badge }}
            </div>

            <div class="flex flex-1 flex-col p-6">
              <!-- Target audience -->
              <p class="text-xs font-medium text-primary-500 dark:text-primary-400">{{ dp.audience }}</p>

              <!-- Plan name -->
              <h3 class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ dp.plan.name }}</h3>

              <!-- Price -->
              <div class="mt-3 flex items-baseline gap-1">
                <span class="text-4xl font-extrabold text-gray-900 dark:text-white">¥{{ (dp.plan.amount_fen / 100).toFixed(0) }}</span>
                <span class="text-sm text-gray-500 dark:text-dark-400">/月</span>
              </div>

              <!-- Highlight text -->
              <p class="mt-3 text-sm font-medium text-primary-600 dark:text-primary-400">{{ dp.highlight }}</p>

              <!-- Divider -->
              <div class="my-4 border-t border-gray-100 dark:border-dark-600"></div>

              <!-- Feature list -->
              <ul class="flex-1 space-y-3">
                <li v-for="(feat, fi) in dp.features" :key="fi" class="flex items-start gap-2.5 text-sm text-gray-600 dark:text-dark-300">
                  <svg class="mt-0.5 h-4 w-4 flex-shrink-0 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                  <span>{{ feat }}</span>
                </li>
              </ul>

              <!-- CTA button -->
              <button
                class="mt-6 w-full rounded-xl py-3 text-sm font-bold transition-colors"
                :class="dp.featured
                  ? 'bg-primary-600 text-white hover:bg-primary-700 shadow-md'
                  : 'bg-gray-900 text-white hover:bg-gray-800 dark:bg-dark-500 dark:hover:bg-dark-400'"
                :disabled="creatingOrder"
              >
                {{ dp.cta }}
              </button>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-12 text-gray-500 dark:text-dark-400">
          {{ t('pricing.noPlans') }}
        </div>
        <p v-if="decoratedPlans.length > 0" class="text-center text-sm text-gray-500 dark:text-dark-400">
          {{ t('pricing.stackable') }}
        </p>
      </template>

      <!-- ==================== Recharge Tab ==================== -->
      <template v-if="!loading && activeTab === 'recharge'">
        <!-- Recharge Plans (if configured) -->
        <div v-if="rechargePlans.length > 0" class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          <div
            v-for="rp in rechargePlans"
            :key="rp.key"
            class="card relative flex flex-col overflow-hidden transition-all hover:shadow-lg hover:-translate-y-1 cursor-pointer"
            :class="{ 'ring-2 ring-primary-500': rp.popular }"
            @click="handleRechargePreset(rp)"
          >
            <div v-if="rp.popular" class="absolute right-4 top-4">
              <span class="badge badge-primary text-xs">{{ t('pricing.mostPopular') }}</span>
            </div>
            <div class="flex flex-1 flex-col p-5">
              <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ rp.name }}</h3>
              <div class="mt-2 flex items-baseline gap-1">
                <span class="text-xs text-gray-500 dark:text-dark-400">{{ t('pricing.rechargeGet') }}</span>
                <span class="text-2xl font-bold text-primary-600">${{ rp.balance_amount }}</span>
              </div>
              <p v-if="rp.description" class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ rp.description }}</p>
              <div class="mt-auto pt-4">
                <div v-if="rp.pay_amount_fen !== rp.balance_amount * 100" class="mb-1 text-center">
                  <span class="text-xs text-gray-400 line-through dark:text-dark-500">¥{{ rp.balance_amount * 100 / 100 }}</span>
                </div>
                <div
                  class="flex items-center justify-center rounded-lg bg-primary-600 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
                  :class="{ 'opacity-50 pointer-events-none': creatingOrder }"
                >
                  {{ t('pricing.rechargeNow') }} ¥{{ (rp.pay_amount_fen / 100).toFixed(rp.pay_amount_fen % 100 === 0 ? 0 : 2) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Custom Amount Input -->
        <div class="card mx-auto max-w-lg p-6 space-y-5">
          <h3 class="text-center text-lg font-semibold text-gray-900 dark:text-white">{{ t('pricing.customRecharge') }}</h3>

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
                :placeholder="rechargeMinAmount > 0 ? t('pricing.minAmountHint', { amount: rechargeMinAmount }) : t('recharge.inputPlaceholder')"
                class="w-full rounded-lg border border-gray-300 bg-white py-3 pl-8 pr-4 text-lg font-medium text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-600 dark:bg-dark-800 dark:text-white dark:placeholder-dark-500"
              />
            </div>
            <p v-if="rechargeMinAmount > 0" class="mt-1 text-xs text-gray-400 dark:text-dark-500">
              {{ t('pricing.minAmountNote', { amount: rechargeMinAmount }) }}
            </p>
          </div>

          <!-- Amount Summary -->
          <div v-if="customFinalAmount > 0" class="rounded-lg bg-gray-50 dark:bg-dark-800 p-4">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600 dark:text-dark-400">{{ t('recharge.rechargeAmount') }}</span>
              <span class="text-lg font-bold text-primary-600">¥{{ customFinalAmount.toFixed(2) }}</span>
            </div>
          </div>

          <!-- Submit -->
          <button
            class="w-full rounded-lg bg-primary-600 py-3 text-base font-medium text-white transition-colors hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="customFinalAmount <= 0 || (rechargeMinAmount > 0 && customFinalAmount < rechargeMinAmount) || creatingOrder"
            @click="handleCustomRecharge"
          >
            <span v-if="creatingOrder" class="flex items-center justify-center gap-2">
              <div class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
              {{ t('recharge.creating') }}
            </span>
            <span v-else>
              {{ customFinalAmount > 0 ? t('recharge.payNow', { amount: customFinalAmount.toFixed(2) }) : t('recharge.enterAmount') }}
            </span>
          </button>
        </div>
      </template>
    </div>

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
            <!-- WeChat Pay -->
            <div
              class="flex items-center gap-3 rounded-lg border-2 p-4 cursor-pointer transition-all"
              :class="selectedPayMethod === 'wechat' ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20' : 'border-gray-200 dark:border-dark-600 hover:border-primary-300'"
              @click="selectedPayMethod = 'wechat'"
            >
              <input type="radio" :checked="selectedPayMethod === 'wechat'" class="h-5 w-5 text-primary-600" />
              <svg class="h-8 w-8 text-green-600" viewBox="0 0 24 24" fill="currentColor">
                <path d="M9.5 4C5.36 4 2 6.69 2 10c0 1.89 1.08 3.56 2.78 4.66L4 17l2.5-1.5C7.55 15.82 8.5 16 9.5 16c.34 0 .68-.02 1-.06A5.95 5.95 0 0110 14c0-3.31 2.69-6 6-6 .34 0 .68.03 1 .08C16.32 5.68 13.17 4 9.5 4zM7 9a1 1 0 110 2 1 1 0 010-2zm5 0a1 1 0 110 2 1 1 0 010-2zm4 3c-2.76 0-5 1.79-5 4s2.24 4 5 4c.71 0 1.39-.11 2-.31L20 21l-.5-1.8C20.45 18.22 21 17.16 21 16c0-2.21-2.24-4-5-4zm-1.5 2.5a.75.75 0 110 1.5.75.75 0 010-1.5zm3 0a.75.75 0 110 1.5.75.75 0 010-1.5z"/>
              </svg>
              <span class="font-medium text-gray-900 dark:text-white">微信支付</span>
            </div>

            <!-- Alipay -->
            <div
              class="flex items-center gap-3 rounded-lg border-2 p-4 cursor-pointer transition-all"
              :class="selectedPayMethod === 'alipay' ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20' : 'border-gray-200 dark:border-dark-600 hover:border-primary-300'"
              @click="selectedPayMethod = 'alipay'"
            >
              <input type="radio" :checked="selectedPayMethod === 'alipay'" class="h-5 w-5 text-primary-600" />
              <svg class="h-8 w-8 text-blue-600" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"/>
              </svg>
              <span class="font-medium text-gray-900 dark:text-white">支付宝</span>
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
          <!-- Header -->
          <div class="mb-6 text-center">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('pricing.payment.scanToPay') }}</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
              {{ t('pricing.payment.amount') }}: <span class="font-bold text-primary-600">¥{{ currentOrderAmount }}</span>
            </p>
          </div>

          <!-- QR Code -->
          <div class="flex justify-center mb-6">
            <div v-if="qrLoading" class="flex h-48 w-48 items-center justify-center">
              <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
            </div>
            <canvas v-else ref="qrCanvas" class="rounded-lg"></canvas>
          </div>

          <!-- Status -->
          <div class="mb-4 text-center">
            <div v-if="paymentStatus === 'pending'" class="flex items-center justify-center gap-2 text-sm text-gray-500 dark:text-dark-400">
              <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
              {{ t('pricing.payment.waitingPayment') }}
            </div>
            <div v-else-if="paymentStatus === 'paid'" class="flex items-center justify-center gap-2 text-sm text-green-600">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ paymentOrderType === 'subscription' ? t('pricing.payment.paymentSuccess') : t('recharge.payment.success') }}
            </div>
            <div v-else-if="paymentStatus === 'closed'" class="text-sm text-red-500">
              {{ t('pricing.payment.orderExpired') }}
            </div>
          </div>

          <!-- Countdown -->
          <div v-if="paymentStatus === 'pending' && countdown > 0" class="mb-4 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('pricing.payment.expiresIn', { minutes: Math.floor(countdown / 60), seconds: countdown % 60 }) }}
          </div>

          <!-- Actions -->
          <div class="flex gap-3">
            <button
              v-if="paymentStatus !== 'paid'"
              class="btn flex-1 border border-gray-300 dark:border-dark-500"
              @click="cancelPayment"
            >
              {{ t('pricing.payment.cancel') }}
            </button>
            <button
              v-if="paymentStatus === 'paid'"
              class="btn btn-primary flex-1"
              @click="goAfterPayment"
            >
              {{ paymentOrderType === 'subscription' ? t('pricing.payment.viewSubscription') : t('recharge.payment.viewDashboard') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { paymentAPI } from '@/api/payment'
import { validatePromoCode } from '@/api/auth'
import type { PaymentPlan, RechargePlan } from '@/api/payment'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

// Tab state
const activeTab = ref<'subscription' | 'recharge'>(
  route.query.tab === 'recharge' ? 'recharge' : 'subscription'
)

const loading = ref(true)
const plans = ref<PaymentPlan[]>([])
const rechargePlans = ref<RechargePlan[]>([])
const rechargeMinAmount = ref(0)
const creatingOrder = ref(false)
const selectedPayMethod = ref<'wechat' | 'alipay'>('wechat')
const showPayMethodModal = ref(false)
const showPaymentModal = ref(false)
const pendingPaymentAction = ref<(() => Promise<void>) | null>(null)
const qrLoading = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)
const paymentStatus = ref<'pending' | 'paid' | 'closed'>('pending')
const paymentOrderType = ref<'subscription' | 'recharge'>('subscription')
const currentOrderNo = ref('')
const currentOrderAmount = ref('')
const countdown = ref(0)
const customInput = ref('')

// Promo code state
const promoCode = ref('')
const promoValidating = ref(false)
const promoValidation = ref<{
  valid: boolean
  invalid: boolean
  discountType: string
  discountAmount: number
  discountFen: number
  finalAmountFen: number
  minOrderAmount: number
  message: string
}>({
  valid: false,
  invalid: false,
  discountType: '',
  discountAmount: 0,
  discountFen: 0,
  finalAmountFen: 0,
  minOrderAmount: 0,
  message: ''
})
let promoValidateTimeout: ReturnType<typeof setTimeout> | null = null

let pollTimer: ReturnType<typeof setInterval> | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null

const customFinalAmount = computed(() => {
  if (customInput.value !== '') {
    const val = parseFloat(customInput.value)
    return isNaN(val) || val <= 0 ? 0 : val
  }
  return 0
})

// === Hardcoded plan decoration data ===
interface DecoratedPlan {
  plan: PaymentPlan
  badge: string
  audience: string
  highlight: string
  features: string[]
  cta: string
  featured: boolean
}

const planDecorations: Record<string, Omit<DecoratedPlan, 'plan'>> = {
  // Match by amount_fen (29900 = ¥299)
  '29900': {
    badge: '',
    audience: '',
    highlight: '',
    features: [
      '$30/day',
      '$600/month',
      '',
      '',
      '',
    ],
    cta: '',
    featured: false,
  },
  // Match by amount_fen (58900 = ¥589)
  '58900': {
    badge: '',
    audience: '',
    highlight: '',
    features: [
      '$60/day',
      '$1,300/month',
      '',
      '',
      '',
      '',
    ],
    cta: '',
    featured: true,
  },
  // Match by amount_fen (118900 = ¥1189)
  '118900': {
    badge: '',
    audience: '',
    highlight: '',
    features: [
      '$120/day',
      '$2,500/month',
      '',
      '',
      '',
    ],
    cta: '',
    featured: false,
  },
}

// Chinese decoration text (hardcoded)
const planDecoZh: Record<string, Omit<DecoratedPlan, 'plan'>> = {
  '29900': {
    badge: '个人首选',
    audience: '适合：独立开发者 / 个人项目',
    highlight: '总额$600：官方价格体系，无需复杂换算',
    features: [
      '每日 $30 高速配额：专属通道，响应丝滑无卡顿',
      '总额 $600/月：充足算力，轻松覆盖日常开发',
      '永不掉线保障：额度用尽自动转按量，业务 7x24 在线',
      '超高性价比：个人开发首选，一次付费全月无忧',
      '微信扫码支付：便捷安全',
    ],
    cta: '购买开发者版',
    featured: false,
  },
  '58900': {
    badge: '生产环境推荐 / 主推套餐',
    audience: '适合：全职开发 / 生产级应用',
    highlight: '折合 $1=¥0.45，立省 55% 成本！',
    features: [
      '每日 $60 高速配额：生产级优先级，拒绝排队等待',
      '总额 $1,300/月：充足算力储备，应对高频调用',
      '无损计费机制：拒绝倍率陷阱，每一分实打实可用',
      '熔断智能防护：异常流量自动隔离，生产环境稳定如磐石',
      '永不掉线保障：额度用尽自动转按量，业务不中断',
      '微信扫码支付：便捷安全',
    ],
    cta: '立即购买，释放生产力',
    featured: true,
  },
  '118900': {
    badge: '团队专属 / SLA 保障',
    audience: '适合：开发团队 / AI 产品项目',
    highlight: '总额$2,500：顶配储备，支撑海量高并发接入',
    features: [
      '每日 $120 高速配额：团队级 VIP 优先级调度，秒级响应',
      '总额 $2,500/月：顶配算力，满足团队级高频需求',
      '高可用服务保障：极致架构，满足企业级稳定运营标准',
      '无感溢出续航：额度海量储备，确保 AI 工作流永不断供',
      '微信扫码支付：便捷安全',
    ],
    cta: '购买旗舰版',
    featured: false,
  },
}

const decoratedPlans = computed<DecoratedPlan[]>(() => {
  // Target plans by amount_fen: 29900, 58900, 118900
  const targetAmounts = [29900, 58900, 118900]
  const matched = plans.value.filter(p => targetAmounts.includes(p.amount_fen))

  // If none of the target amounts match, fall back to showing all plans with generic decoration
  const source = matched.length > 0 ? matched : plans.value

  return source.map(plan => {
    const key = String(plan.amount_fen)
    const deco = planDecoZh[key] || planDecorations[key]
    if (deco) {
      return { ...deco, plan }
    }
    // Fallback for undecorated plans
    return {
      plan,
      badge: '',
      audience: '',
      highlight: plan.description || '',
      features: [
        ...(plan.features || []),
        '微信扫码支付',
      ],
      cta: '立即购买',
      featured: false,
    }
  })
})

// === Promo Code Validation ===
function resetPromoValidation() {
  promoValidation.value = {
    valid: false,
    invalid: false,
    discountType: '',
    discountAmount: 0,
    discountFen: 0,
    finalAmountFen: 0,
    minOrderAmount: 0,
    message: ''
  }
}

function handlePromoCodeInput() {
  const code = promoCode.value.trim()
  resetPromoValidation()

  if (!code) {
    promoValidating.value = false
    return
  }

  // We don't auto-validate here because we need an amount_fen.
  // Validation happens when user clicks buy.
}

async function validatePromoForAmount(code: string, amountFen: number): Promise<boolean> {
  if (!code.trim()) return true // No promo code = ok

  promoValidating.value = true
  try {
    const result = await validatePromoCode(code, amountFen)
    if (result.valid) {
      promoValidation.value = {
        valid: true,
        invalid: false,
        discountType: result.discount_type || '',
        discountAmount: result.discount_amount || 0,
        discountFen: result.discount_fen || 0,
        finalAmountFen: result.final_amount_fen || amountFen,
        minOrderAmount: result.min_order_amount || 0,
        message: ''
      }
      return true
    } else {
      promoValidation.value = {
        valid: false,
        invalid: true,
        discountType: '',
        discountAmount: 0,
        discountFen: 0,
        finalAmountFen: 0,
        minOrderAmount: result.min_order_amount || 0,
        message: getPromoErrorMessage(result.error_code)
      }
      return false
    }
  } catch {
    promoValidation.value = {
      valid: false,
      invalid: true,
      discountType: '',
      discountAmount: 0,
      discountFen: 0,
      finalAmountFen: 0,
      minOrderAmount: 0,
      message: t('pricing.promoCodeInvalid')
    }
    return false
  } finally {
    promoValidating.value = false
  }
}

function getPromoErrorMessage(errorCode?: string): string {
  switch (errorCode) {
    case 'PROMO_CODE_NOT_FOUND':
      return t('pricing.promoCodeNotFound')
    case 'PROMO_CODE_EXPIRED':
      return t('pricing.promoCodeExpired')
    case 'PROMO_CODE_DISABLED':
      return t('pricing.promoCodeDisabled')
    case 'PROMO_CODE_MAX_USED':
      return t('pricing.promoCodeMaxUsed')
    case 'PROMO_CODE_ALREADY_USED':
      return t('pricing.promoCodeAlreadyUsed')
    case 'PROMO_CODE_MIN_ORDER':
      return t('pricing.promoCodeMinOrder', { amount: (promoValidation.value.minOrderAmount / 100).toFixed(0) })
    default:
      return t('pricing.promoCodeInvalid')
  }
}

function formatDiscountText(): string {
  const v = promoValidation.value
  if (!v.valid) return ''
  const amountYuan = (v.discountFen / 100).toFixed(2)
  if (v.discountType === 'percentage') {
    return t('pricing.promoDiscountPercent', { percent: v.discountAmount, amount: amountYuan })
  }
  return t('pricing.promoDiscountFixed', { amount: amountYuan })
}

onMounted(async () => {
  try {
    const [allPlans, rechargeInfo] = await Promise.all([
      paymentAPI.getPlans(),
      paymentAPI.getRechargeInfo()
    ])
    plans.value = allPlans.filter(p => (p.type || 'subscription') === 'subscription')
    rechargePlans.value = rechargeInfo.plans || []
    rechargeMinAmount.value = rechargeInfo.min_amount || 0
  } catch {
    // silently fail
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  clearTimers()
  if (promoValidateTimeout) {
    clearTimeout(promoValidateTimeout)
  }
})

function clearTimers() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  if (countdownTimer) { clearInterval(countdownTimer); countdownTimer = null }
}

// === Payment method selection ===
async function confirmPayMethod() {
  showPayMethodModal.value = false
  if (pendingPaymentAction.value) {
    await pendingPaymentAction.value()
    pendingPaymentAction.value = null
  }
}

// === Subscription purchase ===
async function handleBuy(plan: PaymentPlan) {
  if (creatingOrder.value) return

  // Validate promo code if entered
  const code = promoCode.value.trim()
  if (code) {
    const valid = await validatePromoForAmount(code, plan.amount_fen)
    if (!valid) return
  }

  // Show payment method selection modal
  pendingPaymentAction.value = async () => {
    paymentOrderType.value = 'subscription'
    creatingOrder.value = true
    try {
      const order = await paymentAPI.createOrder(plan.key, code || undefined, selectedPayMethod.value)
      await showQRModal(order)
    } catch (err: any) {
      const msg = err?.message || err?.response?.data?.message || t('pricing.payment.createFailed')
      alert(msg)
    } finally {
      creatingOrder.value = false
    }
  }
  showPayMethodModal.value = true
}

// === Recharge preset ===
async function handleRechargePreset(rp: RechargePlan) {
  if (creatingOrder.value) return

  // Validate promo code if entered
  const code = promoCode.value.trim()
  if (code) {
    const valid = await validatePromoForAmount(code, rp.pay_amount_fen)
    if (!valid) return
  }

  // Show payment method selection modal
  pendingPaymentAction.value = async () => {
    paymentOrderType.value = 'recharge'
    creatingOrder.value = true
    try {
      const payYuan = rp.pay_amount_fen / 100
      const order = await paymentAPI.createRechargeOrder(payYuan, code || undefined, selectedPayMethod.value)
      await showQRModal(order)
    } catch (err: any) {
      const msg = err?.message || err?.response?.data?.message || t('recharge.payment.createFailed')
      alert(msg)
    } finally {
      creatingOrder.value = false
    }
  }
  showPayMethodModal.value = true
}

// === Custom recharge ===
async function handleCustomRecharge() {
  if (customFinalAmount.value <= 0 || creatingOrder.value) return
  if (rechargeMinAmount.value > 0 && customFinalAmount.value < rechargeMinAmount.value) return

  // Validate promo code if entered
  const code = promoCode.value.trim()
  const amountFen = Math.round(customFinalAmount.value * 100)
  if (code) {
    const valid = await validatePromoForAmount(code, amountFen)
    if (!valid) return
  }

  // Show payment method selection modal
  pendingPaymentAction.value = async () => {
    paymentOrderType.value = 'recharge'
    creatingOrder.value = true
    try {
      const order = await paymentAPI.createRechargeOrder(customFinalAmount.value, code || undefined, selectedPayMethod.value)
      await showQRModal(order)
    } catch (err: any) {
      const msg = err?.message || err?.response?.data?.message || t('recharge.payment.createFailed')
      alert(msg)
    } finally {
      creatingOrder.value = false
    }
  }
  showPayMethodModal.value = true
}

// === Shared QR modal logic ===
async function showQRModal(order: { order_no: string; code_url: string | null; amount_fen: number; expired_at: string }) {
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
}

function startPolling() {
  pollTimer = setInterval(async () => {
    try {
      const order = await paymentAPI.queryOrder(currentOrderNo.value)
      if (order.status === 'paid') {
        paymentStatus.value = 'paid'
        clearTimers()
      } else if (order.status === 'closed') {
        paymentStatus.value = 'closed'
        clearTimers()
      }
    } catch {
      // ignore poll errors
    }
  }, 3000)
}

function cancelPayment() {
  showPaymentModal.value = false
  clearTimers()
}

function goAfterPayment() {
  showPaymentModal.value = false
  clearTimers()
  if (paymentOrderType.value === 'subscription') {
    router.push('/subscriptions')
  } else {
    router.push('/dashboard')
  }
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>