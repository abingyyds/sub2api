<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-6">
      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分站管理</h1>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              为不同域名配置独立品牌入口。分站是独立站点，不和代理体系绑定。
            </p>
          </div>
          <button class="btn btn-primary" @click="openCreate">新建分站</button>
        </div>

        <div class="mt-6 grid gap-3 md:grid-cols-[1fr_180px_auto]">
          <input v-model="filters.search" class="input" placeholder="搜索站点名 / slug / 域名 / owner 邮箱" @keyup.enter="loadData(1)" />
          <select v-model="filters.status" class="input" @change="loadData(1)">
            <option value="">全部状态</option>
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
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">入口</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">用户数</th>
                <th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">状态</th>
                <th class="px-4 py-3 text-right text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
              <tr v-if="loading">
                <td colspan="6" class="px-4 py-10 text-center text-sm text-gray-500 dark:text-dark-400">加载中...</td>
              </tr>
              <tr v-else-if="items.length === 0">
                <td colspan="6" class="px-4 py-10 text-center text-sm text-gray-500 dark:text-dark-400">还没有分站</td>
              </tr>
              <tr v-for="item in items" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50">
                <td class="px-4 py-4 align-top">
                  <div class="font-medium text-gray-900 dark:text-white">{{ item.name }}</div>
                  <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">slug: {{ item.slug }}</div>
                  <div v-if="item.custom_domain" class="mt-1 text-xs text-gray-500 dark:text-dark-400">域名: {{ item.custom_domain }}</div>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <div>ID: {{ item.owner_user_id }}</div>
                  <div class="mt-1 break-all text-xs text-gray-500 dark:text-dark-400">{{ item.owner_email || '-' }}</div>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">
                  <a v-if="item.entry_url" :href="item.entry_url" target="_blank" rel="noopener noreferrer" class="break-all text-primary-600 hover:underline dark:text-primary-400">
                    {{ item.entry_url }}
                  </a>
                  <span v-else class="text-xs text-gray-400 dark:text-dark-500">未配置主域名后缀或自定义域名</span>
                </td>
                <td class="px-4 py-4 align-top text-sm text-gray-600 dark:text-dark-300">{{ item.user_count || 0 }}</td>
                <td class="px-4 py-4 align-top">
                  <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium"
                    :class="item.status === 'active' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300' : 'bg-gray-200 text-gray-700 dark:bg-dark-700 dark:text-dark-300'">
                    {{ item.status === 'active' ? '启用' : '停用' }}
                  </span>
                </td>
                <td class="px-4 py-4 align-top text-right">
                  <div class="flex justify-end gap-2">
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
      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label class="input-label">站点名称</label>
          <input v-model="form.name" class="input mt-1 w-full" placeholder="例如：cCoder 华东站" />
        </div>
        <div>
          <label class="input-label">Owner 用户 ID</label>
          <input v-model.number="form.owner_user_id" type="number" min="1" class="input mt-1 w-full" placeholder="例如：123" />
        </div>
        <div>
          <label class="input-label">Slug</label>
          <input v-model="form.slug" class="input mt-1 w-full" placeholder="例如：east-hub" />
        </div>
        <div>
          <label class="input-label">状态</label>
          <select v-model="form.status" class="input mt-1 w-full">
            <option value="active">启用</option>
            <option value="disabled">停用</option>
          </select>
        </div>
        <div class="md:col-span-2">
          <label class="input-label">自定义域名</label>
          <input v-model="form.custom_domain" class="input mt-1 w-full" placeholder="例如：east.ccoder.me" />
        </div>
        <div>
          <label class="input-label">Logo URL</label>
          <input v-model="form.site_logo" class="input mt-1 w-full" placeholder="https://..." />
        </div>
        <div>
          <label class="input-label">Favicon URL</label>
          <input v-model="form.site_favicon" class="input mt-1 w-full" placeholder="https://..." />
        </div>
        <div class="md:col-span-2">
          <label class="input-label">副标题</label>
          <input v-model="form.site_subtitle" class="input mt-1 w-full" placeholder="一句话介绍这个分站" />
        </div>
        <div class="md:col-span-2">
          <label class="input-label">公告</label>
          <input v-model="form.announcement" class="input mt-1 w-full" placeholder="可选，当前先做存储预留" />
        </div>
        <div>
          <label class="input-label">联系信息</label>
          <input v-model="form.contact_info" class="input mt-1 w-full" placeholder="Telegram / WeChat / Email" />
        </div>
        <div>
          <label class="input-label">文档地址</label>
          <input v-model="form.doc_url" class="input mt-1 w-full" placeholder="https://docs.example.com" />
        </div>
        <div class="md:col-span-2">
          <label class="input-label">首页内容</label>
          <textarea v-model="form.home_content" rows="8" class="input mt-1 w-full" placeholder="支持直接填写 HTML，或填一个 https:// 页面地址供 iframe 展示" />
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
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { subSitesAPI, type AdminSubSite, type SaveSubSiteRequest } from '@/api/admin/subsites'
import { useAppStore } from '@/stores'

const appStore = useAppStore()

const loading = ref(false)
const saving = ref(false)
const items = ref<AdminSubSite[]>([])
const page = ref(1)
const pages = ref(1)

const filters = reactive({
  search: '',
  status: ''
})

const showDialog = ref(false)
const showDeleteDialog = ref(false)
const editingId = ref<number | null>(null)
const deleteTarget = ref<AdminSubSite | null>(null)

const form = reactive<SaveSubSiteRequest>({
  owner_user_id: 0,
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
  theme_config: ''
})

onMounted(() => {
  loadData(1)
})

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
  form.theme_config = ''
}

function openCreate() {
  resetForm()
  showDialog.value = true
}

function openEdit(item: AdminSubSite) {
  editingId.value = item.id
  form.owner_user_id = item.owner_user_id
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
  form.theme_config = item.theme_config || ''
  showDialog.value = true
}

function closeDialog() {
  showDialog.value = false
}

async function handleSave() {
  saving.value = true
  try {
    if (editingId.value) {
      await subSitesAPI.update(editingId.value, form)
      appStore.showSuccess('分站已更新')
    } else {
      await subSitesAPI.create(form)
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
</script>
