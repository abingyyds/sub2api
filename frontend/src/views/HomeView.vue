<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Themed Home -->
  <div v-else :class="themeRoot" class="relative flex min-h-screen flex-col overflow-hidden">
    <!-- Background decorations: adapt per-theme -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <template v-if="themeTemplate === 'terminal'">
        <div class="absolute inset-0 bg-[linear-gradient(rgba(16,185,129,0.05)_1px,transparent_1px),linear-gradient(90deg,rgba(16,185,129,0.05)_1px,transparent_1px)] bg-[size:32px_32px]"></div>
      </template>
      <template v-else-if="themeTemplate === 'summit'">
        <div class="absolute -right-40 -top-40 h-96 w-96 rounded-full bg-orange-400/20 blur-3xl"></div>
        <div class="absolute -bottom-40 -left-40 h-96 w-96 rounded-full bg-amber-500/15 blur-3xl"></div>
      </template>
      <template v-else-if="themeTemplate === 'aurora'">
        <div class="absolute -right-40 -top-40 h-[32rem] w-[32rem] animate-pulse rounded-full bg-sky-400/20 blur-3xl"></div>
        <div class="absolute -bottom-40 -left-40 h-[32rem] w-[32rem] animate-pulse rounded-full bg-indigo-500/20 blur-3xl"></div>
        <div class="absolute left-1/3 top-1/4 h-72 w-72 animate-pulse rounded-full bg-cyan-400/20 blur-3xl"></div>
      </template>
      <template v-else>
        <div class="absolute -right-40 -top-40 h-96 w-96 rounded-full bg-primary-400/20 blur-3xl"></div>
        <div class="absolute -bottom-40 -left-40 h-96 w-96 rounded-full bg-primary-500/15 blur-3xl"></div>
        <div class="absolute left-1/3 top-1/4 h-72 w-72 rounded-full bg-primary-300/10 blur-3xl"></div>
        <div class="absolute inset-0 bg-[linear-gradient(rgba(20,184,166,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(20,184,166,0.03)_1px,transparent_1px)] bg-[size:64px_64px]"></div>
      </template>
    </div>

    <component
      :is="ThemeComponent"
      :site-name="siteName"
      :site-logo="siteLogo"
      :site-subtitle="siteSubtitle"
      :doc-url="docUrl"
      :hero-title="heroTitle"
      :hero-description="heroDescription"
      :cta-text="ctaText"
      :feature-tags="featureTags"
      :registration-notice="registrationNotice"
      :allow-sub-site-open="allowSubSiteOpen"
      :sub-site-open-price="subSiteOpenPrice"
      :is-authenticated="isAuthenticated"
      :dashboard-path="dashboardPath"
      :user-initial="userInitial"
      :plans="plans"
      :creating-order="creatingOrder"
      :github-url="githubUrl"
      :current-year="currentYear"
      :is-dark="isDark"
      @toggle-theme="toggleTheme"
      @buy-plan="handleBuyPlan"
    >
      <template #body>
        <!-- Pricing -->
        <div v-if="plans.length > 0" class="mb-16">
          <div class="mb-8 text-center">
            <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">{{ t('pricing.title') }}</h2>
            <p class="text-sm text-gray-600 dark:text-dark-400">{{ t('pricing.subtitle') }}</p>
          </div>
          <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="(plan, index) in plans"
              :key="plan.key"
              class="card relative flex flex-col overflow-hidden transition-all hover:shadow-lg hover:-translate-y-1 cursor-pointer"
              :class="{ 'ring-2 ring-primary-500': index === 0 }"
              @click="handleBuyPlan(plan)"
            >
              <div v-if="index === 0" class="absolute right-4 top-4">
                <span class="badge badge-primary text-xs">{{ t('pricing.recommended') }}</span>
              </div>
              <div class="flex flex-1 flex-col p-6">
                <h3 class="text-xl font-bold text-gray-900 dark:text-white">{{ plan.name }}</h3>
                <p v-if="plan.description" class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ plan.description }}</p>
                <div class="mt-3 text-3xl font-bold text-primary-600">
                  ¥{{ (plan.amount_fen / 100).toFixed(plan.amount_fen % 100 === 0 ? 0 : 2) }}
                </div>
                <ul class="mt-5 flex-1 space-y-2.5">
                  <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                    <svg class="h-4 w-4 flex-shrink-0 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                    </svg>
                    {{ t('pricing.validityDays', { days: plan.validity_days }) }}
                  </li>
                  <li v-for="(feature, fi) in (plan.features || [])" :key="fi" class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                    <svg class="h-4 w-4 flex-shrink-0 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                    </svg>
                    {{ feature }}
                  </li>
                  <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                    <svg class="h-4 w-4 flex-shrink-0 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                    </svg>
                    {{ t('pricing.wechatPay') }}
                  </li>
                </ul>
                <div
                  class="mt-4 flex items-center justify-center rounded-lg bg-primary-600 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
                  :class="{ 'opacity-50 pointer-events-none': creatingOrder }"
                >
                  <template v-if="creatingOrder">
                    <div class="mr-2 h-3.5 w-3.5 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
                    {{ t('common.processing') }}
                  </template>
                  <template v-else>
                    {{ t('pricing.buyNow') }}
                  </template>
                </div>
              </div>
            </div>
          </div>
          <p class="mt-6 text-center text-sm text-gray-500 dark:text-dark-400">
            {{ t('pricing.repurchaseRule') }}
          </p>
        </div>

        <!-- Providers -->
        <div class="mb-8 text-center">
          <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">{{ t('home.providers.title') }}</h2>
          <p class="text-sm text-gray-600 dark:text-dark-400">{{ t('home.providers.description') }}</p>
        </div>
        <div class="mb-16 flex flex-wrap items-center justify-center gap-4">
          <div v-for="(provider, pi) in providers" :key="pi"
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg" :class="provider.badgeClass">
              <span class="text-xs font-bold text-white">{{ provider.initial }}</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ provider.name }}</span>
            <span class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400">
              {{ t('home.providers.supported') }}
            </span>
          </div>
          <div class="flex items-center gap-2 rounded-xl border border-gray-200/50 bg-white/40 px-5 py-3 opacity-60 backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/40">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-gray-500 to-gray-600">
              <span class="text-xs font-bold text-white">+</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.more') }}</span>
            <span class="rounded bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-dark-700 dark:text-dark-400">
              {{ t('home.providers.soon') }}
            </span>
          </div>
        </div>
      </template>
    </component>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left">
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-4">
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white">
            {{ t('home.docs') }}
          </a>
          <a :href="githubUrl" target="_blank" rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white">
            GitHub
          </a>
        </div>
      </div>
    </footer>

    <!-- Payment QR Code Modal -->
    <Teleport to="body">
      <div v-if="showPaymentModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="cancelPayment">
        <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-2xl dark:bg-dark-700 mx-4">
          <div class="mb-6 text-center">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ t('pricing.payment.scanToPay') }}</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
              {{ t('pricing.payment.amount') }}: <span class="font-bold text-primary-600">¥{{ currentOrderAmount }}</span>
            </p>
          </div>
          <div class="flex justify-center mb-6">
            <div class="relative h-48 w-48">
              <canvas ref="qrCanvas" class="rounded-lg"></canvas>
              <div v-if="qrLoading" class="absolute inset-0 flex items-center justify-center rounded-lg bg-white/90 dark:bg-dark-700/90">
                <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
              </div>
            </div>
          </div>
          <div class="mb-4 text-center">
            <div v-if="paymentStatus === 'pending'" class="flex items-center justify-center gap-2 text-sm text-gray-500 dark:text-dark-400">
              <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
              {{ t('pricing.payment.waitingPayment') }}
            </div>
            <div v-else-if="paymentStatus === 'paid'" class="flex items-center justify-center gap-2 text-sm text-green-600">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ t('pricing.payment.paymentSuccess') }}
            </div>
            <div v-else-if="paymentStatus === 'closed'" class="text-sm text-red-500">
              {{ t('pricing.payment.orderExpired') }}
            </div>
          </div>
          <div v-if="paymentStatus === 'pending' && countdown > 0" class="mb-4 text-center text-xs text-gray-400 dark:text-dark-500">
            {{ t('pricing.payment.expiresIn', { minutes: Math.floor(countdown / 60), seconds: countdown % 60 }) }}
          </div>
          <div class="flex gap-3">
            <button v-if="paymentStatus !== 'paid'" class="btn flex-1 border border-gray-300 dark:border-dark-500" @click="cancelPayment">
              {{ t('pricing.payment.cancel') }}
            </button>
            <button v-if="paymentStatus === 'paid'" class="btn btn-primary flex-1" @click="goToSubscriptions">
              {{ t('pricing.payment.viewOrderHistory') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, shallowRef, defineAsyncComponent } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import { paymentAPI } from '@/api/payment'
import type { PaymentPlan } from '@/api/payment'

const StarterHome = defineAsyncComponent(() => import('@/components/themes/StarterHome.vue'))
const AuroraHome = defineAsyncComponent(() => import('@/components/themes/AuroraHome.vue'))
const SummitHome = defineAsyncComponent(() => import('@/components/themes/SummitHome.vue'))
const TerminalHome = defineAsyncComponent(() => import('@/components/themes/TerminalHome.vue'))

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'cCoder.me')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const themeTemplate = computed(() => appStore.cachedPublicSettings?.theme_template || 'starter')

const heroTitle = computed(() => siteName.value)
const heroDescription = computed(() => siteSubtitle.value)
const ctaText = computed(() => t('home.getStarted'))
const featureTags = computed(() => [
  t('home.tags.subscriptionToApi'),
  t('home.tags.stickySession'),
  t('home.tags.realtimeBilling'),
])
const registrationNotice = computed(() => {
  const mode = appStore.cachedPublicSettings?.registration_mode
  if (mode === 'invite') return '当前分站仅支持邀请码注册。'
  if (mode === 'closed') return '当前分站已关闭新用户注册。'
  return ''
})
const allowSubSiteOpen = computed(() => Boolean(
  appStore.cachedPublicSettings?.subsite_entry_enabled
  && appStore.cachedPublicSettings?.allow_sub_site
  && (appStore.cachedPublicSettings?.subsite_price_fen || 0) > 0
))
const subSiteOpenPrice = computed(() => fenToYuan(appStore.cachedPublicSettings?.subsite_price_fen || 0))

const ThemeComponent = shallowRef<any>(StarterHome)
const themeRoot = computed(() => {
  switch (themeTemplate.value) {
    case 'aurora':
      ThemeComponent.value = AuroraHome
      return 'bg-gradient-to-br from-sky-50 via-indigo-50/70 to-cyan-100 dark:from-slate-950 dark:via-indigo-950 dark:to-slate-950'
    case 'summit':
      ThemeComponent.value = SummitHome
      return 'bg-gradient-to-br from-stone-50 via-amber-50/40 to-orange-100 dark:from-neutral-950 dark:via-stone-900 dark:to-neutral-950'
    case 'terminal':
      ThemeComponent.value = TerminalHome
      return 'bg-zinc-950 text-emerald-100'
    default:
      ThemeComponent.value = StarterHome
      return 'bg-gradient-to-br from-gray-50 via-primary-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950'
  }
})

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isDark = ref(document.documentElement.classList.contains('dark'))
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => {
  if (isAdmin.value) return '/admin/dashboard'
  if (appStore.cachedPublicSettings?.is_subsite && authStore.ownedSites[0]?.id) {
    return `/subsite-admin/${authStore.ownedSites[0].id}/dashboard`
  }
  return '/dashboard'
})
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

const currentYear = computed(() => new Date().getFullYear())

const providers = computed(() => [
  { name: t('home.providers.claude'), initial: 'C', badgeClass: 'bg-gradient-to-br from-orange-400 to-orange-500' },
  { name: 'GPT', initial: 'G', badgeClass: 'bg-gradient-to-br from-green-500 to-green-600' },
  { name: t('home.providers.gemini'), initial: 'G', badgeClass: 'bg-gradient-to-br from-blue-500 to-blue-600' },
  { name: t('home.providers.antigravity'), initial: 'A', badgeClass: 'bg-gradient-to-br from-rose-500 to-pink-600' },
])

function fenToYuan(value: number) {
  return (Number(value || 0) / 100).toFixed(value % 100 === 0 ? 0 : 2)
}

const plans = ref<PaymentPlan[]>([])
const creatingOrder = ref(false)
const showPaymentModal = ref(false)
const qrLoading = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)
const paymentStatus = ref<'pending' | 'paid' | 'closed'>('pending')
const currentOrderNo = ref('')
const currentOrderAmount = ref('')
const countdown = ref(0)

let pollTimer: ReturnType<typeof setInterval> | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null

function hideRateFeatures(plan: PaymentPlan): PaymentPlan {
  return {
    ...plan,
    features: (plan.features || []).filter(feature => !/(倍率|费率|\d+(?:\.\d+)?\s*x\b)/i.test(feature))
  }
}

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()
  void authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
  loadPlans()
})

onUnmounted(() => {
  clearTimers()
})

async function loadPlans() {
  try {
    const allPlans = await paymentAPI.getPlans()
    plans.value = allPlans
      .filter(p => (p.type || 'subscription') === 'subscription')
      .map(hideRateFeatures)
  } catch {
    // silently fail
  }
}

function clearTimers() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  if (countdownTimer) { clearInterval(countdownTimer); countdownTimer = null }
}

function getCreateOrderErrorMessage(err: any) {
  if (err?.reason === 'SUBSCRIPTION_REPURCHASE_BLOCKED') {
    return t('pricing.payment.subscriptionRepurchaseBlocked')
  }
  return err?.message || err?.response?.data?.message || t('pricing.payment.createFailed')
}

async function handleBuyPlan(plan: PaymentPlan) {
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  if (creatingOrder.value) return

  creatingOrder.value = true
  try {
    const order = await paymentAPI.createOrder(plan.key)
    currentOrderNo.value = order.order_no
    currentOrderAmount.value = (order.amount_fen / 100).toFixed(order.amount_fen % 100 === 0 ? 0 : 2)
    paymentStatus.value = 'pending'
    qrLoading.value = true
    showPaymentModal.value = true

    const expiresAt = new Date(order.expired_at).getTime()
    countdown.value = Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))

    try {
      if (order.code_url) {
        await nextTick()
        if (qrCanvas.value) {
          const QRCode = (await import('qrcode')).default
          await QRCode.toCanvas(qrCanvas.value, order.code_url, {
            width: 192, margin: 2,
            color: { dark: '#000000', light: '#ffffff' },
          })
        }
      }
    } finally {
      qrLoading.value = false
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
    alert(getCreateOrderErrorMessage(err))
  } finally {
    qrLoading.value = false
    creatingOrder.value = false
  }
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
      // ignore
    }
  }, 3000)
}

async function syncPaymentSuccessState() {
  try {
    await authStore.refreshUser()
    appStore.showSuccess(t('pricing.payment.paymentSuccessToast'), 5000)
  } catch (error) {
    console.error('Failed to refresh home payment success state:', error)
    appStore.showError('支付成功，但页面数据刷新失败，请手动刷新页面查看最新状态')
  }
}

function cancelPayment() {
  showPaymentModal.value = false
  clearTimers()
}

function goToSubscriptions() {
  showPaymentModal.value = false
  clearTimers()
  router.push({
    path: '/order-history',
    query: { highlight: currentOrderNo.value, success: '1' }
  })
}
</script>
