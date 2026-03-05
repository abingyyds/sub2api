<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <select v-model="filters.action" class="input w-40" @change="fetchLogs">
            <option value="">{{ t('org.audit.allActions') }}</option>
            <option value="api_request">{{ t('org.audit.actionApiRequest') }}</option>
            <option value="member.create">{{ t('org.audit.actionMemberCreate') }}</option>
            <option value="member.update">{{ t('org.audit.actionMemberUpdate') }}</option>
            <option value="member.delete">{{ t('org.audit.actionMemberDelete') }}</option>
            <option value="policy.update">{{ t('org.audit.actionPolicyUpdate') }}</option>
          </select>
          <select v-model="flaggedFilter" class="input w-32" @change="fetchLogs">
            <option value="">{{ t('org.audit.allFlags') }}</option>
            <option value="true">{{ t('org.audit.flaggedOnly') }}</option>
            <option value="false">{{ t('org.audit.unflaggedOnly') }}</option>
          </select>
          <input
            v-model="filters.start_date"
            type="date"
            class="input w-40"
            @change="fetchLogs"
          />
          <input
            v-model="filters.end_date"
            type="date"
            class="input w-40"
            @change="fetchLogs"
          />
          <button class="btn btn-ghost btn-sm" @click="showConfigModal = true">
            {{ t('org.audit.config') }}
          </button>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="logs"
          :loading="loading"
          :empty-text="t('org.audit.empty')"
        >
          <template #cell-action="{ value }">
            <span class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300">
              {{ value }}
            </span>
          </template>

          <template #cell-model="{ value }">
            <span v-if="value" class="text-xs font-mono">{{ value }}</span>
            <span v-else class="text-gray-400">-</span>
          </template>

          <template #cell-flagged="{ value }">
            <span v-if="value" class="inline-flex items-center rounded-full bg-red-100 px-2 py-0.5 text-xs font-medium text-red-700 dark:bg-red-900/30 dark:text-red-300">
              {{ t('org.audit.flagged') }}
            </span>
          </template>

          <template #cell-cost_usd="{ value }">
            <span v-if="value != null" class="font-mono text-sm">${{ Number(value).toFixed(6) }}</span>
            <span v-else class="text-gray-400">-</span>
          </template>

          <template #cell-request_summary="{ value }">
            <span v-if="value" class="text-xs text-gray-600 dark:text-gray-400 truncate max-w-xs block">{{ value }}</span>
            <span v-else class="text-gray-400 text-xs">-</span>
          </template>

          <template #cell-created_at="{ value }">
            <span class="text-xs text-gray-500">{{ new Date(value).toLocaleString() }}</span>
          </template>

          <template #cell-actions="{ row }">
            <button class="btn btn-xs btn-ghost" @click="viewDetail(row)">{{ t('common.view') }}</button>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-model:page="page"
          :page-size="pageSize"
          :total="total"
        />
      </template>
    </TablePageLayout>

    <!-- Detail Modal -->
    <Modal v-model="showDetailModal" :title="t('org.audit.detail')">
      <div v-if="selectedLog" class="space-y-4">
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span class="text-gray-500">{{ t('org.audit.action') }}:</span>
            <span class="ml-2 font-medium">{{ selectedLog.action }}</span>
          </div>
          <div>
            <span class="text-gray-500">{{ t('org.audit.auditMode') }}:</span>
            <span class="ml-2">{{ selectedLog.audit_mode }}</span>
          </div>
          <div v-if="selectedLog.model">
            <span class="text-gray-500">{{ t('org.audit.model') }}:</span>
            <span class="ml-2 font-mono">{{ selectedLog.model }}</span>
          </div>
          <div v-if="selectedLog.cost_usd != null">
            <span class="text-gray-500">{{ t('org.audit.cost') }}:</span>
            <span class="ml-2 font-mono">${{ Number(selectedLog.cost_usd).toFixed(6) }}</span>
          </div>
          <div v-if="selectedLog.input_tokens != null">
            <span class="text-gray-500">{{ t('org.audit.inputTokens') }}:</span>
            <span class="ml-2">{{ selectedLog.input_tokens }}</span>
          </div>
          <div v-if="selectedLog.output_tokens != null">
            <span class="text-gray-500">{{ t('org.audit.outputTokens') }}:</span>
            <span class="ml-2">{{ selectedLog.output_tokens }}</span>
          </div>
          <div v-if="selectedLog.ip_address">
            <span class="text-gray-500">IP:</span>
            <span class="ml-2 font-mono text-xs">{{ selectedLog.ip_address }}</span>
          </div>
          <div>
            <span class="text-gray-500">{{ t('org.audit.time') }}:</span>
            <span class="ml-2">{{ new Date(selectedLog.created_at).toLocaleString() }}</span>
          </div>
        </div>
        <div v-if="selectedLog.request_summary" class="border-t pt-3 dark:border-gray-700">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('org.audit.requestSummary') }}</h4>
          <p class="text-sm text-gray-600 dark:text-gray-400 whitespace-pre-wrap">{{ selectedLog.request_summary }}</p>
        </div>
        <div v-if="selectedLog.request_content" class="border-t pt-3 dark:border-gray-700">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('org.audit.requestContent') }}</h4>
          <pre class="text-xs bg-gray-100 dark:bg-gray-800 p-3 rounded overflow-auto max-h-60">{{ selectedLog.request_content }}</pre>
        </div>
        <div v-if="selectedLog.response_summary" class="border-t pt-3 dark:border-gray-700">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('org.audit.responseSummary') }}</h4>
          <p class="text-sm text-gray-600 dark:text-gray-400 whitespace-pre-wrap">{{ selectedLog.response_summary }}</p>
        </div>
        <div v-if="selectedLog.keywords && selectedLog.keywords.length" class="border-t pt-3 dark:border-gray-700">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('org.audit.keywords') }}</h4>
          <div class="flex flex-wrap gap-1">
            <span
              v-for="kw in selectedLog.keywords"
              :key="kw"
              class="inline-flex items-center rounded-full bg-blue-100 px-2 py-0.5 text-xs text-blue-700 dark:bg-blue-900/30 dark:text-blue-300"
            >
              {{ kw }}
            </span>
          </div>
        </div>
        <div v-if="selectedLog.flagged" class="border-t pt-3 dark:border-gray-700">
          <div class="flex items-center gap-2 text-red-600">
            <span class="font-medium">{{ t('org.audit.flagged') }}</span>
            <span v-if="selectedLog.flag_reason" class="text-sm">- {{ selectedLog.flag_reason }}</span>
          </div>
        </div>
      </div>
    </Modal>

    <!-- Audit Config Modal -->
    <Modal v-model="showConfigModal" :title="t('org.audit.configTitle')">
      <form @submit.prevent="handleUpdateConfig" class="space-y-4">
        <FormField :label="t('org.audit.auditMode')">
          <select v-model="configForm.audit_mode" class="input">
            <option value="metadata">{{ t('org.audit.modeMetadata') }}</option>
            <option value="summary">{{ t('org.audit.modeSummary') }}</option>
            <option value="full">{{ t('org.audit.modeFull') }}</option>
          </select>
          <p class="text-xs text-gray-500 mt-1">{{ t('org.audit.modeHint') }}</p>
        </FormField>
        <div class="flex justify-end gap-3">
          <button type="button" class="btn btn-ghost" @click="showConfigModal = false">{{ t('common.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">{{ t('common.save') }}</button>
        </div>
      </form>
    </Modal>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { orgAPI } from '@/api/org'
import type { OrgAuditLog } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import Modal from '@/components/common/Modal.vue'
import FormField from '@/components/common/FormField.vue'
import { useToast } from '@/composables/useToast'

const { t } = useI18n()
const toast = useToast()

const logs = ref<OrgAuditLog[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const submitting = ref(false)

const filters = ref({
  action: '',
  start_date: '',
  end_date: ''
})
const flaggedFilter = ref('')

const showDetailModal = ref(false)
const showConfigModal = ref(false)
const selectedLog = ref<OrgAuditLog | null>(null)
const configForm = ref({ audit_mode: 'metadata' })

const columns = computed(() => [
  { key: 'action', label: t('org.audit.action') },
  { key: 'model', label: t('org.audit.model') },
  { key: 'flagged', label: t('org.audit.flagged') },
  { key: 'cost_usd', label: t('org.audit.cost') },
  { key: 'request_summary', label: t('org.audit.requestSummary') },
  { key: 'created_at', label: t('org.audit.time'), sortable: true },
  { key: 'actions', label: t('common.actions'), align: 'right' as const }
])

async function fetchLogs() {
  loading.value = true
  try {
    const f: Record<string, unknown> = {}
    if (filters.value.action) f.action = filters.value.action
    if (flaggedFilter.value) f.flagged = flaggedFilter.value === 'true'
    if (filters.value.start_date) f.start_date = new Date(filters.value.start_date).toISOString()
    if (filters.value.end_date) f.end_date = new Date(filters.value.end_date).toISOString()

    const res = await orgAPI.listAuditLogs(page.value, pageSize.value, f as any)
    logs.value = res.items
    total.value = res.total
  } catch (e: any) {
    toast.error(e.message || 'Failed to load audit logs')
  } finally {
    loading.value = false
  }
}

function viewDetail(log: OrgAuditLog) {
  selectedLog.value = log
  showDetailModal.value = true
}

async function loadAuditConfig() {
  try {
    const config = await orgAPI.getAuditConfig()
    configForm.value.audit_mode = config.audit_mode
  } catch {
    // ignore
  }
}

async function handleUpdateConfig() {
  submitting.value = true
  try {
    await orgAPI.updateAuditConfig(configForm.value.audit_mode)
    toast.success(t('org.audit.configUpdateSuccess'))
    showConfigModal.value = false
  } catch (e: any) {
    toast.error(e.message || 'Update failed')
  } finally {
    submitting.value = false
  }
}

watch(page, fetchLogs)

onMounted(() => {
  fetchLogs()
  loadAuditConfig()
})
</script>
