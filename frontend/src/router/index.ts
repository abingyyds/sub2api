/**
 * Vue Router configuration for Sub2API frontend
 * Defines all application routes with lazy loading and navigation guards
 */

import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useNavigationLoadingState } from '@/composables/useNavigationLoading'
import { useRoutePrefetch } from '@/composables/useRoutePrefetch'
import {
  CHUNK_RELOAD_STORAGE_KEY,
  installChunkRecoveryListeners,
  isChunkLoadError,
  recoverFromChunkError
} from './chunkLoad'

/**
 * Route definitions with lazy loading
 */
const routes: RouteRecordRaw[] = [
  // ==================== Setup Routes ====================
  {
    path: '/setup',
    name: 'Setup',
    component: () => import('@/views/setup/SetupWizardView.vue'),
    meta: {
      requiresAuth: false,
      title: 'Setup'
    }
  },

  // ==================== Public Routes ====================
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/PublicHomeView.vue'),
    meta: {
      requiresAuth: false,
      title: 'Home'
    }
  },
  {
    path: '/legal/tokushoho',
    name: 'Tokushoho',
    component: () => import('@/views/legal/TokushohoView.vue'),
    meta: {
      requiresAuth: false,
      title: '特定商取引法に基づく表記'
    }
  },
  {
    path: '/legal/disclosure',
    name: 'Disclosure',
    component: () => import('@/views/legal/DisclosureView.vue'),
    meta: {
      requiresAuth: false,
      title: '商業披露'
    }
  },
  {
    path: '/legal/terms',
    name: 'Terms',
    component: () => import('@/views/legal/TermsView.vue'),
    meta: {
      requiresAuth: false,
      title: '用户协议'
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: {
      requiresAuth: false,
      title: 'Login'
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: {
      requiresAuth: false,
      title: 'Register'
    }
  },
  {
    path: '/email-verify',
    name: 'EmailVerify',
    component: () => import('@/views/auth/EmailVerifyView.vue'),
    meta: {
      requiresAuth: false,
      title: 'Verify Email'
    }
  },
  {
    path: '/auth/callback',
    name: 'OAuthCallback',
    component: () => import('@/views/auth/OAuthCallbackView.vue'),
    meta: {
      requiresAuth: false,
      title: 'OAuth Callback'
    }
  },
  {
    path: '/auth/linuxdo/callback',
    name: 'LinuxDoOAuthCallback',
    component: () => import('@/views/auth/LinuxDoCallbackView.vue'),
    meta: {
      requiresAuth: false,
      title: 'LinuxDo OAuth Callback'
    }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPasswordView.vue'),
    meta: {
      requiresAuth: false,
      title: 'Forgot Password'
    }
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/views/auth/ResetPasswordView.vue'),
    meta: {
      requiresAuth: false,
      title: 'Reset Password'
    }
  },

  // ==================== User Routes ====================
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/user/DashboardView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Dashboard',
      titleKey: 'dashboard.title',
      descriptionKey: 'dashboard.welcomeMessage'
    }
  },
  {
    path: '/keys',
    name: 'Keys',
    component: () => import('@/views/user/KeysView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'API Keys',
      titleKey: 'keys.title',
      descriptionKey: 'keys.description'
    }
  },
  {
    path: '/usage',
    name: 'Usage',
    component: () => import('@/views/user/UsageView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Usage Records',
      titleKey: 'usage.title',
      descriptionKey: 'usage.description'
    }
  },
  {
    path: '/redeem',
    name: 'Redeem',
    component: () => import('@/views/user/RedeemView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Redeem Code',
      titleKey: 'redeem.title',
      descriptionKey: 'redeem.description'
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/user/ProfileView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Profile',
      titleKey: 'profile.title',
      descriptionKey: 'profile.description'
    }
  },
  {
    path: '/pricing',
    name: 'Pricing',
    component: () => import('@/views/user/PricingView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Pricing',
      titleKey: 'pricing.title',
      descriptionKey: 'pricing.subtitle'
    }
  },
  {
    path: '/recharge',
    redirect: '/pricing?tab=recharge'
  },
  {
    path: '/tutorial',
    name: 'Tutorial',
    component: () => import('@/views/user/TutorialView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Tutorial',
      titleKey: 'tutorial.title',
      descriptionKey: 'tutorial.subtitle'
    }
  },
  {
    path: '/model-plaza',
    name: 'ModelPlaza',
    component: () => import('@/views/user/ModelPlazaView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Model Plaza',
      titleKey: 'modelPlaza.title',
      descriptionKey: 'modelPlaza.subtitle'
    }
  },
  {
    path: '/subscriptions',
    name: 'Subscriptions',
    component: () => import('@/views/user/SubscriptionsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'My Subscriptions',
      titleKey: 'userSubscriptions.title',
      descriptionKey: 'userSubscriptions.description'
    }
  },
  {
    path: '/referral',
    name: 'Referral',
    component: () => import('@/views/user/ReferralView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Referral',
      titleKey: 'referral.title',
      descriptionKey: 'referral.description'
    }
  },
  {
    path: '/agent',
    name: 'AgentDashboard',
    component: () => import('@/views/user/AgentDashboardView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Agent Dashboard',
      titleKey: 'agent.dashboard.title',
      descriptionKey: 'agent.dashboard.description'
    }
  },
  {
    path: '/subsites',
    name: 'SubSiteCenter',
    component: () => import('@/views/user/SubSiteCenterView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'SubSite Center'
    }
  },
  {
    path: '/subsite-admin',
    name: 'SubSiteAdminEntry',
    redirect: () => {
      // 入口由守卫动态重定向到第一个 owned 分站的 dashboard
      return { path: '/subsite-admin/_entry' }
    },
    meta: { requiresAuth: true, requiresSubSiteOwner: true }
  },
  {
    path: '/subsite-admin/_entry',
    name: 'SubSiteAdminRedirect',
    component: () => import('@/views/subsiteAdmin/EntryRedirectView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Admin' }
  },
  {
    path: '/subsite-admin/:siteId(\\d+)/dashboard',
    name: 'SubSiteAdminDashboard',
    component: () => import('@/views/subsiteAdmin/DashboardView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Dashboard' }
  },
  {
    path: '/subsite-admin/:siteId(\\d+)/users',
    name: 'SubSiteAdminUsers',
    component: () => import('@/views/subsiteAdmin/UsersView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Users' }
  },
  {
    path: '/subsite-admin/:siteId(\\d+)/orders',
    name: 'SubSiteAdminOrders',
    component: () => import('@/views/subsiteAdmin/OrdersView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Orders' }
  },
  {
    path: '/subsite-admin/:siteId(\\d+)/usage',
    name: 'SubSiteAdminUsage',
    component: () => import('@/views/subsiteAdmin/UsageView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Usage' }
  },
  {
    path: '/subsite-admin/:siteId(\\d+)/ledger',
    name: 'SubSiteAdminLedger',
    component: () => import('@/views/subsiteAdmin/LedgerView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Ledger' }
  },
  {
    path: '/subsite-admin/:siteId(\\d+)/settings',
    name: 'SubSiteAdminSettings',
    component: () => import('@/views/subsiteAdmin/SettingsView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Settings' }
  },
  {
    path: '/subsite-admin/:siteId(\\d+)/withdraw',
    name: 'SubSiteAdminWithdraw',
    component: () => import('@/views/subsiteAdmin/WithdrawView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Withdraw' }
  },
  {
    path: '/subsite-admin/:siteId(\\d+)/payment-config',
    name: 'SubSiteAdminPaymentConfig',
    component: () => import('@/views/subsiteAdmin/PaymentConfigView.vue'),
    meta: { requiresAuth: true, requiresSubSiteOwner: true, title: 'SubSite Payment Config' }
  },
  {
    path: '/agent/sub-users',
    name: 'AgentSubUsers',
    component: () => import('@/views/user/AgentSubUsersView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Sub Users',
      titleKey: 'agent.subUsers.title',
      descriptionKey: 'agent.subUsers.description'
    }
  },
  {
    path: '/agent/financial-logs',
    name: 'AgentFinancialLogs',
    component: () => import('@/views/user/AgentFinancialLogsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Financial Logs',
      titleKey: 'agent.financialLogs.title',
      descriptionKey: 'agent.financialLogs.description'
    }
  },
  {
    path: '/agent/commissions',
    name: 'AgentCommissions',
    component: () => import('@/views/user/AgentCommissionsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Commissions',
      titleKey: 'agent.commissions.title',
      descriptionKey: 'agent.commissions.description'
    }
  },
  {
    path: '/changelog',
    name: 'Changelog',
    component: () => import('@/views/user/ChangelogView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Changelog',
      titleKey: 'changelog.title',
      descriptionKey: 'changelog.description'
    }
  },
  {
    path: '/order-history',
    name: 'OrderHistory',
    component: () => import('@/views/user/OrderHistoryView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Order History',
      titleKey: 'nav.orderHistory'
    }
  },
  {
    path: '/invoice',
    redirect: '/order-history'
  },
  {
    path: '/contact',
    redirect: '/dashboard'
  },

  // ==================== Admin Routes ====================
  {
    path: '/admin',
    redirect: '/admin/dashboard'
  },
  {
    path: '/admin/dashboard',
    name: 'AdminDashboard',
    component: () => import('@/views/admin/DashboardView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Admin Dashboard',
      titleKey: 'admin.dashboard.title',
      descriptionKey: 'admin.dashboard.description'
    }
  },
  {
    path: '/admin/ops',
    name: 'AdminOps',
    component: () => import('@/views/admin/ops/OpsDashboard.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Ops Monitoring',
      titleKey: 'admin.ops.title',
      descriptionKey: 'admin.ops.description'
    }
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    component: () => import('@/views/admin/UsersView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'User Management',
      titleKey: 'admin.users.title',
      descriptionKey: 'admin.users.description'
    }
  },
  {
    path: '/admin/groups',
    name: 'AdminGroups',
    component: () => import('@/views/admin/GroupsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Group Management',
      titleKey: 'admin.groups.title',
      descriptionKey: 'admin.groups.description'
    }
  },
  {
    path: '/admin/subscriptions',
    name: 'AdminSubscriptions',
    component: () => import('@/views/admin/SubscriptionsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Subscription Management',
      titleKey: 'admin.subscriptions.title',
      descriptionKey: 'admin.subscriptions.description'
    }
  },
  {
    path: '/admin/accounts',
    name: 'AdminAccounts',
    component: () => import('@/views/admin/AccountsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Account Management',
      titleKey: 'admin.accounts.title',
      descriptionKey: 'admin.accounts.description'
    }
  },
  {
    path: '/admin/proxies',
    name: 'AdminProxies',
    component: () => import('@/views/admin/ProxiesView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Proxy Management',
      titleKey: 'admin.proxies.title',
      descriptionKey: 'admin.proxies.description'
    }
  },
  {
    path: '/admin/redeem',
    name: 'AdminRedeem',
    component: () => import('@/views/admin/RedeemView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Redeem Code Management',
      titleKey: 'admin.redeem.title',
      descriptionKey: 'admin.redeem.description'
    }
  },
  {
    path: '/admin/promo-codes',
    name: 'AdminPromoCodes',
    component: () => import('@/views/admin/PromoCodesView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Promo Code Management',
      titleKey: 'admin.promo.title',
      descriptionKey: 'admin.promo.description'
    }
  },
  {
    path: '/admin/settings',
    name: 'AdminSettings',
    component: () => import('@/views/admin/SettingsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'System Settings',
      titleKey: 'admin.settings.title',
      descriptionKey: 'admin.settings.description'
    }
  },
  {
    path: '/admin/usage',
    name: 'AdminUsage',
    component: () => import('@/views/admin/UsageView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Usage Records',
      titleKey: 'admin.usage.title',
      descriptionKey: 'admin.usage.description'
    }
  },
  {
    path: '/admin/announcements',
    name: 'AdminAnnouncements',
    component: () => import('@/views/admin/AnnouncementsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Announcements',
      titleKey: 'admin.announcements.title',
      descriptionKey: 'admin.announcements.description'
    }
  },
  {
    path: '/admin/referrals',
    name: 'AdminReferrals',
    component: () => import('@/views/admin/ReferralsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Referral Management',
      titleKey: 'admin.referrals.title',
      descriptionKey: 'admin.referrals.description'
    }
  },
  {
    path: '/admin/organizations',
    name: 'AdminOrganizations',
    component: () => import('@/views/admin/OrganizationsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Organization Management',
      titleKey: 'admin.organizations.title',
      descriptionKey: 'admin.organizations.description'
    }
  },
  {
    path: '/admin/invite-codes',
    name: 'AdminInviteCodes',
    component: () => import('@/views/admin/InviteCodesView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Invite Code Management',
      titleKey: 'admin.inviteCodes.title',
      descriptionKey: 'admin.inviteCodes.description'
    }
  },
  {
    path: '/admin/orders',
    name: 'AdminOrders',
    component: () => import('@/views/admin/OrdersView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Order Management',
      titleKey: 'admin.orders.title',
      descriptionKey: 'admin.orders.description'
    }
  },
  {
    path: '/admin/agents',
    name: 'AdminAgents',
    component: () => import('@/views/admin/AgentsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Agent Management',
      titleKey: 'admin.agents.title',
      descriptionKey: 'admin.agents.description'
    }
  },
  {
    path: '/admin/subsites',
    name: 'AdminSubSites',
    component: () => import('@/views/admin/SubSitesView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Sub-site Management',
      descriptionKey: 'Sub-sites'
    }
  },
  {
    path: '/admin/withdraws',
    name: 'AdminWithdraws',
    component: () => import('@/views/admin/WithdrawsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Withdraw Management',
      descriptionKey: 'Withdraws'
    }
  },

  // ==================== Org Admin Routes ====================
  {
    path: '/org',
    redirect: '/org/dashboard'
  },
  {
    path: '/org/dashboard',
    name: 'OrgDashboard',
    component: () => import('@/views/org/OrgDashboardView.vue'),
    meta: {
      requiresAuth: true,
      requiresOrgAdmin: true,
      title: 'Organization Dashboard',
      titleKey: 'org.dashboard.title',
      descriptionKey: 'org.dashboard.description'
    }
  },
  {
    path: '/org/members',
    name: 'OrgMembers',
    component: () => import('@/views/org/OrgMembersView.vue'),
    meta: {
      requiresAuth: true,
      requiresOrgAdmin: true,
      title: 'Organization Members',
      titleKey: 'org.members.title',
      descriptionKey: 'org.members.description'
    }
  },
  {
    path: '/org/projects',
    name: 'OrgProjects',
    component: () => import('@/views/org/OrgProjectsView.vue'),
    meta: {
      requiresAuth: true,
      requiresOrgAdmin: true,
      title: 'Organization Projects',
      titleKey: 'org.projects.title',
      descriptionKey: 'org.projects.description'
    }
  },
  {
    path: '/org/audit-logs',
    name: 'OrgAuditLogs',
    component: () => import('@/views/org/OrgAuditLogView.vue'),
    meta: {
      requiresAuth: true,
      requiresOrgAdmin: true,
      title: 'Audit Logs',
      titleKey: 'org.audit.title',
      descriptionKey: 'org.audit.description'
    }
  },

  // ==================== 404 Not Found ====================
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundView.vue'),
    meta: {
      title: '404 Not Found'
    }
  }
]

/**
 * Create router instance
 */
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    // Scroll to saved position when using browser back/forward
    if (savedPosition) {
      return savedPosition
    }
    // Scroll to top for new routes
    return { top: 0 }
  }
})

/**
 * Navigation guard: Authentication check
 */
let authInitialized = false
let lastRequestedLocation = window.location.pathname + window.location.search + window.location.hash
installChunkRecoveryListeners(() => lastRequestedLocation)

// 初始化导航加载状态和预加载
const navigationLoading = useNavigationLoadingState()
// 延迟初始化预加载，传入 router 实例
let routePrefetch: ReturnType<typeof useRoutePrefetch> | null = null

router.beforeEach(async (to, _from, next) => {
  lastRequestedLocation = to.fullPath
  // 开始导航加载状态
  navigationLoading.startNavigation()

  const authStore = useAuthStore()

  // Restore auth state from localStorage on first navigation (page refresh)
  if (!authInitialized) {
    await authStore.checkAuth()
    authInitialized = true
  }

  // Set page title
  const appStore = useAppStore()
  const siteName = appStore.siteName || 'cCoder.me'
  if (to.meta.title) {
    document.title = `${to.meta.title} - ${siteName}`
  } else {
    document.title = siteName
  }

  // Check if route requires authentication
  const requiresAuth = to.meta.requiresAuth !== false // Default to true
  const requiresAdmin = to.meta.requiresAdmin === true

  // If route doesn't require auth, allow access
  if (!requiresAuth) {
    // If already authenticated and trying to access login/register, redirect to appropriate dashboard
    if (authStore.isAuthenticated && (to.path === '/login' || to.path === '/register')) {
      if (!authStore.isAdmin && !authStore.isOrgAdmin && appStore.cachedPublicSettings?.is_subsite) {
        const sites = authStore.ownedSites.length > 0 ? authStore.ownedSites : await authStore.refreshOwnedSites().catch(() => [])
        if (sites.length > 0) {
          next(`/subsite-admin/${sites[0].id}/dashboard`)
          return
        }
      }
      // Admin users go to admin dashboard, org admins go to org dashboard, regular users go to user dashboard
      if (authStore.isAdmin) {
        next('/admin/dashboard')
      } else if (authStore.isOrgAdmin) {
        next('/org/dashboard')
      } else {
        next('/pricing')
      }
      return
    }
    next()
    return
  }

  // Route requires authentication
  if (!authStore.isAuthenticated) {
    // Not authenticated, redirect to login
    next({
      path: '/login',
      query: { redirect: to.fullPath } // Save intended destination
    })
    return
  }

  // Check admin requirement
  if (requiresAdmin && !authStore.isAdmin) {
    // User is authenticated but not admin, redirect to user dashboard
    next('/dashboard')
    return
  }

  // Check org admin requirement
  const requiresOrgAdmin = to.meta.requiresOrgAdmin === true
  if (requiresOrgAdmin && !authStore.isOrgAdmin) {
    next(authStore.isAdmin ? '/admin/dashboard' : '/dashboard')
    return
  }

  // Check sub-site owner requirement
  const requiresSubSiteOwner = to.meta.requiresSubSiteOwner === true
  if (requiresSubSiteOwner) {
    // 如果 owned sites 尚未加载（例如首次直接访问 URL），尝试拉一次
    if (authStore.ownedSites.length === 0) {
      try {
        await authStore.refreshOwnedSites()
      } catch { /* 忽略，下一步会拒绝 */ }
    }
    if (!authStore.isSubSiteOwner) {
      next('/dashboard')
      return
    }
    const rawSiteId = to.params.siteId
    if (rawSiteId !== undefined && rawSiteId !== '') {
      const siteId = Number(rawSiteId)
      if (!authStore.ownedSites.some((s) => s.id === siteId)) {
        next('/dashboard')
        return
      }
    }
  }

  // Sub-admin cannot access account management
  if (authStore.isAdmin && !authStore.isFullAdmin && to.path.startsWith('/admin/accounts')) {
    next('/admin/dashboard')
    return
  }

  // 简易模式下限制访问某些页面
  if (authStore.isSimpleMode) {
    const restrictedPaths = [
      '/admin/groups',
      '/admin/subscriptions',
      '/admin/redeem',
      '/subscriptions',
      '/redeem'
    ]

    if (restrictedPaths.some((path) => to.path.startsWith(path))) {
      // 简易模式下访问受限页面,重定向到仪表板
      next(authStore.isAdmin ? '/admin/dashboard' : '/dashboard')
      return
    }
  }

  // All checks passed, allow navigation
  next()
})

/**
 * Navigation guard: End loading and trigger prefetch
 */
router.afterEach((to) => {
  // 结束导航加载状态
  navigationLoading.endNavigation()
  sessionStorage.removeItem(CHUNK_RELOAD_STORAGE_KEY)

  // 懒初始化预加载（首次导航时创建，传入 router 实例）
  if (!routePrefetch) {
    routePrefetch = useRoutePrefetch(router)
  }
  // 触发路由预加载（在浏览器空闲时执行）
  routePrefetch.triggerPrefetch(to)
})

/**
 * Navigation guard: Error handling
 * Handles dynamic import failures caused by deployment updates
 */
router.onError((error) => {
  navigationLoading.endNavigation()
  console.error('Router error:', error)

  if (recoverFromChunkError(lastRequestedLocation, error)) {
    console.warn('Chunk load error detected, hard reloading target route to fetch latest version...')
  } else if (isChunkLoadError(error)) {
    console.error('Chunk load error persists after reload. Please clear browser cache.')
  }
})

export default router
