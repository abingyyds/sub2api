<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <SearchInput
            v-model="searchQuery"
            :placeholder="t('org.projects.searchPlaceholder')"
            class="w-64"
          />
          <button class="btn btn-primary" @click="showCreateModal = true">
            {{ t('org.projects.create') }}
          </button>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="projects"
          :loading="loading"
          :empty-text="t('org.projects.empty')"
        >
          <template #cell-name="{ row }">
            <div>
              <div class="font-medium text-gray-900 dark:text-white">{{ row.name }}</div>
              <div v-if="row.description" class="text-xs text-gray-500 truncate max-w-xs">{{ row.description }}</div>
            </div>
          </template>

          <template #cell-allowed_models="{ value }">
            <div v-if="value && value.length" class="flex flex-wrap gap-1">
              <span
                v-for="model in value.slice(0, 3)"
                :key="model"
                class="inline-flex items-center rounded-full bg-blue-100 px-2 py-0.5 text-xs text-blue-700 dark:bg-blue-900/30 dark:text-blue-300"
              >
                {{ model }}
              </span>
              <span v-if="value.length > 3" class="text-xs text-gray-400">+{{ value.length - 3 }}</span>
            </div>
            <span v-else class="text-gray-400 text-xs">{{ t('org.projects.allModels') }}</span>
          </template>

          <template #cell-monthly_budget_usd="{ row }">
            <span class="font-mono text-sm">
              ${{ Number(row.monthly_usage_usd).toFixed(4) }}
              <span class="text-gray-400"> / {{ row.monthly_budget_usd != null ? '$' + Number(row.monthly_budget_usd).toFixed(2) : '∞' }}</span>
            </span>
          </template>

          <template #cell-status="{ value }">
            <span class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
              :class="value === 'active' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300' : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'"
            >
              {{ value === 'active' ? t('common.active') : t('common.disabled') }}
            </span>
          </template>

          <template #cell-actions="{ row }">
            <div class="flex items-center gap-2">
              <button class="btn btn-xs btn-ghost" @click="editProject(row)">{{ t('common.edit') }}</button>
              <button class="btn btn-xs btn-ghost text-red-500" @click="confirmDelete(row)">{{ t('common.delete') }}</button>
            </div>
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

    <!-- Create/Edit Modal -->
    <BaseDialog :show="showCreateModal" :title="editingProject ? t('org.projects.edit') : t('org.projects.create')" @close="showCreateModal = false">
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <FormField :label="t('org.projects.name')" required>
          <input v-model="form.name" type="text" class="input" required />
        </FormField>
        <FormField :label="t('org.projects.description')">
          <textarea v-model="form.description" class="input" rows="2" />
        </FormField>
        <FormField :label="t('org.projects.allowedModels')">
          <input v-model="form.allowed_models_text" type="text" class="input" :placeholder="t('org.projects.allowedModelsPlaceholder')" />
          <p class="text-xs text-gray-500 mt-1">{{ t('org.projects.allowedModelsHint') }}</p>
        </FormField>
        <FormField :label="t('org.projects.monthlyBudget')">
          <input v-model.number="form.monthly_budget_usd" type="number" step="0.01" class="input" :placeholder="t('org.projects.budgetPlaceholder')" />
        </FormField>
        <div v-if="editingProject">
          <FormField :label="t('common.status')">
            <select v-model="form.status" class="input">
              <option value="active">{{ t('common.active') }}</option>
              <option value="disabled">{{ t('common.disabled') }}</option>
            </select>
          </FormField>
        </div>
        <div class="flex justify-end gap-3">
          <button type="button" class="btn btn-ghost" @click="showCreateModal = false">{{ t('common.cancel') }}</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">{{ t('common.save') }}</button>
        </div>
      </form>
    </BaseDialog>

    <!-- Delete Confirmation -->
    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="t('org.projects.deleteTitle')"
      :message="t('org.projects.deleteMessage', { name: deletingProject?.name ?? '' })"
      @confirm="handleDelete"
      @cancel="showDeleteConfirm = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { orgAPI } from '@/api/org'
import type { OrgProject } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import SearchInput from '@/components/common/SearchInput.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import FormField from '@/components/common/FormField.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const appStore = useAppStore()

const projects = ref<OrgProject[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const searchQuery = ref('')

const showCreateModal = ref(false)
const showDeleteConfirm = ref(false)
const editingProject = ref<OrgProject | null>(null)
const deletingProject = ref<OrgProject | null>(null)
const submitting = ref(false)

const form = ref({
  name: '',
  description: '',
  allowed_models_text: '',
  monthly_budget_usd: null as number | null,
  status: 'active' as 'active' | 'disabled'
})

const columns = computed(() => [
  { key: 'name', label: t('org.projects.name'), sortable: true },
  { key: 'allowed_models', label: t('org.projects.allowedModels') },
  { key: 'monthly_budget_usd', label: t('org.projects.monthlyBudget') },
  { key: 'status', label: t('common.status') },
  { key: 'actions', label: t('common.actions'), align: 'right' as const }
])

async function fetchProjects() {
  loading.value = true
  try {
    const res = await orgAPI.listProjects(page.value, pageSize.value)
    projects.value = res.items
    total.value = res.total
  } catch (e: any) {
    appStore.showError(e.message || 'Failed to load projects')
  } finally {
    loading.value = false
  }
}

function editProject(project: OrgProject) {
  editingProject.value = project
  form.value = {
    name: project.name,
    description: project.description || '',
    allowed_models_text: (project.allowed_models || []).join(', '),
    monthly_budget_usd: project.monthly_budget_usd,
    status: project.status
  }
  showCreateModal.value = true
}

function confirmDelete(project: OrgProject) {
  deletingProject.value = project
  showDeleteConfirm.value = true
}

async function handleSubmit() {
  submitting.value = true
  try {
    const allowedModels = form.value.allowed_models_text
      ? form.value.allowed_models_text.split(',').map(s => s.trim()).filter(Boolean)
      : []

    if (editingProject.value) {
      await orgAPI.updateProject(editingProject.value.id, {
        name: form.value.name,
        description: form.value.description || undefined,
        allowed_models: allowedModels,
        monthly_budget_usd: form.value.monthly_budget_usd,
        status: form.value.status
      })
      appStore.showSuccess(t('org.projects.updateSuccess'))
    } else {
      await orgAPI.createProject({
        name: form.value.name,
        description: form.value.description || undefined,
        allowed_models: allowedModels,
        monthly_budget_usd: form.value.monthly_budget_usd
      })
      appStore.showSuccess(t('org.projects.createSuccess'))
    }
    showCreateModal.value = false
    resetForm()
    fetchProjects()
  } catch (e: any) {
    appStore.showError(e.message || 'Operation failed')
  } finally {
    submitting.value = false
  }
}

async function handleDelete() {
  if (!deletingProject.value) return
  try {
    await orgAPI.removeProject(deletingProject.value.id)
    appStore.showSuccess(t('org.projects.deleteSuccess'))
    fetchProjects()
  } catch (e: any) {
    appStore.showError(e.message || 'Delete failed')
  }
}

function resetForm() {
  editingProject.value = null
  form.value = { name: '', description: '', allowed_models_text: '', monthly_budget_usd: null, status: 'active' }
}

watch(showCreateModal, (v) => { if (!v) resetForm() })
watch(page, fetchProjects)

onMounted(fetchProjects)
</script>
