<template>
  <BaseDialog
    :show="show"
    :title="t('admin.users.editUser')"
    width="normal"
    @close="$emit('close')"
  >
    <form v-if="user" id="edit-user-form" @submit.prevent="handleUpdateUser" class="space-y-5">
      <div>
        <label class="input-label">{{ t('admin.users.email') }}</label>
        <input v-model="form.email" type="email" class="input" />
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.password') }}</label>
        <div class="flex gap-2">
          <div class="relative flex-1">
            <input v-model="form.password" type="text" class="input pr-10" :placeholder="t('admin.users.enterNewPassword')" />
            <button v-if="form.password" type="button" @click="copyPassword" class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg p-1 transition-colors hover:bg-gray-100 dark:hover:bg-dark-700" :class="passwordCopied ? 'text-green-500' : 'text-gray-400'">
              <svg v-if="passwordCopied" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
              <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184" /></svg>
            </button>
          </div>
          <button type="button" @click="generatePassword" class="btn btn-secondary px-3">
            <Icon name="refresh" size="md" />
          </button>
        </div>
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.username') }}</label>
        <input v-model="form.username" type="text" class="input" />
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.notes') }}</label>
        <textarea v-model="form.notes" rows="3" class="input"></textarea>
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.columns.concurrency') }}</label>
        <input v-model.number="form.concurrency" type="number" class="input" />
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.inviter') }}</label>
        <input
          v-model.number="form.inviter_id"
          type="number"
          class="input"
          :placeholder="t('admin.users.inviterIdPlaceholder')"
        />
        <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.users.inviterIdHint') }}</p>
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.boundSubSite') }}</label>
        <select v-model.number="form.sub_site_id" class="input">
          <option :value="0">{{ t('admin.users.mainSite') }}</option>
          <option v-for="site in subSites" :key="site.id" :value="site.id">
            #{{ site.id }} · {{ site.name }}{{ site.custom_domain ? ` (${site.custom_domain})` : ` (${site.slug})` }}
          </option>
        </select>
        <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.users.boundSubSiteHint') }}</p>
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.columns.discoverySource') }}</label>
        <div class="flex gap-2">
          <select v-model="form.discovery_source" class="input flex-1">
            <option value="">-</option>
            <option v-for="opt in discoverySourceOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
          <input
            v-if="form.discovery_source === '__custom__'"
            v-model="form.custom_discovery_source"
            type="text"
            class="input flex-1"
            :placeholder="t('admin.users.customSource')"
          />
        </div>
      </div>
      <div v-if="authStore.isFullAdmin">
        <label class="input-label">{{ t('admin.users.columns.role') }}</label>
        <select v-model="form.role" class="input">
          <option value="user">{{ t('admin.users.roles.user') }}</option>
          <option value="sub_admin">{{ t('admin.users.roles.sub_admin') }}</option>
          <option value="admin">{{ t('admin.users.roles.admin') }}</option>
        </select>
      </div>
      <UserAttributeForm v-model="form.customAttributes" :user-id="user?.id" />
    </form>
    <template #footer>
      <div class="flex justify-end gap-3">
        <button @click="$emit('close')" type="button" class="btn btn-secondary">{{ t('common.cancel') }}</button>
        <button type="submit" form="edit-user-form" :disabled="submitting" class="btn btn-primary">
          {{ submitting ? t('admin.users.updating') : t('common.update') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useClipboard } from '@/composables/useClipboard'
import { adminAPI } from '@/api/admin'
import type { AdminUser, UserAttributeValuesMap } from '@/types'
import type { AdminSubSite } from '@/api/admin/subsites'
import BaseDialog from '@/components/common/BaseDialog.vue'
import UserAttributeForm from '@/components/user/UserAttributeForm.vue'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{ show: boolean, user: AdminUser | null }>()
const emit = defineEmits(['close', 'success'])
const { t } = useI18n(); const appStore = useAppStore(); const authStore = useAuthStore(); const { copyToClipboard } = useClipboard()

const submitting = ref(false); const passwordCopied = ref(false)
const subSites = ref<AdminSubSite[]>([])
const form = reactive({ email: '', password: '', username: '', notes: '', concurrency: 1, role: 'user' as string, discovery_source: '' as string, custom_discovery_source: '', customAttributes: {} as UserAttributeValuesMap, inviter_id: null as number | null, sub_site_id: 0 })

const loadSubSites = async () => {
  if (subSites.value.length > 0) return
  const res = await adminAPI.subSites.list(1, 1000)
  subSites.value = res.items
}

const discoverySourceOptions = computed(() => [
  { value: 'douyin', label: t('auth.discoverySource.douyin') },
  { value: 'xiaohongshu', label: t('auth.discoverySource.xiaohongshu') },
  { value: 'bilibili', label: t('auth.discoverySource.bilibili') },
  { value: 'wechat', label: t('auth.discoverySource.wechat') },
  { value: 'twitter', label: t('auth.discoverySource.twitter') },
  { value: 'telegram', label: t('auth.discoverySource.telegram') },
  { value: 'github', label: t('auth.discoverySource.github') },
  { value: 'search', label: t('auth.discoverySource.search') },
  { value: 'friend', label: t('auth.discoverySource.friend') },
  { value: 'other', label: t('auth.discoverySource.other') },
  { value: '__custom__', label: t('admin.users.customSourceOption') }
])

watch(() => props.user, (u) => {
  if (u) {
    const knownSources = ['douyin', 'xiaohongshu', 'bilibili', 'wechat', 'twitter', 'telegram', 'github', 'search', 'friend', 'other', 'skip']
    const src = u.discovery_source || ''
    const isKnown = !src || src === 'skip' || knownSources.includes(src)
    Object.assign(form, {
      email: u.email, password: '', username: u.username || '', notes: u.notes || '',
      concurrency: u.concurrency, role: u.role || 'user',
      discovery_source: isKnown ? (src === 'skip' ? '' : src) : '__custom__',
      custom_discovery_source: isKnown ? '' : src,
      customAttributes: {},
      inviter_id: null,
      sub_site_id: u.bound_sub_site?.id || 0
    })
    passwordCopied.value = false
    // Load current inviter
    if (u.id) {
      adminAPI.users.getUserInviter(u.id).then(inviter => {
        if (inviter) form.inviter_id = inviter.id
      }).catch(() => {})
    }
    loadSubSites().catch(() => {})
  }
}, { immediate: true })

const generatePassword = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$%^&*'
  let p = ''; for (let i = 0; i < 16; i++) p += chars.charAt(Math.floor(Math.random() * chars.length))
  form.password = p
}
const copyPassword = async () => {
  if (form.password && await copyToClipboard(form.password, t('admin.users.passwordCopied'))) {
    passwordCopied.value = true; setTimeout(() => passwordCopied.value = false, 2000)
  }
}
const handleUpdateUser = async () => {
  if (!props.user) return
  if (!form.email.trim()) {
    appStore.showError(t('admin.users.emailRequired'))
    return
  }
  if (form.concurrency < 1) {
    appStore.showError(t('admin.users.concurrencyMin'))
    return
  }
  submitting.value = true
  try {
    const data: any = { email: form.email, username: form.username, notes: form.notes, concurrency: form.concurrency }
    if (form.password.trim()) data.password = form.password.trim()
    if (authStore.isFullAdmin && form.role) data.role = form.role
    // Set discovery_source: custom text if __custom__, selected value otherwise, null to clear
    const resolvedSource = form.discovery_source === '__custom__' ? (form.custom_discovery_source.trim() || null) : (form.discovery_source || null)
    data.discovery_source = resolvedSource
    // Set inviter_id: null to remove, number to set/update
    if (form.inviter_id !== undefined) {
      data.inviter_id = form.inviter_id === null ? 0 : form.inviter_id
    }
    data.sub_site_id = form.sub_site_id || 0
    await adminAPI.users.update(props.user.id, data)
    if (Object.keys(form.customAttributes).length > 0) await adminAPI.userAttributes.updateUserAttributeValues(props.user.id, form.customAttributes)
    appStore.showSuccess(t('admin.users.userUpdated'))
    emit('success'); emit('close')
  } catch (e: any) {
    appStore.showError(e.response?.data?.detail || t('admin.users.failedToUpdate'))
  } finally { submitting.value = false }
}
</script>
