<template>
  <div class="mx-auto max-w-7xl space-y-8">
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
              v-if="sortedGroupModels.length"
              class="inline-flex rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-300"
            >
              {{ t('modelPlaza.pricingGroupsCount', { count: sortedGroupModels.length }) }}
            </span>
          </div>

          <h2 class="mt-4 text-2xl font-semibold text-gray-900 dark:text-white">
            {{ t('modelPlaza.pricingExplorerTitle') }}
          </h2>
          <p class="mt-2 text-sm leading-6 text-gray-600 dark:text-dark-300">
            {{ t('modelPlaza.pricingExplorerDesc') }}
          </p>

          <div class="mt-4 flex flex-wrap gap-2">
            <span class="inline-flex rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-300">
              {{ t('modelPlaza.officialUnitHint') }}
            </span>
            <span class="inline-flex rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-dark-300">
              {{ t('modelPlaza.siteUnitHint') }}
            </span>
            <span class="inline-flex rounded-full bg-primary-50 px-3 py-1 text-xs font-medium text-primary-700 dark:bg-primary-950/30 dark:text-primary-300">
              {{ t('modelPlaza.pricingMultiplierIncluded') }}
            </span>
          </div>

          <p class="mt-3 text-xs leading-5 text-gray-500 dark:text-dark-400">
            {{ t('modelPlaza.pricingSelectionHint') }}
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
            {{ t('modelPlaza.pricingExplorerStats', { groups: filteredGroupSections.length, models: filteredModelCount }) }}
          </p>
        </div>
      </div>

      <div v-if="loading" class="mt-8 flex justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
      </div>

      <div
        v-else-if="!filteredGroupSections.length"
        class="mt-8 rounded-2xl border border-dashed border-gray-200 bg-gray-50/60 px-4 py-12 text-center text-sm text-gray-500 dark:border-dark-700 dark:bg-dark-800/60 dark:text-dark-400"
      >
        {{ pricingSearch.trim() ? t('modelPlaza.noPricingRows') : t('modelPlaza.empty') }}
      </div>

      <div v-else class="mt-8 grid gap-6 xl:grid-cols-[minmax(0,1fr)_360px]">
        <div class="space-y-4">
          <article
            v-for="section in filteredGroupSections"
            :key="section.group.id"
            class="model-plaza-card rounded-2xl border border-gray-200 bg-gray-50/70 p-5 dark:border-dark-700 dark:bg-dark-800/60"
          >
            <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ section.group.name }}</h3>
                  <span class="rounded-full px-2.5 py-1 text-xs font-medium" :class="platformBadgeClass(section.group.platform)">
                    {{ t(`admin.groups.platforms.${section.group.platform}`) }}
                  </span>
                </div>

                <p v-if="section.group.description" class="mt-1 text-sm text-gray-500 dark:text-dark-400">
                  {{ section.group.description }}
                </p>
              </div>

              <div class="flex flex-wrap gap-2 text-xs">
                <span class="rounded-full bg-white px-3 py-1 text-gray-600 dark:bg-dark-900 dark:text-dark-300">
                  {{ groupSummary(section.pricingGroup ?? section.group) }}
                </span>
                <span
                  v-if="section.group.subscription_type === 'subscription'"
                  class="rounded-full bg-primary-50 px-3 py-1 text-primary-700 dark:bg-primary-950/30 dark:text-primary-300"
                >
                  {{ t('modelPlaza.subscriptionType') }}
                </span>
                <span
                  v-else
                  class="rounded-full bg-emerald-50 px-3 py-1 text-emerald-700 dark:bg-emerald-950/30 dark:text-emerald-300"
                >
                  {{ t('modelPlaza.standardType') }}
                </span>
              </div>
            </div>

            <div class="mt-4 flex flex-wrap gap-2">
              <button
                v-for="model in visibleModels(section)"
                :key="model"
                type="button"
                class="inline-flex items-center rounded-xl border px-3 py-2 text-left text-xs font-mono transition-colors"
                :class="isSelectedModel(section.group.id, model)
                  ? 'border-primary-300 bg-primary-50 text-primary-700 shadow-sm dark:border-primary-700 dark:bg-primary-950/30 dark:text-primary-300'
                  : 'border-gray-200 bg-white text-gray-700 hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300'"
                :title="t('modelPlaza.clickToCopy')"
                @click="handleModelSelect(section.group.id, model)"
              >
                {{ model }}
              </button>
            </div>

            <div v-if="hiddenModelCount(section) > 0" class="mt-3">
              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:text-dark-300 dark:hover:border-primary-700 dark:hover:text-primary-300"
                @click="toggleExpanded(section.group.id)"
              >
                {{ isExpanded(section.group.id) ? t('common.collapse') : t('common.more') }}
                <span v-if="!isExpanded(section.group.id)" class="text-gray-400 dark:text-dark-400">
                  +{{ hiddenModelCount(section) }}
                </span>
              </button>
            </div>
          </article>
        </div>

        <aside class="xl:sticky xl:top-24 xl:self-start">
          <div class="rounded-3xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900 xl:max-h-[calc(100vh-7rem)] xl:overflow-y-auto">
            <p class="text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">
              {{ t('modelPlaza.pricingPanelTitle') }}
            </p>

            <template v-if="selectedGroup && selectedRequestedModel">
              <h3 class="mt-3 break-all font-mono text-sm font-semibold leading-6 text-gray-900 dark:text-white">
                {{ selectedRequestedModel }}
              </h3>

              <p class="mt-2 text-xs leading-5 text-gray-500 dark:text-dark-400">
                {{ t('modelPlaza.pricingPanelCopiedHint') }}
              </p>

              <div class="mt-4 rounded-2xl border border-gray-200 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-800/70">
                <div class="text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">
                  {{ t('modelPlaza.pricingCurrentGroup') }}
                </div>
                <div class="mt-2 flex flex-wrap items-center gap-2">
                  <span class="text-sm font-semibold text-gray-900 dark:text-white">
                    {{ selectedGroup.name }}
                  </span>
                  <span class="rounded-full px-2 py-0.5 text-[10px] font-medium" :class="platformBadgeClass(selectedGroup.platform)">
                    {{ t(`admin.groups.platforms.${selectedGroup.platform}`) }}
                  </span>
                </div>
                <p class="mt-2 text-xs leading-5 text-gray-500 dark:text-dark-400">
                  {{ groupSummary(selectedPricingGroup ?? selectedGroup) }}
                </p>
              </div>

              <div
                v-if="selectedPricingRow && normalizeModelKey(selectedRequestedModel) !== normalizeModelKey(selectedPricingRow.model)"
                class="mt-4 rounded-2xl border border-amber-200 bg-amber-50/80 p-4 text-xs text-amber-800 dark:border-amber-900/40 dark:bg-amber-950/20 dark:text-amber-200"
              >
                <div class="font-semibold">{{ t('modelPlaza.pricingMatchedModelLabel') }}</div>
                <div class="mt-1 break-all font-mono">{{ selectedPricingRow.model }}</div>
              </div>

              <div class="mt-4 space-y-4">
                <div class="overflow-hidden rounded-3xl bg-gradient-to-br from-primary-600 via-primary-600 to-sky-500 p-5 text-white shadow-lg shadow-primary-900/10">
                  <div class="flex items-start justify-between gap-3">
                    <div>
                      <div class="text-[11px] font-semibold uppercase tracking-[0.18em] text-white/75">
                        {{ t('modelPlaza.sitePriceCardTitle') }}
                      </div>
                      <div class="mt-2 text-base font-semibold">
                        {{ selectedGroup.name }}
                      </div>
                    </div>
                    <div class="rounded-full border border-white/15 bg-white/10 px-3 py-1 text-[11px] font-medium text-white/85">
                      {{ t('modelPlaza.siteUnitHint') }}
                    </div>
                  </div>

                  <div class="mt-3 flex flex-wrap items-center gap-2 text-xs text-white/80">
                    <span class="rounded-full border border-white/15 bg-white/10 px-3 py-1">
                      {{ groupSummary(selectedPricingGroup ?? selectedGroup) }}
                    </span>
                    <span class="rounded-full border border-white/15 bg-white/10 px-3 py-1">
                      {{ t('modelPlaza.pricingMultiplierIncluded') }}
                    </span>
                  </div>

                  <div v-if="selectedPrimaryComparisonRows.length" class="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-1 2xl:grid-cols-2">
                    <div
                      v-for="metric in selectedPrimaryComparisonRows"
                      :key="metric.key"
                      class="rounded-2xl border border-white/15 bg-white/10 px-4 py-4 backdrop-blur-sm"
                    >
                      <div class="text-xs font-medium text-white/75">{{ metric.label }}</div>
                      <div class="mt-2 font-mono text-3xl font-semibold leading-none text-white">
                        {{ metric.siteValue }}
                      </div>
                      <div class="mt-3 text-xs text-white/70">
                        {{ t('modelPlaza.pricingCompareOfficialShort') }} {{ metric.officialValue }}
                      </div>
                    </div>
                  </div>

                  <div v-if="selectedSecondaryComparisonRows.length" class="mt-3 grid gap-3 sm:grid-cols-2 xl:grid-cols-1 2xl:grid-cols-2">
                    <div
                      v-for="metric in selectedSecondaryComparisonRows"
                      :key="metric.key"
                      class="rounded-2xl border border-white/10 bg-black/10 px-4 py-3"
                    >
                      <div class="text-xs font-medium text-white/70">{{ metric.label }}</div>
                      <div class="mt-1 font-mono text-lg font-semibold text-white">
                        {{ metric.siteValue }}
                      </div>
                      <div class="mt-2 text-xs text-white/65">
                        {{ t('modelPlaza.pricingCompareOfficialShort') }} {{ metric.officialValue }}
                      </div>
                    </div>
                  </div>

                  <div v-if="!selectedComparisonRows.length" class="mt-4 text-sm text-white/75">
                    {{ t('modelPlaza.notSupportedInGroup') }}
                  </div>
                </div>

                <div class="rounded-2xl border border-gray-200 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-800/70">
                  <div class="flex items-start justify-between gap-3">
                    <div>
                      <div class="text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-dark-400">
                        {{ t('modelPlaza.pricingCompareTitle') }}
                      </div>
                      <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">
                        {{ t('modelPlaza.pricingCompareDesc') }}
                      </div>
                    </div>
                    <div class="rounded-full bg-white px-3 py-1 text-[11px] text-gray-500 dark:bg-dark-900 dark:text-dark-400">
                      {{ t('modelPlaza.pricingCompareRange') }}
                    </div>
                  </div>

                  <div v-if="selectedComparisonRows.length" class="mt-4 overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
                    <div class="grid grid-cols-[minmax(0,1fr)_110px_110px] border-b border-gray-200 bg-gray-50 px-4 py-3 text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-400">
                      <div>{{ t('modelPlaza.pricingCompareMetricLabel') }}</div>
                      <div class="text-right">{{ t('modelPlaza.pricingCompareOfficialShort') }}</div>
                      <div class="text-right text-primary-700 dark:text-primary-300">{{ t('modelPlaza.pricingCompareSiteShort') }}</div>
                    </div>
                    <div
                      v-for="metric in selectedComparisonRows"
                      :key="metric.key"
                      class="grid grid-cols-[minmax(0,1fr)_110px_110px] items-center gap-3 border-t border-gray-100 px-4 py-3 first:border-t-0 dark:border-dark-800"
                    >
                      <div class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ metric.label }}</div>
                      <div class="text-right font-mono text-sm text-gray-500 dark:text-dark-400">{{ metric.officialValue }}</div>
                      <div class="text-right font-mono text-sm font-semibold text-primary-700 dark:text-primary-300">{{ metric.siteValue }}</div>
                    </div>
                  </div>

                  <div v-else class="mt-3 text-xs text-gray-400 dark:text-dark-500">
                    {{ t('modelPlaza.notSupportedInGroup') }}
                  </div>
                </div>
              </div>
            </template>

            <div v-else class="mt-4 rounded-2xl border border-dashed border-gray-200 bg-gray-50/70 px-4 py-10 text-center dark:border-dark-700 dark:bg-dark-800/60">
              <div class="text-sm font-semibold text-gray-900 dark:text-white">
                {{ t('modelPlaza.pricingPanelEmptyTitle') }}
              </div>
              <p class="mt-2 text-sm leading-6 text-gray-500 dark:text-dark-400">
                {{ t('modelPlaza.pricingPanelEmptyDesc') }}
              </p>
            </div>
          </div>
        </aside>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
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
type MetricField = 'input_per_million' | 'output_per_million' | 'cache_write_per_million' | 'cache_read_per_million'

interface GroupSection {
  group: GroupModels['group']
  pricingGroup: ModelPlazaPricingGroup | null
  models: string[]
}

interface MetricComparisonRow {
  key: string
  label: string
  officialValue: string
  siteValue: string
}

const { t } = useI18n()
const { copyToClipboard } = useClipboard()
const INITIAL_VISIBLE_MODELS = 18

const loading = ref(true)
const groupModels = ref<GroupModels[]>([])
const pricingTable = ref<ModelPlazaPricingTable | null>(null)
const expandedGroupIds = ref<number[]>([])
const pricingSearch = ref('')
const selectedGroupId = ref<number | null>(null)
const selectedRequestedModel = ref('')

const platformOrder: Record<string, number> = {
  anthropic: 0,
  openai: 1,
  gemini: 2,
  antigravity: 3,
  multi: 4,
}

const pricingGroups = computed(() => pricingTable.value?.groups || [])
const pricingRows = computed(() => pricingTable.value?.items || [])
const pricingGroupMap = computed(() => new Map(pricingGroups.value.map((group) => [group.id, group])))

const pricingRowLookup = computed(() => {
  const lookup = new Map<string, ModelPlazaPricingItem>()

  for (const row of pricingRows.value) {
    lookup.set(normalizeModelKey(row.model), row)
    for (const alias of row.aliases ?? []) {
      const key = normalizeModelKey(alias)
      if (!lookup.has(key)) {
        lookup.set(key, row)
      }
    }
  }

  return lookup
})

const pricingModelsByGroup = computed(() => {
  const map = new Map<number, string[]>()

  for (const row of pricingRows.value) {
    for (const groupId of Object.keys(row.group_prices)) {
      const numericGroupId = Number(groupId)
      if (!Number.isFinite(numericGroupId)) {
        continue
      }
      const list = map.get(numericGroupId) ?? []
      list.push(row.model)
      map.set(numericGroupId, list)
    }
  }

  return map
})

const sortedGroupModels = computed(() => {
  return [...groupModels.value].sort((a, b) => {
    const platformDelta = (platformOrder[a.group.platform] ?? 999) - (platformOrder[b.group.platform] ?? 999)
    if (platformDelta !== 0) {
      return platformDelta
    }

    return a.group.name.localeCompare(b.group.name)
  })
})

const filteredGroupSections = computed<GroupSection[]>(() => {
  const query = pricingSearch.value.trim().toLowerCase()

  return sortedGroupModels.value
    .map((item) => {
      const pricingGroup = pricingGroupMap.value.get(item.group.id) ?? null
      const baseModels = getBaseGroupModels(item.group.id, getGroupModels(item))
      const groupMatches = query ? groupSearchText(item, pricingGroup).includes(query) : false
      const models = query
        ? (groupMatches ? baseModels : baseModels.filter((model) => modelSearchText(item, model, pricingGroup).includes(query)))
        : baseModels

      return {
        group: item.group,
        pricingGroup,
        models,
      }
    })
    .filter((section) => section.models.length > 0)
})

const filteredModelCount = computed(() => {
  return filteredGroupSections.value.reduce((count, section) => count + section.models.length, 0)
})

const selectedGroup = computed(() => {
  if (selectedGroupId.value === null) {
    return null
  }
  return sortedGroupModels.value.find((item) => item.group.id === selectedGroupId.value)?.group ?? null
})

const selectedPricingGroup = computed(() => {
  if (selectedGroupId.value === null) {
    return null
  }
  return pricingGroupMap.value.get(selectedGroupId.value) ?? null
})

const selectedPricingRow = computed(() => {
  if (!selectedRequestedModel.value) {
    return null
  }
  return pricingRowLookup.value.get(normalizeModelKey(selectedRequestedModel.value)) ?? null
})

const selectedSiteMetrics = computed(() => {
  if (!selectedPricingRow.value || selectedGroupId.value === null) {
    return null
  }
  return selectedPricingRow.value.group_prices[selectedGroupId.value] ?? null
})

const selectedDisplayMultiplier = computed(() => {
  return (
    extractDisplayPriceMultiplier(selectedPricingGroup.value?.display_price) ??
    extractDisplayPriceMultiplier(selectedGroup.value?.display_price) ??
    normalizePositiveNumber(selectedPricingGroup.value?.display_rate_multiplier) ??
    normalizePositiveNumber(selectedPricingGroup.value?.rate_multiplier) ??
    normalizePositiveNumber(selectedGroup.value?.rate_multiplier) ??
    1
  )
})

const selectedDisplayedSiteMetrics = computed(() => {
  if (!selectedPricingRow.value) {
    return null
  }

  if (metricsDifferFromOfficial(selectedSiteMetrics.value, selectedPricingRow.value.official)) {
    return selectedSiteMetrics.value
  }

  return scaleMetrics(selectedPricingRow.value.official, selectedDisplayMultiplier.value)
})

const selectedComparisonRows = computed(() => {
  return buildMetricComparisonRows(selectedPricingRow.value?.official, selectedDisplayedSiteMetrics.value)
})

const selectedPrimaryComparisonRows = computed(() => {
  return selectedComparisonRows.value.filter((row) => row.key === 'input' || row.key === 'output')
})

const selectedSecondaryComparisonRows = computed(() => {
  return selectedComparisonRows.value.filter((row) => row.key !== 'input' && row.key !== 'output')
})

function normalizeModelKey(model: string) {
  return model.trim().toLowerCase()
}

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

function dedupeModels(models: string[]) {
  const seen = new Set<string>()
  const result: string[] = []

  for (const model of models) {
    const normalized = model.trim()
    if (!normalized || seen.has(normalized)) {
      continue
    }
    seen.add(normalized)
    result.push(normalized)
  }

  return result
}

function getGroupModels(item: GroupModels) {
  return Array.isArray(item.models) ? item.models : []
}

function getBaseGroupModels(groupId: number, explicitModels: string[]) {
  if (explicitModels.length > 0) {
    return dedupeModels(explicitModels)
  }

  const derivedModels = pricingModelsByGroup.value.get(groupId) ?? []
  return dedupeModels(derivedModels)
}

function getPricingRowByModel(model: string) {
  return pricingRowLookup.value.get(normalizeModelKey(model)) ?? null
}

function groupSearchText(item: GroupModels, pricingGroup: ModelPlazaPricingGroup | null) {
  return [
    item.group.name,
    item.group.description,
    item.group.platform,
    item.group.display_price,
    item.group.display_discount,
    pricingGroup?.display_price,
    pricingGroup?.display_discount,
  ]
    .filter(Boolean)
    .join(' ')
    .toLowerCase()
}

function modelSearchText(item: GroupModels, model: string, pricingGroup: ModelPlazaPricingGroup | null) {
  const pricingRow = getPricingRowByModel(model)

  return [
    model,
    pricingRow?.model,
    ...(pricingRow?.aliases ?? []),
    pricingRow?.provider,
    pricingRow?.mode,
    pricingRow?.platform,
    item.group.name,
    item.group.platform,
    pricingGroup?.display_price,
  ]
    .filter(Boolean)
    .join(' ')
    .toLowerCase()
}

function formatMetricValue(value: number | null | undefined, unit: MetricUnit) {
  if (value === null || value === undefined || Number.isNaN(value)) {
    return '-'
  }

  const digits = value >= 100 ? 2 : 3
  const formatted = value.toFixed(digits).replace(/\.?0+$/, '')
  return unit === 'usd' ? `$${formatted}` : `¥${formatted}`
}

function formatMultiplierLabel(value: number | null | undefined) {
  if (value === null || value === undefined || Number.isNaN(value)) {
    return '-'
  }
  const digits = value >= 100 ? 2 : 3
  return `${value.toFixed(digits).replace(/\.?0+$/, '')}x`
}

function normalizePositiveNumber(value: number | null | undefined) {
  if (value === null || value === undefined || Number.isNaN(value) || value <= 0) {
    return null
  }

  return value
}

function hasPositiveMetric(value: number | null | undefined) {
  return value !== null && value !== undefined && !Number.isNaN(value) && value > 0
}

function extractDisplayPriceMultiplier(displayPrice: string | null | undefined) {
  if (!displayPrice) {
    return null
  }

  const normalized = displayPrice.trim().toLowerCase()
  if (!normalized) {
    return null
  }

  const match = normalized.match(/(\d+(?:\.\d+)?)/)
  if (!match) {
    return null
  }

  const parsed = Number(match[1])
  return normalizePositiveNumber(parsed)
}

function scaleMetricValue(value: number | null | undefined, multiplier: number) {
  if (value === null || value === undefined || Number.isNaN(value)) {
    return 0
  }

  return value * multiplier
}

function scaleOptionalMetricValue(value: number | null | undefined, multiplier: number) {
  if (value === null || value === undefined || Number.isNaN(value)) {
    return null
  }

  return value * multiplier
}

function scaleMetrics(metrics: ModelPlazaPricingMetrics, multiplier: number): ModelPlazaPricingMetrics {
  return {
    input_per_million: scaleMetricValue(metrics.input_per_million, multiplier),
    output_per_million: scaleMetricValue(metrics.output_per_million, multiplier),
    cache_write_per_million: scaleMetricValue(metrics.cache_write_per_million, multiplier),
    cache_read_per_million: scaleMetricValue(metrics.cache_read_per_million, multiplier),
    input_per_million_above_200k: scaleOptionalMetricValue(metrics.input_per_million_above_200k, multiplier),
    output_per_million_above_200k: scaleOptionalMetricValue(metrics.output_per_million_above_200k, multiplier),
    cache_write_per_million_above_200k: scaleOptionalMetricValue(metrics.cache_write_per_million_above_200k, multiplier),
    cache_read_per_million_above_200k: scaleOptionalMetricValue(metrics.cache_read_per_million_above_200k, multiplier),
  }
}

function metricsDifferFromOfficial(
  site: ModelPlazaPricingMetrics | null | undefined,
  official: ModelPlazaPricingMetrics | null | undefined
) {
  if (!site || !official) {
    return false
  }

  const fields: MetricField[] = [
    'input_per_million',
    'output_per_million',
    'cache_write_per_million',
    'cache_read_per_million',
  ]

  return fields.some((field) => {
    const siteValue = getMetricValue(site, field)
    const officialValue = getMetricValue(official, field)
    if (siteValue === null || officialValue === null) {
      return false
    }

    return Math.abs(siteValue - officialValue) > 0.000001
  })
}

function getMetricValue(metrics: ModelPlazaPricingMetrics | null | undefined, field: MetricField) {
  if (!metrics) {
    return null
  }

  return metrics[field]
}

function buildMetricComparisonRows(
  official: ModelPlazaPricingMetrics | null | undefined,
  site: ModelPlazaPricingMetrics | null | undefined
): MetricComparisonRow[] {
  const definitions: Array<{ key: MetricComparisonRow['key']; label: string; field: MetricField; always?: boolean }> = [
    { key: 'input', label: t('modelPlaza.metricInput'), field: 'input_per_million', always: true },
    { key: 'output', label: t('modelPlaza.metricOutput'), field: 'output_per_million', always: true },
    { key: 'cache-write', label: t('modelPlaza.metricCacheWrite'), field: 'cache_write_per_million' },
    { key: 'cache-read', label: t('modelPlaza.metricCacheRead'), field: 'cache_read_per_million' },
  ]

  return definitions
    .filter((definition) => {
      if (definition.always) {
        return Boolean(official || site)
      }

      return hasPositiveMetric(getMetricValue(official, definition.field)) || hasPositiveMetric(getMetricValue(site, definition.field))
    })
    .map((definition) => ({
      key: definition.key,
      label: definition.label,
      officialValue: formatMetricValue(getMetricValue(official, definition.field), 'usd'),
      siteValue: formatMetricValue(getMetricValue(site, definition.field), 'balance'),
    }))
}

function groupSummary(group: {
  display_discount?: string | null
  display_price?: string | null
  rate_multiplier: number
  display_rate_multiplier?: number | null
}) {
  const parts = [
    group.display_discount || formatMultiplierLabel(group.display_rate_multiplier ?? group.rate_multiplier)
  ]

  if (group.display_price) {
    parts.push(group.display_price)
  }

  return parts.join(' · ')
}

function isExpanded(groupId: number) {
  return expandedGroupIds.value.includes(groupId)
}

function toggleExpanded(groupId: number) {
  expandedGroupIds.value = isExpanded(groupId)
    ? expandedGroupIds.value.filter((id) => id !== groupId)
    : [...expandedGroupIds.value, groupId]
}

function hiddenModelCount(section: GroupSection) {
  if (isExpanded(section.group.id) || section.models.length <= INITIAL_VISIBLE_MODELS) {
    return 0
  }

  return section.models.length - INITIAL_VISIBLE_MODELS
}

function visibleModels(section: GroupSection) {
  if (isExpanded(section.group.id) || section.models.length <= INITIAL_VISIBLE_MODELS) {
    return section.models
  }

  return section.models.slice(0, INITIAL_VISIBLE_MODELS)
}

function isSelectedModel(groupId: number, model: string) {
  return selectedGroupId.value === groupId && selectedRequestedModel.value === model
}

function handleModelSelect(groupId: number, model: string) {
  selectedGroupId.value = groupId
  selectedRequestedModel.value = model
  copyToClipboard(model, t('modelPlaza.modelCopied'))
}

watch(filteredGroupSections, (sections) => {
  if (!sections.length) {
    selectedGroupId.value = null
    selectedRequestedModel.value = ''
    return
  }

  const currentSection = sections.find((section) => section.group.id === selectedGroupId.value)
  if (currentSection && currentSection.models.includes(selectedRequestedModel.value)) {
    return
  }

  selectedGroupId.value = null
  selectedRequestedModel.value = ''
}, { immediate: true })

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
