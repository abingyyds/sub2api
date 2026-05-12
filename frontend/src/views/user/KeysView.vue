<template>
  <div class="mx-auto max-w-7xl space-y-8 px-4 py-6 sm:px-6 lg:px-8">

        <!-- Section 1: 新建 API Key - Group Cards -->
        <section>
            <div class="mb-4">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('keys.newApiKey') }}</h2>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('keys.newApiKeyDesc') }}</p>
              <router-link
                to="/model-plaza"
                class="mt-3 inline-flex max-w-3xl items-start gap-3 rounded-2xl border border-primary-100 bg-primary-50/80 px-4 py-3 text-left transition-colors hover:border-primary-200 hover:bg-primary-50 dark:border-primary-900/40 dark:bg-primary-900/20 dark:hover:border-primary-800 dark:hover:bg-primary-900/30"
              >
                <span class="mt-0.5 flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-100 text-primary-700 dark:bg-primary-900/50 dark:text-primary-300">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.042-.02a.75.75 0 011.037.852l-.708 2.836a.75.75 0 001.124.824l.042-.02M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
                  </svg>
                </span>
                <span class="min-w-0">
                  <span class="block text-sm leading-6 text-gray-600 dark:text-gray-300">{{ t('keys.newApiKeyBillingHint') }}</span>
                  <span class="mt-1 inline-flex items-center gap-1 text-sm font-medium text-primary-700 dark:text-primary-300">
                    {{ t('keys.newApiKeyBillingLink') }}
                    <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                    </svg>
                  </span>
                </span>
              </router-link>
            </div>

            <!-- Groups by Platform -->
            <div
              v-for="platformGroup in groupedByPlatform"
              :key="platformGroup.platform"
              class="keys-platform-group mb-6"
            >
              <div class="mb-3 flex items-center gap-2">
                <PlatformIcon :platform="platformGroup.platform" size="md" />
                <h3 class="font-medium text-gray-800 dark:text-gray-200">
                  {{ t('admin.groups.platforms.' + platformGroup.platform) }}
                </h3>
                <span class="text-xs text-gray-400 dark:text-gray-500">
                  {{ t('keys.groupCount', { count: platformGroup.groups.length }) }}
                </span>
              </div>

              <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                <div
                  v-for="group in visiblePlatformGroups(platformGroup)"
                  :key="group.id"
                  :class="[
                    'relative flex flex-col rounded-xl border p-4 transition-all duration-200',
                    isGroupUnavailable(group)
                      ? 'cursor-not-allowed border-gray-200 bg-gray-50 opacity-60 dark:border-dark-600 dark:bg-dark-800/50'
                      : 'cursor-pointer border-gray-200 bg-white hover:border-primary-300 hover:shadow-md dark:border-dark-600 dark:bg-dark-800 dark:hover:border-primary-600'
                  ]"
                >
                  <!-- Header: Name + Tags -->
                  <div class="mb-2 flex flex-wrap items-start gap-1.5">
                    <span class="font-medium text-gray-900 dark:text-white">{{ group.name }}</span>
                    <span
                      v-for="tag in getGroupSpecialTags(group)"
                      :key="tag.text"
                      :class="[
                        'rounded px-1.5 py-0.5 text-[10px] font-medium',
                        tag.class
                      ]"
                    >
                      {{ tag.text }}
                    </span>
                  </div>

                  <!-- Price & Discount -->
                  <div v-if="group.display_price" class="mb-2">
                    <span class="text-lg font-bold text-gray-900 dark:text-white">{{ group.display_price }}</span>
                    <span v-if="group.display_discount" class="ml-2 text-sm font-medium text-green-600 dark:text-green-400">
                      {{ group.display_discount }}
                    </span>
                  </div>

                  <!-- Regular Tags -->
                  <div v-if="getGroupRegularTags(group).length > 0" class="mb-2 flex flex-wrap gap-1">
                    <span
                      v-for="tag in getGroupRegularTags(group)"
                      :key="tag"
                      class="rounded bg-gray-100 px-1.5 py-0.5 text-[11px] text-gray-600 dark:bg-dark-600 dark:text-gray-400"
                    >
                      {{ tag }}
                    </span>
                  </div>

                  <!-- Description -->
                  <p v-if="group.description" class="mb-3 line-clamp-2 text-xs text-gray-500 dark:text-gray-400">
                    {{ group.description }}
                  </p>
                  <div class="flex-1"></div>

                  <!-- Create Button -->
                  <div class="mt-2 flex items-center justify-end">
                    <button
                      v-if="!isGroupUnavailable(group)"
                      @click="quickCreateKey(group)"
                      class="flex items-center gap-1 rounded-lg bg-primary-50 px-3 py-1.5 text-sm font-medium text-primary-700 transition-colors hover:bg-primary-100 dark:bg-primary-900/20 dark:text-primary-400 dark:hover:bg-primary-900/40"
                    >
                      {{ t('keys.create') }}
                      <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                      </svg>
                    </button>
                    <span v-else class="text-xs font-medium text-gray-400 dark:text-gray-500">
                      {{ t('keys.unavailable') }}
                    </span>
                  </div>
                </div>

                <!-- 购买新套餐 card (only for multi/通用套餐 platform) -->
                <router-link
                  v-if="platformGroup.platform === 'multi'"
                  to="/pricing"
                  class="relative flex flex-col items-center justify-center rounded-xl border-2 border-dashed border-primary-300 bg-primary-50/50 p-4 transition-all duration-200 hover:border-primary-400 hover:bg-primary-50 hover:shadow-md dark:border-primary-700 dark:bg-primary-900/10 dark:hover:border-primary-600 dark:hover:bg-primary-900/20"
                >
                  <svg class="mb-2 h-8 w-8 text-primary-500 dark:text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 00-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 00-16.536-1.84M7.5 14.25L5.106 5.272M6 20.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm12.75 0a.75.75 0 11-1.5 0 .75.75 0 011.5 0z" />
                  </svg>
                  <span class="text-sm font-medium text-primary-700 dark:text-primary-400">{{ t('keys.buyNewPlan') }}</span>
                  <span class="mt-1 text-[11px] text-primary-500/70 dark:text-primary-500/50">{{ t('keys.buyNewPlanDesc') }}</span>
                </router-link>

                <!-- 充值 card (last card in every platform group) -->
                <router-link
                  to="/pricing?tab=recharge"
                  class="relative flex flex-col items-center justify-center rounded-xl border-2 border-dashed border-amber-300 bg-amber-50/50 p-4 transition-all duration-200 hover:border-amber-400 hover:bg-amber-50 hover:shadow-md dark:border-amber-700 dark:bg-amber-900/10 dark:hover:border-amber-600 dark:hover:bg-amber-900/20"
                >
                  <svg class="mb-2 h-8 w-8 text-amber-500 dark:text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span class="text-sm font-medium text-amber-700 dark:text-amber-400">{{ t('keys.recharge') }}</span>
                  <span class="mt-1 text-[11px] text-amber-500/70 dark:text-amber-500/50">{{ t('keys.rechargeDesc') }}</span>
                </router-link>
              </div>

              <div v-if="shouldShowPlatformToggle(platformGroup)" class="mt-3 flex justify-end">
                <button
                  class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:text-dark-300 dark:hover:border-primary-700 dark:hover:text-primary-400"
                  @click="togglePlatformExpanded(platformGroup.platform)"
                >
                  {{ isPlatformExpanded(platformGroup.platform) ? t('common.collapse') : t('common.more') }}
                  <span
                    v-if="!isPlatformExpanded(platformGroup.platform)"
                    class="text-gray-400 dark:text-dark-400"
                  >
                    +{{ hiddenPlatformGroupCount(platformGroup) }}
                  </span>
                </button>
              </div>
            </div>

            <!-- Empty state when no groups -->
            <div v-if="groupedByPlatform.length === 0 && !loading" class="rounded-xl border border-dashed border-gray-300 p-8 text-center dark:border-dark-600">
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('keys.noGroup') }}</p>
            </div>
        </section>

        <!-- Section 2: Keys 管理 -->
        <section ref="keysManagementSectionRef" class="keys-table-section">
            <div class="mb-4 flex items-center justify-between">
              <div>
                <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('keys.keysManagement') }}</h2>
                <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('keys.keysCount', { count: pagination.total }) }}</p>
              </div>
              <div class="flex gap-2">
                <button
                  @click="loadApiKeys"
                  :disabled="loading"
                  class="btn btn-secondary"
                  :title="t('common.refresh')"
                >
                  <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
                </button>
                <button @click="openCreateModal()" class="btn btn-primary" data-tour="keys-create-btn">
                  <Icon name="plus" size="md" class="mr-2" />
                  {{ t('keys.createKey') }}
                </button>
              </div>
            </div>

        <template v-if="renderKeysManagementSection">
        <DataTable :columns="columns" :data="apiKeys" :loading="loading">
          <template #cell-key="{ value, row }">
            <div class="flex items-center gap-2">
              <code class="code text-xs">
                {{ maskKey(value) }}
              </code>
              <button
                @click="copyToClipboard(value, row.id)"
                class="rounded-lg p-1 transition-colors hover:bg-gray-100 dark:hover:bg-dark-700"
                :class="
                  copiedKeyId === row.id
                    ? 'text-green-500'
                    : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'
                "
                :title="copiedKeyId === row.id ? t('keys.copied') : t('keys.copyToClipboard')"
              >
                <Icon
                  v-if="copiedKeyId === row.id"
                  name="check"
                  size="sm"
                  :stroke-width="2"
                />
                <Icon v-else name="clipboard" size="sm" />
              </button>
            </div>
          </template>

          <template #cell-name="{ value, row }">
            <div class="flex items-center gap-1.5">
              <span class="font-medium text-gray-900 dark:text-white">{{ value }}</span>
              <Icon
                v-if="row.ip_whitelist?.length > 0 || row.ip_blacklist?.length > 0"
                name="shield"
                size="sm"
                class="text-blue-500"
                :title="t('keys.ipRestrictionEnabled')"
              />
            </div>
          </template>

          <template #cell-group="{ row }">
            <div class="group/dropdown relative">
              <button
                :ref="(el) => setGroupButtonRef(row.id, el)"
                @click="openGroupSelector(row)"
                class="-mx-2 -my-1 flex cursor-pointer items-center gap-2 rounded-lg px-2 py-1 transition-all duration-200 hover:bg-gray-100 dark:hover:bg-dark-700"
                :title="t('keys.clickToChangeGroup')"
              >
                <GroupBadge
                  v-if="row.group"
                  :name="row.group.name"
                  :platform="row.group.platform"
                  :subscription-type="row.group.subscription_type"
                />
                <span v-else class="text-sm text-gray-400 dark:text-dark-500">{{
                  t('keys.noGroup')
                }}</span>
                <svg
                  class="h-3.5 w-3.5 text-gray-400 opacity-0 transition-opacity group-hover/dropdown:opacity-100"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M8.25 15L12 18.75 15.75 15m-7.5-6L12 5.25 15.75 9"
                  />
                </svg>
              </button>
            </div>
          </template>

          <template #cell-usage="{ row }">
            <div class="text-sm">
              <div class="flex items-center gap-1.5">
                <span class="text-gray-500 dark:text-gray-400">{{ t('keys.today') }}:</span>
                <span class="font-medium text-gray-900 dark:text-white">
                  {{ formatUsageValue(row.id, 'today_actual_cost') }}
                </span>
              </div>
              <div class="mt-0.5 flex items-center gap-1.5">
                <span class="text-gray-500 dark:text-gray-400">{{ t('keys.total') }}:</span>
                <span class="font-medium text-gray-900 dark:text-white">
                  {{ formatUsageValue(row.id, 'total_actual_cost') }}
                </span>
              </div>
            </div>
          </template>

          <template #cell-status="{ value }">
            <span :class="['badge', value === 'active' ? 'badge-success' : 'badge-gray']">
              {{ t('admin.accounts.status.' + value) }}
            </span>
          </template>

          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
          </template>

          <template #cell-actions="{ row }">
            <div class="flex items-center gap-1">
              <!-- Use Key Button -->
              <button
                @click="openUseKeyModal(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-green-50 hover:text-green-600 dark:hover:bg-green-900/20 dark:hover:text-green-400"
              >
                <Icon name="terminal" size="sm" />
                <span class="text-xs">{{ t('keys.useKey') }}</span>
              </button>
              <!-- Import to CC Switch Button -->
              <button
                @click="importToCcswitch(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-blue-50 hover:text-blue-600 dark:hover:bg-blue-900/20 dark:hover:text-blue-400"
              >
                <Icon name="upload" size="sm" />
                <span class="text-xs">{{ getCcSwitchButtonLabel(row) }}</span>
              </button>
              <!-- Toggle Status Button -->
              <button
                @click="toggleKeyStatus(row)"
                :class="[
                  'flex flex-col items-center gap-0.5 rounded-lg p-1.5 transition-colors',
                  row.status === 'active'
                    ? 'text-gray-500 hover:bg-yellow-50 hover:text-yellow-600 dark:hover:bg-yellow-900/20 dark:hover:text-yellow-400'
                    : 'text-gray-500 hover:bg-green-50 hover:text-green-600 dark:hover:bg-green-900/20 dark:hover:text-green-400'
                ]"
              >
                <Icon v-if="row.status === 'active'" name="ban" size="sm" />
                <Icon v-else name="checkCircle" size="sm" />
                <span class="text-xs">{{ row.status === 'active' ? t('keys.disable') : t('keys.enable') }}</span>
              </button>
              <!-- Edit Button -->
              <button
                @click="editKey(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-primary-600 dark:hover:bg-dark-700 dark:hover:text-primary-400"
              >
                <Icon name="edit" size="sm" />
                <span class="text-xs">{{ t('common.edit') }}</span>
              </button>
              <!-- Delete Button -->
              <button
                @click="confirmDelete(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
              >
                <Icon name="trash" size="sm" />
                <span class="text-xs">{{ t('common.delete') }}</span>
              </button>
            </div>
          </template>

          <template #empty>
            <EmptyState
              :title="t('keys.noKeysYet')"
              :description="t('keys.createFirstKey')"
              :action-text="t('keys.createKey')"
              @action="openCreateModal()"
            />
          </template>
        </DataTable>

        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
          class="mt-4"
        />
        </template>
        <div
          v-else
          class="rounded-xl border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-900"
        >
          <div class="space-y-3">
            <div class="h-4 w-40 animate-pulse rounded bg-gray-200 dark:bg-dark-700"></div>
            <div class="h-4 w-full animate-pulse rounded bg-gray-100 dark:bg-dark-800"></div>
            <div class="h-4 w-11/12 animate-pulse rounded bg-gray-100 dark:bg-dark-800"></div>
            <div class="h-4 w-10/12 animate-pulse rounded bg-gray-100 dark:bg-dark-800"></div>
          </div>
        </div>
        </section>

        <!-- Section 3: FAQ -->
        <section ref="faqSectionRef" class="keys-faq-section">
            <div class="mb-4">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('keys.faq.title') }}</h2>
            </div>
            <div
              v-if="renderFaqSection"
              class="divide-y divide-gray-200 rounded-xl border border-gray-200 bg-white dark:divide-dark-700 dark:border-dark-700 dark:bg-dark-900"
            >
              <details v-for="(item, idx) in faqItems" :key="idx" class="group">
                <summary class="flex cursor-pointer items-center justify-between px-5 py-4 text-sm font-medium text-gray-900 transition-colors hover:bg-gray-50 dark:text-white dark:hover:bg-dark-800 [&::-webkit-details-marker]:hidden">
                  {{ item.q }}
                  <svg class="h-4 w-4 flex-shrink-0 text-gray-400 transition-transform group-open:rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </summary>
                <div class="px-5 pb-4 text-sm leading-relaxed text-gray-600 dark:text-gray-400" v-html="item.a"></div>
              </details>
            </div>
            <div
              v-else
              class="rounded-xl border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-900"
            >
              <div class="space-y-3">
                <div class="h-4 w-2/3 animate-pulse rounded bg-gray-200 dark:bg-dark-700"></div>
                <div class="h-4 w-5/6 animate-pulse rounded bg-gray-100 dark:bg-dark-800"></div>
                <div class="h-4 w-3/4 animate-pulse rounded bg-gray-100 dark:bg-dark-800"></div>
              </div>
            </div>
        </section>

    <KeyEditorDialog
      v-if="showCreateModal || showEditModal"
      :show="showCreateModal || showEditModal"
      :groups="groups"
      :selected-key="showEditModal ? selectedKey : null"
      :initial-group-id="createModalGroupId"
      @close="closeModals"
      @saved="handleKeySaved"
    />

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      v-if="showDeleteDialog"
      :show="showDeleteDialog"
      :title="t('keys.deleteKey')"
      :message="t('keys.deleteConfirmMessage', { name: selectedKey?.name })"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="handleDelete"
      @cancel="showDeleteDialog = false"
    />

    <!-- Use Key Modal -->
    <UseKeyModal
      v-if="showUseKeyModal && selectedKey"
      :show="showUseKeyModal"
      :api-key="selectedKey.key"
      :base-url="publicSettings?.api_base_url || ''"
      :platform="selectedKey.group?.platform || null"
      @close="closeUseKeyModal"
    />

    <!-- CCS Client Selection Dialog for Antigravity -->
    <BaseDialog
      v-if="showCcsClientSelect"
      :show="showCcsClientSelect"
      :title="t('keys.ccsClientSelect.title')"
      width="wide"
      @close="closeCcsClientSelect"
    >
      <div class="space-y-4">
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ t('keys.ccsClientSelect.description') }}
        </p>
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <button
            v-for="option in availablePendingCcsApps"
            :key="option.app"
            @click="handleCcsClientSelect(option.app)"
            class="flex flex-col items-center gap-2 rounded-xl border-2 border-gray-200 p-4 transition-all hover:border-primary-500 hover:bg-primary-50 dark:border-dark-600 dark:hover:border-primary-500 dark:hover:bg-primary-900/20"
          >
            <Icon :name="option.icon" size="xl" class="text-gray-600 dark:text-gray-400" />
            <span class="font-medium text-gray-900 dark:text-white">{{ option.label }}</span>
            <span class="text-xs text-center text-gray-500 dark:text-gray-400">{{ option.description }}</span>
          </button>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end">
          <button @click="closeCcsClientSelect" class="btn btn-secondary">
            {{ t('common.cancel') }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Group Selector Dropdown (Teleported to body to avoid overflow clipping) -->
    <Teleport to="body">
      <div
        v-if="groupSelectorKeyId !== null && dropdownPosition"
        ref="dropdownRef"
        class="animate-in fade-in slide-in-from-top-2 fixed z-[100000020] w-64 overflow-hidden rounded-xl bg-white shadow-lg ring-1 ring-black/5 duration-200 dark:bg-dark-800 dark:ring-white/10"
        style="pointer-events: auto !important;"
        :style="{ top: dropdownPosition.top + 'px', left: dropdownPosition.left + 'px' }"
      >
        <div class="max-h-64 overflow-y-auto p-1.5">
          <button
            v-for="option in groupOptions"
            :key="option.value ?? 'null'"
            @click="changeGroup(selectedKeyForGroup!, option.value)"
            :class="[
              'flex w-full items-center justify-between rounded-lg px-3 py-2 text-sm transition-colors',
              selectedKeyForGroup?.group_id === option.value ||
              (!selectedKeyForGroup?.group_id && option.value === null)
                ? 'bg-primary-50 dark:bg-primary-900/20'
                : 'hover:bg-gray-100 dark:hover:bg-dark-700'
            ]"
            :title="option.description || undefined"
          >
            <GroupOptionItem
              :name="option.label"
              :platform="option.platform"
              :subscription-type="option.subscriptionType"
              :description="option.description"
              :selected="
                selectedKeyForGroup?.group_id === option.value ||
                (!selectedKeyForGroup?.group_id && option.value === null)
              "
            />
          </button>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
	import { ref, computed, onMounted, onUnmounted, nextTick, type ComponentPublicInstance } from 'vue'
	import { useI18n } from 'vue-i18n'
	import { useAppStore } from '@/stores/app'
	import { useClipboard } from '@/composables/useClipboard'
const { t } = useI18n()
import { keysAPI, usageAPI, userGroupsAPI } from '@/api'
import { defineChunkResilientAsyncComponent } from '@/router/asyncComponent'
	import Icon from '@/components/icons/Icon.vue'
	import GroupBadge from '@/components/common/GroupBadge.vue'
	import GroupOptionItem from '@/components/common/GroupOptionItem.vue'
	import PlatformIcon from '@/components/common/PlatformIcon.vue'

	// 弹窗组件异步加载，用户交互时才加载
	const DataTable = defineChunkResilientAsyncComponent(() => import('@/components/common/DataTable.vue'))
	const Pagination = defineChunkResilientAsyncComponent(() => import('@/components/common/Pagination.vue'))
	const EmptyState = defineChunkResilientAsyncComponent(() => import('@/components/common/EmptyState.vue'))
	const BaseDialog = defineChunkResilientAsyncComponent(() => import('@/components/common/BaseDialog.vue'))
	const ConfirmDialog = defineChunkResilientAsyncComponent(() => import('@/components/common/ConfirmDialog.vue'))
	const KeyEditorDialog = defineChunkResilientAsyncComponent(() => import('@/components/keys/KeyEditorDialog.vue'))
	const UseKeyModal = defineChunkResilientAsyncComponent(() => import('@/components/keys/UseKeyModal.vue'))

import type { ApiKey, Group, GroupPlatform, SelectOption, SubscriptionType } from '@/types'
import type { Column } from '@/components/common/types'
import type { BatchApiKeyUsageStats } from '@/api/usage'
import { formatDateTime } from '@/utils/format'
import {
  buildCcSwitchImportUrl,
  resolveCcSwitchImportTarget,
  type SupportedCcSwitchApp,
} from '@/utils/clientImport'

interface GroupOption extends SelectOption {
  value: number | null
  label: string
  description: string | null
  subscriptionType: SubscriptionType
  platform: GroupPlatform
}

const AVAILABLE_GROUPS_CACHE_TTL_MS = 5 * 60 * 1000
const AVAILABLE_GROUPS_SESSION_CACHE_KEY = 'sub2api:available-groups'
const INITIAL_VISIBLE_GROUPS_PER_PLATFORM = 6

let cachedAvailableGroups: Group[] | null = null
let cachedAvailableGroupsExpiresAt = 0
let availableGroupsPromise: Promise<Group[]> | null = null
const supportedAppCandidates: SupportedCcSwitchApp[] = ['claude', 'gemini', 'codex']

interface AvailableGroupsCachePayload {
  expiresAt: number
  groups: Group[]
}

const readAvailableGroupsSessionCache = (): AvailableGroupsCachePayload | null => {
  if (typeof window === 'undefined') {
    return null
  }

  try {
    const raw = window.sessionStorage.getItem(AVAILABLE_GROUPS_SESSION_CACHE_KEY)
    if (!raw) {
      return null
    }

    const payload = JSON.parse(raw) as AvailableGroupsCachePayload
    if (!Array.isArray(payload.groups) || payload.expiresAt <= Date.now()) {
      window.sessionStorage.removeItem(AVAILABLE_GROUPS_SESSION_CACHE_KEY)
      return null
    }

    return payload
  } catch {
    return null
  }
}

const writeAvailableGroupsSessionCache = (groups: Group[], expiresAt: number) => {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.sessionStorage.setItem(
      AVAILABLE_GROUPS_SESSION_CACHE_KEY,
      JSON.stringify({ groups, expiresAt })
    )
  } catch {
    // Ignore storage quota errors and keep the in-memory cache.
  }
}

async function getCachedAvailableGroups(force = false): Promise<Group[]> {
  const now = Date.now()

  if (!force && cachedAvailableGroups && cachedAvailableGroupsExpiresAt > now) {
    return cachedAvailableGroups
  }

  if (!force && !cachedAvailableGroups) {
    const cachedPayload = readAvailableGroupsSessionCache()
    if (cachedPayload) {
      cachedAvailableGroups = cachedPayload.groups
      cachedAvailableGroupsExpiresAt = cachedPayload.expiresAt
      return cachedAvailableGroups
    }
  }

  if (!force && availableGroupsPromise) {
    return availableGroupsPromise
  }

  availableGroupsPromise = userGroupsAPI.getAvailable()
    .then((data) => {
      cachedAvailableGroups = data
      cachedAvailableGroupsExpiresAt = Date.now() + AVAILABLE_GROUPS_CACHE_TTL_MS
      writeAvailableGroupsSessionCache(data, cachedAvailableGroupsExpiresAt)
      return data
    })
    .finally(() => {
      availableGroupsPromise = null
    })

  return availableGroupsPromise
}

const appStore = useAppStore()
const { copyToClipboard: clipboardCopy } = useClipboard()

const columns = computed<Column[]>(() => [
  { key: 'name', label: t('common.name'), sortable: true },
  { key: 'key', label: t('keys.apiKey'), sortable: false },
  { key: 'group', label: t('keys.group'), sortable: false },
  { key: 'usage', label: t('keys.usage'), sortable: false },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'created_at', label: t('keys.created'), sortable: true },
  { key: 'actions', label: t('common.actions'), sortable: false }
])

const apiKeys = ref<ApiKey[]>([])
const groups = ref<Group[]>([])
const loading = ref(false)
const usageStats = ref<Record<string, BatchApiKeyUsageStats>>({})

const pagination = ref({
  page: 1,
  page_size: 10,
  total: 0,
  pages: 0
})

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteDialog = ref(false)
const showUseKeyModal = ref(false)
const showCcsClientSelect = ref(false)
const createModalGroupId = ref<number | null>(null)
const renderKeysManagementSection = ref(false)
const renderFaqSection = ref(false)
const expandedPlatformSections = ref<GroupPlatform[]>([])
const pendingCcsRow = ref<ApiKey | null>(null)
const selectedKey = ref<ApiKey | null>(null)
const copiedKeyId = ref<number | null>(null)
const groupSelectorKeyId = ref<number | null>(null)
const publicSettings = computed(() => appStore.cachedPublicSettings)
const dropdownRef = ref<HTMLElement | null>(null)
const dropdownPosition = ref<{ top: number; left: number } | null>(null)
const groupButtonRefs = ref<Map<number, HTMLElement>>(new Map())
const keysManagementSectionRef = ref<HTMLElement | null>(null)
const faqSectionRef = ref<HTMLElement | null>(null)
let abortController: AbortController | null = null
let loadRequestId = 0
let usageStatsIdleHandle: number | ReturnType<typeof setTimeout> | null = null
let keysManagementObserver: IntersectionObserver | null = null
let faqObserver: IntersectionObserver | null = null

type IdleCallbackDeadline = {
  didTimeout: boolean
  timeRemaining: () => number
}
type IdleCallbackFn = (deadline: IdleCallbackDeadline) => void
type IdleCallbackOptions = {
  timeout?: number
}

const scheduleIdleCallback = (
  callback: IdleCallbackFn,
  options?: IdleCallbackOptions
): number | ReturnType<typeof setTimeout> => {
  if (typeof window.requestIdleCallback === 'function') {
    return window.requestIdleCallback(callback, options)
  }

  return window.setTimeout(() => {
    callback({ didTimeout: false, timeRemaining: () => 50 })
  }, 120)
}

const cancelIdleCallbackHandle = (handle: number | ReturnType<typeof setTimeout>) => {
  if (typeof window.cancelIdleCallback === 'function' && typeof handle === 'number') {
    window.cancelIdleCallback(handle)
  } else {
    clearTimeout(handle)
  }
}

// Get the currently selected key for group change
const selectedKeyForGroup = computed(() => {
  if (groupSelectorKeyId.value === null) return null
  return apiKeys.value.find((k) => k.id === groupSelectorKeyId.value) || null
})

const setGroupButtonRef = (keyId: number, el: Element | ComponentPublicInstance | null) => {
  if (el instanceof HTMLElement) {
    groupButtonRefs.value.set(keyId, el)
  } else {
    groupButtonRefs.value.delete(keyId)
  }
}

// Convert groups to Select options format with subscription type
const groupOptions = computed<GroupOption[]>(() =>
  groups.value.map((group) => ({
    value: group.id,
    label: group.name,
    description: group.description,
    subscriptionType: group.subscription_type,
    platform: group.platform
  }))
)

// Group cards organized by platform
const platformOrder: GroupPlatform[] = ['anthropic', 'openai', 'gemini', 'antigravity', 'multi']

interface PlatformGrouping {
  platform: GroupPlatform
  groups: Group[]
}

const groupedByPlatform = computed<PlatformGrouping[]>(() => {
  const platformMap = new Map<GroupPlatform, Group[]>()
  for (const group of groups.value) {
    const platform = group.platform || 'anthropic'
    if (!platformMap.has(platform)) {
      platformMap.set(platform, [])
    }
    platformMap.get(platform)!.push(group)
  }
  return platformOrder
    .filter(p => platformMap.has(p))
    .map(p => ({ platform: p, groups: platformMap.get(p)! }))
})

const isPlatformExpanded = (platform: GroupPlatform) => {
  return expandedPlatformSections.value.includes(platform)
}

const togglePlatformExpanded = (platform: GroupPlatform) => {
  expandedPlatformSections.value = isPlatformExpanded(platform)
    ? expandedPlatformSections.value.filter((item) => item !== platform)
    : [...expandedPlatformSections.value, platform]
}

const visiblePlatformGroups = (platformGroup: PlatformGrouping) => {
  if (
    isPlatformExpanded(platformGroup.platform) ||
    platformGroup.groups.length <= INITIAL_VISIBLE_GROUPS_PER_PLATFORM
  ) {
    return platformGroup.groups
  }

  return platformGroup.groups.slice(0, INITIAL_VISIBLE_GROUPS_PER_PLATFORM)
}

const hiddenPlatformGroupCount = (platformGroup: PlatformGrouping) => {
  return Math.max(platformGroup.groups.length - INITIAL_VISIBLE_GROUPS_PER_PLATFORM, 0)
}

const shouldShowPlatformToggle = (platformGroup: PlatformGrouping) => {
  return platformGroup.groups.length > INITIAL_VISIBLE_GROUPS_PER_PLATFORM
}

// Special tags that get highlighted styling (推荐, 暂不可用, 备用)
const specialTagKeywords: Record<string, string> = {
  '推荐': 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400',
  'Recommended': 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400',
  '暂不可用': 'bg-gray-200 text-gray-500 dark:bg-dark-500 dark:text-gray-400',
  'Unavailable': 'bg-gray-200 text-gray-500 dark:bg-dark-500 dark:text-gray-400',
  '备用': 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
  'Backup': 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
}

const getGroupSpecialTags = (group: Group) => {
  const tags = group.tags || []
  return tags
    .filter(tag => specialTagKeywords[tag])
    .map(tag => ({ text: tag, class: specialTagKeywords[tag] }))
}

const getGroupRegularTags = (group: Group) => {
  const tags = group.tags || []
  return tags.filter(tag => !specialTagKeywords[tag])
}

const isGroupUnavailable = (group: Group) => {
  const tags = group.tags || []
  return tags.some(tag => tag === '暂不可用' || tag === 'Unavailable')
}

const openCreateModal = (groupId: number | null = null) => {
  createModalGroupId.value = groupId
  selectedKey.value = null
  showEditModal.value = false
  showCreateModal.value = true
}

const quickCreateKey = (group: Group) => {
  openCreateModal(group.id)
}

// FAQ items
const faqItems = computed(() => [
  { q: t('keys.faq.q1'), a: t('keys.faq.a1') },
  { q: t('keys.faq.q2'), a: t('keys.faq.a2') },
  { q: t('keys.faq.q3'), a: t('keys.faq.a3') },
  { q: t('keys.faq.q4'), a: t('keys.faq.a4') },
  { q: t('keys.faq.q5'), a: t('keys.faq.a5') }
])

const maskKey = (key: string): string => {
  if (key.length <= 12) return key
  return `${key.slice(0, 8)}...${key.slice(-4)}`
}

const copyToClipboard = async (text: string, keyId: number) => {
  const success = await clipboardCopy(text, t('keys.copied'))
  if (success) {
    copiedKeyId.value = keyId
    setTimeout(() => {
      copiedKeyId.value = null
    }, 800)
  }
}

const formatUsageValue = (
  keyId: number,
  field: keyof Pick<BatchApiKeyUsageStats, 'today_actual_cost' | 'total_actual_cost'>
) => {
  const stats = usageStats.value[keyId]
  if (!stats) {
    return '...'
  }

  return `$${stats[field].toFixed(4)}`
}

const isAbortError = (error: unknown) => {
  if (!error || typeof error !== 'object') return false
  const { name, code } = error as { name?: string; code?: string }
  return name === 'AbortError' || code === 'ERR_CANCELED'
}

const loadUsageStatsForKeys = async (
  keyList: ApiKey[],
  signal: AbortSignal,
  requestId: number
) => {
  if (keyList.length === 0) {
    if (requestId === loadRequestId) {
      usageStats.value = {}
    }
    return
  }

  try {
    const usageResponse = await usageAPI.getDashboardApiKeysUsage(
      keyList.map((key) => key.id),
      { signal }
    )
    if (signal.aborted || requestId !== loadRequestId) return
    usageStats.value = usageResponse.stats
  } catch (error) {
    if (!isAbortError(error)) {
      console.error('Failed to load usage stats:', error)
    }
  }
}

const loadApiKeys = async () => {
  abortController?.abort()
  const controller = new AbortController()
  abortController = controller
  const { signal } = controller
  const requestId = ++loadRequestId
  loading.value = true
  usageStats.value = {}
  try {
    const response = await keysAPI.list(pagination.value.page, pagination.value.page_size, {
      signal
    })
    if (signal.aborted || requestId !== loadRequestId) return
    apiKeys.value = response.items
    pagination.value.total = response.total
    pagination.value.pages = response.pages

    if (usageStatsIdleHandle !== null) {
      cancelIdleCallbackHandle(usageStatsIdleHandle)
    }

    usageStatsIdleHandle = scheduleIdleCallback(() => {
      usageStatsIdleHandle = null
      void loadUsageStatsForKeys(response.items, signal, requestId)
    }, { timeout: 800 })
  } catch (error) {
    if (isAbortError(error)) {
      return
    }
    appStore.showError(t('keys.failedToLoad'))
  } finally {
    if (abortController === controller) {
      loading.value = false
    }
  }
}

const loadGroups = async () => {
  try {
    groups.value = await getCachedAvailableGroups()
  } catch (error) {
    console.error('Failed to load groups:', error)
  }
}

const loadPublicSettings = async () => {
  try {
    await appStore.fetchPublicSettings()
  } catch (error) {
    console.error('Failed to load public settings:', error)
  }
}

const openUseKeyModal = (key: ApiKey) => {
  selectedKey.value = key
  showUseKeyModal.value = true
  if (!publicSettings.value) {
    void loadPublicSettings()
  }
}

const closeUseKeyModal = () => {
  showUseKeyModal.value = false
  selectedKey.value = null
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  loadApiKeys()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.value.page_size = pageSize
  pagination.value.page = 1
  loadApiKeys()
}

const editKey = (key: ApiKey) => {
  selectedKey.value = key
  createModalGroupId.value = null
  showCreateModal.value = false
  showEditModal.value = true
}

const toggleKeyStatus = async (key: ApiKey) => {
  const newStatus = key.status === 'active' ? 'inactive' : 'active'
  try {
    await keysAPI.toggleStatus(key.id, newStatus)
    appStore.showSuccess(
      newStatus === 'active' ? t('keys.keyEnabledSuccess') : t('keys.keyDisabledSuccess')
    )
    loadApiKeys()
  } catch (error) {
    appStore.showError(t('keys.failedToUpdateStatus'))
  }
}

const openGroupSelector = (key: ApiKey) => {
  if (groupSelectorKeyId.value === key.id) {
    groupSelectorKeyId.value = null
    dropdownPosition.value = null
  } else {
    const buttonEl = groupButtonRefs.value.get(key.id)
    if (buttonEl) {
      const rect = buttonEl.getBoundingClientRect()
      dropdownPosition.value = {
        top: rect.bottom + 4,
        left: rect.left
      }
    }
    groupSelectorKeyId.value = key.id
  }
}

const changeGroup = async (key: ApiKey, newGroupId: number | null) => {
  groupSelectorKeyId.value = null
  dropdownPosition.value = null
  if (key.group_id === newGroupId) return

  try {
    await keysAPI.update(key.id, { group_id: newGroupId })
    appStore.showSuccess(t('keys.groupChangedSuccess'))
    loadApiKeys()
  } catch (error) {
    appStore.showError(t('keys.failedToChangeGroup'))
  }
}

const closeGroupSelector = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  // Check if click is inside the dropdown or the trigger button
  if (!target.closest('.group\\/dropdown') && !dropdownRef.value?.contains(target)) {
    groupSelectorKeyId.value = null
    dropdownPosition.value = null
  }
}

const confirmDelete = (key: ApiKey) => {
  selectedKey.value = key
  showDeleteDialog.value = true
}

/**
 * 处理删除 API Key 的操作
 * 优化：错误处理改进，优先显示后端返回的具体错误消息（如权限不足等），
 * 若后端未返回消息则显示默认的国际化文本
 */
const handleDelete = async () => {
  if (!selectedKey.value) return

  try {
    await keysAPI.delete(selectedKey.value.id)
    appStore.showSuccess(t('keys.keyDeletedSuccess'))
    showDeleteDialog.value = false
    loadApiKeys()
  } catch (error: any) {
    // 优先使用后端返回的错误消息，提供更具体的错误信息给用户
    const errorMsg = error?.message || t('keys.failedToDelete')
    appStore.showError(errorMsg)
  }
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  createModalGroupId.value = null
  selectedKey.value = null
}

const handleKeySaved = () => {
  void loadApiKeys()
}

const availablePendingCcsApps = computed(() => {
  if (!pendingCcsRow.value) {
    return []
  }

  const baseUrl = publicSettings.value?.api_base_url || window.location.origin
  const candidates: Array<{ app: SupportedCcSwitchApp; label: string; description: string; icon: 'terminal' | 'sparkles' | 'cpu' }> = [
    {
      app: 'claude',
      label: t('keys.ccsClientSelect.claudeCode'),
      description: t('keys.ccsClientSelect.claudeCodeDesc'),
      icon: 'terminal',
    },
    {
      app: 'gemini',
      label: t('keys.ccsClientSelect.geminiCli'),
      description: t('keys.ccsClientSelect.geminiCliDesc'),
      icon: 'sparkles',
    },
    {
      app: 'codex',
      label: t('keys.ccsClientSelect.codexCli'),
      description: t('keys.ccsClientSelect.codexCliDesc'),
      icon: 'cpu',
    },
  ]

  return candidates.filter(candidate =>
    Boolean(resolveCcSwitchImportTarget(pendingCcsRow.value as ApiKey, candidate.app, baseUrl))
  )
})

const importToCcswitch = (row: ApiKey) => {
  const baseUrl = publicSettings.value?.api_base_url || window.location.origin
  const supportedApps = supportedAppCandidates.filter((app) =>
    Boolean(resolveCcSwitchImportTarget(row, app, baseUrl))
  )

  if (supportedApps.length === 0) {
    appStore.showError(t('tutorial.configExport.ccSwitchUnsupported'))
    return
  }

  if (supportedApps.length === 1) {
    executeCcsImport(row, supportedApps[0])
    return
  }

  pendingCcsRow.value = row
  showCcsClientSelect.value = true
}

const getCcSwitchButtonLabel = (row: ApiKey) => {
  const baseUrl = publicSettings.value?.api_base_url || window.location.origin
  const supportedApps = supportedAppCandidates.filter((app) =>
    Boolean(resolveCcSwitchImportTarget(row, app, baseUrl))
  )

  if (supportedApps.length !== 1) {
    return t('keys.chooseCcSwitchImport')
  }

  const app = supportedApps[0]
  if (app === 'claude') {
    return t('keys.importToCcSwitchApp', { app: t('keys.ccsClientSelect.claudeCode') })
  }
  if (app === 'gemini') {
    return t('keys.importToCcSwitchApp', { app: t('keys.ccsClientSelect.geminiCli') })
  }
  return t('keys.importToCcSwitchApp', { app: t('keys.ccsClientSelect.codexCli') })
}

const executeCcsImport = (row: ApiKey, clientType: SupportedCcSwitchApp) => {
  const baseUrl = publicSettings.value?.api_base_url || window.location.origin
  const deeplink = buildCcSwitchImportUrl(row, clientType, baseUrl)
  if (!deeplink) {
    appStore.showError(t('tutorial.configExport.ccSwitchUnsupported'))
    return
  }

  try {
    window.open(deeplink, '_self')

    // Check if the protocol handler worked by detecting if we're still focused
    setTimeout(() => {
      if (document.hasFocus()) {
        // Still focused means the protocol handler likely failed
        appStore.showError(t('keys.ccSwitchNotInstalled'))
      }
    }, 100)
  } catch (error) {
    appStore.showError(t('keys.ccSwitchNotInstalled'))
  }
}

const handleCcsClientSelect = (clientType: SupportedCcSwitchApp) => {
  if (pendingCcsRow.value) {
    executeCcsImport(pendingCcsRow.value, clientType)
  }
  showCcsClientSelect.value = false
  pendingCcsRow.value = null
}

const closeCcsClientSelect = () => {
  showCcsClientSelect.value = false
  pendingCcsRow.value = null
}

const observeDeferredSection = (
  target: HTMLElement | null,
  onVisible: () => void,
  rootMargin: string
) => {
  if (!target) {
    onVisible()
    return null
  }

  if (typeof window === 'undefined' || typeof window.IntersectionObserver !== 'function') {
    onVisible()
    return null
  }

  const observer = new window.IntersectionObserver(
    (entries) => {
      if (entries.some((entry) => entry.isIntersecting)) {
        onVisible()
        observer.disconnect()
      }
    },
    { rootMargin }
  )

  observer.observe(target)
  return observer
}

onMounted(async () => {
  loadApiKeys()
  loadGroups()
  document.addEventListener('click', closeGroupSelector)

  await nextTick()

  keysManagementObserver = observeDeferredSection(
    keysManagementSectionRef.value,
    () => {
      renderKeysManagementSection.value = true
    },
    '320px 0px'
  )

  faqObserver = observeDeferredSection(
    faqSectionRef.value,
    () => {
      renderFaqSection.value = true
    },
    '240px 0px'
  )
})

onUnmounted(() => {
  if (usageStatsIdleHandle !== null) {
    cancelIdleCallbackHandle(usageStatsIdleHandle)
  }
  abortController?.abort()
  keysManagementObserver?.disconnect()
  faqObserver?.disconnect()
  document.removeEventListener('click', closeGroupSelector)
})
</script>

<style scoped>
.keys-platform-group,
.keys-table-section,
.keys-faq-section {
  content-visibility: auto;
}

.keys-platform-group {
  contain-intrinsic-size: 420px;
}

.keys-table-section {
  contain-intrinsic-size: 720px;
}

.keys-faq-section {
  contain-intrinsic-size: 360px;
}
</style>
