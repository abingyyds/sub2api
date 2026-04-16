<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-6">
      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站中心</h1>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              分站默认继承主站的套餐、配置和界面。这里只需要维护分站品牌信息，并设置消耗倍率。
            </p>
          </div>
          <div class="flex flex-wrap gap-2">
            <span class="inline-flex rounded-full bg-primary-50 px-3 py-1 text-xs font-semibold text-primary-700 dark:bg-primary-900/20 dark:text-primary-300">
              {{ openInfo.scope === 'subsite' ? `当前挂靠：${openInfo.parent_sub_site_name || '当前分站'}` : '平台直属开通' }}
            </span>
            <span class="inline-flex rounded-full bg-gray-100 px-3 py-1 text-xs font-semibold text-gray-700 dark:bg-dark-800 dark:text-dark-300">
              ￥{{ fenToYuan(openInfo.price_fen) }}/{{ openInfo.validity_days }} 天
            </span>
          </div>
        </div>
      </section>

      <section class="grid gap-6 lg:grid-cols-[1.2fr_1fr]">
        <div class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between gap-3">
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">开通新分站</h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
                支付成功后自动创建分站，默认继承主站全部配置，只保留你的品牌信息和消耗倍率。
              </p>
            </div>
            <span
              class="inline-flex rounded-full px-3 py-1 text-xs font-semibold"
              :class="openInfo.enabled ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300' : 'bg-gray-100 text-gray-700 dark:bg-dark-800 dark:text-dark-300'"
            >
              {{ openInfo.enabled ? '可自助开通' : '当前不可开通' }}
            </span>
          </div>

          <div v-if="!openInfo.enabled" class="mt-4 rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-700 dark:border-amber-900/40 dark:bg-amber-950/20 dark:text-amber-300">
            当前入口没有开启分站自助开通。你可以切到允许开站的分站域名，或联系管理员开启平台自助开通。
          </div>

          <div class="mt-6 grid gap-4 md:grid-cols-2">
            <div>
              <label class="input-label">分站名称</label>
              <input v-model="createForm.name" class="input mt-1 w-full" placeholder="例如：cCoder 华东站" />
            </div>
            <div>
              <label class="input-label">Slug</label>
              <input v-model="createForm.slug" class="input mt-1 w-full" placeholder="例如：east-hub" />
            </div>
            <div class="md:col-span-2">
              <label class="input-label">自定义域名</label>
              <input v-model="createForm.custom_domain" class="input mt-1 w-full" placeholder="例如：east.ccoder.me" />
            </div>
            <div>
              <label class="input-label">Logo URL</label>
              <input v-model="createForm.site_logo" class="input mt-1 w-full" placeholder="https://..." />
            </div>
            <div>
              <label class="input-label">Favicon URL</label>
              <input v-model="createForm.site_favicon" class="input mt-1 w-full" placeholder="https://..." />
            </div>
            <div class="md:col-span-2">
              <label class="input-label">副标题</label>
              <input v-model="createForm.site_subtitle" class="input mt-1 w-full" placeholder="一句话介绍你的分站" />
            </div>
            <div class="md:col-span-2">
              <label class="input-label">公告</label>
              <input v-model="createForm.announcement" class="input mt-1 w-full" placeholder="分站首页顶部公告" />
            </div>
            <div>
              <label class="input-label">联系信息</label>
              <input v-model="createForm.contact_info" class="input mt-1 w-full" placeholder="Telegram / WeChat / Email" />
            </div>
            <div>
              <label class="input-label">文档地址</label>
              <input v-model="createForm.doc_url" class="input mt-1 w-full" placeholder="https://docs.example.com" />
            </div>
            <div class="md:col-span-2">
              <label class="input-label">首页内容</label>
              <textarea v-model="createForm.home_content" rows="5" class="input mt-1 w-full" placeholder="支持直接填 HTML，或填一个 https:// 页面地址用于 iframe 展示"></textarea>
            </div>
          </div>

          <div class="mt-6 grid gap-4 rounded-2xl border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/60 md:grid-cols-3">
            <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
              <input v-model="createForm.allow_sub_site" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
              允许发展下级分站
            </label>
            <div>
              <label class="input-label">下级分站售价</label>
              <input v-model.number="createSubSitePriceYuan" type="number" min="0" step="1" class="input mt-1 w-full" placeholder="例如：399" />
            </div>
            <div>
              <label class="input-label">消耗倍率</label>
              <input v-model.number="createForm.consume_rate_multiplier" type="number" min="1" step="0.1" class="input mt-1 w-full" placeholder="例如：1.5" />
            </div>
          </div>

          <div class="mt-6 flex justify-end">
            <button class="btn btn-primary" :disabled="creatingOrder || !openInfo.enabled" @click="handleCreateActivationOrder">
              {{ creatingOrder ? '创建订单中...' : `支付 ￥${fenToYuan(openInfo.price_fen)} 开通分站` }}
            </button>
          </div>
        </div>

        <div class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between gap-3">
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">我的分站</h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">已开通的分站可以继续修改模板、域名和售价。</p>
            </div>
            <button class="btn btn-secondary btn-sm" :disabled="loadingSites" @click="loadOwnedSites">
              {{ loadingSites ? '刷新中...' : '刷新' }}
            </button>
          </div>

          <div v-if="loadingSites" class="py-10 text-center text-sm text-gray-500 dark:text-dark-400">加载中...</div>
          <div v-else-if="ownedSites.length === 0" class="py-10 text-center text-sm text-gray-500 dark:text-dark-400">你还没有自己的分站</div>
          <div v-else class="mt-4 space-y-4">
            <div v-for="site in ownedSites" :key="site.id" class="rounded-2xl border border-gray-100 p-4 dark:border-dark-700">
              <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ site.name }}</h3>
                    <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium" :class="site.status === 'active' ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300' : site.status === 'pending' ? 'bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300' : 'bg-gray-100 text-gray-700 dark:bg-dark-800 dark:text-dark-300'">
                      {{ site.status }}
                    </span>
                    <span class="inline-flex rounded-full bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-700 dark:bg-dark-800 dark:text-dark-300">
                      L{{ site.level }}
                    </span>
                  </div>
                  <div class="mt-2 space-y-1 text-sm text-gray-500 dark:text-dark-400">
                    <div>slug：{{ site.slug }}</div>
                    <div v-if="site.custom_domain">域名：{{ site.custom_domain }}</div>
                    <div v-if="site.entry_url">入口：<a :href="site.entry_url" target="_blank" rel="noopener noreferrer" class="text-primary-600 hover:underline dark:text-primary-400">{{ site.entry_url }}</a></div>
                    <div>用户数：{{ site.user_count || 0 }}，下级分站：{{ site.child_site_count || 0 }}</div>
                    <div>消耗倍率：{{ Number(site.consume_rate_multiplier || 1).toFixed(2) }}x，下级分站售价：￥{{ fenToYuan(site.sub_site_price_fen || 0) }}</div>
                    <div class="font-medium text-gray-900 dark:text-white">
                      池余额：￥{{ fenToYuan(site.balance_fen || 0) }}
                      <span class="ml-2 text-xs text-gray-500 dark:text-dark-400">累充 ￥{{ fenToYuan(site.total_topup_fen || 0) }} / 累消 ￥{{ fenToYuan(site.total_consumed_fen || 0) }}</span>
                    </div>
                  </div>
                </div>
                <div class="flex flex-wrap gap-2">
                  <button class="btn btn-secondary btn-sm" @click="openEdit(site)">编辑分站</button>
                  <button class="btn btn-secondary btn-sm" @click="openOfflineTopup(site)">线下加余额</button>
                  <button class="btn btn-secondary btn-sm" @click="openOwnedLedger(site)">查看流水</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>

    <BaseDialog :show="showEditDialog" title="编辑分站" width="extra-wide" @close="closeEdit">
      <div v-if="editForm" class="space-y-5">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">分站名称</label>
            <input v-model="editForm.name" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Slug</label>
            <input v-model="editForm.slug" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">自定义域名</label>
            <input v-model="editForm.custom_domain" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Logo URL</label>
            <input v-model="editForm.site_logo" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Favicon URL</label>
            <input v-model="editForm.site_favicon" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">副标题</label>
            <input v-model="editForm.site_subtitle" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">公告</label>
            <input v-model="editForm.announcement" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">联系信息</label>
            <input v-model="editForm.contact_info" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">文档地址</label>
            <input v-model="editForm.doc_url" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">首页内容</label>
            <textarea v-model="editForm.home_content" rows="5" class="input mt-1 w-full"></textarea>
          </div>
        </div>

        <div class="grid gap-4 rounded-2xl border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/60 md:grid-cols-3">
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
            <input v-model="editForm.allow_sub_site" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            允许发展下级分站
          </label>
          <div>
            <label class="input-label">下级分站售价</label>
            <input v-model.number="editSubSitePriceYuan" type="number" min="0" step="1" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">消耗倍率</label>
            <input v-model.number="editForm.consume_rate_multiplier" type="number" min="1" step="0.1" class="input mt-1 w-full" />
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="closeEdit">取消</button>
          <button class="btn btn-primary" :disabled="savingSite" @click="handleSaveOwnedSite">
            {{ savingSite ? '保存中...' : '保存分站' }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <Teleport to="body">
      <div v-if="showPaymentModal" class="fixed inset-0 z-50 flex items-center justify-center px-4">
        <div class="absolute inset-0 bg-black/50" @click="closePaymentModal"></div>
        <div class="relative w-full max-w-md rounded-3xl bg-white p-6 shadow-2xl dark:bg-dark-900">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">扫码支付分站开通费</h3>
          <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">订单号：{{ activationOrder.order_no || '-' }}</p>
          <div class="mt-2 text-sm text-gray-500 dark:text-dark-400">金额：￥{{ activationOrder.amount_fen ? fenToYuan(activationOrder.amount_fen) : '0' }}</div>
          <div class="mt-6 flex justify-center">
            <canvas ref="qrCanvas" class="rounded-2xl border border-gray-200 bg-white p-3 dark:border-dark-700"></canvas>
          </div>
          <p class="mt-4 text-center text-sm text-gray-500 dark:text-dark-400">支付成功后会自动刷新你的分站列表</p>
          <div class="mt-6 flex justify-end gap-3">
            <button class="btn btn-secondary" @click="closePaymentModal">关闭</button>
            <button class="btn btn-primary" @click="loadOwnedSites">手动刷新</button>
          </div>
        </div>
      </div>
    </Teleport>

    <BaseDialog :show="showOfflineDialog" title="线下给用户加余额" @close="showOfflineDialog = false">
      <div v-if="offlineTarget" class="space-y-4">
        <div class="rounded-2xl border border-gray-200 px-4 py-3 text-sm dark:border-dark-700">
          <div class="font-medium text-gray-900 dark:text-white">{{ offlineTarget.name }}</div>
          <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">当前池余额：￥{{ fenToYuan(offlineTarget.balance_fen || 0) }}</div>
          <div class="text-xs text-gray-500 dark:text-dark-400">说明：给用户加余额的同时会从分站池按等额扣除。</div>
        </div>
        <div>
          <label class="input-label">用户 ID</label>
          <input v-model.number="offlineUserID" type="number" min="1" class="input mt-1 w-full" />
        </div>
        <div>
          <label class="input-label">加余额金额（元）</label>
          <input v-model.number="offlineAmountYuan" type="number" min="0" step="1" class="input mt-1 w-full" />
        </div>
        <div>
          <label class="input-label">备注（可选）</label>
          <input v-model="offlineNote" class="input mt-1 w-full" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="showOfflineDialog = false">取消</button>
          <button class="btn btn-primary" :disabled="offlineSubmitting" @click="handleOfflineTopup">{{ offlineSubmitting ? '提交中...' : '确认加余额' }}</button>
        </div>
      </template>
    </BaseDialog>

    <BaseDialog :show="showOwnedLedgerDialog" title="分站池流水" width="wide" @close="showOwnedLedgerDialog = false">
      <div v-if="ownedLedgerTarget" class="space-y-3">
        <div class="text-sm text-gray-600 dark:text-dark-300">
          {{ ownedLedgerTarget.name }} · 当前余额 ￥{{ fenToYuan(ownedLedgerTarget.balance_fen || 0) }}
        </div>
        <div class="max-h-96 overflow-y-auto">
          <table class="min-w-full divide-y divide-gray-200 text-sm dark:divide-dark-700">
            <thead class="bg-gray-50 dark:bg-dark-800/60">
              <tr>
                <th class="px-3 py-2 text-left text-xs font-semibold uppercase">时间</th>
                <th class="px-3 py-2 text-left text-xs font-semibold uppercase">类型</th>
                <th class="px-3 py-2 text-right text-xs font-semibold uppercase">变动</th>
                <th class="px-3 py-2 text-right text-xs font-semibold uppercase">变动后</th>
                <th class="px-3 py-2 text-left text-xs font-semibold uppercase">备注</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-for="row in ownedLedgerRows" :key="row.id">
                <td class="px-3 py-2 text-xs text-gray-500 dark:text-dark-400">{{ formatLedgerDate(row.created_at) }}</td>
                <td class="px-3 py-2 text-xs">{{ row.tx_type }}</td>
                <td class="px-3 py-2 text-right text-xs" :class="row.delta_fen >= 0 ? 'text-emerald-600' : 'text-red-600'">
                  {{ row.delta_fen >= 0 ? '+' : '' }}{{ fenToYuan(row.delta_fen) }}
                </td>
                <td class="px-3 py-2 text-right text-xs">{{ fenToYuan(row.balance_after_fen) }}</td>
                <td class="px-3 py-2 text-xs text-gray-500 dark:text-dark-400">{{ row.note || '-' }}</td>
              </tr>
              <tr v-if="ownedLedgerRows.length === 0">
                <td colspan="5" class="px-3 py-6 text-center text-xs text-gray-400">暂无流水</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { paymentAPI, type CreateOrderResponse, type CreateSubSiteActivationInput } from '@/api/payment'
import { subSiteAPI, type OwnedSubSite, type OwnedSubSiteLedgerEntry, type SubSiteOpenInfo } from '@/api/subsite'
import { useAppStore } from '@/stores'

const appStore = useAppStore()

const loadingSites = ref(false)
const creatingOrder = ref(false)
const savingSite = ref(false)
const showEditDialog = ref(false)
const showPaymentModal = ref(false)
const qrCanvas = ref<HTMLCanvasElement | null>(null)
const pollTimer = ref<number | null>(null)

const ownedSites = ref<OwnedSubSite[]>([])

const openInfo = ref<SubSiteOpenInfo>({
  enabled: false,
  scope: 'platform',
  level: 1,
  max_level: 2,
  price_fen: 0,
  validity_days: 365,
  currency: 'CNY',
  allow_custom_domain: true,
  default_theme_template: 'starter',
  default_custom_config: '',
  theme_templates: [],
})

const activationOrder = ref<CreateOrderResponse>({
  order_no: '',
  code_url: null,
  amount_fen: 0,
  discount_amount: 0,
  expired_at: '',
})

const createForm = reactive<CreateSubSiteActivationInput>({
  name: '',
  slug: '',
  custom_domain: '',
  site_logo: '',
  site_favicon: '',
  site_subtitle: '',
  announcement: '',
  contact_info: '',
  doc_url: '',
  home_content: '',
  allow_sub_site: false,
  sub_site_price_fen: 0,
  consume_rate_multiplier: 1,
})

const createSubSitePriceYuan = ref<number | null>(null)

const editForm = ref<OwnedSubSite | null>(null)
const editSubSitePriceYuan = ref<number | null>(null)

onMounted(async () => {
  await Promise.all([loadOpenInfo(), loadOwnedSites()])
})

onBeforeUnmount(() => {
  stopPolling()
})

async function loadOpenInfo() {
  try {
    openInfo.value = await subSiteAPI.getOpenInfo()
  } catch (error: any) {
    appStore.showError(error?.message || '加载分站开通信息失败')
  }
}

async function loadOwnedSites() {
  loadingSites.value = true
  try {
    ownedSites.value = await subSiteAPI.listOwnedSites()
  } catch (error: any) {
    appStore.showError(error?.message || '加载我的分站失败')
  } finally {
    loadingSites.value = false
  }
}

function fenToYuan(value: number) {
  return (Number(value || 0) / 100).toFixed(value % 100 === 0 ? 0 : 2)
}

function parseInputYuan(value: string | number | null | undefined): number | null {
  if (value === '' || value === null || value === undefined) {
    return null
  }
  const num = Number(value)
  if (!Number.isFinite(num) || num < 0) {
    return null
  }
  return Math.round(num * 100)
}

function resetCreateForm() {
  createForm.name = ''
  createForm.slug = ''
  createForm.custom_domain = ''
  createForm.site_logo = ''
  createForm.site_favicon = ''
  createForm.site_subtitle = ''
  createForm.announcement = ''
  createForm.contact_info = ''
  createForm.doc_url = ''
  createForm.home_content = ''
  createForm.allow_sub_site = false
  createForm.sub_site_price_fen = 0
  createForm.consume_rate_multiplier = 1
  createSubSitePriceYuan.value = null
}

async function handleCreateActivationOrder() {
  const subSitePriceFen = parseInputYuan(createSubSitePriceYuan.value)
  createForm.sub_site_price_fen = subSitePriceFen ?? 0

  creatingOrder.value = true
  try {
    const order = await paymentAPI.createSubSiteActivationOrder(createForm)
    activationOrder.value = order
    showPaymentModal.value = true
    await nextTick()
    if (order.code_url && qrCanvas.value) {
      await QRCode.toCanvas(qrCanvas.value, order.code_url, { width: 220, margin: 1 })
    }
    startPolling(order.order_no)
    appStore.showSuccess('分站开通订单已创建，请扫码支付')
  } catch (error: any) {
    appStore.showError(error?.message || '创建分站开通订单失败')
  } finally {
    creatingOrder.value = false
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
        resetCreateForm()
        await loadOwnedSites()
        appStore.showSuccess('分站开通成功')
      }
    } catch {
      // ignore poll errors
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

function openEdit(site: OwnedSubSite) {
  editForm.value = JSON.parse(JSON.stringify(site)) as OwnedSubSite
  editForm.value.consume_rate_multiplier = Number(site.consume_rate_multiplier || 1)
  editSubSitePriceYuan.value = site.sub_site_price_fen ? site.sub_site_price_fen / 100 : null
  showEditDialog.value = true
}

function closeEdit() {
  showEditDialog.value = false
  editForm.value = null
}

async function handleSaveOwnedSite() {
  if (!editForm.value) return
  savingSite.value = true
  try {
    const payload: CreateSubSiteActivationInput = {
      name: editForm.value.name,
      slug: editForm.value.slug,
      custom_domain: editForm.value.custom_domain || '',
      site_logo: editForm.value.site_logo || '',
      site_favicon: editForm.value.site_favicon || '',
      site_subtitle: editForm.value.site_subtitle || '',
      announcement: editForm.value.announcement || '',
      contact_info: editForm.value.contact_info || '',
      doc_url: editForm.value.doc_url || '',
      home_content: editForm.value.home_content || '',
      allow_sub_site: editForm.value.allow_sub_site,
      sub_site_price_fen: parseInputYuan(editSubSitePriceYuan.value) || 0,
      consume_rate_multiplier: Number(editForm.value.consume_rate_multiplier || 1),
    }
    await subSiteAPI.updateOwnedSite(editForm.value.id, payload)
    appStore.showSuccess('分站已更新')
    closeEdit()
    await loadOwnedSites()
  } catch (error: any) {
    appStore.showError(error?.message || '保存分站失败')
  } finally {
    savingSite.value = false
  }
}

const showOfflineDialog = ref(false)
const offlineTarget = ref<OwnedSubSite | null>(null)
const offlineUserID = ref<number>(0)
const offlineAmountYuan = ref<number>(0)
const offlineNote = ref('')
const offlineSubmitting = ref(false)

function openOfflineTopup(site: OwnedSubSite) {
  offlineTarget.value = site
  offlineUserID.value = 0
  offlineAmountYuan.value = 0
  offlineNote.value = ''
  showOfflineDialog.value = true
}

async function handleOfflineTopup() {
  if (!offlineTarget.value) return
  const amountFen = Math.round((offlineAmountYuan.value || 0) * 100)
  if (!offlineUserID.value || offlineUserID.value <= 0) {
    appStore.showError('请输入有效的用户 ID')
    return
  }
  if (amountFen <= 0) {
    appStore.showError('请输入有效金额')
    return
  }
  offlineSubmitting.value = true
  try {
    await subSiteAPI.offlineTopupUser(offlineTarget.value.id, {
      user_id: offlineUserID.value,
      amount_fen: amountFen,
      note: offlineNote.value
    })
    appStore.showSuccess('加余额成功')
    showOfflineDialog.value = false
    offlineTarget.value = null
    await loadOwnedSites()
  } catch (error: any) {
    appStore.showError(error?.message || '加余额失败')
  } finally {
    offlineSubmitting.value = false
  }
}

const showOwnedLedgerDialog = ref(false)
const ownedLedgerTarget = ref<OwnedSubSite | null>(null)
const ownedLedgerRows = ref<OwnedSubSiteLedgerEntry[]>([])

async function openOwnedLedger(site: OwnedSubSite) {
  ownedLedgerTarget.value = site
  ownedLedgerRows.value = []
  showOwnedLedgerDialog.value = true
  try {
    const res = await subSiteAPI.listOwnedLedger(site.id, 1, 50)
    ownedLedgerRows.value = res.items || []
  } catch (error: any) {
    appStore.showError(error?.message || '加载流水失败')
  }
}

function formatLedgerDate(v: string) {
  if (!v) return '-'
  try {
    return new Date(v).toLocaleString('zh-CN', { hour12: false })
  } catch {
    return v
  }
}
</script>
