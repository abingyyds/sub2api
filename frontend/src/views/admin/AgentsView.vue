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
          <template #cell-agent_status="{ value }">
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
          <template #cell-agent_approved_at="{ value }">
            <span v-if="value" class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
            <span v-else class="text-sm text-gray-400">-</span>
          </template>
          <template #cell-identity_status="{ value }">
            <span :class="['badge', value === 'submitted' ? 'badge-success' : 'badge-gray']">
              {{ value === 'submitted' ? '已提交' : '未提交' }}
            </span>
          </template>
          <template #cell-contract_status="{ value }">
            <span :class="['badge', value === 'signed' ? 'badge-success' : 'badge-gray']">
              {{ value === 'signed' ? '已确认' : '未确认' }}
            </span>
          </template>
          <template #cell-activation_fee_paid_at="{ value }">
            <span v-if="value" class="text-sm text-emerald-600 dark:text-emerald-400">{{ formatDateTime(value) }}</span>
            <span v-else class="text-sm text-gray-400">未支付</span>
          </template>
          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-2">
              <button @click="openDetail(row)" class="btn btn-xs btn-secondary" :title="t('admin.agents.viewDetail')">
                <Icon name="document" size="sm" />
              </button>
              <template v-if="row.agent_status === 'pending'">
                <button @click="handleApprove(row)" class="btn btn-xs btn-success">{{ t('admin.agents.approve') }}</button>
                <button @click="handleReject(row)" class="btn btn-xs btn-danger">{{ t('admin.agents.reject') }}</button>
              </template>
              <button @click="openEditRate(row)" class="btn btn-xs btn-secondary" :title="t('admin.agents.editRate')">
                <Icon name="edit" size="sm" />
              </button>
              <button @click="handleToggleFrozen(row)" class="btn btn-xs" :class="row.is_frozen ? 'btn-success' : 'btn-danger'">
                {{ row.is_frozen ? '解冻' : '冻结' }}
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
    <Teleport to="body">
      <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="fixed inset-0 bg-black/50" @click="showEditModal = false"></div>
        <div class="relative w-full max-w-md rounded-xl bg-white p-6 shadow-xl dark:bg-dark-800">
          <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.agents.editRateTitle') }}</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-dark-300 mb-1">{{ t('admin.agents.commissionRate') }} (%)</label>
              <input v-model.number="editRate" type="number" min="0" max="100" step="0.1" class="input w-full" />
            </div>
          </div>
          <div class="mt-6 flex justify-end gap-3">
            <button @click="showEditModal = false" class="btn btn-secondary">{{ t('common.cancel') }}</button>
            <button @click="saveRate" :disabled="saving" class="btn btn-primary">{{ saving ? t('common.saving') : t('common.save') }}</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Agent Detail Modal -->
    <Teleport to="body">
      <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="fixed inset-0 bg-black/50" @click="showDetailModal = false"></div>
        <div class="relative w-full max-w-md rounded-xl bg-white p-6 shadow-xl dark:bg-dark-800">
          <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.agents.viewDetail') }}</h3>
          <div v-if="detailLoading" class="py-10 text-center text-sm text-gray-500 dark:text-dark-400">加载中...</div>
          <div v-else-if="detailAgent" class="space-y-3">
            <div class="text-sm text-gray-500 dark:text-dark-400">{{ detailAgent.email }}</div>
            <div class="grid grid-cols-2 gap-3 rounded-2xl bg-gray-50 p-3 text-sm dark:bg-dark-800">
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">实名状态</label>
                <p class="text-gray-900 dark:text-white">{{ detailAgent.identity_status === 'submitted' ? '已提交' : '未提交' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">合同状态</label>
                <p class="text-gray-900 dark:text-white">{{ detailAgent.contract_status === 'signed' ? '已确认' : '未确认' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">真实姓名</label>
                <p class="text-gray-900 dark:text-white break-all">{{ detailAgent.real_name || '-' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">身份证号</label>
                <p class="text-gray-900 dark:text-white break-all">{{ detailAgent.id_card_no || '-' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">手机号</label>
                <p class="text-gray-900 dark:text-white break-all">{{ detailAgent.phone || '-' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">开通费支付</label>
                <p class="text-gray-900 dark:text-white">{{ detailAgent.activation_fee_paid_at ? formatDateTime(detailAgent.activation_fee_paid_at) : '未支付' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">冻结状态</label>
                <p class="text-gray-900 dark:text-white">{{ detailAgent.is_frozen ? `已冻结：${detailAgent.frozen_reason || '无原因'}` : '正常' }}</p>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3 rounded-2xl bg-gray-50 p-3 text-sm dark:bg-dark-800">
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">实名提交时间</label>
                <p class="text-gray-900 dark:text-white">{{ detailAgent.identity_submitted_at ? formatDateTime(detailAgent.identity_submitted_at) : '-' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">合同版本</label>
                <p class="text-gray-900 dark:text-white">{{ detailAgent.contract_version || '-' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">签署时间</label>
                <p class="text-gray-900 dark:text-white">{{ detailAgent.contract_signed_at ? formatDateTime(detailAgent.contract_signed_at) : '-' }}</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">签署 IP</label>
                <p class="text-gray-900 dark:text-white break-all">{{ detailAgent.contract_ip || '-' }}</p>
              </div>
            </div>
            <div>
              <label class="mb-2 block text-xs font-medium text-gray-500 dark:text-dark-400">合同签字</label>
              <div v-if="detailAgent.contract_signature_data" class="rounded-2xl border border-gray-200 bg-white p-3 dark:border-dark-700 dark:bg-dark-900">
                <img :src="detailAgent.contract_signature_data" alt="合同签字" class="max-h-28 rounded-lg" />
              </div>
              <p v-else class="text-sm text-gray-500 dark:text-dark-400">未签字</p>
            </div>
            <div v-if="detailAgent.agent_note">
              <label class="block text-xs font-medium text-gray-500 dark:text-dark-400 mb-1">{{ t('admin.agents.agentNote') }}</label>
              <p class="text-sm text-gray-900 dark:text-white whitespace-pre-wrap break-all">{{ detailAgent.agent_note }}</p>
            </div>
          </div>
          <div class="mt-6 flex justify-end">
            <button @click="showDetailModal = false" class="btn btn-secondary">{{ t('common.close') }}</button>
          </div>
        </div>
      </div>
    </Teleport>

  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { agentsAPI } from '@/api/admin'
import type { AdminAgent, AdminAgentDetail } from '@/api/admin/agents'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
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

// Detail modal
const showDetailModal = ref(false)
const detailLoading = ref(false)
const detailAgent = ref<AdminAgentDetail | null>(null)

const columns = computed<Column[]>(() => [
  { key: 'id', label: 'ID', sortable: false },
  { key: 'email', label: t('admin.agents.email'), sortable: false },
  { key: 'real_name', label: '实名', sortable: false },
  { key: 'username', label: t('admin.agents.username'), sortable: false },
  { key: 'agent_status', label: t('common.status'), sortable: false },
  { key: 'identity_status', label: '实名资料', sortable: false },
  { key: 'contract_status', label: '合同', sortable: false },
  { key: 'activation_fee_paid_at', label: '开通费', sortable: false },
  { key: 'agent_commission_rate', label: t('admin.agents.commissionRate'), sortable: false },
  { key: 'sub_user_count', label: t('admin.agents.subUserCount'), sortable: false },
  { key: 'total_commission', label: t('admin.agents.totalCommission'), sortable: false },
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

async function openDetail(agent: AdminAgent) {
  showDetailModal.value = true
  detailLoading.value = true
  try {
    detailAgent.value = await agentsAPI.getDetail(agent.id)
  } catch {
    detailAgent.value = null
    appStore.showError('加载代理详情失败')
  } finally {
    detailLoading.value = false
  }
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

async function handleToggleFrozen(agent: AdminAgent) {
  try {
    await agentsAPI.setFrozen(agent.id, !agent.is_frozen, agent.is_frozen ? '' : '管理员手动冻结')
    appStore.showSuccess(agent.is_frozen ? '代理已解冻' : '代理已冻结')
    await loadData()
  } catch {
    appStore.showError(agent.is_frozen ? '代理解冻失败' : '代理冻结失败')
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

onMounted(loadData)
</script>
