<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('admin.inviteCodes.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.description') }}</p>
        </div>
        <button @click="showCreateModal = true" class="btn-primary flex items-center gap-2">
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {{ t('admin.inviteCodes.create') }}
        </button>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
        <div class="card p-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.inviteCodes.totalCodes') }}</p>
          <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ codes.length }}</p>
        </div>
        <div class="card p-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.inviteCodes.activeCodes') }}</p>
          <p class="text-2xl font-bold text-green-600 dark:text-green-400">{{ codes.filter(c => c.enabled).length }}</p>
        </div>
        <div class="card p-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.inviteCodes.totalUsed') }}</p>
          <p class="text-2xl font-bold text-blue-600 dark:text-blue-400">{{ codes.reduce((sum, c) => sum + c.used_count, 0) }}</p>
        </div>
        <div class="card p-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.inviteCodes.sources') }}</p>
          <p class="text-2xl font-bold text-purple-600 dark:text-purple-400">{{ new Set(codes.map(c => c.source_name)).size }}</p>
        </div>
      </div>

      <!-- Table -->
      <div class="card overflow-hidden">
        <div v-if="loading" class="flex items-center justify-center py-12">
          <LoadingSpinner />
        </div>
        <table v-else class="min-w-full divide-y divide-gray-200 dark:divide-dark-700">
          <thead class="bg-gray-50 dark:bg-dark-800">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.code') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.sourceName') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.usedCount') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.maxUses') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">{{ t('common.status') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.notes') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">{{ t('common.createdAt') }}</th>
              <th class="px-4 py-3 text-right text-xs font-medium uppercase text-gray-500 dark:text-dark-400">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-dark-700">
            <tr v-if="codes.length === 0">
              <td colspan="8" class="px-4 py-8 text-center text-sm text-gray-500 dark:text-dark-400">
                {{ t('admin.inviteCodes.noCodes') }}
              </td>
            </tr>
            <tr v-for="code in codes" :key="code.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <code class="rounded bg-gray-100 px-2 py-0.5 text-xs font-mono text-gray-800 dark:bg-dark-700 dark:text-dark-200">
                    {{ code.code }}
                  </code>
                  <button @click="copyCode(code.code)" class="text-gray-400 hover:text-gray-600 dark:hover:text-dark-200" :title="t('common.copy')">
                    <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184" />
                    </svg>
                  </button>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="inline-flex items-center rounded-full bg-purple-100 px-2.5 py-0.5 text-xs font-medium text-purple-800 dark:bg-purple-900/30 dark:text-purple-300">
                  {{ code.source_name }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-900 dark:text-white">{{ code.used_count }}</td>
              <td class="px-4 py-3 text-sm text-gray-500 dark:text-dark-400">
                {{ code.max_uses !== null ? code.max_uses : t('common.unlimited') }}
              </td>
              <td class="px-4 py-3">
                <span
                  class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
                  :class="code.enabled
                    ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'
                    : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'"
                >
                  {{ code.enabled ? t('common.enabled') : t('common.disabled') }}
                </span>
              </td>
              <td class="max-w-[200px] truncate px-4 py-3 text-sm text-gray-500 dark:text-dark-400" :title="code.notes">
                {{ code.notes || '-' }}
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dark-400">
                {{ formatDate(code.created_at) }}
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button @click="copyInviteLink(code.code)" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300" :title="t('admin.inviteCodes.copyLink')">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" />
                    </svg>
                  </button>
                  <button @click="openEditModal(code)" class="text-gray-400 hover:text-gray-600 dark:hover:text-dark-200">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Z" />
                    </svg>
                  </button>
                  <button @click="toggleEnabled(code)" class="text-gray-400 hover:text-gray-600 dark:hover:text-dark-200">
                    <svg v-if="code.enabled" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                    </svg>
                    <svg v-else class="h-4 w-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </button>
                  <button @click="confirmDelete(code)" class="text-red-400 hover:text-red-600">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Create/Edit Modal -->
      <Teleport to="body">
        <div v-if="showCreateModal || showEditModal" class="fixed inset-0 z-50 flex items-center justify-center">
          <div class="fixed inset-0 bg-black/50" @click="closeModal"></div>
          <div class="relative w-full max-w-md rounded-xl bg-white p-6 shadow-xl dark:bg-dark-800">
            <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
              {{ showEditModal ? t('admin.inviteCodes.edit') : t('admin.inviteCodes.create') }}
            </h3>

            <div class="space-y-4">
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-dark-300">
                  {{ t('admin.inviteCodes.sourceName') }} *
                </label>
                <input
                  v-model="form.source_name"
                  type="text"
                  class="input w-full"
                  :placeholder="t('admin.inviteCodes.sourceNamePlaceholder')"
                  maxlength="100"
                />
                <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.sourceNameHint') }}</p>
              </div>

              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-dark-300">
                  {{ t('admin.inviteCodes.maxUses') }}
                </label>
                <input
                  v-model.number="form.max_uses"
                  type="number"
                  min="1"
                  class="input w-full"
                  :placeholder="t('admin.inviteCodes.maxUsesPlaceholder')"
                />
                <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.maxUsesHint') }}</p>
              </div>

              <div v-if="showEditModal">
                <label class="flex items-center gap-2">
                  <input v-model="form.enabled" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700" />
                  <span class="text-sm font-medium text-gray-700 dark:text-dark-300">{{ t('common.enabled') }}</span>
                </label>
              </div>

              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-dark-300">
                  {{ t('admin.inviteCodes.notes') }}
                </label>
                <textarea
                  v-model="form.notes"
                  rows="2"
                  class="input w-full"
                  :placeholder="t('admin.inviteCodes.notesPlaceholder')"
                ></textarea>
              </div>
            </div>

            <div class="mt-6 flex justify-end gap-3">
              <button @click="closeModal" class="btn-secondary">{{ t('common.cancel') }}</button>
              <button @click="submitForm" :disabled="!form.source_name || submitting" class="btn-primary">
                {{ submitting ? t('common.saving') : t('common.save') }}
              </button>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Delete Confirmation Modal -->
      <Teleport to="body">
        <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center">
          <div class="fixed inset-0 bg-black/50" @click="showDeleteModal = false"></div>
          <div class="relative w-full max-w-sm rounded-xl bg-white p-6 shadow-xl dark:bg-dark-800">
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">{{ t('common.confirmDelete') }}</h3>
            <p class="mb-4 text-sm text-gray-500 dark:text-dark-400">{{ t('admin.inviteCodes.deleteConfirm') }}</p>
            <div class="flex justify-end gap-3">
              <button @click="showDeleteModal = false" class="btn-secondary">{{ t('common.cancel') }}</button>
              <button @click="executeDelete" :disabled="submitting" class="btn-danger">
                {{ submitting ? t('common.deleting') : t('common.delete') }}
              </button>
            </div>
          </div>
        </div>
      </Teleport>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { adminAPI } from '@/api/admin'
import { useClipboard } from '@/composables/useClipboard'
import type { AdminInviteCode } from '@/api/admin/inviteCodes'

const { t } = useI18n()
const appStore = useAppStore()
const { copyToClipboard } = useClipboard()

const loading = ref(false)
const submitting = ref(false)
const codes = ref<AdminInviteCode[]>([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const editingCode = ref<AdminInviteCode | null>(null)
const deletingCode = ref<AdminInviteCode | null>(null)

const form = ref({
  source_name: '',
  max_uses: null as number | null,
  enabled: true,
  notes: ''
})

const loadCodes = async () => {
  loading.value = true
  try {
    const resp = await adminAPI.inviteCodes.list(1, 200)
    codes.value = resp.items || []
  } catch (error) {
    appStore.showError(t('admin.inviteCodes.loadFailed'))
    console.error('Error loading invite codes:', error)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  form.value = { source_name: '', max_uses: null, enabled: true, notes: '' }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingCode.value = null
  resetForm()
}

const openEditModal = (code: AdminInviteCode) => {
  editingCode.value = code
  form.value = {
    source_name: code.source_name,
    max_uses: code.max_uses,
    enabled: code.enabled,
    notes: code.notes || ''
  }
  showEditModal.value = true
}

const submitForm = async () => {
  if (!form.value.source_name) return
  submitting.value = true
  try {
    if (showEditModal.value && editingCode.value) {
      await adminAPI.inviteCodes.update(editingCode.value.id, {
        source_name: form.value.source_name,
        max_uses: form.value.max_uses,
        enabled: form.value.enabled,
        notes: form.value.notes
      })
      appStore.showSuccess(t('admin.inviteCodes.updateSuccess'))
    } else {
      await adminAPI.inviteCodes.create({
        source_name: form.value.source_name,
        max_uses: form.value.max_uses,
        notes: form.value.notes
      })
      appStore.showSuccess(t('admin.inviteCodes.createSuccess'))
    }
    closeModal()
    await loadCodes()
  } catch (error) {
    appStore.showError(showEditModal.value ? t('admin.inviteCodes.updateFailed') : t('admin.inviteCodes.createFailed'))
    console.error('Error saving invite code:', error)
  } finally {
    submitting.value = false
  }
}

const confirmDelete = (code: AdminInviteCode) => {
  deletingCode.value = code
  showDeleteModal.value = true
}

const executeDelete = async () => {
  if (!deletingCode.value) return
  submitting.value = true
  try {
    await adminAPI.inviteCodes.remove(deletingCode.value.id)
    appStore.showSuccess(t('admin.inviteCodes.deleteSuccess'))
    showDeleteModal.value = false
    deletingCode.value = null
    await loadCodes()
  } catch (error) {
    appStore.showError(t('admin.inviteCodes.deleteFailed'))
    console.error('Error deleting invite code:', error)
  } finally {
    submitting.value = false
  }
}

const toggleEnabled = async (code: AdminInviteCode) => {
  try {
    await adminAPI.inviteCodes.update(code.id, { enabled: !code.enabled })
    code.enabled = !code.enabled
    appStore.showSuccess(code.enabled ? t('admin.inviteCodes.enabled') : t('admin.inviteCodes.disabled'))
  } catch (error) {
    appStore.showError(t('admin.inviteCodes.updateFailed'))
  }
}

const copyCode = (code: string) => {
  copyToClipboard(code, t('common.copied'))
}

const copyInviteLink = (code: string) => {
  const link = `${window.location.origin}/register?invite=${code}`
  copyToClipboard(link, t('admin.inviteCodes.linkCopied'))
}

const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

onMounted(loadCodes)
</script>
