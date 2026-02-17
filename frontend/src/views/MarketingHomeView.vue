<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950">
    <!-- Navigation -->
    <nav class="sticky top-0 z-50 border-b border-gray-200/50 bg-white/80 backdrop-blur-lg dark:border-dark-800/50 dark:bg-dark-900/80">
      <div class="mx-auto max-w-7xl px-6 py-4">
        <div class="flex items-center justify-between">
          <!-- Brand -->
          <div class="flex items-center gap-3">
            <Logo :size="40" theme="auto" />
            <div>
              <div class="text-xl font-bold text-gray-900 dark:text-white">
                Sub<span class="text-primary-600">Router</span><span class="text-primary-600">.ai</span>
              </div>
              <div class="text-xs text-gray-500 dark:text-dark-400">AI API Routing & Billing Platform</div>
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

    <!-- Hero Section -->
    <section class="px-6 py-20">
      <div class="mx-auto max-w-7xl text-center">
        <div class="mb-8">
          <h1 class="mb-4 text-6xl font-bold text-gray-900 dark:text-white md:text-7xl lg:text-8xl">
            Sub<span class="text-primary-600">Router</span><span class="text-primary-600">.ai</span>
          </h1>
          <p class="text-2xl font-semibold text-gray-700 dark:text-dark-200 mb-3">
            AI API Routing & Billing Platform
          </p>
          <p class="text-lg text-gray-600 dark:text-dark-300">
            智能路由 &middot; 灵活计费 &middot; 多模型支持
          </p>
        </div>

        <div class="mb-4 text-sm font-medium text-primary-600 dark:text-primary-400">
          Powered by Claude AI - 让 AI 成为你的编程伙伴
        </div>
        <h2 class="mb-6 text-4xl font-bold text-gray-900 dark:text-white md:text-5xl">
          Claude 先进模型使用平台
        </h2>
        <p class="mb-8 text-xl text-gray-600 dark:text-dark-300">
          支持 Opus 4.6、Opus 4.5、Sonnet 4.5、GPT-5.2-Codex 等先进 AI 模型
        </p>
        <div class="flex flex-col items-center gap-4">
          <div class="flex flex-wrap items-center justify-center gap-4">
            <a
              href="#pricing"
              class="inline-flex items-center gap-2 rounded-full bg-primary-600 px-8 py-4 text-lg font-medium text-white shadow-lg shadow-primary-500/30 transition-all hover:bg-primary-700 hover:shadow-xl"
            >
              查看套餐价格
              <Icon name="arrowRight" size="md" />
            </a>
            <router-link
              :to="isAuthenticated ? dashboardPath : '/login'"
              class="inline-flex items-center gap-2 rounded-full border-2 border-primary-600 px-8 py-4 text-lg font-medium text-primary-600 transition-all hover:bg-primary-50 dark:hover:bg-primary-900/20"
            >
              {{ isAuthenticated ? '进入控制台' : '注册 / 登录' }}
            </router-link>
          </div>
          <a href="https://fk.subrouter.ai" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-1 text-sm font-medium text-primary-600 transition-colors hover:text-primary-700 dark:text-primary-400">
            直接购买套餐 &rarr; fk.subrouter.ai
          </a>
        </div>
      </div>
    </section>

    <!-- New Model Announcement Banner -->
    <section class="px-6 py-4">
      <div class="mx-auto max-w-7xl">
        <div class="rounded-xl border border-green-300 bg-gradient-to-r from-green-50 to-emerald-50 p-4 text-center dark:border-green-700 dark:from-green-900/30 dark:to-emerald-900/30">
          <div class="flex items-center justify-center gap-2">
            <span class="inline-flex items-center rounded-full bg-green-500 px-2.5 py-0.5 text-xs font-semibold text-white">NEW</span>
            <p class="text-lg font-semibold text-green-800 dark:text-green-200">已支持 claude-opus-4-6 模型</p>
          </div>
        </div>
      </div>
    </section>


    <!-- Pricing Section -->
    <section id="pricing" class="px-6 py-20">
      <div class="mx-auto max-w-7xl">
        <h2 class="mb-12 text-center text-4xl font-bold text-gray-900 dark:text-white">套餐价格</h2>

        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <!-- Pay As You Go -->
          <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
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
            <a href="https://fk.subrouter.ai" target="_blank" rel="noopener noreferrer" class="mt-auto block rounded-xl bg-gray-100 py-3 text-center text-sm font-semibold text-gray-700 transition-colors hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-200 dark:hover:bg-dark-600">
              立即购买
            </a>
          </div>

          <!-- Monthly Plans -->
          <div class="relative flex flex-col rounded-2xl border-2 border-primary-500 bg-white p-8 shadow-xl ring-2 ring-primary-500/20 dark:bg-dark-800">
            <div class="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-primary-600 px-4 py-1 text-xs font-bold text-white shadow-md">最受欢迎</div>
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 289</h3>
            <div class="mb-1 text-4xl font-bold text-primary-600">¥289</div>
            <p class="mb-4 text-xs text-gray-500 dark:text-dark-400">约 ¥9.6/天</p>
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
            <a href="https://fk.subrouter.ai" target="_blank" rel="noopener noreferrer" class="mt-auto block rounded-xl bg-primary-600 py-3 text-center text-sm font-semibold text-white shadow-md shadow-primary-500/30 transition-all hover:bg-primary-700 hover:shadow-lg">
              立即购买
            </a>
          </div>

          <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 389</h3>
            <div class="mb-1 text-4xl font-bold text-primary-600">¥389</div>
            <p class="mb-4 text-xs text-gray-500 dark:text-dark-400">约 ¥13/天</p>
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
            <a href="https://fk.subrouter.ai" target="_blank" rel="noopener noreferrer" class="mt-auto block rounded-xl bg-gray-100 py-3 text-center text-sm font-semibold text-gray-700 transition-colors hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-200 dark:hover:bg-dark-600">
              立即购买
            </a>
          </div>

          <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 459</h3>
            <div class="mb-1 text-4xl font-bold text-primary-600">¥459</div>
            <p class="mb-4 text-xs text-gray-500 dark:text-dark-400">约 ¥15.3/天</p>
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
            <a href="https://fk.subrouter.ai" target="_blank" rel="noopener noreferrer" class="mt-auto block rounded-xl bg-gray-100 py-3 text-center text-sm font-semibold text-gray-700 transition-colors hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-200 dark:hover:bg-dark-600">
              立即购买
            </a>
          </div>

          <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 559</h3>
            <div class="mb-1 text-4xl font-bold text-primary-600">¥559</div>
            <p class="mb-4 text-xs text-gray-500 dark:text-dark-400">约 ¥18.6/天</p>
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
            <a href="https://fk.subrouter.ai" target="_blank" rel="noopener noreferrer" class="mt-auto block rounded-xl bg-gray-100 py-3 text-center text-sm font-semibold text-gray-700 transition-colors hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-200 dark:hover:bg-dark-600">
              立即购买
            </a>
          </div>

          <div class="flex flex-col rounded-2xl border border-gray-200 bg-white p-8 shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-2 inline-block w-fit rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-700 dark:bg-amber-900/30 dark:text-amber-300">大额用户</div>
            <h3 class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">月卡 1180</h3>
            <div class="mb-1 text-4xl font-bold text-primary-600">¥1180</div>
            <p class="mb-4 text-xs text-gray-500 dark:text-dark-400">约 ¥39.3/天</p>
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
            <a href="https://fk.subrouter.ai" target="_blank" rel="noopener noreferrer" class="mt-auto block rounded-xl bg-gray-100 py-3 text-center text-sm font-semibold text-gray-700 transition-colors hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-200 dark:hover:bg-dark-600">
              立即购买
            </a>
          </div>
        </div>

        <p class="mt-8 text-center text-sm text-gray-600 dark:text-dark-400">
          月卡支持叠加使用 / 支持企业组团购买
        </p>
      </div>
    </section>

    <!-- Purchase Method -->
    <section id="purchase" class="px-6 py-16">
      <div class="mx-auto max-w-4xl">
        <div class="overflow-hidden rounded-2xl border-2 border-primary-200 bg-gradient-to-br from-primary-50 via-white to-blue-50 shadow-xl dark:border-primary-800 dark:from-primary-900/30 dark:via-dark-800 dark:to-blue-900/20">
          <div class="p-8 text-center md:p-12">
            <h2 class="mb-3 text-3xl font-bold text-gray-900 dark:text-white">立即购买</h2>
            <p class="mb-8 text-gray-600 dark:text-dark-300">选好套餐？点击下方按钮前往购买页面</p>
            <a
              href="https://fk.subrouter.ai"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center gap-3 rounded-full bg-primary-600 px-10 py-4 text-lg font-bold text-white shadow-lg shadow-primary-500/30 transition-all hover:bg-primary-700 hover:shadow-xl hover:scale-105"
            >
              <Icon name="creditCard" size="md" />
              前往 fk.subrouter.ai 购买
              <Icon name="arrowRight" size="md" />
            </a>
            <div class="mt-8 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
              <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-green-500" />
                支持支付宝/微信
              </div>
              <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-green-500" />
                即买即用
              </div>
              <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-dark-300">
                <Icon name="check" size="sm" class="text-green-500" />
                支持企业团购
              </div>
            </div>
          </div>
          <div class="border-t border-primary-200 bg-yellow-50/80 px-8 py-4 dark:border-primary-800 dark:bg-yellow-900/20">
            <div class="flex items-center justify-center gap-3">
              <Icon name="exclamationCircle" size="md" class="flex-shrink-0 text-yellow-600 dark:text-yellow-400" />
              <p class="text-sm font-medium text-yellow-800 dark:text-yellow-300">购买后请务必保存下单详情里的兑换码，然后到控制台兑换区激活</p>
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
              <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">登录 subrouter.ai 注册账号</h3>
              <p class="text-sm text-gray-600 dark:text-dark-300">首次使用需要注册一个账号</p>
            </div>
          </div>

          <div class="flex gap-4 rounded-2xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <div class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white">2</div>
            <div>
              <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">在兑换区填入从 fk.subrouter.ai 购买的兑换码</h3>
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
                <div><span class="text-gray-500">URL:</span> <span class="text-primary-600">https://subrouter.ai</span></div>
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
          <a href="https://rcnwwid75eku.feishu.cn/wiki/WDWFwhUZiiNApaknf5BcqZsgndg" target="_blank" rel="noopener noreferrer" class="group rounded-2xl border border-gray-200 bg-white p-6 transition-all hover:border-primary-500 hover:shadow-lg dark:border-dark-700 dark:bg-dark-800">
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-100 text-blue-600 transition-transform group-hover:scale-110 dark:bg-blue-900/30">
              <Icon name="book" size="lg" />
            </div>
            <h3 class="mb-2 font-semibold text-gray-900 dark:text-white">新手指南</h3>
            <p class="text-sm text-gray-600 dark:text-dark-300">Claude Code 小白无门槛宝典</p>
          </a>

          <a href="https://rcnwwid75eku.feishu.cn/wiki/WDWFwhUZiiNApaknf5BcqZsgndg" target="_blank" rel="noopener noreferrer" class="group rounded-2xl border border-gray-200 bg-white p-6 transition-all hover:border-primary-500 hover:shadow-lg dark:border-dark-700 dark:bg-dark-800">
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
        <div class="flex flex-col items-center gap-6">
          <div class="flex items-center gap-3">
            <Logo :size="32" theme="auto" />
            <span class="text-xl font-bold text-gray-900 dark:text-white">Sub<span class="text-primary-600">Router</span><span class="text-primary-600">.ai</span></span>
          </div>
          <p class="text-sm text-gray-500 dark:text-dark-400">AI API Routing & Billing Platform</p>
          <div class="flex flex-wrap items-center justify-center gap-4 text-xs text-gray-400 dark:text-dark-500">
            <router-link to="/legal/tokushoho" class="transition-colors hover:text-gray-600 dark:hover:text-dark-300">特定商取引法に基づく表記</router-link>
            <span>|</span>
            <router-link to="/legal/disclosure" class="transition-colors hover:text-gray-600 dark:hover:text-dark-300">商業披露</router-link>
          </div>
          <p class="text-xs text-gray-400 dark:text-dark-500">
            &copy; {{ currentYear }} SubRouter. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import Logo from '@/components/Logo.vue'

const authStore = useAuthStore()

const isDark = ref(document.documentElement.classList.contains('dark'))
const currentYear = computed(() => new Date().getFullYear())

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  authStore.checkAuth()
})
</script>
