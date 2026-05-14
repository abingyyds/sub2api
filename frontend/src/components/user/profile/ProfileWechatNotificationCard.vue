<template>
  <div class="rounded-2xl border-2 border-gray-200 bg-white p-6 shadow-soft dark:border-dark-700 dark:bg-dark-900">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-start gap-4">
        <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400">
          <Icon name="chat" size="lg" />
        </div>
        <div>
          <h3 class="font-semibold text-gray-900 dark:text-white">微信公众号通知</h3>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
            {{ description }}
          </p>
          <p v-if="status?.openid_masked" class="mt-2 text-xs text-gray-500 dark:text-dark-400">
            OpenID：{{ status.openid_masked }}
          </p>
        </div>
      </div>

      <div class="flex flex-wrap gap-2">
        <button
          v-if="canBind"
          type="button"
          class="btn btn-primary"
          :disabled="operating"
          @click="bindWechat"
        >
          绑定公众号
        </button>
        <button
          v-if="status?.bound"
          type="button"
          class="btn btn-secondary"
          :disabled="operating"
          @click="unbindWechat"
        >
          解绑
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores'
import { Icon } from '@/components/icons'
import {
  getWechatBindURL,
  getWechatNotificationStatus,
  unbindWechatNotification,
  type WechatNotificationStatus
} from '@/api/notifications'

const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

const status = ref<WechatNotificationStatus | null>(null)
const operating = ref(false)

const canBind = computed(() => Boolean(status.value?.enabled && status.value?.configured && !status.value?.bound))

const description = computed(() => {
  if (!status.value) return '正在读取绑定状态'
  if (!status.value.enabled) return '管理员尚未启用公众号通知'
  if (!status.value.configured) return '公众号通知配置尚未完成'
  if (status.value.bound) return '余额不足、额度包不足和套餐额度不足时会通过公众号提醒'
  return '绑定后可接收余额不足、额度包不足和套餐额度不足提醒'
})

async function loadStatus() {
  try {
    status.value = await getWechatNotificationStatus()
  } catch (error: any) {
    appStore.showError(error?.message || '读取公众号绑定状态失败')
  }
}

async function bindWechat() {
  operating.value = true
  try {
    const url = await getWechatBindURL('/profile')
    window.location.href = url
  } catch (error: any) {
    appStore.showError(error?.message || '获取公众号绑定链接失败')
  } finally {
    operating.value = false
  }
}

async function unbindWechat() {
  operating.value = true
  try {
    await unbindWechatNotification()
    appStore.showSuccess('已解绑微信公众号')
    await loadStatus()
  } catch (error: any) {
    appStore.showError(error?.message || '解绑失败')
  } finally {
    operating.value = false
  }
}

async function clearWechatBindQuery() {
  const query = { ...route.query }
  delete query.wechat_bind
  delete query.reason
  await router.replace({ path: route.path, query })
}

onMounted(async () => {
  const bindResult = route.query.wechat_bind
  if (bindResult === 'success') {
    appStore.showSuccess('微信公众号绑定成功')
    await clearWechatBindQuery()
  } else if (bindResult === 'error') {
    appStore.showError('微信公众号绑定失败')
    await clearWechatBindQuery()
  }
  await loadStatus()
})
</script>
