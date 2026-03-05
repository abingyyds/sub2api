<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <SearchInput
            v-model="searchQuery"
            :placeholder="t('org.members.searchPlaceholder')"
            class="w-64"
          />
          <button class="btn btn-primary" @click="showCreateModal = true">
            {{ t('org.members.createEmployee') }}
          </button>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="members"
          :loading="loading"
          :empty-text="t('org.members.empty')"
        >
          <template #cell-username="{ row }">
            <div>
              <div class="font-medium text-gray-900 dark:text-white">{{ row.username }}</div>
              <div class="text-xs text-gray-500">{{ row.email }}</div>
            </div>
          </template>

          <template #cell-role="{ value }">
            <span class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
              :class="value === 'org_admin' ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300' : 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'"
            >
              {{ value === 'org_admin' ? t('org.members.roleAdmin') : t('org.members.roleMember') }}
            </span>
          </template>

          <template #cell-daily_quota_usd="{ row }">
            <span class="font-mono text-sm">
              ${{ Number(row.daily_usage_usd).toFixed(4) }}
              <span class="text-gray-400"> / {{ row.daily_quota_usd != null ? '$' + Number(row.daily_quota_usd).toFixed(2) : '∞' }}</span>
            </span>
          </template>

          <template #cell-monthly_quota_usd="{ row }">
            <span class="font-mono text-sm">
              ${{ Number(row.monthly_usage_usd).toFixed(4) }}
              <span class="text-gray-400"> / {{ row.monthly_quota_usd != null ? '$' + Number(row.monthly_quota_usd).toFixed(2) : '∞' }}</span>
            </span>
          </template>

          <template #cell-status="{ value }">
            <StatusBadge :status="value" />
          </template>

          <template #cell-actions="{ row }">
            <div class="flex items-center gap-2">
              <button class="btn btn-sm btn-ghost" @click="openEdit(row)">{{ t('common.edit') }}</button>
              <button
                v-if="row.status === 'active'"
                class="btn btn-sm btn-ghost text-orange-600"
                @click="handleSuspend(row)"
              >
                {{ t('org.members.suspend') }}
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
          :current-page="pagination.page"
          :total-pages="pagination.pages"
          :total-items="pagination.total"
          @page-change="handlePageChange"
        />
      </template>
    </TablePageLayout>

    <!-- Create Employee Modal -->
    <BaseDialog v-model="showCreateModal" :title="t('org.members.createTitle')">
      <form @submit.prevent="handleCreate">
        <div class="space-y-4">
          <Input v-model="createForm.email" :label="t('org.members.email')" type="email" required />
          <Input v-model="createForm.password" :label="t('org.members.password')" type="password" required />
          <Input v-model="createForm.username" :label="t('org.members.username')" />
          <Select v-model="createForm.role" :label="t('org.members.role')" :options="roleOptions" />
          <Input v-model.number="createForm.daily_quota_usd" :label="t('org.members.dailyQuota')" type="number" step="0.01" :placeholder="t('org.members.quotaUnlimited')" />
          <Input v-model.number="createForm.monthly_quota_usd" :label="t('org.members.monthlyQuota')" type="number" step="0.01" :placeholder="t('org.members.quotaUnlimited')" />
        </div>
        <div class="mt-6 flex justify-end gap-3">
          <button type="button" class="btn btn-ghost" @click="showCreateModal = false">{{ t('common.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="creating">{{ t('common.save') }}</button>
        </div>
      </form>
    </BaseDialog>

    <!-- Edit Member Modal -->
    <BaseDialog v-model="showEditModal" :title="t('org.members.editTitle')">
      <form @submit.prevent="handleUpdate">
        <div class="space-y-4">
          <Select v-model="editForm.role" :label="t('org.members.role')" :options="roleOptions" />
          <Input v-model.number="editForm.daily_quota_usd" :label="t('org.members.dailyQuota')" type="number" step="0.01" :placeholder="t('org.members.quotaUnlimited')" />
          <Input v-model.number="editForm.monthly_quota_usd" :label="t('org.members.monthlyQuota')" type="number" step="0.01" :placeholder="t('org.members.quotaUnlimited')" />
          <TextArea v-model="editForm.notes" :label="t('org.members.notes')" />
        </div>
        <div class="mt-6 flex justify-end gap-3">
          <button type="button" class="btn btn-ghost" @click="showEditModal = false">{{ t('common.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="updating">{{ t('common.save') }}</button>
        </div>
      </form>
    </BaseDialog>

    <!-- Delete Confirm -->
    <ConfirmDialog
      v-model="showDeleteConfirm"
      :title="t('org.members.deleteTitle')"
      :message="t('org.members.deleteConfirm', { email: selectedMember?.email ?? '' })"
      :confirm-text="t('common.delete')"
      danger
      :loading="deleting"
      @confirm="handleDelete"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { orgAPI } from '@/api/org'
import type { OrgMember } from '@/types'
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

const members = ref<OrgMember[]>([])
const loading = ref(false)
const searchQuery = ref('')
const pagination = reactive({ page: 1, pages: 1, total: 0 })
let abortController: AbortController | null = null

const columns = computed(() => [
  { key: 'username', label: t('org.members.name') },
  { key: 'role', label: t('org.members.role') },
  { key: 'daily_quota_usd', label: t('org.members.dailyUsage') },
  { key: 'monthly_quota_usd', label: t('org.members.monthlyUsage') },
  { key: 'status', label: t('common.status') },
  { key: 'actions', label: t('common.actions') }
])

const roleOptions = [
  { label: t('org.members.roleMember'), value: 'member' },
  { label: t('org.members.roleAdmin'), value: 'org_admin' }
]

// Create
const showCreateModal = ref(false)
const creating = ref(false)
const createForm = reactive({
  email: '',
  password: '',
  username: '',
  role: 'member',
  daily_quota_usd: null as number | null,
  monthly_quota_usd: null as number | null
})

// Edit
const showEditModal = ref(false)
const updating = ref(false)
const selectedMember = ref<OrgMember | null>(null)
const editForm = reactive({
  role: 'member',
  daily_quota_usd: null as number | null,
  monthly_quota_usd: null as number | null,
  notes: ''
})

// Delete
const showDeleteConfirm = ref(false)
const deleting = ref(false)

async function loadMembers() {
  abortController?.abort()
  abortController = new AbortController()
  loading.value = true
  try {
    const res = await orgAPI.listMembers(pagination.page, 20, { signal: abortController.signal })
    members.value = res.items
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
  loadMembers()
}

async function handleCreate() {
  creating.value = true
  try {
    await orgAPI.createMember({
      email: createForm.email,
      password: createForm.password,
      username: createForm.username || undefined,
      role: createForm.role as 'org_admin' | 'member',
      daily_quota_usd: createForm.daily_quota_usd,
      monthly_quota_usd: createForm.monthly_quota_usd
    })
    showCreateModal.value = false
    Object.assign(createForm, { email: '', password: '', username: '', role: 'member', daily_quota_usd: null, monthly_quota_usd: null })
    appStore.showSuccess(t('common.createSuccess'))
    loadMembers()
  } catch {
    appStore.showError(t('common.createError'))
  } finally {
    creating.value = false
  }
}

function openEdit(member: OrgMember) {
  selectedMember.value = member
  Object.assign(editForm, {
    role: member.role,
    daily_quota_usd: member.daily_quota_usd,
    monthly_quota_usd: member.monthly_quota_usd,
    notes: member.notes || ''
  })
  showEditModal.value = true
}

async function handleUpdate() {
  if (!selectedMember.value) return
  updating.value = true
  try {
    await orgAPI.updateMember(selectedMember.value.id, {
      role: editForm.role as 'org_admin' | 'member',
      daily_quota_usd: editForm.daily_quota_usd,
      monthly_quota_usd: editForm.monthly_quota_usd,
      notes: editForm.notes
    })
    showEditModal.value = false
    appStore.showSuccess(t('common.updateSuccess'))
    loadMembers()
  } catch {
    appStore.showError(t('common.updateError'))
  } finally {
    updating.value = false
  }
}

async function handleSuspend(member: OrgMember) {
  try {
    await orgAPI.suspendMember(member.id)
    appStore.showSuccess(t('org.members.suspendSuccess'))
    loadMembers()
  } catch {
    appStore.showError(t('common.updateError'))
  }
}

function confirmDelete(member: OrgMember) {
  selectedMember.value = member
  showDeleteConfirm.value = true
}

async function handleDelete() {
  if (!selectedMember.value) return
  deleting.value = true
  try {
    await orgAPI.removeMember(selectedMember.value.id)
    showDeleteConfirm.value = false
    appStore.showSuccess(t('common.deleteSuccess'))
    loadMembers()
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
    loadMembers()
  }, 300)
})

onMounted(() => {
  loadMembers()
})
</script>
