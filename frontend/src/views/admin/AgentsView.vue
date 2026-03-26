<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap gap-3">
          <Select v-model="filters.status" :options="statusOptions" :placeholder="t('admin.agents.allStatus')" class="w-36" @update:modelValue="onFilterChange" />
          <input v-model="filters.search" :placeholder="t('admin.agents.searchPlaceholder')" class="input w-48" @keyup.enter="onFilterChange" />
          <button @click="loadData" :disabled="loading" class="btn btn-secondary" :title="t('common.refresh')">
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="items" :loading="loading">
          <template #cell-agent_status="{ value, row }">
            <span :class="['badge', statusBadgeClass(value)]">
              {{ statusLabel(value) }}
            </span>
          </template>
          <template #cell-agent_commission_rate="{ value }">
            <span class="font-medium">{{ (value * 100).toFixed(1) }}%</span>
          </template>
          <template #cell-total_commission="{ value }">
            <span class="font-medium text-primary-600 dark:text-primary-400">${{ (value || 0).toFixed(2) }}</span>
          </template>
          <template #cell-pending_commission="{ value }">
            <span class="font-medium text-orange-600 dark:text-orange-400">${{ (value || 0).toFixed(2) }}</span>
          </template>
          <template #cell-agent_approved_at="{ value }">
            <span v-if="value" class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
            <span v-else class="text-sm text-gray-400">-</span>
          </template>
          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-2">
              <template v-if="row.agent_status === 'pending'">
                <button @click="handleApprove(row)" class="btn btn-xs btn-success">{{ t('admin.agents.approve') }}</button>
                <button @click="handleReject(row)" class="btn btn-xs btn-danger">{{ t('admin.agents.reject') }}</button>
              </template>
              <button @click="openEditRate(row)" class="btn btn-xs btn-secondary" :title="t('admin.agents.editRate')">
                <Icon name="edit" size="sm" />
              </button>
              <button v-if="row.agent_status === 'approved' && row.pending_commission > 0" @click="handleSettle(row)" class="btn btn-xs btn-primary">
                {{ t('admin.agents.settle') }}
              </button>
            </div>
          </template>
          <template #empty>
            <EmptyState :title="t('admin.agents.noData')" :description="t('admin.agents.noDataDesc')" />
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination v-if="pagination.total > 0" :page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="p => { pagination.page = p; loadData() }" @update:pageSize="s => { pagination.page_size = s; pagination.page = 1; loadData() }" />
      </template>
    </TablePageLayout>

    <!-- Edit Commission Rate Modal -->
    <Modal v-model="showEditModal" :title="t('admin.agents.editRateTitle')">
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1">{{ t('admin.agents.commissionRate') }} (%)</label>
          <input v-model.number="editRate" type="number" min="0" max="100" step="0.1" class="input w-full" />
        </div>
      </div>
      <template #footer>
        <button @click="showEditModal = false" class="btn btn-secondary">{{ t('common.cancel') }}</button>
        <button @click="saveRate" :disabled="saving" class="btn btn-primary">{{ saving ? t('common.saving') : t('common.save') }}</button>
      </template>
    </Modal>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { agentsAPI } from '@/api/admin'
import type { AdminAgent } from '@/api/admin/agents'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import Modal from '@/components/common/Modal.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateTime } from '@/utils/format'
import type { Column } from '@/components/common/types'

const { t } = useI18n()
const appStore = useAppStore()

const items = ref<AdminAgent[]>([])
const loading = ref(false)
const saving = ref(false)
const pagination = ref({ page: 1, page_size: 20, total: 0 })
const filters = ref({ status: '', search: '' })

// Edit modal
const showEditModal = ref(false)
const editRate = ref(0)
const editingAgent = ref<AdminAgent | null>(null)

const columns = computed<Column[]>(() => [
  { key: 'id', label: 'ID', sortable: false },
  { key: 'email', label: t('admin.agents.email'), sortable: false },
  { key: 'username', label: t('admin.agents.username'), sortable: false },
  { key: 'agent_status', label: t('common.status'), sortable: false },
  { key: 'agent_commission_rate', label: t('admin.agents.commissionRate'), sortable: false },
  { key: 'sub_user_count', label: t('admin.agents.subUserCount'), sortable: false },
  { key: 'total_commission', label: t('admin.agents.totalCommission'), sortable: false },
  { key: 'pending_commission', label: t('admin.agents.pendingCommission'), sortable: false },
  { key: 'created_at', label: t('common.createdAt'), sortable: false },
  { key: 'actions', label: t('common.actions'), sortable: false }
])

const statusOptions = computed(() => [
  { value: '', label: t('admin.agents.allStatus') },
  { value: 'pending', label: t('admin.agents.statusPending') },
  { value: 'approved', label: t('admin.agents.statusApproved') },
  { value: 'rejected', label: t('admin.agents.statusRejected') }
])

function statusBadgeClass(status: string) {
  switch (status) {
    case 'approved': return 'badge-success'
    case 'pending': return 'badge-warning'
    case 'rejected': return 'badge-danger'
    default: return 'badge-gray'
  }
}

function statusLabel(status: string) {
  switch (status) {
    case 'approved': return t('admin.agents.statusApproved')
    case 'pending': return t('admin.agents.statusPending')
    case 'rejected': return t('admin.agents.statusRejected')
    default: return status
  }
}

function onFilterChange() {
  pagination.value.page = 1
  loadData()
}

async function loadData() {
  loading.value = true
  try {
    const res = await agentsAPI.list(pagination.value.page, pagination.value.page_size, {
      status: filters.value.status,
      search: filters.value.search
    })
    items.value = res.items || []
    pagination.value.total = res.total
  } catch {
    appStore.showError(t('admin.agents.loadError'))
  } finally {
    loading.value = false
  }
}

async function handleApprove(agent: AdminAgent) {
  try {
    await agentsAPI.approve(agent.id)
    appStore.showSuccess(t('admin.agents.approveSuccess'))
    loadData()
  } catch {
    appStore.showError(t('admin.agents.approveError'))
  }
}

async function handleReject(agent: AdminAgent) {
  try {
    await agentsAPI.reject(agent.id)
    appStore.showSuccess(t('admin.agents.rejectSuccess'))
    loadData()
  } catch {
    appStore.showError(t('admin.agents.rejectError'))
  }
}

function openEditRate(agent: AdminAgent) {
  editingAgent.value = agent
  editRate.value = agent.agent_commission_rate * 100
  showEditModal.value = true
}

async function saveRate() {
  if (!editingAgent.value) return
  saving.value = true
  try {
    await agentsAPI.update(editingAgent.value.id, { commission_rate: editRate.value / 100 })
    appStore.showSuccess(t('admin.agents.updateSuccess'))
    showEditModal.value = false
    loadData()
  } catch {
    appStore.showError(t('admin.agents.updateError'))
  } finally {
    saving.value = false
  }
}

async function handleSettle(agent: AdminAgent) {
  try {
    const res = await agentsAPI.settle(agent.id)
    appStore.showSuccess(t('admin.agents.settleSuccess', { count: res.settled_count, amount: res.settled_amount.toFixed(2) }))
    loadData()
  } catch {
    appStore.showError(t('admin.agents.settleError'))
  }
}

onMounted(loadData)
</script>
