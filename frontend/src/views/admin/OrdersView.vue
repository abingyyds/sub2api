<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap gap-3">
          <Select v-model="filters.status" :options="statusOptions" :placeholder="t('admin.orders.allStatus')" class="w-36" @update:modelValue="onFilterChange" />
          <Select v-model="filters.order_type" :options="typeOptions" :placeholder="t('admin.orders.allTypes')" class="w-36" @update:modelValue="onFilterChange" />
          <button @click="loadData" :disabled="loading" class="btn btn-secondary" :title="t('common.refresh')">
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="items" :loading="loading">
          <template #cell-amount_fen="{ value }">
            <span class="font-medium">¥{{ (value / 100).toFixed(value % 100 === 0 ? 0 : 2) }}</span>
          </template>
          <template #cell-order_type="{ value }">
            <span :class="['badge', value === 'balance' ? 'badge-primary' : 'badge-gray']">
              {{ value === 'balance' ? t('admin.orders.typeBalance') : t('admin.orders.typeSubscription') }}
            </span>
          </template>
          <template #cell-status="{ value }">
            <span :class="['badge', statusBadgeClass(value)]">
              {{ statusLabel(value) }}
            </span>
          </template>
          <template #cell-invoice_status="{ row }">
            <span :class="['badge', invoiceBadgeClass(row)]">
              {{ invoiceStatusLabel(row) }}
            </span>
          </template>
          <template #cell-invoice_request="{ row }">
            <div v-if="row.invoice_requested_at" class="max-w-xs whitespace-normal text-xs leading-5 text-gray-600 dark:text-dark-300">
              <div><span class="font-medium">{{ t('admin.orders.invoiceCompanyName') }}:</span> {{ row.invoice_company_name }}</div>
              <div><span class="font-medium">{{ t('admin.orders.invoiceTaxId') }}:</span> {{ row.invoice_tax_id }}</div>
              <div><span class="font-medium">{{ t('admin.orders.invoiceEmail') }}:</span> {{ row.invoice_email }}</div>
              <div><span class="font-medium">{{ t('admin.orders.invoiceRequestedAt') }}:</span> {{ formatDateTime(row.invoice_requested_at) }}</div>
              <div v-if="row.invoice_remark"><span class="font-medium">{{ t('admin.orders.invoiceRemark') }}:</span> {{ row.invoice_remark }}</div>
            </div>
            <span v-else class="text-sm text-gray-400">{{ t('admin.orders.invoiceNone') }}</span>
          </template>
          <template #cell-paid_at="{ value }">
            <span v-if="value" class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
            <span v-else class="text-sm text-gray-400">-</span>
          </template>
          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex justify-end">
              <button
                v-if="row.invoice_requested_at && !row.invoice_processed_at"
                class="rounded-lg bg-primary-600 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-primary-700 disabled:opacity-50"
                :disabled="processingOrderId === row.id"
                @click="markProcessed(row)"
              >
                {{ processingOrderId === row.id ? t('common.processing') : t('admin.orders.markProcessed') }}
              </button>
              <span v-else class="text-sm text-gray-400">-</span>
            </div>
          </template>
          <template #empty>
            <EmptyState :title="t('admin.orders.noData')" :description="t('admin.orders.noDataDesc')" />
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination v-if="pagination.total > 0" :page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="p => { pagination.page = p; loadData() }" @update:pageSize="s => { pagination.page_size = s; pagination.page = 1; loadData() }" />
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { ordersAPI } from '@/api/admin'
import type { AdminPaymentOrder } from '@/api/admin/orders'
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

const items = ref<AdminPaymentOrder[]>([])
const loading = ref(false)
const processingOrderId = ref<number | null>(null)
const pagination = ref({ page: 1, page_size: 20, total: 0 })
const filters = ref({ status: '', order_type: '' })

const columns = computed<Column[]>(() => [
  { key: 'order_no', label: t('admin.orders.orderNo'), sortable: false },
  { key: 'user_id', label: t('admin.orders.userId'), sortable: false },
  { key: 'order_type', label: t('admin.orders.orderType'), sortable: false },
  { key: 'amount_fen', label: t('admin.orders.amount'), sortable: false },
  { key: 'status', label: t('common.status'), sortable: false },
  { key: 'invoice_status', label: t('admin.orders.invoiceStatus'), sortable: false },
  { key: 'invoice_request', label: t('admin.orders.invoiceInfo'), sortable: false },
  { key: 'plan_key', label: t('admin.orders.planKey'), sortable: false },
  { key: 'paid_at', label: t('admin.orders.paidAt'), sortable: false },
  { key: 'created_at', label: t('common.createdAt'), sortable: false },
  { key: 'actions', label: t('common.actions'), sortable: false }
])

const statusOptions = computed(() => [
  { value: '', label: t('admin.orders.allStatus') },
  { value: 'pending', label: t('admin.orders.statusPending') },
  { value: 'paid', label: t('admin.orders.statusPaid') },
  { value: 'closed', label: t('admin.orders.statusClosed') }
])

const typeOptions = computed(() => [
  { value: '', label: t('admin.orders.allTypes') },
  { value: 'subscription', label: t('admin.orders.typeSubscription') },
  { value: 'balance', label: t('admin.orders.typeBalance') }
])

function statusBadgeClass(status: string) {
  switch (status) {
    case 'paid': return 'badge-success'
    case 'pending': return 'badge-warning'
    case 'closed': return 'badge-gray'
    default: return 'badge-gray'
  }
}

function statusLabel(status: string) {
  switch (status) {
    case 'paid': return t('admin.orders.statusPaid')
    case 'pending': return t('admin.orders.statusPending')
    case 'closed': return t('admin.orders.statusClosed')
    default: return status
  }
}

function invoiceStatusLabel(order: AdminPaymentOrder) {
  if (order.invoice_processed_at) return t('admin.orders.invoiceProcessed')
  if (order.invoice_requested_at) return t('admin.orders.invoicePending')
  return t('admin.orders.invoiceNone')
}

function invoiceBadgeClass(order: AdminPaymentOrder) {
  if (order.invoice_processed_at) return 'badge-success'
  if (order.invoice_requested_at) return 'badge-warning'
  return 'badge-gray'
}

function onFilterChange() {
  pagination.value.page = 1
  loadData()
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await ordersAPI.list(pagination.value.page, pagination.value.page_size, {
      status: filters.value.status,
      order_type: filters.value.order_type
    })
    items.value = res.items
    pagination.value.total = res.total
  } catch {
    appStore.showError(t('admin.orders.loadError'))
  } finally {
    loading.value = false
  }
}

const markProcessed = async (row: AdminPaymentOrder) => {
  processingOrderId.value = row.id
  try {
    await ordersAPI.markInvoiceProcessed(row.id)
    appStore.showSuccess(t('admin.orders.markProcessedSuccess'))
    await loadData()
  } catch (err: any) {
    appStore.showError(err?.message || t('admin.orders.markProcessedError'))
  } finally {
    processingOrderId.value = null
  }
}

onMounted(loadData)
</script>
