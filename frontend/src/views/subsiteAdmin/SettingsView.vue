<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl space-y-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站设置</h1>

      <div v-if="loading" class="flex justify-center py-16">
        <LoadingSpinner />
      </div>

      <div v-else-if="site" class="space-y-5 rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">分站名称</label>
            <input v-model="site.name" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Slug</label>
            <input v-model="site.slug" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">自定义域名</label>
            <input v-model="site.custom_domain" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Logo URL</label>
            <input v-model="site.site_logo" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">Favicon URL</label>
            <input v-model="site.site_favicon" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">副标题</label>
            <input v-model="site.site_subtitle" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">公告</label>
            <input v-model="site.announcement" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">联系信息</label>
            <input v-model="site.contact_info" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">文档地址</label>
            <input v-model="site.doc_url" class="input mt-1 w-full" />
          </div>
          <div class="md:col-span-2">
            <label class="input-label">首页内容</label>
            <textarea v-model="site.home_content" rows="5" class="input mt-1 w-full"></textarea>
          </div>
          <div class="md:col-span-2">
            <label class="input-label">主题模板</label>
            <input v-model="site.theme_template" class="input mt-1 w-full" placeholder="例如：starter / aurora / summit / terminal" />
          </div>
        </div>

        <div class="grid gap-4 rounded-2xl border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/60 md:grid-cols-3">
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-dark-200">
            <input v-model="site.allow_sub_site" type="checkbox" class="h-4 w-4 rounded border-gray-300" />
            允许发展下级分站
          </label>
          <div>
            <label class="input-label">下级分站售价（元）</label>
            <input v-model.number="subSitePriceYuan" type="number" min="0" step="1" class="input mt-1 w-full" />
          </div>
          <div>
            <label class="input-label">消耗倍率</label>
            <input v-model.number="site.consume_rate_multiplier" type="number" min="1" step="0.1" class="input mt-1 w-full" />
          </div>
        </div>

        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" :disabled="saving" @click="loadSite">重置</button>
          <button class="btn btn-primary" :disabled="saving" @click="handleSave">
            {{ saving ? '保存中...' : '保存分站' }}
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { subSiteAdminAPI } from '@/api/subsiteAdmin'
import type { OwnedSubSite } from '@/api/subsite'
import type { CreateSubSiteActivationInput } from '@/api/payment'
import { useAppStore, useAuthStore } from '@/stores'

const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteId = computed(() => Number(route.params.siteId))

const loading = ref(true)
const saving = ref(false)
const site = ref<OwnedSubSite | null>(null)
const subSitePriceYuan = ref<number | null>(null)

async function loadSite() {
  if (!siteId.value) return
  loading.value = true
  try {
    const data = await subSiteAdminAPI.getSite(siteId.value)
    site.value = JSON.parse(JSON.stringify(data)) as OwnedSubSite
    site.value.consume_rate_multiplier = Number(data.consume_rate_multiplier || 1)
    subSitePriceYuan.value = data.sub_site_price_fen ? data.sub_site_price_fen / 100 : null
  } catch (error: any) {
    appStore.showError(error?.message || '加载分站设置失败')
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!site.value || !siteId.value) return
  saving.value = true
  try {
    const payload: CreateSubSiteActivationInput = {
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
      sub_site_price_fen: subSitePriceYuan.value ? Math.round(subSitePriceYuan.value * 100) : 0,
      consume_rate_multiplier: Number(site.value.consume_rate_multiplier || 1),
    }
    await subSiteAdminAPI.updateSite(siteId.value, payload)
    appStore.showSuccess('分站已更新')
    await authStore.refreshOwnedSites().catch(() => { /* ignore */ })
    await loadSite()
  } catch (error: any) {
    appStore.showError(error?.message || '保存分站失败')
  } finally {
    saving.value = false
  }
}

onMounted(loadSite)
watch(siteId, loadSite)
</script>
