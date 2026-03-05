<template>
  <AppLayout>
    <div class="space-y-6 p-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('org.dashboard.title') }}
        </h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t('org.dashboard.description') }}
        </p>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Skeleton v-for="i in 4" :key="i" class="h-28 rounded-xl" />
      </div>

      <!-- Stats Cards -->
      <div v-else-if="dashboard" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          :title="t('org.dashboard.balance')"
          :value="'$' + Number(dashboard.organization.balance).toFixed(4)"
        />
        <StatCard
          :title="t('org.dashboard.memberCount')"
          :value="String(dashboard.member_count)"
        />
        <StatCard
          :title="t('org.dashboard.todayUsage')"
          :value="'$' + Number(dashboard.today_usage_usd).toFixed(4)"
        />
        <StatCard
          :title="t('org.dashboard.totalUsage')"
          :value="'$' + Number(dashboard.total_usage_usd).toFixed(4)"
        />
      </div>

      <!-- Organization Info -->
      <div v-if="dashboard" class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
        <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
          {{ t('org.dashboard.orgInfo') }}
        </h2>
        <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('admin.organizations.name') }}</dt>
            <dd class="mt-1 font-medium text-gray-900 dark:text-white">{{ dashboard.organization.name }}</dd>
          </div>
          <div>
            <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('admin.organizations.billingMode') }}</dt>
            <dd class="mt-1 font-medium text-gray-900 dark:text-white">{{ dashboard.organization.billing_mode }}</dd>
          </div>
          <div>
            <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('admin.organizations.maxMembers') }}</dt>
            <dd class="mt-1 font-medium text-gray-900 dark:text-white">{{ dashboard.member_count }} / {{ dashboard.organization.max_members }}</dd>
          </div>
          <div>
            <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('admin.organizations.maxApiKeys') }}</dt>
            <dd class="mt-1 font-medium text-gray-900 dark:text-white">{{ dashboard.api_key_count }} / {{ dashboard.organization.max_api_keys }}</dd>
          </div>
        </dl>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { orgAPI } from '@/api/org'
import type { OrgDashboard } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import Skeleton from '@/components/common/Skeleton.vue'
import StatCard from '@/components/common/StatCard.vue'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const dashboard = ref<OrgDashboard | null>(null)

async function loadDashboard() {
  loading.value = true
  try {
    dashboard.value = await orgAPI.getDashboard()
  } catch {
    appStore.showError(t('common.loadError'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboard()
})
</script>
