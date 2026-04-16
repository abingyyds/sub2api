<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-6">
      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站管理</h1>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              分站是独立站点体系，不属于代理系统。这里可以配置平台自助开通、模板、层级和分站售卖价格。
            </p>
          </div>
          <button class="btn btn-primary" @click="openCreate">新建分站</button>
        </div>
      </section>

      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">平台自助开通配置</h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">主站用户从平台直接购买分站时，会使用这里的开通价格和默认模板。</p>
          </div>
          <button class="btn btn-secondary btn-sm" :disabled="savingPlatform" @click="handleSavePlatform">
            {{ savingPlatform ? '保存中...' : '保存平台配置' }}
          </button>
        </div>
        <div class="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          <label class="flex items-center gap-3 rounded-2xl border border-gray-100 px-4 py-3 text-sm text-gray-700 dark:border-dark-700 dark:text-dark-200">
            <input v-model="platformForm.entry_enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            展示分站系统入口
          </label>
          <label class="flex items-center gap-3 rounded-2xl border border-gray-100 px-4 py-3 text-sm text-gray-700 dark:border-dark-700 dark:text-dark-200">
            <input v-model="platformForm.enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            开启平台自助开通
          </label>
          <div>
            <label class="input-label">平台开通价</label>
            <input v-model.number="platformPriceYuan" type="number" min="0" step="1" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">有效期（天）</label>
            <input v-model.number="platformForm.validity_days" type="number" min="1" step="1" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">最大分站层级</label>
            <input v-model.number="platformForm.max_level" type="number" min="1" max="10" step="1" class="input mt-1 w-full" />
            <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">上限 10 层；修改后已存在的超层分站仍保留，但无法继续向下开新站。</p>
          </div>
          <div>
            <label class="input-label">默认模板</label>
            <select v-model="platformForm.default_theme_template" class="input mt-1 w-full">
              <option v-for="item in platformForm.theme_templates" :key="item.key" :value="item.key">{{ item.label }}</option>
            </select>
          </div>
        </div>
      </section>

      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="grid gap-3 md:grid-cols-[1fr_180px_auto]">
          <input v-model="filters.search" class="input" placeholder="搜索站点名 / slug / 域名 / owner 邮箱" @keyup.enter="loadData(1)" />
          <select v-model="filters.status" class="input" @change="loadData(1)">
            <option value="">全部状态</option>
            <option value="pending">待激活</option>
            <option value="active">启用</option>
            <option value="disabled">停用</option>
          </select>
          <button class="btn btn-secondary" @click="loadData(1)">刷新</button>
        </div>
      </section>

      <section class="overflow-hidden rounded-3xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-dark-700">
            <thead class="bg-gray-50 dark:bg-dark-800/60">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">站点</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">Owner</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">层级 / 入口</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">用户 / 下级</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">下级分站售价</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">余额池（元）</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">状态</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-if="loading">
                <td colspan="8" class="px-4 py-10 text-center text-sm text-gray-500 dark:text-dark-400">加载中...</td>
              </tr>
              <tr v-else-if="items.length === 0">
                <td colspan="8" class="px-4 py-10 text-center text-sm text-gray-500 dark:text-dark-400">还没有分站</td>
              </tr>
              <tr v-for="item in items" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50">
                <td class="px-4 py-4 align-top">
                  <div class="font-medium text-gray-900 dark:text-white">{{ item.name }}</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">slug: {{ item.slug }}</div>
                  <div v-if="item.parent_sub_site_name" class="mt-1 text-xs text-gray-500 dark:text-dark-400">上级：{{ item.parent_sub_site_name }}</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">模板：{{ item.theme_template || 'starter' }}</div>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <div>ID: {{ item.owner_user_id }}</div>
                  <div class="mt-1 break-all text-xs text-gray-500 dark:text-dark-400">{{ item.owner_email || '-' }}</div>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <div>L{{ item.level }}</div>
                  <a v-if="item.entry_url" :href="item.entry_url" target="_blank" rel="noopener noreferrer" class="mt-1 block break-all text-primary-600 hover:underline dark:text-primary-400">
                    {{ item.entry_url }}
                  </a>
                  <span v-else class="mt-1 block text-xs text-gray-400 dark:text-dark-500">未配置入口域名</span>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <div>{{ item.user_count || 0 }} 个用户</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ item.child_site_count || 0 }} 个下级分站</div>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <div>￥{{ fenToYuan(item.sub_site_price_fen || 0) }}</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ item.allow_sub_site ? '允许发展下级' : '未开放下级分站' }}</div>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <div class="font-medium text-gray-900 dark:text-white">￥{{ fenToYuan(item.balance_fen || 0) }}</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">累充 ￥{{ fenToYuan(item.total_topup_fen || 0) }}</div>
                  <div class="text-xs text-gray-500 dark:text-dark-400">累消 ￥{{ fenToYuan(item.total_consumed_fen || 0) }}</div>
                </td>
                <td class="px-4 py-4 align-top">
                  <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium"
                    :class="item.status === 'active'
                      ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
                      : item.status === 'pending'
                        ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
                        : 'bg-gray-200 text-gray-700 dark:bg-dark-700 dark:text-dark-300'">
                    {{ item.status }}
                  </span>
                </td>
                <td class="px-4 py-4 align-top text-right">
                  <div class="flex justify-end gap-2">
                    <button class="btn btn-secondary btn-sm" @click="openTopup(item)">充值</button>
                    <button class="btn btn-secondary btn-sm" @click="openLedger(item)">流水</button>
                    <button class="btn btn-secondary btn-sm" @click="openEdit(item)">编辑</button>
                    <button class="btn btn-secondary btn-sm text-red-600 dark:text-red-300" @click="askDelete(item)">删除</button>
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

    <BaseDialog :show="showDialog" :title="editingId ? '编辑分站' : '新建分站'" width="extra-wide" @close="closeDialog">
      <div class="space-y-5">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">站点名称</label>
            <input v-model="form.name" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Owner 用户 ID</label>
            <input v-model.number="form.owner_user_id" type="number" min="1" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">上级分站 ID</label>
            <input v-model.number="form.parent_sub_site_id" type="number" min="0" class="input mt-1 w-full" placeholder="顶级分站留空" />
          </div>
          <div>
            <label class="input-label">状态</label>
            <select v-model="form.status" class="input mt-1 w-full">
              <option value="pending">待激活</option>
              <option value="active">启用</option>
              <option value="disabled">停用</option>
            </select>
          </div>
          <div>
            <label class="input-label">Slug</label>
            <input v-model="form.slug" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">模板</label>
            <select v-model="form.theme_template" class="input mt-1 w-full">
              <option v-for="item in platformForm.theme_templates" :key="item.key" :value="item.key">{{ item.label }}</option>
            </select>
          </div>
          <div class="md:col-span-2">
            <label class="input-label">自定义域名</label>
            <input v-model="form.custom_domain" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Logo URL</label>
            <input v-model="form.site_logo" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Favicon URL</label>
            <input v-model="form.site_favicon" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">副标题</label>
            <input v-model="form.site_subtitle" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">公告</label>
            <input v-model="form.announcement" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">联系信息</label>
            <input v-model="form.contact_info" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">文档地址</label>
            <input v-model="form.doc_url" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">注册模式</label>
            <select v-model="form.registration_mode" class="input mt-1 w-full">
              <option value="open">开放注册</option>
              <option value="invite">邀请码注册</option>
              <option value="closed">关闭注册</option>
            </select>
          </div>
          <div>
            <label class="input-label">到期时间</label>
            <input v-model="form.subscription_expired_at" type="datetime-local" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">首页内容</label>
            <textarea v-model="form.home_content" rows="5" class="input mt-1 w-full"></textarea>
          </div>
        </div>

        <div class="grid gap-4 rounded-2xl border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/60 md:grid-cols-3">
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
            <input v-model="form.enable_topup" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            启用余额充值
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
            <input v-model="form.allow_sub_site" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            允许发展下级分站
          </label>
          <div>
            <label class="input-label">下级分站售价</label>
            <input v-model.number="formSubSitePriceYuan" type="number" min="0" step="1" class="input mt-1 w-full" />
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="closeDialog">取消</button>
          <button class="btn btn-primary" :disabled="saving" @click="handleSave">{{ saving ? '保存中...' : '保存' }}</button>
        </div>
      </template>
    </BaseDialog>

    <ConfirmDialog
      :show="showDeleteDialog"
      title="删除分站"
      :message="deleteTarget ? `确认删除 ${deleteTarget.name} 吗？` : '确认删除该分站吗？'"
      confirm-text="删除"
      cancel-text="取消"
      :danger="true"
      @confirm="handleDelete"
      @cancel="showDeleteDialog = false"
    />

    <BaseDialog :show="showTopupDialog" title="给分站池充值" @close="showTopupDialog = false">
      <div v-if="topupTarget" class="space-y-4">
        <div class="rounded-2xl border border-gray-200 px-4 py-3 text-sm dark:border-dark-700">
          <div class="font-medium text-gray-900 dark:text-white">{{ topupTarget.name }}</div>
          <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">当前池余额：￥{{ fenToYuan(topupTarget.balance_fen || 0) }}</div>
        </div>
        <div>
          <label class="input-label">充值金额（元）</label>
          <input v-model.number="topupAmountYuan" type="number" min="0" step="1" class="input mt-1 w-full" />
        </div>
        <div>
          <label class="input-label">备注（可选）</label>
          <input v-model="topupNote" class="input mt-1 w-full" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="showTopupDialog = false">取消</button>
          <button class="btn btn-primary" :disabled="topupSubmitting" @click="handleTopup">{{ topupSubmitting ? '提交中...' : '确认充值' }}</button>
        </div>
      </template>
    </BaseDialog>

    <BaseDialog :show="showLedgerDialog" title="分站池流水" width="wide" @close="showLedgerDialog = false">
      <div v-if="ledgerTarget" class="space-y-3">
        <div class="text-sm text-gray-600 dark:text-dark-300">
          {{ ledgerTarget.name }} · 当前余额 ￥{{ fenToYuan(ledgerTarget.balance_fen || 0) }}
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
              <tr v-for="row in ledgerRows" :key="row.id">
                <td class="px-3 py-2 text-xs text-gray-500 dark:text-dark-400">{{ formatDate(row.created_at) }}</td>
                <td class="px-3 py-2 text-xs">{{ row.tx_type }}</td>
                <td class="px-3 py-2 text-right text-xs" :class="row.delta_fen >= 0 ? 'text-emerald-600' : 'text-red-600'">
                  {{ row.delta_fen >= 0 ? '+' : '' }}{{ fenToYuan(row.delta_fen) }}
                </td>
                <td class="px-3 py-2 text-right text-xs">{{ fenToYuan(row.balance_after_fen) }}</td>
                <td class="px-3 py-2 text-xs text-gray-500 dark:text-dark-400">{{ row.note || '-' }}</td>
              </tr>
              <tr v-if="ledgerRows.length === 0">
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
import { onMounted, reactive, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { subSitesAPI, type AdminSubSite, type PlatformSubSiteConfig, type SaveSubSiteRequest, type SubSiteLedgerEntry } from '@/api/admin/subsites'
import { useAppStore } from '@/stores'

const appStore = useAppStore()

const loading = ref(false)
const saving = ref(false)
const loadingPlatform = ref(false)
const savingPlatform = ref(false)
const items = ref<AdminSubSite[]>([])
const page = ref(1)
const pages = ref(1)

const filters = reactive({
  search: '',
  status: ''
})

const platformForm = reactive<PlatformSubSiteConfig>({
  entry_enabled: false,
  enabled: true,
  activation_price_fen: 38800,
  validity_days: 365,
  max_level: 2,
  default_theme_template: 'starter',
  theme_templates: []
})
const platformPriceYuan = ref<number | null>(388)

const showDialog = ref(false)
const showDeleteDialog = ref(false)
const editingId = ref<number | null>(null)
const deleteTarget = ref<AdminSubSite | null>(null)
const formSubSitePriceYuan = ref<number | null>(null)

const form = reactive<SaveSubSiteRequest>({
  owner_user_id: 0,
  parent_sub_site_id: undefined,
  name: '',
  slug: '',
  custom_domain: '',
  status: 'active',
  site_logo: '',
  site_favicon: '',
  site_subtitle: '',
  announcement: '',
  contact_info: '',
  doc_url: '',
  home_content: '',
  theme_template: 'starter',
  registration_mode: 'open',
  enable_topup: true,
  allow_sub_site: false,
  sub_site_price_fen: 0,
  subscription_expired_at: null,
})

onMounted(async () => {
  await Promise.all([loadPlatformConfig(), loadData(1)])
})

function fenToYuan(value: number) {
  return (Number(value || 0) / 100).toFixed(value % 100 === 0 ? 0 : 2)
}

function parseYuanToFen(value: string | number | null | undefined): number | null {
  if (value === '' || value === null || value === undefined) return null
  const num = Number(value)
  if (!Number.isFinite(num) || num < 0) return null
  return Math.round(num * 100)
}

function normalizePositiveInt(value: number | string | null | undefined): number | undefined {
  if (value === '' || value === null || value === undefined) return undefined
  const num = Number(value)
  if (!Number.isFinite(num) || num <= 0) return undefined
  return Math.trunc(num)
}

function formatDateTimeLocal(value?: string | null): string | null {
  if (!value) return null
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return null
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

function parseDateTimeLocal(value?: string | null): string | null {
  if (!value) return null
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return null
  return date.toISOString()
}

async function loadPlatformConfig() {
  loadingPlatform.value = true
  try {
    const config = await subSitesAPI.getPlatformConfig()
    Object.assign(platformForm, config)
    platformPriceYuan.value = config.activation_price_fen / 100
  } catch (error: any) {
    appStore.showError(error?.message || '加载平台分站配置失败')
  } finally {
    loadingPlatform.value = false
  }
}

async function handleSavePlatform() {
  savingPlatform.value = true
  try {
    platformForm.activation_price_fen = parseYuanToFen(platformPriceYuan.value) || 0
    const next = await subSitesAPI.updatePlatformConfig(platformForm)
    Object.assign(platformForm, next)
    platformPriceYuan.value = next.activation_price_fen / 100
    await appStore.fetchPublicSettings(true)
    appStore.showSuccess('平台分站配置已保存')
  } catch (error: any) {
    appStore.showError(error?.message || '保存平台分站配置失败')
  } finally {
    savingPlatform.value = false
  }
}

async function loadData(targetPage = page.value) {
  loading.value = true
  try {
    const res = await subSitesAPI.list(targetPage, 20, filters)
    items.value = res.items || []
    page.value = res.page || targetPage
    pages.value = res.pages || 1
  } catch (error: any) {
    appStore.showError(error?.message || '加载分站失败')
  } finally {
    loading.value = false
  }
}

function resetForm() {
  editingId.value = null
  form.owner_user_id = 0
  form.parent_sub_site_id = undefined
  form.name = ''
  form.slug = ''
  form.custom_domain = ''
  form.status = 'active'
  form.site_logo = ''
  form.site_favicon = ''
  form.site_subtitle = ''
  form.announcement = ''
  form.contact_info = ''
  form.doc_url = ''
  form.home_content = ''
  form.theme_template = platformForm.default_theme_template || 'starter'
  form.registration_mode = 'open'
  form.enable_topup = true
  form.allow_sub_site = false
  form.sub_site_price_fen = 0
  form.subscription_expired_at = null
  formSubSitePriceYuan.value = null
}

function openCreate() {
  resetForm()
  showDialog.value = true
}

function openEdit(item: AdminSubSite) {
  editingId.value = item.id
  form.owner_user_id = item.owner_user_id
  form.parent_sub_site_id = item.parent_sub_site_id
  form.name = item.name
  form.slug = item.slug
  form.custom_domain = item.custom_domain || ''
  form.status = item.status
  form.site_logo = item.site_logo || ''
  form.site_favicon = item.site_favicon || ''
  form.site_subtitle = item.site_subtitle || ''
  form.announcement = item.announcement || ''
  form.contact_info = item.contact_info || ''
  form.doc_url = item.doc_url || ''
  form.home_content = item.home_content || ''
  form.theme_template = item.theme_template || 'starter'
  form.registration_mode = item.registration_mode || 'open'
  form.enable_topup = item.enable_topup
  form.allow_sub_site = item.allow_sub_site
  form.sub_site_price_fen = item.sub_site_price_fen || 0
  form.subscription_expired_at = formatDateTimeLocal(item.subscription_expired_at)
  formSubSitePriceYuan.value = item.sub_site_price_fen ? item.sub_site_price_fen / 100 : null
  showDialog.value = true
}

function closeDialog() {
  showDialog.value = false
}

async function handleSave() {
  saving.value = true
  try {
    form.sub_site_price_fen = parseYuanToFen(formSubSitePriceYuan.value) || 0
    const payload: SaveSubSiteRequest = {
      owner_user_id: form.owner_user_id,
      parent_sub_site_id: normalizePositiveInt(form.parent_sub_site_id),
      name: form.name,
      slug: form.slug,
      custom_domain: form.custom_domain || '',
      status: form.status,
      site_logo: form.site_logo || '',
      site_favicon: form.site_favicon || '',
      site_subtitle: form.site_subtitle || '',
      announcement: form.announcement || '',
      contact_info: form.contact_info || '',
      doc_url: form.doc_url || '',
      home_content: form.home_content || '',
      theme_template: form.theme_template || 'starter',
      registration_mode: form.registration_mode || 'open',
      enable_topup: form.enable_topup,
      allow_sub_site: form.allow_sub_site,
      sub_site_price_fen: form.sub_site_price_fen,
      subscription_expired_at: parseDateTimeLocal(form.subscription_expired_at),
    }
    if (editingId.value) {
      await subSitesAPI.update(editingId.value, payload)
      appStore.showSuccess('分站已更新')
    } else {
      await subSitesAPI.create(payload)
      appStore.showSuccess('分站已创建')
    }
    showDialog.value = false
    await loadData(editingId.value ? page.value : 1)
  } catch (error: any) {
    appStore.showError(error?.message || '保存分站失败')
  } finally {
    saving.value = false
  }
}

function askDelete(item: AdminSubSite) {
  deleteTarget.value = item
  showDeleteDialog.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  try {
    await subSitesAPI.remove(deleteTarget.value.id)
    appStore.showSuccess('分站已删除')
    showDeleteDialog.value = false
    deleteTarget.value = null
    await loadData(page.value)
  } catch (error: any) {
    appStore.showError(error?.message || '删除分站失败')
  }
}

const showTopupDialog = ref(false)
const topupTarget = ref<AdminSubSite | null>(null)
const topupAmountYuan = ref<number>(0)
const topupNote = ref('')
const topupSubmitting = ref(false)

function openTopup(item: AdminSubSite) {
  topupTarget.value = item
  topupAmountYuan.value = 0
  topupNote.value = ''
  showTopupDialog.value = true
}

async function handleTopup() {
  if (!topupTarget.value) return
  const amountFen = Math.round((topupAmountYuan.value || 0) * 100)
  if (amountFen <= 0) {
    appStore.showError('请输入有效金额')
    return
  }
  topupSubmitting.value = true
  try {
    await subSitesAPI.topupPool(topupTarget.value.id, amountFen, topupNote.value)
    appStore.showSuccess('充值成功')
    showTopupDialog.value = false
    topupTarget.value = null
    await loadData(page.value)
  } catch (error: any) {
    appStore.showError(error?.message || '充值失败')
  } finally {
    topupSubmitting.value = false
  }
}

const showLedgerDialog = ref(false)
const ledgerTarget = ref<AdminSubSite | null>(null)
const ledgerRows = ref<SubSiteLedgerEntry[]>([])

async function openLedger(item: AdminSubSite) {
  ledgerTarget.value = item
  ledgerRows.value = []
  showLedgerDialog.value = true
  try {
    const res = await subSitesAPI.listLedger(item.id, 1, 50)
    ledgerRows.value = res.items || []
  } catch (error: any) {
    appStore.showError(error?.message || '加载流水失败')
  }
}

function formatDate(v: string) {
  if (!v) return '-'
  try {
    return new Date(v).toLocaleString('zh-CN', { hour12: false })
  } catch {
    return v
  }
}
</script>
