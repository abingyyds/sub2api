import type { PaymentPlan } from '@/api/payment'

export interface ThemeHomeProps {
  siteName: string
  siteLogo: string
  siteSubtitle: string
  docUrl: string
  heroTitle: string
  heroDescription: string
  ctaText: string
  featureTags: string[]
  registrationNotice: string
  allowSubSiteOpen: boolean
  subSiteOpenPrice: string
  isAuthenticated: boolean
  dashboardPath: string
  userInitial: string
  plans: PaymentPlan[]
  creatingOrder: boolean
  githubUrl: string
  currentYear: number
  isDark: boolean
}

export interface ThemeHomeEmits {
  (e: 'toggle-theme'): void
  (e: 'buy-plan', plan: PaymentPlan): void
}
