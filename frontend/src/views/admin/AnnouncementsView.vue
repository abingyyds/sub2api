<template>
  <AppLayout>
    <TablePageLayout>
      <template #actions>
        <div class="flex justify-end gap-3">
          <button @click="loadData" :disabled="loading" class="btn btn-secondary" :title="t('common.refresh')">
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
          <button @click="openCreate" class="btn btn-primary">
            <Icon name="plus" size="md" class="mr-2" />
            {{ t('admin.announcements.create') }}
          </button>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="items" :loading="loading">
          <template #cell-status="{ value }">
            <span :class="['badge', value === 'active' ? 'badge-success' : 'badge-gray']">
              {{ value === 'active' ? t('common.active') : t('common.inactive') }}
            </span>
          </template>
          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-1">
              <button @click="openEdit(row)" class="btn-icon">
                <Icon name="edit" size="sm" />
              </button>
              <button @click="confirmDelete(row)" class="btn-icon text-red-500 hover:text-red-700">
                <Icon name="trash" size="sm" />
              </button>
            </div>
          </template>
          <template #empty>
            <EmptyState :title="t('admin.announcements.noData')" :description="t('admin.announcements.noDataDesc')" :action-text="t('admin.announcements.create')" @action="openCreate" />
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination v-if="pagination.total > 0" :page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="p => { pagination.page = p; loadData() }" @update:pageSize="s => { pagination.page_size = s; pagination.page = 1; loadData() }" />
      </template>
    </TablePageLayout>

    <!-- Create/Edit Dialog -->
    <BaseDialog :show="showDialog" :title="editing ? t('admin.announcements.edit') : t('admin.announcements.create')" width="normal" @close="showDialog = false">
      <form id="ann-form" @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="input-label">{{ t('admin.announcements.titleLabel') }}</label>
          <input v-model="form.title" type="text" required class="input" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.announcements.contentLabel') }}</label>
          <textarea v-model="form.content" rows="4" class="input" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="input-label">{{ t('common.status') }}</label>
            <Select v-model="form.status" :options="statusOptions" />
          </div>
          <div>
            <label class="input-label">{{ t('admin.announcements.priority') }}</label>
            <input v-model.number="form.priority" type="number" class="input" min="0" />
          </div>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button @click="showDialog = false" class="btn btn-secondary">{{ t('common.cancel') }}</button>
          <button form="ann-form" type="submit" :disabled="submitting" class="btn btn-primary">
            {{ submitting ? t('common.saving') : editing ? t('common.update') : t('common.create') }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <ConfirmDialog :show="showDeleteDialog" :title="t('admin.announcements.delete')" :message="t('admin.announcements.deleteConfirm')" :confirm-text="t('common.delete')" :cancel-text="t('common.cancel')" :danger="true" @confirm="handleDelete" @cancel="showDeleteDialog = false" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { announcementsAPI } from '@/api/admin'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateTime } from '@/utils/format'
import type { Announcement } from '@/types'
import type { Column } from '@/components/common/types'

const { t } = useI18n()
const appStore = useAppStore()

const items = ref<Announcement[]>([])
const loading = ref(false)
const submitting = ref(false)
const showDialog = ref(false)
const showDeleteDialog = ref(false)
const editing = ref<Announcement | null>(null)
const deleting = ref<Announcement | null>(null)
const pagination = ref({ page: 1, page_size: 20, total: 0 })

const form = ref({ title: '', content: '', status: 'active', priority: 0 })

const columns = computed<Column[]>(() => [
  { key: 'id', label: 'ID', sortable: true },
  { key: 'title', label: t('admin.announcements.titleLabel'), sortable: true },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'priority', label: t('admin.announcements.priority'), sortable: true },
  { key: 'created_at', label: t('common.createdAt'), sortable: true },
  { key: 'actions', label: t('common.actions'), sortable: false }
])

const statusOptions = computed(() => [
  { value: 'active', label: t('common.active') },
  { value: 'inactive', label: t('common.inactive') }
])

const loadData = async () => {
  loading.value = true
  try {
    const res = await announcementsAPI.list(pagination.value.page, pagination.value.page_size)
    items.value = res.items
    pagination.value.total = res.total
  } catch { appStore.showError(t('admin.announcements.loadError')) }
  finally { loading.value = false }
}

const openCreate = () => {
  editing.value = null
  form.value = { title: '', content: '', status: 'active', priority: 0 }
  showDialog.value = true
}

const openEdit = (row: Announcement) => {
  editing.value = row
  form.value = { title: row.title, content: row.content, status: row.status, priority: row.priority }
  showDialog.value = true
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (editing.value) {
      await announcementsAPI.update(editing.value.id, form.value)
      appStore.showSuccess(t('admin.announcements.updated'))
    } else {
      await announcementsAPI.create(form.value)
      appStore.showSuccess(t('admin.announcements.created'))
    }
    showDialog.value = false
    loadData()
  } catch { appStore.showError(t('admin.announcements.saveError')) }
  finally { submitting.value = false }
}

const confirmDelete = (row: Announcement) => { deleting.value = row; showDeleteDialog.value = true }

const handleDelete = async () => {
  if (!deleting.value) return
  try {
    await announcementsAPI.remove(deleting.value.id)
    appStore.showSuccess(t('admin.announcements.deleted'))
    showDeleteDialog.value = false
    loadData()
  } catch { appStore.showError(t('admin.announcements.deleteError')) }
}

onMounted(loadData)
</script>
