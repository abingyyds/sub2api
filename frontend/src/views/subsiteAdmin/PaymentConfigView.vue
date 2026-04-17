<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl space-y-6">
      <div class="flex items-center gap-3">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">自有收款账号</h1>
        <span v-if="site?.mode" class="inline-flex rounded-full px-3 py-1 text-xs font-medium"
          :class="site.mode === 'pool'
            ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
            : 'bg-gray-200 text-gray-700 dark:bg-dark-700 dark:text-dark-300'">
          {{ site.mode === 'pool' ? '资金池模式' : '倍率分成模式' }}
        </span>
      </div>

      <div v-if="site?.mode !== 'pool'" class="rounded-2xl border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800 dark:border-amber-800/50 dark:bg-amber-900/20 dark:text-amber-200">
        自有收款账号仅对资金池模式生效：用户通过分站域名充值时，款项会直接入账到您配置的商户号，同时从分站池扣等额作为"自动进货"。
        当前分站为倍率分成模式，无需配置。
      </div>

      <template v-else-if="site">
        <div class="rounded-2xl border border-blue-200 bg-blue-50 p-4 text-sm text-blue-800 dark:border-blue-800/50 dark:bg-blue-900/20 dark:text-blue-200">
          <p class="font-medium">工作原理</p>
          <ul class="mt-2 list-disc space-y-1 pl-4">
            <li>用户通过分站域名发起余额充值时，订单使用下方配置的支付商户收款</li>
            <li>付款成功后：用户余额增加，分站池自动扣等额（<code>auto_restock</code> 流水）</li>
            <li>未填写的渠道会回退到平台主站配置</li>
            <li>分站激活 / 下级分站购买订单始终走主站收款，不受此处影响</li>
          </ul>
        </div>

        <section class="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">微信 Native 支付</h2>
            <label class="flex items-center gap-2 text-sm">
              <input v-model="config.wechat.enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
              启用
            </label>
          </div>
          <div class="mt-4 grid gap-4 md:grid-cols-2">
            <div>
              <label class="input-label">AppID</label>
              <input v-model="config.wechat.app_id" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">商户号 MchID</label>
              <input v-model="config.wechat.mch_id" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">APIv3 Key</label>
              <input v-model="config.wechat.apiv3_key" type="password" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">商户证书序列号</label>
              <input v-model="config.wechat.mch_serial_no" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">平台公钥 ID</label>
              <input v-model="config.wechat.public_key_id" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">Notify URL</label>
              <input v-model="config.wechat.notify_url" class="input mt-1 w-full" />
            </div>
            <div class="md:col-span-2">
              <label class="input-label">商户私钥（PEM）</label>
              <textarea v-model="config.wechat.private_key" rows="4" class="input mt-1 w-full font-mono text-xs"></textarea>
            </div>
            <div class="md:col-span-2">
              <label class="input-label">平台公钥（PEM）</label>
              <textarea v-model="config.wechat.public_key" rows="4" class="input mt-1 w-full font-mono text-xs"></textarea>
            </div>
          </div>
        </section>

        <section class="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">支付宝 Native 支付</h2>
            <label class="flex items-center gap-2 text-sm">
              <input v-model="config.alipay.enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
              启用
            </label>
          </div>
          <div class="mt-4 grid gap-4 md:grid-cols-2">
            <div>
              <label class="input-label">AppID</label>
              <input v-model="config.alipay.app_id" class="input mt-1 w-full" />
            </div>
            <div class="flex items-center gap-2 md:col-span-1 md:pt-6">
              <input v-model="config.alipay.is_production" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
              <label class="text-sm text-gray-700 dark:text-dark-200">生产环境</label>
            </div>
            <div>
              <label class="input-label">Notify URL</label>
              <input v-model="config.alipay.notify_url" class="input mt-1 w-full" />
            </div>
            <div class="md:col-span-2">
              <label class="input-label">应用私钥（PEM）</label>
              <textarea v-model="config.alipay.private_key" rows="4" class="input mt-1 w-full font-mono text-xs"></textarea>
            </div>
            <div class="md:col-span-2">
              <label class="input-label">支付宝公钥（PEM）</label>
              <textarea v-model="config.alipay.public_key" rows="4" class="input mt-1 w-full font-mono text-xs"></textarea>
            </div>
          </div>
        </section>

        <section class="rounded-2xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">易支付</h2>
            <label class="flex items-center gap-2 text-sm">
              <input v-model="config.epay.enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
              启用
            </label>
          </div>
          <div class="mt-4 grid gap-4 md:grid-cols-2">
            <div class="md:col-span-2">
              <label class="input-label">网关地址</label>
              <input v-model="config.epay.gateway" class="input mt-1 w-full" placeholder="https://pay.example.com" />
            </div>
            <div>
              <label class="input-label">商户 PID</label>
              <input v-model="config.epay.pid" class="input mt-1 w-full" />
            </div>
            <div>
              <label class="input-label">商户密钥 PKEY</label>
              <input v-model="config.epay.pkey" type="password" class="input mt-1 w-full" />
            </div>
            <div class="md:col-span-2">
              <label class="input-label">Notify URL</label>
              <input v-model="config.epay.notify_url" class="input mt-1 w-full" />
            </div>
          </div>
        </section>

        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" :disabled="saving" @click="loadSite">重置</button>
          <button class="btn btn-primary" :disabled="saving" @click="handleSave">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { subSiteAdminAPI } from '@/api/subsiteAdmin'
import type { OwnedSubSite, OwnerPaymentConfig } from '@/api/subsite'
import type { CreateSubSiteActivationInput } from '@/api/payment'
import { useAppStore } from '@/stores'

const route = useRoute()
const appStore = useAppStore()

const siteId = computed(() => Number(route.params.siteId))

const site = ref<OwnedSubSite | null>(null)
const saving = ref(false)

const config = reactive<Required<OwnerPaymentConfig>>({
  wechat: { enabled: false, app_id: '', mch_id: '', apiv3_key: '', mch_serial_no: '', public_key_id: '', public_key: '', private_key: '', notify_url: '' },
  alipay: { enabled: false, app_id: '', private_key: '', public_key: '', notify_url: '', is_production: true },
  epay: { enabled: false, gateway: '', pid: '', pkey: '', notify_url: '' }
})

function populateConfig(source?: OwnerPaymentConfig | null) {
  config.wechat = {
    enabled: source?.wechat?.enabled ?? false,
    app_id: source?.wechat?.app_id ?? '',
    mch_id: source?.wechat?.mch_id ?? '',
    apiv3_key: source?.wechat?.apiv3_key ?? '',
    mch_serial_no: source?.wechat?.mch_serial_no ?? '',
    public_key_id: source?.wechat?.public_key_id ?? '',
    public_key: source?.wechat?.public_key ?? '',
    private_key: source?.wechat?.private_key ?? '',
    notify_url: source?.wechat?.notify_url ?? ''
  }
  config.alipay = {
    enabled: source?.alipay?.enabled ?? false,
    app_id: source?.alipay?.app_id ?? '',
    private_key: source?.alipay?.private_key ?? '',
    public_key: source?.alipay?.public_key ?? '',
    notify_url: source?.alipay?.notify_url ?? '',
    is_production: source?.alipay?.is_production ?? true
  }
  config.epay = {
    enabled: source?.epay?.enabled ?? false,
    gateway: source?.epay?.gateway ?? '',
    pid: source?.epay?.pid ?? '',
    pkey: source?.epay?.pkey ?? '',
    notify_url: source?.epay?.notify_url ?? ''
  }
}

async function loadSite() {
  if (!siteId.value) return
  try {
    const data = await subSiteAdminAPI.getSite(siteId.value)
    site.value = data
    populateConfig(data.owner_payment_config)
  } catch (e: any) {
    appStore.showError(e?.message || '加载分站失败')
  }
}

async function handleSave() {
  if (!site.value || !siteId.value) return
  saving.value = true
  try {
    const payload: CreateSubSiteActivationInput & { owner_payment_config?: OwnerPaymentConfig } = {
      name: site.value.name,
      slug: site.value.slug,
      custom_domain: site.value.custom_domain || '',
      site_logo: site.value.site_logo || '',
      site_favicon: site.value.site_favicon || '',
      site_subtitle: site.value.site_subtitle || '',
      announcement: site.value.announcement || '',
      contact_info: site.value.contact_info || '',
      doc_url: site.value.doc_url || '',
      home_content: site.value.home_content || '',
      theme_template: site.value.theme_template || 'starter',
      allow_sub_site: site.value.allow_sub_site,
      sub_site_price_fen: site.value.sub_site_price_fen || 0,
      consume_rate_multiplier: Number(site.value.consume_rate_multiplier || 1),
      owner_payment_config: config as OwnerPaymentConfig
    }
    await subSiteAdminAPI.updateSite(siteId.value, payload as any)
    appStore.showSuccess('已保存')
    await loadSite()
  } catch (e: any) {
    appStore.showError(e?.response?.data?.error || e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(loadSite)
watch(siteId, loadSite)
</script>
