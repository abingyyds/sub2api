<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950">
    <!-- Navigation -->
    <nav class="sticky top-0 z-50 border-b border-gray-200/50 bg-white/80 backdrop-blur-lg dark:border-dark-800/50 dark:bg-dark-900/80">
      <div class="mx-auto max-w-7xl px-6 py-4">
        <div class="flex items-center justify-between">
          <!-- Brand with Animation -->
          <div class="flex items-center gap-3">
            <Logo :size="40" theme="auto" />
            <div>
              <div class="text-xl font-bold text-gray-900 dark:text-white">
                <span class="inline-block">c</span><span class="inline-block">{{ animatedText }}</span><span class="inline-block">Coder.me</span>
              </div>
              <div class="text-xs text-gray-500 dark:text-dark-400">新一代 AI 代码工作平台</div>
            </div>
          </div>

          <!-- Nav Links -->
          <div class="hidden items-center gap-8 md:flex">
            <a href="#pricing" class="text-sm font-medium text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-300 dark:hover:text-white">套餐价格</a>
            <a href="#purchase" class="text-sm font-medium text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-300 dark:hover:text-white">购买方式</a>
            <a href="#guide" class="text-sm font-medium text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-300 dark:hover:text-white">使用指南</a>
            <a href="#help" class="text-sm font-medium text-gray-600 transition-colors hover:text-gray-900 dark:text-dark-300 dark:hover:text-white">帮助文档</a>
            <router-link
              :to="isAuthenticated ? dashboardPath : '/login'"
              class="rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700"
            >
              {{ isAuthenticated ? '控制台' : '控制台 / 注册' }}
            </router-link>
          </div>

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
        </div>
      </div>
    </nav>

    <!-- Hero Section with Brand Animation -->
    <section class="px-6 py-20">
      <div class="mx-auto max-w-7xl text-center">
        <!-- Animated Brand Title -->
        <div class="mb-8">
          <h1 class="mb-4 text-6xl font-bold text-gray-900 dark:text-white md:text-7xl lg:text-8xl">
            <span class="inline-block">c</span><span class="inline-block">{{ animatedText }}</span><span class="inline-block">Coder</span><span class="text-primary-600">.me</span>
          </h1>
          <p class="text-2xl font-semibold text-gray-700 dark:text-dark-200 mb-3">
            新一代代码大师就是我
          </p>
          <p class="text-lg text-gray-600 dark:text-dark-300">
            我就是 <span class="font-bold text-primary-600">claude</span> <span class="font-bold text-primary-600">coder</span> <span class="text-primary-600">=</span> <span class="font-bold text-primary-600">me</span>
          </p>
        </div>

        <div class="mb-4 text-sm font-medium text-primary-600 dark:text-primary-400">
          Powered by Claude AI - 让 AI 成为你的编程伙伴
        </div>
        <h2 class="mb-6 text-4xl font-bold text-gray-900 dark:text-white md:text-5xl">
          Claude 先进模型使用平台
        </h2>
        <p class="mb-8 text-xl text-gray-600 dark:text-dark-300">
          支持 Opus 4.5、Sonnet 4.5、GPT-5.2-Codex 等先进 AI 模型
        </p>
        <div class="flex flex-col items-center gap-4">
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex items-center gap-2 rounded-full bg-primary-600 px-8 py-4 text-lg font-medium text-white shadow-lg shadow-primary-500/30 transition-all hover:bg-primary-700 hover:shadow-xl"
          >
            进入控制台
            <Icon name="arrowRight" size="md" />
          </router-link>
          <p class="text-sm text-gray-500 dark:text-dark-400">^ 新用户注册也点这里 ^</p>
        </div>
      </div>
    </section>

    <!-- Claude API Service Banner -->
    <section class="px-6 py-12">
      <div class="mx-auto max-w-7xl">
        <div class="rounded-2xl border border-primary-200 bg-gradient-to-r from-primary-50 to-blue-50 p-8 text-center dark:border-primary-800 dark:from-primary-900/20 dark:to-blue-900/20">
          <h2 class="mb-4 text-3xl font-bold text-gray-900 dark:text-white">Claude API 服务</h2>
          <p class="text-lg text-gray-600 dark:text-dark-300">支持 Opus 4.5、Sonnet 4.5、GPT-5.2-Codex 等模型</p>
        </div>
      </div>
    </section>

    <!-- Pricing Section -->
    <section id="pricing" class="px-6 py-20">
      <div class="mx-auto max-w-7xl">
        <h2 class="mb-12 text-center text-4xl font-bold text-gray-900 dark:text-white">套餐价格</h2>

        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <!-- Pay As You Go -->
          <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">按量付费</h3>
            <div class="mb-4 text-4xl font-bold text-primary-600">1:1</div>
            <p class="mb-6 text-gray-600 dark:text-dark-300">1$ = 1RMB</p>
            <ul class="mb-8 space-y-3">
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                Opus 4.5
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                Sonnet 4.5
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                GPT-5.2-Codex
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                用多少付多少
              </li>
            </ul>
          </div>

          <!-- Monthly Plans -->
          <div class="rounded-2xl border-2 border-primary-500 bg-white p-8 shadow-xl dark:bg-dark-800">
            <div class="mb-2 inline-block rounded-full bg-primary-100 px-3 py-1 text-xs font-medium text-primary-600 dark:bg-primary-900/30">推荐</div>
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 289</h3>
            <div class="mb-4 text-4xl font-bold text-primary-600">¥289</div>
            <ul class="mb-8 space-y-3">
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                每日 30 美金
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                无周限额
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                30 天有效期
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                支持叠加
              </li>
            </ul>
          </div>

          <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 389</h3>
            <div class="mb-4 text-4xl font-bold text-primary-600">¥389</div>
            <ul class="mb-8 space-y-3">
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                每日 40 美金
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                无周限额
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                30 天有效期
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                支持叠加
              </li>
            </ul>
          </div>

          <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 459</h3>
            <div class="mb-4 text-4xl font-bold text-primary-600">¥459</div>
            <ul class="mb-8 space-y-3">
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                每日 50 美金
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                无周限额
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                30 天有效期
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                支持叠加
              </li>
            </ul>
          </div>

          <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 559</h3>
            <div class="mb-4 text-4xl font-bold text-primary-600">¥559</div>
            <ul class="mb-8 space-y-3">
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                每日 60 美金
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                无周限额
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                30 天有效期
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                支持叠加
              </li>
            </ul>
          </div>

          <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 1180</h3>
            <div class="mb-4 text-4xl font-bold text-primary-600">¥1180</div>
            <ul class="mb-8 space-y-3">
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                每日 120 美金
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                无周限额
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                30 天有效期
              </li>
              <li class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-primary-600" />
                支持叠加
              </li>
            </ul>
          </div>
        </div>

        <p class="mt-8 text-center text-sm text-gray-600 dark:text-dark-400">
          月卡支持叠加使用 / 支持企业组团购买
        </p>
      </div>
    </section>

    <!-- Purchase Method -->
    <section id="purchase" class="bg-gray-50 px-6 py-20 dark:bg-dark-900/50">
      <div class="mx-auto max-w-4xl">
        <h2 class="mb-12 text-center text-4xl font-bold text-gray-900 dark:text-white">购买方式</h2>
        <div class="rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
          <div class="mb-6 flex items-start gap-4">
            <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-primary-100 text-primary-600 dark:bg-primary-900/30">
              <Icon name="creditCard" size="md" />
            </div>
            <div>
              <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">浏览器打开下单链接，支付后保存兑换码</h3>
              <a href="#" class="text-primary-600 hover:text-primary-700 dark:text-primary-400">fk.ccoder.me</a>
            </div>
          </div>
          <div class="rounded-lg bg-yellow-50 p-4 dark:bg-yellow-900/20">
            <div class="flex items-start gap-3">
              <Icon name="exclamationCircle" size="md" class="mt-0.5 flex-shrink-0 text-yellow-600 dark:text-yellow-400" />
              <div>
                <h4 class="mb-1 font-semibold text-yellow-900 dark:text-yellow-200">重要提示</h4>
                <p class="text-sm text-yellow-800 dark:text-yellow-300">购买后请保存下单详情里的兑换码</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Usage Guide -->
    <section id="guide" class="px-6 py-20">
      <div class="mx-auto max-w-4xl">
        <h2 class="mb-12 text-center text-4xl font-bold text-gray-900 dark:text-white">使用指南</h2>
        <div class="space-y-6">
          <div class="flex gap-4 rounded-2xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">1</div>
            <div>
              <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">登录 ccoder.me 注册账号</h3>
              <p class="text-sm text-gray-600 dark:text-dark-300">首次使用需要注册一个账号</p>
            </div>
          </div>

          <div class="flex gap-4 rounded-2xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">2</div>
            <div>
              <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">在兑换区填入从 fk.ccoder.me 购买的兑换码</h3>
              <p class="text-sm text-gray-600 dark:text-dark-300">输入购买时获得的兑换码激活服务</p>
            </div>
          </div>

          <div class="flex gap-4 rounded-2xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">3</div>
            <div>
              <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">选择对应的分组（月卡或按量）</h3>
              <p class="text-sm text-gray-600 dark:text-dark-300">根据购买的套餐类型选择分组</p>
            </div>
          </div>

          <div class="flex gap-4 rounded-2xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">4</div>
            <div>
              <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">配置 API</h3>
              <div class="mt-3 space-y-2 rounded-lg bg-gray-50 p-4 font-mono text-sm dark:bg-dark-900">
                <div><span class="text-gray-500">URL:</span> <span class="text-primary-600">https://ccoder.me</span></div>
                <div><span class="text-gray-500">Key:</span> <span class="text-gray-600 dark:text-dark-300">您的 API 密钥</span></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Help Documentation -->
    <section id="help" class="bg-gray-50 px-6 py-20 dark:bg-dark-900/50">
      <div class="mx-auto max-w-4xl">
        <h2 class="mb-12 text-center text-4xl font-bold text-gray-900 dark:text-white">帮助文档</h2>
        <div class="grid gap-6 md:grid-cols-3">
          <a href="#" class="group rounded-2xl border border-gray-200 bg-white p-6 transition-all hover:border-primary-500 hover:shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-100 text-blue-600 transition-transform group-hover:scale-110 dark:bg-blue-900/30">
              <Icon name="book" size="lg" />
            </div>
            <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">新手指南</h3>
            <p class="text-sm text-gray-600 dark:text-dark-300">Claude Code 小白无门槛宝典</p>
          </a>

          <a href="#" class="group rounded-2xl border border-gray-200 bg-white p-6 transition-all hover:border-primary-500 hover:shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-green-100 text-green-600 transition-transform group-hover:scale-110 dark:bg-green-900/30">
              <Icon name="search" size="lg" />
            </div>
            <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">问题排查</h3>
            <p class="text-sm text-gray-600 dark:text-dark-300">使用问题速查</p>
          </a>

          <a href="#" class="group rounded-2xl border border-gray-200 bg-white p-6 transition-all hover:border-primary-500 hover:shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-purple-100 text-purple-600 transition-transform group-hover:scale-110 dark:bg-purple-900/30">
              <Icon name="document" size="lg" />
            </div>
            <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">开具发票</h3>
            <p class="text-sm text-gray-600 dark:text-dark-300">填写发票申请表单</p>
          </a>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-gray-200 px-6 py-12 dark:border-dark-800">
      <div class="mx-auto max-w-7xl">
        <div class="text-center">
          <div class="mb-4">
            <div class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">
              <span class="inline-block">c</span><span class="inline-block text-primary-600">Coder</span><span class="text-primary-600">.me</span>
            </div>
            <div class="text-sm text-gray-600 dark:text-dark-400">新一代 AI 代码工作平台</div>
          </div>
          <p class="text-sm text-gray-500 dark:text-dark-400">
            &copy; {{ currentYear }} cCoder.me. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import Logo from '@/components/Logo.vue'

const authStore = useAuthStore()

const isDark = ref(document.documentElement.classList.contains('dark'))
const currentYear = computed(() => new Date().getFullYear())

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')

// Brand animation state - typewriter effect
const animatedText = ref('laude')
const fullText = 'laude'
let animationTimer: ReturnType<typeof setTimeout> | null = null

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Typewriter animation: delete and type "laude" repeatedly
function startTypewriterAnimation() {
  let currentIndex = fullText.length
  let isDeleting = false
  let isPaused = false

  function animate() {
    if (isPaused) {
      // Pause phase
      isPaused = false
      animationTimer = setTimeout(animate, isDeleting ? 1000 : 3000)
      return
    }

    if (isDeleting) {
      // Deleting phase
      if (currentIndex > 0) {
        currentIndex--
        animatedText.value = fullText.substring(0, currentIndex)
        animationTimer = setTimeout(animate, 150)
      } else {
        // Finished deleting, pause then start typing
        isDeleting = false
        isPaused = true
        animationTimer = setTimeout(animate, 1000)
      }
    } else {
      // Typing phase
      if (currentIndex < fullText.length) {
        currentIndex++
        animatedText.value = fullText.substring(0, currentIndex)
        animationTimer = setTimeout(animate, 150)
      } else {
        // Finished typing, pause then start deleting
        isDeleting = true
        isPaused = true
        animationTimer = setTimeout(animate, 3000)
      }
    }
  }

  // Start animation after initial display
  animationTimer = setTimeout(() => {
    isDeleting = true
    isPaused = true
    animate()
  }, 3000)
}

onMounted(() => {
  authStore.checkAuth()
  startTypewriterAnimation()
})

onUnmounted(() => {
  if (animationTimer) {
    clearTimeout(animationTimer)
    animationTimer = null
  }
})
</script>
