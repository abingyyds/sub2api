<template>
  <BaseDialog
    :show="show"
    :title="isEditMode ? t('keys.editKey') : t('keys.createKey')"
    width="normal"
    @close="handleClose"
  >
    <form id="key-form" class="space-y-5" @submit.prevent="handleSubmit">
      <div
        v-if="!isEditMode && preSelectedGroup"
        class="rounded-xl border border-gray-200 bg-gray-50 p-4 dark:border-dark-600 dark:bg-dark-800/50"
      >
        <p class="mb-2 text-xs font-medium text-gray-500 dark:text-gray-400">
          {{ t('keys.groupInfo') }}
        </p>
        <div class="flex items-center gap-3">
          <GroupBadge
            :name="preSelectedGroup.name"
            :platform="preSelectedGroup.platform"
            :subscription-type="preSelectedGroup.subscription_type"
          />
        </div>
        <p v-if="preSelectedGroup.description" class="mt-2 text-xs text-gray-500 dark:text-gray-400">
          {{ preSelectedGroup.description }}
        </p>
      </div>

      <div>
        <label class="input-label">{{ t('keys.nameLabel') }}</label>
        <input
          v-model="formData.name"
          type="text"
          class="input"
          :placeholder="t('keys.namePlaceholder')"
          data-tour="key-form-name"
        />
      </div>

      <div>
        <label class="input-label">{{ t('keys.groupLabel') }}</label>
        <Select
          v-model="formData.group_id"
          :options="groupOptions"
          :placeholder="t('keys.selectGroup')"
          data-tour="key-form-group"
        >
          <template #selected="{ option }">
            <GroupBadge
              v-if="option"
              :name="(option as unknown as GroupOption).label"
              :platform="(option as unknown as GroupOption).platform"
              :subscription-type="(option as unknown as GroupOption).subscriptionType"
            />
            <span v-else class="text-gray-400">{{ t('keys.selectGroup') }}</span>
          </template>
          <template #option="{ option, selected }">
            <GroupOptionItem
              :name="(option as unknown as GroupOption).label"
              :platform="(option as unknown as GroupOption).platform"
              :subscription-type="(option as unknown as GroupOption).subscriptionType"
              :description="(option as unknown as GroupOption).description"
              :selected="selected"
            />
          </template>
        </Select>
      </div>

      <div v-if="!isEditMode" class="space-y-3">
        <div class="flex items-center justify-between">
          <label class="input-label mb-0">{{ t('keys.customKeyLabel') }}</label>
          <button
            type="button"
            :class="[
              'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
              formData.use_custom_key ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
            ]"
            @click="formData.use_custom_key = !formData.use_custom_key"
          >
            <span
              :class="[
                'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                formData.use_custom_key ? 'translate-x-4' : 'translate-x-0'
              ]"
            />
          </button>
        </div>
        <div v-if="formData.use_custom_key">
          <input
            v-model="formData.custom_key"
            type="text"
            class="input font-mono"
            :placeholder="t('keys.customKeyPlaceholder')"
            :class="{ 'border-red-500 dark:border-red-500': customKeyError }"
          />
          <p v-if="customKeyError" class="mt-1 text-sm text-red-500">{{ customKeyError }}</p>
          <p v-else class="input-hint">{{ t('keys.customKeyHint') }}</p>
        </div>
      </div>

      <div v-if="isEditMode">
        <label class="input-label">{{ t('keys.statusLabel') }}</label>
        <Select
          v-model="formData.status"
          :options="statusOptions"
          :placeholder="t('keys.selectStatus')"
        />
      </div>

      <div>
        <label class="input-label">{{ t('keys.usageLimitLabel') }}</label>
        <input
          v-model.number="formData.usage_limit"
          type="number"
          min="0"
          step="0.01"
          class="input"
          :placeholder="t('keys.usageLimitPlaceholder')"
        />
        <p class="input-hint">{{ t('keys.usageLimitHint') }}</p>
      </div>

      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <label class="input-label mb-0">{{ t('keys.ipRestriction') }}</label>
          <button
            type="button"
            :class="[
              'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
              formData.enable_ip_restriction ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
            ]"
            @click="formData.enable_ip_restriction = !formData.enable_ip_restriction"
          >
            <span
              :class="[
                'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                formData.enable_ip_restriction ? 'translate-x-4' : 'translate-x-0'
              ]"
            />
          </button>
        </div>

        <div v-if="formData.enable_ip_restriction" class="space-y-4 pt-2">
          <div>
            <label class="input-label">{{ t('keys.ipWhitelist') }}</label>
            <textarea
              v-model="formData.ip_whitelist"
              rows="3"
              class="input font-mono text-sm"
              :placeholder="t('keys.ipWhitelistPlaceholder')"
            />
            <p class="input-hint">{{ t('keys.ipWhitelistHint') }}</p>
          </div>

          <div>
            <label class="input-label">{{ t('keys.ipBlacklist') }}</label>
            <textarea
              v-model="formData.ip_blacklist"
              rows="3"
              class="input font-mono text-sm"
              :placeholder="t('keys.ipBlacklistPlaceholder')"
            />
            <p class="input-hint">{{ t('keys.ipBlacklistHint') }}</p>
          </div>
        </div>
      </div>

      <label
        v-if="!isEditMode"
        class="flex cursor-pointer items-start gap-3 rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-700 dark:border-dark-700 dark:bg-dark-900/60 dark:text-dark-200"
      >
        <input
          v-model="formData.legal_accepted"
          type="checkbox"
          class="mt-0.5 h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600"
        />
        <span>{{ t('keys.legal.commitmentLabel') }}</span>
      </label>
    </form>

    <template #footer>
      <div class="flex justify-end gap-3">
        <button type="button" class="btn btn-secondary" @click="handleClose">
          {{ t('common.cancel') }}
        </button>
        <button
          form="key-form"
          type="submit"
          :disabled="submitting || (!isEditMode && !formData.legal_accepted)"
          class="btn btn-primary"
          data-tour="key-form-submit"
        >
          <svg
            v-if="submitting"
            class="-ml-1 mr-2 h-4 w-4 animate-spin"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          {{
            submitting
              ? t('keys.saving')
              : isEditMode
                ? t('common.update')
                : t('common.create')
          }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { keysAPI } from '@/api'
import BaseDialog from '@/components/common/BaseDialog.vue'
import GroupBadge from '@/components/common/GroupBadge.vue'
import GroupOptionItem from '@/components/common/GroupOptionItem.vue'
import Select from '@/components/common/Select.vue'
import { useAppStore } from '@/stores/app'
import { useOnboardingStore } from '@/stores/onboarding'
import type { ApiKey, Group, GroupPlatform, SelectOption, SubscriptionType } from '@/types'

interface Props {
  show: boolean
  groups: Group[]
  selectedKey: ApiKey | null
  initialGroupId?: number | null
}

interface Emits {
  (e: 'close'): void
  (e: 'saved', payload: { created: boolean }): void
}

interface GroupOption extends SelectOption {
  value: number | null
  label: string
  description: string | null
  subscriptionType: SubscriptionType
  platform: GroupPlatform
}

interface KeyFormData {
  name: string
  group_id: number | null
  status: 'active' | 'inactive'
  use_custom_key: boolean
  custom_key: string
  enable_ip_restriction: boolean
  ip_whitelist: string
  ip_blacklist: string
  usage_limit: number | null
  legal_accepted: boolean
}

const props = withDefaults(defineProps<Props>(), {
  initialGroupId: null
})

const emit = defineEmits<Emits>()

const { t } = useI18n()
const appStore = useAppStore()
const onboardingStore = useOnboardingStore()

const isEditMode = computed(() => props.selectedKey !== null)
const submitting = ref(false)

const createEmptyFormData = (groupId: number | null = null): KeyFormData => ({
  name: '',
  group_id: groupId,
  status: 'active',
  use_custom_key: false,
  custom_key: '',
  enable_ip_restriction: false,
  ip_whitelist: '',
  ip_blacklist: '',
  usage_limit: null,
  legal_accepted: false
})

const createFormDataFromKey = (key: ApiKey): KeyFormData => {
  const hasIPRestriction = key.ip_whitelist.length > 0 || key.ip_blacklist.length > 0

  return {
    name: key.name,
    group_id: key.group_id,
    status: key.status,
    use_custom_key: false,
    custom_key: '',
    enable_ip_restriction: hasIPRestriction,
    ip_whitelist: key.ip_whitelist.join('\n'),
    ip_blacklist: key.ip_blacklist.join('\n'),
    usage_limit: key.usage_limit,
    legal_accepted: true
  }
}

const formData = ref<KeyFormData>(createEmptyFormData())

const syncFormFromProps = () => {
  formData.value = props.selectedKey
    ? createFormDataFromKey(props.selectedKey)
    : createEmptyFormData(props.initialGroupId)
}

watch(
  () => [props.show, props.selectedKey, props.initialGroupId],
  ([show]) => {
    if (show) {
      syncFormFromProps()
      return
    }

    if (!submitting.value) {
      formData.value = createEmptyFormData()
    }
  },
  { immediate: true }
)

const customKeyError = computed(() => {
  if (!formData.value.use_custom_key || !formData.value.custom_key) {
    return ''
  }

  const key = formData.value.custom_key
  if (key.length < 16) {
    return t('keys.customKeyTooShort')
  }

  if (!/^[a-zA-Z0-9_-]+$/.test(key)) {
    return t('keys.customKeyInvalidChars')
  }

  return ''
})

const statusOptions = computed(() => [
  { value: 'active', label: t('common.active') },
  { value: 'inactive', label: t('common.inactive') }
])

const groupOptions = computed<GroupOption[]>(() =>
  props.groups.map((group) => ({
    value: group.id,
    label: group.name,
    description: group.description,
    subscriptionType: group.subscription_type,
    platform: group.platform
  }))
)

const preSelectedGroup = computed(() => {
  if (isEditMode.value || formData.value.group_id === null) {
    return null
  }

  return props.groups.find((group) => group.id === formData.value.group_id) || null
})

const parseIPList = (text: string): string[] =>
  text
    .split('\n')
    .map((ip) => ip.trim())
    .filter((ip) => ip.length > 0)

const handleClose = () => {
  if (!submitting.value) {
    emit('close')
  }
}

const handleSubmit = async () => {
  if (formData.value.group_id === null) {
    appStore.showError(t('keys.groupRequired'))
    return
  }

  if (!isEditMode.value && formData.value.use_custom_key) {
    if (!formData.value.custom_key) {
      appStore.showError(t('keys.customKeyRequired'))
      return
    }

    if (customKeyError.value) {
      appStore.showError(customKeyError.value)
      return
    }
  }

  if (!isEditMode.value && !formData.value.legal_accepted) {
    appStore.showError(t('keys.legal.required'))
    return
  }

  const ipWhitelist = formData.value.enable_ip_restriction
    ? parseIPList(formData.value.ip_whitelist)
    : []
  const ipBlacklist = formData.value.enable_ip_restriction
    ? parseIPList(formData.value.ip_blacklist)
    : []

  submitting.value = true

  try {
    if (isEditMode.value && props.selectedKey) {
      await keysAPI.update(props.selectedKey.id, {
        name: formData.value.name,
        group_id: formData.value.group_id,
        status: formData.value.status,
        ip_whitelist: ipWhitelist,
        ip_blacklist: ipBlacklist,
        usage_limit: formData.value.usage_limit
      })
      appStore.showSuccess(t('keys.keyUpdatedSuccess'))
      emit('saved', { created: false })
    } else {
      const customKey = formData.value.use_custom_key ? formData.value.custom_key : undefined
      await keysAPI.create(
        formData.value.name,
        formData.value.group_id,
        customKey,
        ipWhitelist,
        ipBlacklist,
        formData.value.usage_limit,
        formData.value.legal_accepted
      )
      appStore.showSuccess(t('keys.keyCreatedSuccess'))

      if (onboardingStore.isCurrentStep('[data-tour="key-form-submit"]')) {
        onboardingStore.nextStep(500)
      }

      emit('saved', { created: true })
    }

    emit('close')
  } catch (error: any) {
    const errorMsg = error.response?.data?.detail || t('keys.failedToSave')
    appStore.showError(errorMsg)
  } finally {
    submitting.value = false
  }
}
</script>
