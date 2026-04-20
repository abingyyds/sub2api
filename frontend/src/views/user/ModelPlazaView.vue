<template>
  <div class="mx-auto max-w-6xl space-y-8">
    <div class="text-center">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('modelPlaza.title') }}</h1>
      <p class="mt-2 text-gray-500 dark:text-dark-400">{{ t('modelPlaza.subtitle') }}</p>
    </div>

    <section
      id="pricing-standard"
      class="scroll-mt-24 rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900"
    >
      <div class="flex flex-col gap-5 xl:flex-row xl:items-start xl:justify-between">
        <div class="max-w-3xl">
          <div class="flex flex-wrap items-center gap-2">
            <span class="inline-flex items-center rounded-full bg-primary-100 px-3 py-1 text-xs font-semibold tracking-wide text-primary-700 dark:bg-primary-900/50 dark:text-primary-300">
              {{ t('modelPlaza.pricingStandardBadge') }}
            </span>
            <span
              v-if="pricingGroups.length"
              class="inline-flex rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-300"
            >
              {{ t('modelPlaza.pricingGroupsCount', { count: pricingGroups.length }) }}
            </span>
          </div>

          <h2 class="mt-4 text-2xl font-semibold text-gray-900 dark:text-white">{{ t('modelPlaza.pricingTableTitle') }}</h2>
          <p class="mt-2 text-sm leading-6 text-gray-600 dark:text-dark-300">
            {{ t('modelPlaza.pricingTableDesc') }}
          </p>

          <div class="mt-4 flex flex-wrap gap-2">
            <span class="inline-flex rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-300">
              {{ t('modelPlaza.officialUnitHint') }}
            </span>
            <span class="inline-flex rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-300">
              {{ t('modelPlaza.siteUnitHint') }}
            </span>
            <span class="inline-flex rounded-full bg-primary-50 px-3 py-1 text-xs font-medium text-primary-700 dark:bg-primary-950/30 dark:text-primary-300">
              {{ t('modelPlaza.pricingAvailableOnly') }}
            </span>
          </div>

          <p v-if="pricingGroups.length" class="mt-3 text-xs leading-5 text-gray-500 dark:text-dark-400">
            {{ t('modelPlaza.pricingCompactHint') }}
          </p>
        </div>

        <div class="w-full max-w-sm space-y-3">
          <router-link
            to="/pricing"
            class="inline-flex items-center justify-center rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300"
          >
            {{ t('pricing.title') }}
          </router-link>

          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-dark-300">
              {{ t('common.search') }}
            </label>
            <input
              v-model="pricingSearch"
              type="text"
              class="input w-full"
              :placeholder="t('modelPlaza.pricingSearchPlaceholder')"
            />
          </div>

          <p class="text-xs leading-5 text-gray-500 dark:text-dark-400">
            {{ t('modelPlaza.pricingResults', { count: filteredPricingRows.length, total: pricingRows.length }) }}
          </p>
        </div>
      </div>

      <div v-if="pricingGroups.length" class="mt-6 space-y-3">
        <div v-if="filteredPricingRows.length" class="space-y-3">
          <article
            v-for="row in paginatedPricingRows"
            :key="row.model"
            class="rounded-2xl border border-gray-200 bg-gray-50/70 p-4 dark:border-dark-700 dark:bg-dark-800/70"
          >
            <div class="flex flex-col gap-4 xl:flex-row">
              <div class="xl:w-72 xl:flex-shrink-0">
                <div class="flex flex-wrap items-center gap-2">
                  <span class="inline-flex rounded-full px-2.5 py-1 text-xs font-medium" :class="platformBadgeClass(row.platform)">
                    {{ t(`admin.groups.platforms.${row.platform}`) }}
                  </span>
                </div>

                <h3 class="mt-3 break-all font-mono text-sm font-semibold leading-6 text-gray-900 dark:text-white">
                  {{ row.model }}
                </h3>

                <p v-if="row.aliases?.length" class="mt-2 text-xs leading-5 text-gray-500 dark:text-dark-400">
                  {{ t('modelPlaza.aliasesLabel') }}: {{ aliasSummary(row.aliases) }}
                </p>
              </div>

              <div class="flex-1 space-y-3">
                <div class="rounded-2xl border border-white/80 bg-white p-4 dark:border-dark-700 dark:bg-dark-900/80">
                  <div class="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
                    <div class="text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">
                      {{ t('modelPlaza.pricingColOfficial') }}
                    </div>
                    <div class="text-[11px] text-gray-400 dark:text-dark-500">
                      {{ t('modelPlaza.officialUnitHint') }}
                    </div>
                  </div>

                  <div class="mt-3 flex flex-wrap gap-x-4 gap-y-2">
                    <div
                      v-for="metric in metricRows(row.official, 'usd')"
                      :key="metric.key"
                      class="inline-flex items-center gap-2 rounded-full bg-gray-100 px-3 py-1.5 text-xs dark:bg-dark-800"
                    >
                      <span class="text-gray-500 dark:text-dark-400">{{ metric.label }}</span>
                      <span class="font-mono font-semibold text-gray-900 dark:text-white">{{ metric.value }}</span>
                    </div>
                  </div>
                </div>

                <div class="rounded-2xl border border-white/80 bg-white p-4 dark:border-dark-700 dark:bg-dark-900/80">
                  <div class="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
                    <div class="text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">
                      {{ t('modelPlaza.pricingSupportedGroupsTitle') }}
                    </div>
                    <div class="text-[11px] text-gray-400 dark:text-dark-500">
                      {{ t('modelPlaza.pricingSupportedGroupsCount', { count: supportedGroupEntries(row).length }) }}
                    </div>
                  </div>

                  <div
                    v-if="supportedGroupEntries(row).length"
                    class="mt-3 grid gap-2 lg:grid-cols-2 2xl:grid-cols-3"
                  >
                    <div
                      v-for="entry in supportedGroupEntries(row)"
                      :key="entry.group.id"
                      class="rounded-xl border border-gray-200 bg-gray-50 px-3 py-3 dark:border-dark-700 dark:bg-dark-800/80"
                    >
                      <div class="flex flex-wrap items-center gap-2">
                        <span class="inline-flex rounded-full px-2 py-0.5 text-[10px] font-medium" :class="platformBadgeClass(entry.group.platform)">
                          {{ t(`admin.groups.platforms.${entry.group.platform}`) }}
                        </span>
                        <span class="truncate text-sm font-semibold text-gray-900 dark:text-white">
                          {{ entry.group.name }}
                        </span>
                      </div>

                      <p class="mt-1 text-[11px] leading-5 text-gray-500 dark:text-dark-400">
                        {{ groupSummary(entry.group) }}
                      </p>

                      <div class="mt-3 flex flex-wrap gap-x-3 gap-y-2">
                        <div
                          v-for="metric in metricRows(entry.metrics, 'balance')"
                          :key="metric.key"
                          class="inline-flex items-center gap-1.5 text-xs"
                        >
                          <span class="text-gray-500 dark:text-dark-400">{{ metric.label }}</span>
                          <span class="font-mono font-semibold text-gray-900 dark:text-white">{{ metric.value }}</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div v-else class="mt-3 rounded-xl border border-dashed border-gray-200 bg-gray-50 px-3 py-4 text-center text-xs text-gray-400 dark:border-dark-700 dark:bg-dark-800/80 dark:text-dark-500">
                    {{ t('modelPlaza.notSupportedInGroup') }}
                  </div>
                </div>
              </div>
            </div>
          </article>
        </div>

        <div v-else class="rounded-2xl border border-dashed border-gray-200 bg-gray-50/60 px-4 py-12 text-center text-sm text-gray-500 dark:border-dark-700 dark:bg-dark-800/60 dark:text-dark-400">
          {{ t('modelPlaza.noPricingRows') }}
        </div>

        <Pagination
          v-if="filteredPricingRows.length"
          :page="pricingPage"
          :page-size="pricingPageSize"
          :total="filteredPricingRows.length"
          :page-size-options="[10, 20, 50, 100]"
          @update:page="handlePricingPageChange"
          @update:pageSize="handlePricingPageSizeChange"
        />
      </div>

      <div v-else-if="!loading" class="mt-6 rounded-2xl border border-dashed border-gray-200 bg-gray-50/60 px-4 py-10 text-center text-sm text-gray-500 dark:border-dark-700 dark:bg-dark-800/60 dark:text-dark-400">
        {{ t('modelPlaza.noPricingRows') }}
      </div>
    </section>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
    </div>

    <div v-else-if="!groupModels.length" class="py-12 text-center text-gray-500 dark:text-dark-400">
      {{ t('modelPlaza.empty') }}
    </div>

    <div v-else class="space-y-4">
      <div>
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('modelPlaza.groupListTitle') }}</h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('modelPlaza.groupListDesc') }}</p>
      </div>

      <div class="grid gap-6 md:grid-cols-2">
        <div
          v-for="item in sortedGroupModels"
          :key="item.group.id"
          class="model-plaza-card rounded-2xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900"
        >
          <div class="mb-4 flex items-center justify-between">
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ item.group.name }}</h2>
              <p v-if="item.group.description" class="mt-1 text-sm text-gray-500 dark:text-dark-400">
                {{ item.group.description }}
              </p>
            </div>
            <span class="rounded-full px-3 py-1 text-xs font-medium" :class="platformBadgeClass(item.group.platform)">
              {{ t(`admin.groups.platforms.${item.group.platform}`) }}
            </span>
          </div>

          <div class="mb-4 flex flex-wrap gap-3 text-xs text-gray-500 dark:text-dark-400">
            <span v-if="item.group.rate_multiplier !== 1">
              {{ t('modelPlaza.rate') }}: {{ item.group.rate_multiplier }}x
            </span>
            <span v-if="item.group.subscription_type === 'subscription'" class="text-primary-600 dark:text-primary-400">
              {{ t('modelPlaza.subscriptionType') }}
            </span>
            <span v-else class="text-green-600 dark:text-green-400">
              {{ t('modelPlaza.standardType') }}
            </span>
            <span v-if="item.group.claude_code_only" class="text-amber-600 dark:text-amber-400">
              Claude Code Only
            </span>
          </div>

          <div v-if="getGroupModels(item).length">
            <div class="mb-2 text-sm font-medium text-gray-700 dark:text-dark-300">
              {{ t('modelPlaza.availableModels') }} ({{ getGroupModels(item).length }})
            </div>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="model in visibleModels(item)"
                :key="model"
                @click="copyModel(model)"
                class="inline-block cursor-pointer rounded-md bg-gray-100 px-2.5 py-1 text-xs font-mono text-gray-700 transition-colors hover:bg-primary-100 hover:text-primary-700 dark:bg-dark-700 dark:text-dark-300 dark:hover:bg-primary-900/30 dark:hover:text-primary-400"
                :title="t('modelPlaza.clickToCopy')"
              >
                {{ model }}
              </button>
            </div>
            <div v-if="hiddenModelCount(item) > 0" class="mt-3">
              <button
                class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:text-dark-300 dark:hover:border-primary-700 dark:hover:text-primary-400"
                @click="toggleExpanded(item.group.id)"
              >
                {{ isExpanded(item.group.id) ? t('common.collapse') : t('common.more') }}
                <span v-if="!isExpanded(item.group.id)" class="text-gray-400 dark:text-dark-400">
                  +{{ hiddenModelCount(item) }}
                </span>
              </button>
            </div>
          </div>
          <div v-else class="text-sm italic text-gray-400 dark:text-dark-500">
            {{ t('modelPlaza.allModels') }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Pagination from '@/components/common/Pagination.vue'
import {
  getModelPlaza,
  getModelPlazaPricingTable,
  type GroupModels,
  type ModelPlazaPricingGroup,
  type ModelPlazaPricingItem,
  type ModelPlazaPricingMetrics,
  type ModelPlazaPricingTable
} from '@/api/model-plaza'
import { useClipboard } from '@/composables/useClipboard'

type MetricUnit = 'usd' | 'balance'

const { t } = useI18n()
const { copyToClipboard } = useClipboard()
const INITIAL_VISIBLE_MODELS = 18

const loading = ref(true)
const groupModels = ref<GroupModels[]>([])
const pricingTable = ref<ModelPlazaPricingTable | null>(null)
const expandedGroupIds = ref<number[]>([])
const pricingSearch = ref('')
const pricingPage = ref(1)
const pricingPageSize = ref(20)

const platformOrder: Record<string, number> = {
  anthropic: 0,
  openai: 1,
  gemini: 2,
  antigravity: 3,
  multi: 4,
}

const pricingGroups = computed(() => pricingTable.value?.groups || [])
const pricingRows = computed(() => pricingTable.value?.items || [])
const sortedGroupModels = computed(() => {
  return [...groupModels.value].sort((a, b) => {
    const platformDelta = (platformOrder[a.group.platform] ?? 999) - (platformOrder[b.group.platform] ?? 999)
    if (platformDelta !== 0) {
      return platformDelta
    }

    return a.group.name.localeCompare(b.group.name)
  })
})

const filteredPricingRows = computed(() => {
  const query = pricingSearch.value.trim().toLowerCase()
  if (!query) {
    return pricingRows.value
  }

  return pricingRows.value.filter((row) => {
    const haystack = [
      row.model,
      row.platform,
      row.provider,
      row.mode,
      ...(row.aliases ?? []),
    ]
      .join(' ')
      .toLowerCase()

    return haystack.includes(query)
  })
})

const paginatedPricingRows = computed(() => {
  const start = (pricingPage.value - 1) * pricingPageSize.value
  return filteredPricingRows.value.slice(start, start + pricingPageSize.value)
})

function platformBadgeClass(platform: string) {
  const map: Record<string, string> = {
    anthropic: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
    openai: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    gemini: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    antigravity: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
    multi: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
  }
  return map[platform] || map.multi
}

function formatMetricValue(value: number | null | undefined, unit: MetricUnit) {
  if (value === null || value === undefined || Number.isNaN(value)) {
    return '-'
  }

  const digits = value >= 100 ? 2 : 3
  const formatted = value.toFixed(digits).replace(/\.?0+$/, '')
  return unit === 'usd' ? `$${formatted}` : formatted
}

function formatMultiplierLabel(value: number | null | undefined) {
  if (value === null || value === undefined || Number.isNaN(value)) {
    return '-'
  }
  const digits = value >= 100 ? 2 : 3
  return `${value.toFixed(digits).replace(/\.?0+$/, '')}x`
}

function metricRows(metrics: ModelPlazaPricingMetrics | undefined, unit: MetricUnit) {
  if (!metrics) {
    return []
  }

  const rows = [
    {
      key: 'input',
      label: t('modelPlaza.metricInput'),
      value: formatMetricValue(metrics.input_per_million, unit)
    },
    {
      key: 'output',
      label: t('modelPlaza.metricOutput'),
      value: formatMetricValue(metrics.output_per_million, unit)
    }
  ]

  if (hasPositiveMetric(metrics.cache_write_per_million)) {
    rows.push({
      key: 'cache-write',
      label: t('modelPlaza.metricCacheWrite'),
      value: formatMetricValue(metrics.cache_write_per_million, unit)
    })
  }

  if (hasPositiveMetric(metrics.cache_read_per_million)) {
    rows.push({
      key: 'cache-read',
      label: t('modelPlaza.metricCacheRead'),
      value: formatMetricValue(metrics.cache_read_per_million, unit)
    })
  }

  return rows
}

function hasPositiveMetric(value: number | null | undefined) {
  return value !== null && value !== undefined && !Number.isNaN(value) && value > 0
}

function supportedGroupEntries(row: ModelPlazaPricingItem) {
  return pricingGroups.value
    .filter((group) => Boolean(row.group_prices[group.id]))
    .map((group) => ({
      group,
      metrics: row.group_prices[group.id]
    }))
}

function groupSummary(group: ModelPlazaPricingGroup) {
  const parts = [
    group.display_discount || formatMultiplierLabel(group.display_rate_multiplier ?? group.rate_multiplier)
  ]

  if (group.display_price) {
    parts.push(group.display_price)
  }

  return parts.join(' · ')
}

function aliasSummary(aliases: string[]) {
  if (!aliases.length) {
    return ''
  }

  const visibleAliases = aliases.slice(0, 2)
  if (aliases.length <= visibleAliases.length) {
    return visibleAliases.join(', ')
  }

  return `${visibleAliases.join(', ')} +${aliases.length - visibleAliases.length}`
}

function handlePricingPageChange(page: number) {
  pricingPage.value = page
}

function handlePricingPageSizeChange(pageSize: number) {
  pricingPageSize.value = pageSize
  pricingPage.value = 1
}

function copyModel(model: string) {
  copyToClipboard(model, t('modelPlaza.modelCopied'))
}

function isExpanded(groupId: number) {
  return expandedGroupIds.value.includes(groupId)
}

function toggleExpanded(groupId: number) {
  expandedGroupIds.value = isExpanded(groupId)
    ? expandedGroupIds.value.filter((id) => id !== groupId)
    : [...expandedGroupIds.value, groupId]
}

function getGroupModels(item: GroupModels) {
  return Array.isArray(item.models) ? item.models : []
}

function hiddenModelCount(item: GroupModels) {
  const models = getGroupModels(item)
  if (isExpanded(item.group.id) || models.length <= INITIAL_VISIBLE_MODELS) {
    return 0
  }

  return models.length - INITIAL_VISIBLE_MODELS
}

function visibleModels(item: GroupModels) {
  const models = getGroupModels(item)
  if (isExpanded(item.group.id) || models.length <= INITIAL_VISIBLE_MODELS) {
    return models
  }

  return models.slice(0, INITIAL_VISIBLE_MODELS)
}

watch(pricingSearch, () => {
  pricingPage.value = 1
})

watch(filteredPricingRows, (rows) => {
  const totalPages = Math.max(1, Math.ceil(rows.length / pricingPageSize.value))
  if (pricingPage.value > totalPages) {
    pricingPage.value = totalPages
  }
})

onMounted(async () => {
  try {
    const [groupsResult, pricingResult] = await Promise.allSettled([
      getModelPlaza(),
      getModelPlazaPricingTable(),
    ])

    if (groupsResult.status === 'fulfilled') {
      groupModels.value = groupsResult.value
    } else {
      console.error('Failed to load model plaza groups:', groupsResult.reason)
    }

    if (pricingResult.status === 'fulfilled') {
      pricingTable.value = pricingResult.value
    } else {
      console.error('Failed to load model plaza pricing table:', pricingResult.reason)
    }
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.model-plaza-card {
  content-visibility: auto;
  contain-intrinsic-size: 320px;
}
</style>
