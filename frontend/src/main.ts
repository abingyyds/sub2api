import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n, { getLocale, loadLocaleMessages } from './i18n'
import { useAppStore } from '@/stores/app'
import './style.css'
import './styles/animations.css'

async function bootstrap() {
  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)

  // Initialize settings from injected config BEFORE mounting (prevents flash)
  // This must happen after pinia is installed but before router and i18n
  const appStore = useAppStore()
  appStore.initFromInjectedConfig()

  // Load only the active locale on the critical path.
  const initialLocale = getLocale()
  await loadLocaleMessages(initialLocale)
  document.documentElement.setAttribute('lang', initialLocale)

  // Set document title immediately after config is loaded
  if (appStore.siteName && appStore.siteName !== 'Sub2API') {
    document.title = `${appStore.siteName} - AI API Gateway`
  }

  app.use(router)
  app.use(i18n)

  // 立即挂载应用，不等待路由就绪（提升首屏速度）
  app.mount('#app')
}

void bootstrap()
