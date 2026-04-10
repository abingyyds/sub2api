<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-6">
      <div v-if="loading" class="flex items-center justify-center py-16">
        <div class="h-9 w-9 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <template v-else>
        <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white">代理中心</h1>
              <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
                先完成实名、合同确认和开通费支付，再提交代理申请。
              </p>
            </div>
            <div class="flex flex-wrap gap-2">
              <span class="badge" :class="statusBadgeClass(status.agent_status)">
                {{ statusLabel(status.agent_status) }}
              </span>
              <span class="badge" :class="profileBadgeClass(status.profile?.identity_status)">
                实名：{{ identityLabel(status.profile?.identity_status) }}
              </span>
              <span class="badge" :class="profileBadgeClass(status.profile?.contract_status)">
                合同：{{ contractLabel(status.profile?.contract_status) }}
              </span>
            </div>
          </div>

          <div v-if="status.profile?.is_frozen" class="mt-4 rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900/40 dark:bg-red-950/30 dark:text-red-300">
            当前代理已被冻结：{{ status.profile?.frozen_reason || '请联系管理员处理' }}
          </div>
          <div v-else-if="status.agent_status === 'pending'" class="mt-4 rounded-2xl border border-yellow-200 bg-yellow-50 px-4 py-3 text-sm text-yellow-700 dark:border-yellow-900/40 dark:bg-yellow-950/30 dark:text-yellow-300">
            资料已提交，等待管理员审核。
          </div>
          <div v-else-if="status.agent_status === 'rejected'" class="mt-4 rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900/40 dark:bg-red-950/30 dark:text-red-300">
            申请已被驳回。你可以补充资料后重新提交。
          </div>
        </section>

        <template v-if="status.agent_status === 'approved'">
          <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            <div class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
              <p class="text-sm text-gray-500 dark:text-dark-400">站内消费余额</p>
              <p class="mt-2 text-3xl font-bold text-gray-900 dark:text-white">{{ formatMoney(status.wallet.site_balance) }}</p>
              <p class="mt-2 text-xs text-gray-500 dark:text-dark-400">邀请代理返佣会直接进入这里，只能站内消费。</p>
            </div>
            <div class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
              <p class="text-sm text-gray-500 dark:text-dark-400">冻结余额</p>
              <p class="mt-2 text-3xl font-bold text-orange-600 dark:text-orange-400">{{ formatMoney(status.wallet.frozen_balance) }}</p>
              <p class="mt-2 text-xs text-gray-500 dark:text-dark-400">可提现收益会先冻结 {{ status.withdraw_freeze_days }} 天。</p>
            </div>
            <div class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
              <p class="text-sm text-gray-500 dark:text-dark-400">可提现余额</p>
              <p class="mt-2 text-3xl font-bold text-emerald-600 dark:text-emerald-400">{{ formatMoney(status.wallet.withdrawable_balance) }}</p>
              <p class="mt-2 text-xs text-gray-500 dark:text-dark-400">仅这部分余额允许申请提现。</p>
            </div>
            <div class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900">
              <p class="text-sm text-gray-500 dark:text-dark-400">累计已提现</p>
              <p class="mt-2 text-3xl font-bold text-indigo-600 dark:text-indigo-400">{{ formatMoney(status.wallet.total_withdrawn) }}</p>
              <p class="mt-2 text-xs text-gray-500 dark:text-dark-400">提现开放时间：每周 {{ weekdayLabel(status.withdraw_window.weekday) }} {{ windowLabel }}</p>
            </div>
          </section>

          <section class="grid gap-6 lg:grid-cols-[1.2fr_1fr]">
            <div class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
              <p class="text-sm text-gray-500 dark:text-dark-400">邀请链接</p>
              <div class="mt-3 break-all rounded-2xl bg-gray-50 px-4 py-3 text-sm text-gray-800 dark:bg-dark-800 dark:text-dark-100">
                {{ inviteLink }}
              </div>
              <div class="mt-4 flex gap-3">
                <button class="btn btn-primary" @click="copyLink">复制邀请链接</button>
                <router-link class="btn btn-secondary" to="/agent/commissions">查看返佣明细</router-link>
              </div>
            </div>

            <div class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
              <p class="text-sm text-gray-500 dark:text-dark-400">经营概览</p>
              <div class="mt-4 grid grid-cols-2 gap-4">
                <div class="rounded-2xl bg-gray-50 p-4 dark:bg-dark-800">
                  <p class="text-xs text-gray-500 dark:text-dark-400">直属用户</p>
                  <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">{{ dashboard.total_sub_users }}</p>
                </div>
                <div class="rounded-2xl bg-gray-50 p-4 dark:bg-dark-800">
                  <p class="text-xs text-gray-500 dark:text-dark-400">累计返佣</p>
                  <p class="mt-2 text-2xl font-bold text-primary-600 dark:text-primary-400">{{ formatMoney(dashboard.total_commission) }}</p>
                </div>
                <div class="rounded-2xl bg-gray-50 p-4 dark:bg-dark-800">
                  <p class="text-xs text-gray-500 dark:text-dark-400">今日新增</p>
                  <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">{{ dashboard.today_new_users }}</p>
                </div>
                <div class="rounded-2xl bg-gray-50 p-4 dark:bg-dark-800">
                  <p class="text-xs text-gray-500 dark:text-dark-400">今日开通费成交</p>
                  <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">{{ formatMoney(dashboard.today_recharge) }}</p>
                </div>
              </div>
              <div class="mt-4 grid gap-3 sm:grid-cols-2">
                <router-link to="/agent/sub-users" class="rounded-2xl border border-gray-200 px-4 py-3 text-sm font-medium text-gray-700 transition hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-400">
                  查看直属用户
                </router-link>
                <router-link to="/agent/financial-logs" class="rounded-2xl border border-gray-200 px-4 py-3 text-sm font-medium text-gray-700 transition hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-400">
                  查看资金动态
                </router-link>
              </div>
            </div>
          </section>
        </template>

        <template v-else>
          <section class="grid gap-6 xl:grid-cols-3">
            <div class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900 xl:col-span-2">
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-2xl bg-primary-50 dark:bg-primary-900/20">
                  <Icon name="users" size="md" class="text-primary-600 dark:text-primary-400" />
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-gray-900 dark:text-white">步骤 1：填写实名资料并确认合同</h2>
                  <p class="text-sm text-gray-500 dark:text-dark-400">一期先采用站内留痕，不接第三方实名和电子签。</p>
                </div>
              </div>

              <div class="mt-6 grid gap-4 md:grid-cols-2">
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-dark-300">真实姓名</label>
                  <input v-model="profileForm.real_name" class="input w-full" placeholder="请输入真实姓名" />
                </div>
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-dark-300">手机号</label>
                  <input v-model="profileForm.phone" class="input w-full" placeholder="请输入本人手机号" />
                </div>
                <div class="md:col-span-2">
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-dark-300">身份证号</label>
                  <input v-model="profileForm.id_card_no" class="input w-full" placeholder="请输入身份证号" />
                </div>
              </div>

              <label class="mt-5 flex items-start gap-3 rounded-2xl bg-gray-50 px-4 py-3 text-sm text-gray-700 dark:bg-dark-800 dark:text-dark-200">
                <input v-model="profileForm.contract_accepted" type="checkbox" class="mt-1 h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
                <span>我已阅读并同意当前代理合作条款，并将在下方完成电子签字。</span>
              </label>

              <div class="mt-5 rounded-2xl border border-gray-200 p-4 dark:border-dark-700">
                <div class="flex items-center justify-between gap-3">
                  <div>
                    <h3 class="text-sm font-medium text-gray-900 dark:text-white">合同 PDF</h3>
                    <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">请先阅读管理员上传的合同模板，再完成签字。</p>
                  </div>
                  <a
                    v-if="status.contract_template"
                    :href="status.contract_template"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="btn btn-secondary btn-sm"
                  >
                    新窗口查看
                  </a>
                </div>

                <div v-if="status.contract_template" class="mt-4 overflow-hidden rounded-2xl border border-gray-200 dark:border-dark-700">
                  <iframe :src="status.contract_template" class="h-72 w-full bg-white"></iframe>
                </div>
                <p v-else class="mt-4 text-sm text-amber-600 dark:text-amber-400">管理员还没有上传合同 PDF，当前无法完成签署。</p>
              </div>

              <div class="mt-5 rounded-2xl border border-gray-200 p-4 dark:border-dark-700">
                <div class="flex items-center justify-between gap-3">
                  <div>
                    <h3 class="text-sm font-medium text-gray-900 dark:text-white">合同签字</h3>
                    <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">请使用鼠标或手指在签字板内签名，系统会保存签字图、时间和 IP 作为签署留痕。</p>
                  </div>
                  <button type="button" class="btn btn-secondary btn-sm" @click="clearSignature">清空重签</button>
                </div>

                <div class="mt-4 overflow-hidden rounded-2xl border border-dashed border-gray-300 bg-white dark:border-dark-600">
                  <canvas
                    ref="signatureCanvas"
                    class="block h-40 w-full touch-none"
                    @pointerdown="startSignature"
                    @pointermove="moveSignature"
                    @pointerup="endSignature"
                    @pointerleave="endSignature"
                    @pointercancel="endSignature"
                  ></canvas>
                </div>

                <p v-if="signatureError" class="mt-2 text-xs text-red-500">{{ signatureError }}</p>

                <div v-if="status.profile?.contract_signature_data" class="mt-4">
                  <p class="mb-2 text-xs text-gray-500 dark:text-dark-400">当前已保存签字</p>
                  <img :src="status.profile.contract_signature_data" alt="合同签字" class="max-h-24 rounded-lg border border-gray-200 bg-white p-2 dark:border-dark-700" />
                </div>
              </div>

              <div class="mt-5 flex items-center gap-3">
                <button class="btn btn-primary" :disabled="savingProfile || !canSaveProfile" @click="handleSaveProfile">
                  {{ savingProfile ? '保存中...' : '保存资料' }}
                </button>
                <span v-if="status.profile?.identity_status === 'submitted'" class="text-sm text-emerald-600 dark:text-emerald-400">
                  资料已保存
                </span>
              </div>
            </div>

            <div class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-2xl bg-orange-50 dark:bg-orange-900/20">
                  <Icon name="dollar" size="md" class="text-orange-600 dark:text-orange-400" />
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-gray-900 dark:text-white">步骤 2：支付开通费</h2>
                  <p class="text-sm text-gray-500 dark:text-dark-400">当前开通费 {{ formatMoney(status.activation_fee) }}</p>
                </div>
              </div>

              <div class="mt-5 rounded-2xl bg-gray-50 p-4 text-sm text-gray-700 dark:bg-dark-800 dark:text-dark-200">
                <p>支付状态：<span class="font-semibold">{{ status.profile?.activation_fee_paid_at ? '已支付' : '未支付' }}</span></p>
                <p class="mt-2">返佣只针对这笔代理开通费，且仅结算一级邀请关系。</p>
              </div>

              <button class="btn btn-primary mt-5 w-full" :disabled="creatingOrder || !!status.profile?.activation_fee_paid_at" @click="handleCreateActivationOrder">
                {{ status.profile?.activation_fee_paid_at ? '已完成支付' : creatingOrder ? '创建订单中...' : '去支付开通费' }}
              </button>
            </div>
          </section>

          <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <h2 class="text-lg font-semibold text-gray-900 dark:text-white">步骤 3：提交代理申请</h2>
                <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
                  只有实名资料、合同确认和开通费支付都完成后，才可以提交申请。
                </p>
              </div>
              <button class="btn btn-primary" :disabled="applying || !status.can_apply" @click="handleApply">
                {{ applying ? '提交中...' : '提交代理申请' }}
              </button>
            </div>
          </section>
        </template>
      </template>
    </div>

    <Teleport to="body">
      <div v-if="showPaymentModal" class="fixed inset-0 z-50 flex items-center justify-center px-4">
        <div class="absolute inset-0 bg-black/50" @click="closePaymentModal"></div>
        <div class="relative w-full max-w-md rounded-3xl bg-white p-6 shadow-2xl dark:bg-dark-900">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">扫码支付代理开通费</h3>
          <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
            订单号：{{ activationOrder.order_no || '-' }}
          </p>
          <div class="mt-6 flex justify-center">
            <canvas ref="qrCanvas" class="rounded-2xl border border-gray-200 bg-white p-3 dark:border-dark-700"></canvas>
          </div>
          <p class="mt-4 text-center text-sm text-gray-500 dark:text-dark-400">
            支付成功后会自动刷新状态
          </p>
          <div class="mt-6 flex justify-end gap-3">
            <button class="btn btn-secondary" @click="closePaymentModal">关闭</button>
            <button class="btn btn-primary" @click="refreshStatus">手动刷新</button>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { agentAPI, type AgentDashboardStats, type AgentStatus } from '@/api/agent'
import { paymentAPI, type CreateOrderResponse } from '@/api/payment'
import { useAppStore } from '@/stores'

const appStore = useAppStore()

const loading = ref(true)
const savingProfile = ref(false)
const creatingOrder = ref(false)
const applying = ref(false)
const showPaymentModal = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)
const signatureCanvas = ref<HTMLCanvasElement | null>(null)
const pollTimer = ref<number | null>(null)
const signatureError = ref('')
const isSigning = ref(false)
let signatureContext: CanvasRenderingContext2D | null = null

const activationOrder = ref<CreateOrderResponse & { order_no?: string }>({
  order_no: '',
  code_url: null,
  amount_fen: 0,
  discount_amount: 0,
  expired_at: ''
})

const status = ref<AgentStatus>({
  enabled: true,
  is_agent: false,
  agent_status: '',
  commission_rate: 0,
  invite_code: '',
  can_apply: false,
  activation_fee: 0,
  contract_template: '',
  profile: {
    user_id: 0,
    real_name: '',
    id_card_no: '',
    phone: '',
    identity_status: 'unsubmitted',
    identity_submitted_at: null,
    contract_status: 'unsigned',
    contract_version: 'v1',
    contract_signed_at: null,
    contract_ip: '',
    contract_signature_data: '',
    activation_order_id: null,
    activation_fee_paid_at: null,
    frozen_balance: 0,
    withdrawable_balance: 0,
    total_withdrawn: 0,
    is_frozen: false,
    frozen_reason: ''
  },
  wallet: {
    site_balance: 0,
    frozen_balance: 0,
    withdrawable_balance: 0,
    total_withdrawn: 0
  },
  withdraw_freeze_days: 7,
  withdraw_window: {
    weekday: 5,
    start_hour: 14,
    end_hour: 24,
    label: ''
  }
})

const dashboard = ref<AgentDashboardStats>({
  total_sub_users: 0,
  total_recharge: 0,
  total_consumed: 0,
  total_commission: 0,
  pending_commission: 0,
  settled_commission: 0,
  today_new_users: 0,
  today_recharge: 0,
  site_balance: 0,
  frozen_balance: 0,
  withdrawable_balance: 0,
  total_withdrawn: 0
})

const inviteLink = ref('')
const profileForm = ref({
  real_name: '',
  id_card_no: '',
  phone: '',
  contract_accepted: false,
  contract_signature_data: ''
})

const canSaveProfile = computed(() =>
  profileForm.value.real_name.trim() &&
  profileForm.value.id_card_no.trim() &&
  profileForm.value.phone.trim() &&
  status.value.contract_template &&
  profileForm.value.contract_accepted &&
  profileForm.value.contract_signature_data.trim()
)

const windowLabel = computed(() => `${String(status.value.withdraw_window.start_hour).padStart(2, '0')}:00-${String(status.value.withdraw_window.end_hour).padStart(2, '0')}:00`)

onMounted(async () => {
  await refreshStatus()
})

onBeforeUnmount(() => {
  stopPolling()
})

async function refreshStatus() {
  loading.value = true
  try {
    status.value = await agentAPI.getStatus()
    profileForm.value.real_name = status.value.profile?.real_name || ''
    profileForm.value.id_card_no = status.value.profile?.id_card_no || ''
    profileForm.value.phone = status.value.profile?.phone || ''
    profileForm.value.contract_accepted = status.value.profile?.contract_status === 'signed'
    profileForm.value.contract_signature_data = status.value.profile?.contract_signature_data || ''
    await nextTick()
    renderSignature(profileForm.value.contract_signature_data)

    if (status.value.agent_status === 'approved') {
      const [dashData, linkData] = await Promise.all([
        agentAPI.getDashboard(),
        agentAPI.getLink()
      ])
      dashboard.value = dashData
      inviteLink.value = linkData.invite_url || `${window.location.origin}/register?invite=${linkData.invite_code}`
    }
  } catch (err) {
    console.error('Failed to load agent status:', err)
    appStore.showError('加载代理状态失败')
  } finally {
    loading.value = false
  }
}

async function handleSaveProfile() {
  if (!canSaveProfile.value) return
  savingProfile.value = true
  try {
    await agentAPI.saveProfile({
      real_name: profileForm.value.real_name.trim(),
      id_card_no: profileForm.value.id_card_no.trim(),
      phone: profileForm.value.phone.trim(),
      contract_accepted: profileForm.value.contract_accepted,
      contract_signature_data: profileForm.value.contract_signature_data
    })
    appStore.showSuccess('资料已保存')
    await refreshStatus()
  } catch (err) {
    console.error('Failed to save profile:', err)
    appStore.showError('保存资料失败')
  } finally {
    savingProfile.value = false
  }
}

async function handleCreateActivationOrder() {
  creatingOrder.value = true
  try {
    const order = await paymentAPI.createAgentActivationOrder()
    activationOrder.value = order
    showPaymentModal.value = true
    await nextTick()
    if (order.code_url && qrCanvas.value) {
      await QRCode.toCanvas(qrCanvas.value, order.code_url, { width: 220, margin: 1 })
    }
    startPolling(order.order_no)
  } catch (err) {
    console.error('Failed to create activation order:', err)
    appStore.showError('创建支付订单失败')
  } finally {
    creatingOrder.value = false
  }
}

async function handleApply() {
  applying.value = true
  try {
    await agentAPI.apply()
    appStore.showSuccess('代理申请已提交')
    await refreshStatus()
  } catch (err) {
    console.error('Failed to apply agent:', err)
    appStore.showError('提交申请失败')
  } finally {
    applying.value = false
  }
}

function startPolling(orderNo: string) {
  stopPolling()
  pollTimer.value = window.setInterval(async () => {
    try {
      const order = await paymentAPI.queryOrder(orderNo)
      if (order.status === 'paid') {
        stopPolling()
        showPaymentModal.value = false
        appStore.showSuccess('开通费支付成功')
        await refreshStatus()
      }
    } catch (err) {
      console.error('Failed to poll activation order:', err)
    }
  }, 3000)
}

function stopPolling() {
  if (pollTimer.value !== null) {
    window.clearInterval(pollTimer.value)
    pollTimer.value = null
  }
}

function closePaymentModal() {
  showPaymentModal.value = false
  stopPolling()
}

function copyLink() {
  navigator.clipboard.writeText(inviteLink.value)
  appStore.showSuccess('邀请链接已复制')
}

function formatMoney(value: number) {
  return `¥${Number(value || 0).toFixed(2)}`
}

function statusLabel(value: string) {
  switch (value) {
    case 'approved': return '已通过'
    case 'pending': return '待审核'
    case 'rejected': return '已驳回'
    default: return '未申请'
  }
}

function identityLabel(value?: string) {
  return value === 'submitted' ? '已提交' : '未提交'
}

function contractLabel(value?: string) {
  return value === 'signed' ? '已确认' : '未确认'
}

function weekdayLabel(value: number) {
  return ['一', '二', '三', '四', '五', '六', '日'][Math.max(1, value) - 1] || '五'
}

function statusBadgeClass(value: string) {
  if (value === 'approved') return 'badge-success'
  if (value === 'pending') return 'badge-warning'
  if (value === 'rejected') return 'badge-danger'
  return 'badge-gray'
}

function profileBadgeClass(value?: string) {
  return value === 'signed' || value === 'submitted' ? 'badge-success' : 'badge-gray'
}

function ensureSignatureCanvas() {
  const canvas = signatureCanvas.value
  if (!canvas) return null

  const rect = canvas.getBoundingClientRect()
  const ratio = window.devicePixelRatio || 1
  const width = Math.max(1, Math.floor(rect.width))
  const height = Math.max(1, Math.floor(rect.height))

  if (canvas.width !== Math.floor(width * ratio) || canvas.height !== Math.floor(height * ratio)) {
    canvas.width = Math.floor(width * ratio)
    canvas.height = Math.floor(height * ratio)
  }

  const ctx = canvas.getContext('2d')
  if (!ctx) return null

  ctx.setTransform(ratio, 0, 0, ratio, 0, 0)
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.lineWidth = 2
  ctx.strokeStyle = '#111827'
  signatureContext = ctx
  return ctx
}

function clearCanvasSurface() {
  const canvas = signatureCanvas.value
  const ctx = ensureSignatureCanvas()
  if (!canvas || !ctx) return
  const width = canvas.width / (window.devicePixelRatio || 1)
  const height = canvas.height / (window.devicePixelRatio || 1)
  ctx.clearRect(0, 0, width, height)
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, width, height)
  ctx.strokeStyle = '#111827'
}

function renderSignature(dataUrl?: string) {
  clearCanvasSurface()
  if (!dataUrl || !signatureCanvas.value || !signatureContext) return

  const canvas = signatureCanvas.value
  const ctx = signatureContext
  const img = new Image()
  img.onload = () => {
    const width = canvas.width / (window.devicePixelRatio || 1)
    const height = canvas.height / (window.devicePixelRatio || 1)
    ctx.clearRect(0, 0, width, height)
    ctx.fillStyle = '#ffffff'
    ctx.fillRect(0, 0, width, height)
    ctx.drawImage(img, 0, 0, width, height)
  }
  img.src = dataUrl
}

function pointerPosition(event: PointerEvent) {
  const canvas = signatureCanvas.value
  if (!canvas) return null
  const rect = canvas.getBoundingClientRect()
  return {
    x: event.clientX - rect.left,
    y: event.clientY - rect.top
  }
}

function startSignature(event: PointerEvent) {
  const canvas = signatureCanvas.value
  const ctx = ensureSignatureCanvas()
  const point = pointerPosition(event)
  if (!canvas || !ctx || !point) return

  signatureError.value = ''
  isSigning.value = true
  canvas.setPointerCapture(event.pointerId)
  ctx.beginPath()
  ctx.moveTo(point.x, point.y)
}

function moveSignature(event: PointerEvent) {
  if (!isSigning.value || !signatureContext) return
  const point = pointerPosition(event)
  if (!point) return
  signatureContext.lineTo(point.x, point.y)
  signatureContext.stroke()
  profileForm.value.contract_signature_data = signatureCanvas.value?.toDataURL('image/png') || ''
}

function endSignature(event: PointerEvent) {
  const canvas = signatureCanvas.value
  if (canvas?.hasPointerCapture(event.pointerId)) {
    canvas.releasePointerCapture(event.pointerId)
  }
  if (!isSigning.value) return
  isSigning.value = false
  profileForm.value.contract_signature_data = signatureCanvas.value?.toDataURL('image/png') || ''
}

function clearSignature() {
  profileForm.value.contract_signature_data = ''
  signatureError.value = ''
  clearCanvasSurface()
}
</script>
