<template>
  <AppLayout>
    <div class="flex min-h-[50vh] items-center justify-center text-sm text-gray-500 dark:text-dark-400">
      {{ loading ? '加载中...' : message }}
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAuthStore } from '@/stores'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const message = ref('你还没有自己的分站')

onMounted(async () => {
  try {
    if (authStore.ownedSites.length === 0) {
      await authStore.refreshOwnedSites()
    }
    const first = authStore.ownedSites[0]
    if (first) {
      await router.replace(`/subsite-admin/${first.id}/dashboard`)
      return
    }
  } finally {
    loading.value = false
  }
})
</script>
