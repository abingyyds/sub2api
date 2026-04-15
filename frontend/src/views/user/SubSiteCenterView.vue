<template>
  <AppLayout>
    <div class="mx-auto max-w-6xl space-y-6">
      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站中心</h1>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              分站是独立站点体系，不和代理中心混用。这里可以自助开通分站、设置模板、设置对外售卖价格。
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
                支付成功后自动创建分站，并带上你设置好的模板、开站规则和售价。
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
            <div>
              <label class="input-label">模板</label>
              <select v-model="createForm.theme_template" class="input mt-1 w-full">
                <option v-for="item in openInfo.theme_templates" :key="item.key" :value="item.key">
                  {{ item.label }}
                </option>
              </select>
            </div>
            <div>
              <label class="input-label">注册模式</label>
              <select v-model="createForm.registration_mode" class="input mt-1 w-full">
                <option value="open">开放注册</option>
                <option value="invite">邀请码注册</option>
                <option value="closed">关闭注册</option>
              </select>
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
            <div class="md:col-span-2">
              <label class="input-label">主题配置 JSON</label>
              <textarea v-model="createForm.theme_config" rows="4" class="input mt-1 w-full font-mono text-xs" placeholder='例如：{"accent":"aurora","hero_badge":"官方合作站"}'></textarea>
            </div>
            <div class="md:col-span-2">
              <label class="input-label">自定义配置 JSON</label>
              <textarea v-model="createForm.custom_config" rows="5" class="input mt-1 w-full font-mono text-xs" placeholder='例如：{"hero_title":"你的品牌标题","feature_tags":["稳定","高并发"]}'></textarea>
            </div>
          </div>

          <div class="mt-6 grid gap-4 rounded-2xl border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/60 md:grid-cols-3">
            <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
              <input v-model="createForm.enable_topup" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
              启用余额充值
            </label>
            <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
              <input v-model="createForm.allow_sub_site" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
              允许发展下级分站
            </label>
            <div>
              <label class="input-label">下级分站售价</label>
              <input v-model.number="createSubSitePriceYuan" type="number" min="0" step="1" class="input mt-1 w-full" placeholder="例如：399" />
            </div>
          </div>

          <div class="mt-6 space-y-4">
            <div>
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white">订阅套餐售价覆盖</h3>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">留空表示沿用平台默认售价。</p>
              <div class="mt-3 space-y-3">
                <div v-for="plan in baseSubscriptionPlans" :key="plan.key" class="grid gap-3 rounded-2xl border border-gray-100 px-4 py-3 dark:border-dark-700 md:grid-cols-[1fr_140px]">
                  <div>
                    <div class="font-medium text-gray-900 dark:text-white">{{ plan.name }}</div>
                    <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">默认价：￥{{ fenToYuan(plan.amount_fen) }}</div>
                  </div>
                  <input
                    :value="groupPriceInputs[plan.group_id] ?? ''"
                    type="number"
                    min="0"
                    step="1"
                    class="input w-full"
                    placeholder="自定义售价"
                    @input="updateGroupPrice(plan.group_id, $event)"
                  />
                </div>
              </div>
            </div>

            <div>
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white">充值套餐售价覆盖</h3>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">留空表示沿用平台默认售价。</p>
              <div class="mt-3 space-y-3">
                <div v-for="plan in baseRechargePlans" :key="plan.key" class="grid gap-3 rounded-2xl border border-gray-100 px-4 py-3 dark:border-dark-700 md:grid-cols-[1fr_140px]">
                  <div>
                    <div class="font-medium text-gray-900 dark:text-white">{{ plan.name }}</div>
                    <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">默认价：￥{{ fenToYuan(plan.pay_amount_fen) }} / 到账 ${{ plan.balance_amount.toFixed(2) }}</div>
                  </div>
                  <input
                    :value="rechargePriceInputs[plan.key] ?? ''"
                    type="number"
                    min="0"
                    step="1"
                    class="input w-full"
                    placeholder="自定义售价"
                    @input="updateRechargePrice(plan.key, $event)"
                  />
                </div>
              </div>
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
                    <div>分站模板：{{ site.theme_template || 'starter' }}，下级分站售价：￥{{ fenToYuan(site.sub_site_price_fen || 0) }}</div>
                  </div>
                </div>
                <div class="flex gap-2">
                  <button class="btn btn-secondary btn-sm" @click="openEdit(site)">编辑分站</button>
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
          <div>
            <label class="input-label">模板</label>
            <select v-model="editForm.theme_template" class="input mt-1 w-full">
              <option v-for="item in openInfo.theme_templates" :key="item.key" :value="item.key">
                {{ item.label }}
              </option>
            </select>
          </div>
          <div>
            <label class="input-label">注册模式</label>
            <select v-model="editForm.registration_mode" class="input mt-1 w-full">
              <option value="open">开放注册</option>
              <option value="invite">邀请码注册</option>
              <option value="closed">关闭注册</option>
            </select>
          </div>
          <div class="md:col-span-2">
            <label class="input-label">首页内容</label>
            <textarea v-model="editForm.home_content" rows="5" class="input mt-1 w-full"></textarea>
          </div>
          <div class="md:col-span-2">
            <label class="input-label">主题配置 JSON</label>
            <textarea v-model="editForm.theme_config" rows="4" class="input mt-1 w-full font-mono text-xs"></textarea>
          </div>
          <div class="md:col-span-2">
            <label class="input-label">自定义配置 JSON</label>
            <textarea v-model="editForm.custom_config" rows="5" class="input mt-1 w-full font-mono text-xs"></textarea>
          </div>
        </div>

        <div class="grid gap-4 rounded-2xl border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/60 md:grid-cols-3">
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
            <input v-model="editForm.enable_topup" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            启用余额充值
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
            <input v-model="editForm.allow_sub_site" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            允许发展下级分站
          </label>
          <div>
            <label class="input-label">下级分站售价</label>
            <input v-model.number="editSubSitePriceYuan" type="number" min="0" step="1" class="input mt-1 w-full" />
          </div>
        </div>

        <div>
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">订阅套餐售价覆盖</h3>
          <div class="mt-3 space-y-3">
            <div v-for="plan in baseSubscriptionPlans" :key="plan.key" class="grid gap-3 rounded-2xl border border-gray-100 px-4 py-3 dark:border-dark-700 md:grid-cols-[1fr_140px]">
              <div>
                <div class="font-medium text-gray-900 dark:text-white">{{ plan.name }}</div>
                <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">默认价：￥{{ fenToYuan(plan.amount_fen) }}</div>
              </div>
              <input
                :value="editGroupPriceInputs[plan.group_id] ?? ''"
                type="number"
                min="0"
                step="1"
                class="input w-full"
                placeholder="自定义售价"
                @input="updateEditGroupPrice(plan.group_id, $event)"
              />
            </div>
          </div>
        </div>

        <div>
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">充值套餐售价覆盖</h3>
          <div class="mt-3 space-y-3">
            <div v-for="plan in baseRechargePlans" :key="plan.key" class="grid gap-3 rounded-2xl border border-gray-100 px-4 py-3 dark:border-dark-700 md:grid-cols-[1fr_140px]">
              <div>
                <div class="font-medium text-gray-900 dark:text-white">{{ plan.name }}</div>
                <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">默认价：￥{{ fenToYuan(plan.pay_amount_fen) }} / 到账 ${{ plan.balance_amount.toFixed(2) }}</div>
              </div>
              <input
                :value="editRechargePriceInputs[plan.key] ?? ''"
                type="number"
                min="0"
                step="1"
                class="input w-full"
                placeholder="自定义售价"
                @input="updateEditRechargePrice(plan.key, $event)"
              />
            </div>
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
  </AppLayout>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { paymentAPI, type CreateOrderResponse, type CreateSubSiteActivationInput, type PaymentPlan, type RechargePlan } from '@/api/payment'
import { subSiteAPI, type OwnedSubSite, type SubSiteOpenInfo } from '@/api/subsite'
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
const baseSubscriptionPlans = ref<PaymentPlan[]>([])
const baseRechargePlans = ref<RechargePlan[]>([])

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
  theme_template: 'starter',
  theme_config: '',
  custom_config: '',
  registration_mode: 'open',
  enable_topup: true,
  allow_sub_site: false,
  sub_site_price_fen: 0,
  group_price_overrides: [],
  recharge_price_overrides: [],
})

const createSubSitePriceYuan = ref<number | null>(null)
const groupPriceInputs = reactive<Record<number, string>>({})
const rechargePriceInputs = reactive<Record<string, string>>({})

const editForm = ref<OwnedSubSite | null>(null)
const editSubSitePriceYuan = ref<number | null>(null)
const editGroupPriceInputs = reactive<Record<number, string>>({})
const editRechargePriceInputs = reactive<Record<string, string>>({})

onMounted(async () => {
  await Promise.all([loadOpenInfo(), loadBaseCatalog(), loadOwnedSites()])
})

onBeforeUnmount(() => {
  stopPolling()
})

async function loadOpenInfo() {
  try {
    openInfo.value = await subSiteAPI.getOpenInfo()
    if (!createForm.theme_template) {
      createForm.theme_template = openInfo.value.default_theme_template || 'starter'
    }
    if (!createForm.custom_config) {
      createForm.custom_config = openInfo.value.default_custom_config || ''
    }
  } catch (error: any) {
    appStore.showError(error?.message || '加载分站开通信息失败')
  }
}

async function loadBaseCatalog() {
  try {
    const [plans, rechargeInfo] = await Promise.all([
      paymentAPI.getPlans(),
      paymentAPI.getRechargeInfo(),
    ])
    baseSubscriptionPlans.value = plans.filter(plan => (plan.type || 'subscription') === 'subscription')
    baseRechargePlans.value = rechargeInfo.plans || []
  } catch (error: any) {
    appStore.showError(error?.message || '加载售价基准失败')
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

function buildGroupOverrides(source: Record<number, string>) {
  return Object.entries(source)
    .map(([groupId, value]) => {
      const priceFen = parseInputYuan(value)
      if (priceFen === null) return null
      return { group_id: Number(groupId), price_fen: priceFen }
    })
    .filter(Boolean) as Array<{ group_id: number; price_fen: number }>
}

function buildRechargeOverrides(source: Record<string, string>) {
  return Object.entries(source)
    .map(([planKey, value]) => {
      const priceFen = parseInputYuan(value)
      if (priceFen === null) return null
      return { plan_key: planKey, pay_amount_fen: priceFen }
    })
    .filter(Boolean) as Array<{ plan_key: string; pay_amount_fen: number }>
}

function updateGroupPrice(groupId: number, event: Event) {
  groupPriceInputs[groupId] = (event.target as HTMLInputElement).value
}

function updateRechargePrice(planKey: string, event: Event) {
  rechargePriceInputs[planKey] = (event.target as HTMLInputElement).value
}

function updateEditGroupPrice(groupId: number, event: Event) {
  editGroupPriceInputs[groupId] = (event.target as HTMLInputElement).value
}

function updateEditRechargePrice(planKey: string, event: Event) {
  editRechargePriceInputs[planKey] = (event.target as HTMLInputElement).value
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
  createForm.theme_template = openInfo.value.default_theme_template || 'starter'
  createForm.theme_config = ''
  createForm.custom_config = openInfo.value.default_custom_config || ''
  createForm.registration_mode = 'open'
  createForm.enable_topup = true
  createForm.allow_sub_site = false
  createForm.sub_site_price_fen = 0
  createSubSitePriceYuan.value = null
  Object.keys(groupPriceInputs).forEach((key) => delete groupPriceInputs[Number(key)])
  Object.keys(rechargePriceInputs).forEach((key) => delete rechargePriceInputs[key])
}

async function handleCreateActivationOrder() {
  const subSitePriceFen = parseInputYuan(createSubSitePriceYuan.value)
  createForm.sub_site_price_fen = subSitePriceFen ?? 0
  createForm.group_price_overrides = buildGroupOverrides(groupPriceInputs)
  createForm.recharge_price_overrides = buildRechargeOverrides(rechargePriceInputs)

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
  editSubSitePriceYuan.value = site.sub_site_price_fen ? site.sub_site_price_fen / 100 : null
  Object.keys(editGroupPriceInputs).forEach((key) => delete editGroupPriceInputs[Number(key)])
  Object.keys(editRechargePriceInputs).forEach((key) => delete editRechargePriceInputs[key])
  ;(site.group_price_overrides || []).forEach(item => {
    editGroupPriceInputs[item.group_id] = item.price_fen ? String(item.price_fen / 100) : ''
  })
  ;(site.recharge_price_overrides || []).forEach(item => {
    editRechargePriceInputs[item.plan_key] = item.pay_amount_fen ? String(item.pay_amount_fen / 100) : ''
  })
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
      theme_template: editForm.value.theme_template || 'starter',
      theme_config: editForm.value.theme_config || '',
      custom_config: editForm.value.custom_config || '',
      registration_mode: editForm.value.registration_mode || 'open',
      enable_topup: editForm.value.enable_topup,
      allow_sub_site: editForm.value.allow_sub_site,
      sub_site_price_fen: parseInputYuan(editSubSitePriceYuan.value) || 0,
      group_price_overrides: buildGroupOverrides(editGroupPriceInputs),
      recharge_price_overrides: buildRechargeOverrides(editRechargePriceInputs),
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
</script>
