<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <SearchInput
            v-model="searchQuery"
            :placeholder="t('admin.organizations.searchPlaceholder')"
            class="w-64"
          />
          <button
            class="btn btn-primary"
            @click="showCreateModal = true"
          >
            {{ t('admin.organizations.create') }}
          </button>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="organizations"
          :loading="loading"
          :empty-text="t('admin.organizations.empty')"
        >
          <template #cell-name="{ row }">
            <div>
              <div class="font-medium text-gray-900 dark:text-white">{{ row.name }}</div>
              <div class="text-xs text-gray-500">{{ row.slug }}</div>
            </div>
          </template>

          <template #cell-status="{ value }">
            <StatusBadge :status="value" :label="t(`common.${value}`)" />
          </template>

          <template #cell-billing_mode="{ value }">
            <span class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
              :class="value === 'balance' ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300' : 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300'"
            >
              {{ value === 'balance' ? t('admin.organizations.billingBalance') : t('admin.organizations.billingSubscription') }}
            </span>
          </template>

          <template #cell-balance="{ value }">
            <span class="font-mono">${{ Number(value).toFixed(4) }}</span>
          </template>

          <template #cell-member_count="{ row }">
            {{ row.member_count ?? 0 }} / {{ row.max_members }}
          </template>

          <template #cell-created_at="{ value }">
            {{ new Date(value).toLocaleDateString() }}
          </template>

          <template #cell-actions="{ row }">
            <div class="flex items-center gap-2">
              <button class="btn btn-sm btn-ghost" @click="openEdit(row)">
                {{ t('common.edit') }}
              </button>
              <button class="btn btn-sm btn-ghost" @click="openBalance(row)">
                {{ t('admin.organizations.updateBalance') }}
              </button>
              <button class="btn btn-sm btn-ghost text-red-600" @click="confirmDelete(row)">
                {{ t('common.delete') }}
              </button>
            </div>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.pages > 1"
          :page="pagination.page"
          :page-size="20"
          :total="pagination.total"
          @update:page="handlePageChange"
        />
      </template>
    </TablePageLayout>

    <!-- Create Modal -->
    <BaseDialog :show="showCreateModal" :title="t('admin.organizations.createTitle')" @close="showCreateModal = false">
      <form @submit.prevent="handleCreate">
        <div class="space-y-4">
          <Input v-model="createForm.name" :label="t('admin.organizations.name')" required />
          <Input v-model="createForm.slug" :label="t('admin.organizations.slug')" required />
          <Input v-model.number="createForm.owner_user_id" :label="t('admin.organizations.ownerUserId')" type="number" required />
          <TextArea v-model="createForm.description" :label="t('admin.organizations.descriptionField')" />
          <Select v-model="createForm.billing_mode" :label="t('admin.organizations.billingMode')" :options="billingModeOptions" />
          <Input v-model.number="createForm.balance" :label="t('admin.organizations.initialBalance')" type="number" step="0.0001" />
          <Input v-model.number="createForm.max_members" :label="t('admin.organizations.maxMembers')" type="number" />
          <Input v-model.number="createForm.max_api_keys" :label="t('admin.organizations.maxApiKeys')" type="number" />
        </div>
        <div class="mt-6 flex justify-end gap-3">
          <button type="button" class="btn btn-ghost" @click="showCreateModal = false">{{ t('common.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="creating">{{ t('common.save') }}</button>
        </div>
      </form>
    </BaseDialog>

    <!-- Edit Modal -->
    <BaseDialog :show="showEditModal" :title="t('admin.organizations.editTitle')" @close="showEditModal = false">
      <form @submit.prevent="handleUpdate">
        <div class="space-y-4">
          <Input v-model="editForm.name" :label="t('admin.organizations.name')" />
          <Input v-model="editForm.slug" :label="t('admin.organizations.slug')" />
          <TextArea v-model="editForm.description" :label="t('admin.organizations.descriptionField')" />
          <Select v-model="editForm.billing_mode" :label="t('admin.organizations.billingMode')" :options="billingModeOptions" />
          <Input v-model.number="editForm.max_members" :label="t('admin.organizations.maxMembers')" type="number" />
          <Input v-model.number="editForm.max_api_keys" :label="t('admin.organizations.maxApiKeys')" type="number" />
          <Select v-model="editForm.status" :label="t('common.status')" :options="statusOptions" />
        </div>
        <div class="mt-6 flex justify-end gap-3">
          <button type="button" class="btn btn-ghost" @click="showEditModal = false">{{ t('common.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="updating">{{ t('common.save') }}</button>
        </div>
      </form>
    </BaseDialog>

    <!-- Balance Modal -->
    <BaseDialog :show="showBalanceModal" :title="t('admin.organizations.updateBalanceTitle')" @close="showBalanceModal = false">
      <form @submit.prevent="handleUpdateBalance">
        <div class="space-y-4">
          <div class="text-sm text-gray-500">
            {{ t('admin.organizations.currentBalance') }}: <span class="font-mono font-medium">${{ selectedOrg?.balance?.toFixed(4) ?? '0.0000' }}</span>
          </div>
          <Select v-model="balanceForm.action" :label="t('admin.organizations.balanceAction')" :options="balanceActionOptions" />
          <Input v-model.number="balanceForm.amount" :label="t('admin.organizations.amount')" type="number" step="0.0001" required />
        </div>
        <div class="mt-6 flex justify-end gap-3">
          <button type="button" class="btn btn-ghost" @click="showBalanceModal = false">{{ t('common.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="updatingBalance">{{ t('common.save') }}</button>
        </div>
      </form>
    </BaseDialog>

    <!-- Delete Confirm -->
    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="t('admin.organizations.deleteTitle')"
      :message="t('admin.organizations.deleteConfirm', { name: selectedOrg?.name ?? '' })"
      :confirm-text="t('common.delete')"
      danger
      @confirm="handleDelete"
      @cancel="showDeleteConfirm = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { adminAPI } from '@/api/admin'
import type { Organization } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import SearchInput from '@/components/common/SearchInput.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import Input from '@/components/common/Input.vue'
import TextArea from '@/components/common/TextArea.vue'
import Select from '@/components/common/Select.vue'

const { t } = useI18n()
const appStore = useAppStore()

const organizations = ref<Organization[]>([])
const loading = ref(false)
const searchQuery = ref('')
const pagination = reactive({ page: 1, pages: 1, total: 0 })
let abortController: AbortController | null = null

const columns = computed(() => [
  { key: 'name', label: t('admin.organizations.name'), sortable: true },
  { key: 'status', label: t('common.status') },
  { key: 'billing_mode', label: t('admin.organizations.billingMode') },
  { key: 'balance', label: t('admin.organizations.balance') },
  { key: 'member_count', label: t('admin.organizations.members') },
  { key: 'created_at', label: t('common.createdAt') },
  { key: 'actions', label: t('common.actions') }
])

const billingModeOptions = [
  { label: t('admin.organizations.billingBalance'), value: 'balance' },
  { label: t('admin.organizations.billingSubscription'), value: 'subscription' }
]
const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Suspended', value: 'suspended' },
  { label: 'Disabled', value: 'disabled' }
]
const balanceActionOptions = [
  { label: t('admin.organizations.balanceAdd'), value: 'add' },
  { label: t('admin.organizations.balanceSet'), value: 'set' }
]

// Create
const showCreateModal = ref(false)
const creating = ref(false)
const createForm = reactive({
  name: '',
  slug: '',
  owner_user_id: 0,
  description: '',
  billing_mode: 'balance',
  balance: 0,
  max_members: 50,
  max_api_keys: 100
})

// Edit
const showEditModal = ref(false)
const updating = ref(false)
const selectedOrg = ref<Organization | null>(null)
const editForm = reactive({
  name: '',
  slug: '',
  description: '',
  billing_mode: 'balance',
  max_members: 50,
  max_api_keys: 100,
  status: 'active'
})

// Balance
const showBalanceModal = ref(false)
const updatingBalance = ref(false)
const balanceForm = reactive({ action: 'add' as 'add' | 'set', amount: 0 })

// Delete
const showDeleteConfirm = ref(false)
const deleting = ref(false)

async function loadOrganizations() {
  abortController?.abort()
  abortController = new AbortController()
  loading.value = true
  try {
    const res = await adminAPI.organizations.list(pagination.page, 20, { signal: abortController.signal })
    organizations.value = res.items
    pagination.total = res.total
    pagination.pages = res.pages
  } catch (e: unknown) {
    if ((e as Error).name !== 'AbortError') {
      appStore.showError(t('common.loadError'))
    }
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  pagination.page = page
  loadOrganizations()
}

async function handleCreate() {
  creating.value = true
  try {
    await adminAPI.organizations.create({
      name: createForm.name,
      slug: createForm.slug,
      owner_user_id: createForm.owner_user_id,
      description: createForm.description || undefined,
      billing_mode: createForm.billing_mode as 'balance' | 'subscription',
      balance: createForm.balance,
      max_members: createForm.max_members,
      max_api_keys: createForm.max_api_keys
    })
    showCreateModal.value = false
    Object.assign(createForm, { name: '', slug: '', owner_user_id: 0, description: '', billing_mode: 'balance', balance: 0, max_members: 50, max_api_keys: 100 })
    appStore.showSuccess(t('common.createSuccess'))
    loadOrganizations()
  } catch {
    appStore.showError(t('common.createError'))
  } finally {
    creating.value = false
  }
}

function openEdit(org: Organization) {
  selectedOrg.value = org
  Object.assign(editForm, {
    name: org.name,
    slug: org.slug,
    description: org.description,
    billing_mode: org.billing_mode,
    max_members: org.max_members,
    max_api_keys: org.max_api_keys,
    status: org.status
  })
  showEditModal.value = true
}

async function handleUpdate() {
  if (!selectedOrg.value) return
  updating.value = true
  try {
    await adminAPI.organizations.update(selectedOrg.value.id, {
      name: editForm.name,
      slug: editForm.slug,
      description: editForm.description,
      billing_mode: editForm.billing_mode as 'balance' | 'subscription',
      max_members: editForm.max_members,
      max_api_keys: editForm.max_api_keys,
      status: editForm.status as 'active' | 'suspended' | 'disabled'
    })
    showEditModal.value = false
    appStore.showSuccess(t('common.updateSuccess'))
    loadOrganizations()
  } catch {
    appStore.showError(t('common.updateError'))
  } finally {
    updating.value = false
  }
}

function openBalance(org: Organization) {
  selectedOrg.value = org
  balanceForm.action = 'add'
  balanceForm.amount = 0
  showBalanceModal.value = true
}

async function handleUpdateBalance() {
  if (!selectedOrg.value) return
  updatingBalance.value = true
  try {
    await adminAPI.organizations.updateBalance(selectedOrg.value.id, {
      action: balanceForm.action,
      amount: balanceForm.amount
    })
    showBalanceModal.value = false
    appStore.showSuccess(t('common.updateSuccess'))
    loadOrganizations()
  } catch {
    appStore.showError(t('common.updateError'))
  } finally {
    updatingBalance.value = false
  }
}

function confirmDelete(org: Organization) {
  selectedOrg.value = org
  showDeleteConfirm.value = true
}

async function handleDelete() {
  if (!selectedOrg.value) return
  deleting.value = true
  try {
    await adminAPI.organizations.remove(selectedOrg.value.id)
    showDeleteConfirm.value = false
    appStore.showSuccess(t('common.deleteSuccess'))
    loadOrganizations()
  } catch {
    appStore.showError(t('common.deleteError'))
  } finally {
    deleting.value = false
  }
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null
watch(searchQuery, () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    pagination.page = 1
    loadOrganizations()
  }, 300)
})

onMounted(() => {
  loadOrganizations()
})
</script>
