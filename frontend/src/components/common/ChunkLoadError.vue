<template>
  <div class="flex min-h-[50vh] items-center justify-center px-4 py-12">
    <div
      class="w-full max-w-md rounded-xl border border-gray-200 bg-white p-6 text-center shadow-sm dark:border-dark-700 dark:bg-dark-900"
      role="alert"
      aria-live="assertive"
    >
      <div
        class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300"
      >
        <Icon name="refresh" size="lg" />
      </div>
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">版本已更新</h2>
      <p class="mt-2 text-sm leading-6 text-gray-600 dark:text-gray-300">
        当前页面资源已失效，请刷新页面加载最新版本。
      </p>
      <button class="btn btn-primary mt-6" type="button" @click="reload">
        刷新页面
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'
import { forceReloadForChunkError, getCurrentLocationPath } from '@/router/chunkLoad'

const props = defineProps<{
  path?: string
}>()

const reload = () => {
  forceReloadForChunkError(props.path || getCurrentLocationPath())
}
</script>
