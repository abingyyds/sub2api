<template>
  <div class="relative min-h-screen bg-[#FAFAFA] dark:bg-dark-950 overflow-x-hidden antialiased">
    <!-- 顶部公告栏 -->
    <div class="fixed top-0 left-0 w-full bg-indigo-50 dark:bg-indigo-900/20 border-b border-indigo-100 dark:border-indigo-800 text-indigo-700 dark:text-indigo-300 text-center py-2.5 text-sm font-medium z-50 transition-all hover:bg-indigo-100 dark:hover:bg-indigo-900/30">
      <span class="mr-2">🎉</span> 新一代企业级多节点路由已上线 | 新用户注册即送 <strong>$10 无门槛纯净体验金</strong>
    </div>

    <!-- 导航栏 -->
    <nav class="max-w-7xl mx-auto px-6 py-6 flex justify-between items-center relative z-10 bg-[#FAFAFA]/80 dark:bg-dark-950/80 backdrop-blur-md sticky top-[40px]">
      <div class="text-2xl font-extrabold tracking-tighter flex items-center gap-2">
        <div class="bg-primary-500/10 text-primary-600 dark:text-primary-400 w-8 h-8 rounded-lg flex items-center justify-center">
          <Icon name="terminal" size="sm" />
        </div>
        <span class="text-gray-900 dark:text-white">C</span>
        <span class="text-primary-600 dark:text-primary-400">coder</span>
        <span class="text-gray-400 dark:text-dark-500">.me</span>
      </div>
      <div class="flex items-center gap-6">
        <a href="#pricing" class="text-sm font-medium text-gray-500 dark:text-dark-400 hover:text-gray-900 dark:hover:text-white transition hidden md:inline-block">定价套餐</a>
        <a href="#faq" class="text-sm font-medium text-gray-500 dark:text-dark-400 hover:text-gray-900 dark:hover:text-white transition hidden md:inline-block">常见问题</a>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="text-sm px-5 py-2 rounded-full bg-white dark:bg-dark-800 border border-gray-200 dark:border-dark-700 shadow-sm hover:border-primary-500 hover:text-primary-600 dark:hover:border-primary-500 dark:hover:text-primary-400 transition font-medium"
        >
          {{ isAuthenticated ? '控制台' : '登录控制台' }}
        </router-link>
      </div>
    </nav>

    <!-- Hero Section -->
    <header class="max-w-5xl mx-auto px-6 pt-16 pb-24 text-center relative">
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-gradient-to-tr from-orange-200/40 to-indigo-200/40 dark:from-orange-900/20 dark:to-indigo-900/20 rounded-full blur-[100px] -z-10 pointer-events-none"></div>
      <FadeIn :delay="100">
        <h1 class="text-5xl md:text-7xl font-extrabold tracking-tight leading-tight mb-6 text-gray-900 dark:text-white">
          突破算力封锁，<br/>你的 <span class="bg-gradient-to-r from-primary-600 to-indigo-600 bg-clip-text text-transparent">10x 生产力引擎</span>
        </h1>
      </FadeIn>
      <SlideIn direction="up" :delay="200">
        <p class="text-lg md:text-xl text-gray-500 dark:text-dark-400 mb-10 max-w-2xl mx-auto leading-relaxed">
          直连 Claude / Codex / Gemini。无损计费机制，企业级专线调度，告别一切排队与 429 报错。
        </p>
      </SlideIn>
      <SlideIn direction="up" :delay="300">
        <div class="flex flex-col items-center">
          <MagneticButton>
            <router-link
              :to="isAuthenticated ? dashboardPath : '/register'"
              class="bg-gray-900 dark:bg-white hover:bg-gray-800 dark:hover:bg-gray-100 text-white dark:text-gray-900 text-lg font-bold py-4 px-10 rounded-full shadow-xl transform hover:-translate-y-1 transition duration-300 flex items-center gap-3"
            >
              <Icon name="sparkles" size="md" /> 领取 $10 体验金，即刻接入
            </router-link>
          </MagneticButton>
        </div>
      </SlideIn>
    </header>

    <!-- 痛点对比模块 -->
    <section class="max-w-5xl mx-auto px-6 mb-20">
      <SlideIn direction="up" :delay="100">
        <GlowCard glow-color="rgb(99, 102, 241)">
          <div class="bg-white dark:bg-dark-900 border border-gray-200 dark:border-dark-800 rounded-3xl p-8 md:p-10 relative overflow-hidden shadow-soft">
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-8 text-center">算一笔账：为什么全职开发者都在用 Ccoder？</h3>

            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse">
                <thead>
                  <tr class="border-b-2 border-gray-200 dark:border-dark-700">
                    <th class="py-4 px-4 text-sm font-bold text-gray-700 dark:text-dark-300">对比维度</th>
                    <th class="py-4 px-4 text-sm font-bold text-red-600 dark:text-red-400">❌ 传统方案</th>
                    <th class="py-4 px-4 text-sm font-bold text-green-600 dark:text-green-400">✅ Ccoder.me</th>
                  </tr>
                </thead>
                <tbody class="text-sm">
                  <tr class="border-b border-gray-100 dark:border-dark-800">
                    <td class="py-4 px-4 font-medium text-gray-900 dark:text-white">月度成本</td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400">$20 订阅 + 超额按量 = <strong class="text-red-600 dark:text-red-400">$50-80</strong></td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400"><strong class="text-green-600 dark:text-green-400">¥299 = $42</strong>，纯净无损</td>
                  </tr>
                  <tr class="border-b border-gray-100 dark:border-dark-800">
                    <td class="py-4 px-4 font-medium text-gray-900 dark:text-white">速度体验</td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400">高峰期排队 <strong class="text-red-600 dark:text-red-400">3-10 分钟</strong></td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400"><strong class="text-green-600 dark:text-green-400">秒级响应</strong>，专线直连</td>
                  </tr>
                  <tr class="border-b border-gray-100 dark:border-dark-800">
                    <td class="py-4 px-4 font-medium text-gray-900 dark:text-white">稳定性</td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400">频繁 <strong class="text-red-600 dark:text-red-400">429 限流</strong></td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400"><strong class="text-green-600 dark:text-green-400">99.9% 可用</strong>，多节点容灾</td>
                  </tr>
                  <tr class="border-b border-gray-100 dark:border-dark-800">
                    <td class="py-4 px-4 font-medium text-gray-900 dark:text-white">计费透明度</td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400">黑盒计费，<strong class="text-red-600 dark:text-red-400">不知道花了多少</strong></td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400"><strong class="text-green-600 dark:text-green-400">实时账单</strong>，每笔可查</td>
                  </tr>
                  <tr>
                    <td class="py-4 px-4 font-medium text-gray-900 dark:text-white">团队协作</td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400">需额外购买 <strong class="text-red-600 dark:text-red-400">Team 版</strong></td>
                    <td class="py-4 px-4 text-gray-600 dark:text-dark-400"><strong class="text-green-600 dark:text-green-400">免费多用户</strong>，统一管理</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="mt-8 p-6 bg-gradient-to-r from-green-50 to-blue-50 dark:from-green-900/20 dark:to-blue-900/20 rounded-2xl border border-green-200 dark:border-green-800">
              <p class="text-center text-lg font-bold text-gray-900 dark:text-white mb-2">💰 每月节省 <span class="text-green-600 dark:text-green-400">$10-40</span>，一年省下 <span class="text-green-600 dark:text-green-400">$120-480</span></p>
              <p class="text-center text-sm text-gray-600 dark:text-dark-400">更快、更稳、更省心 —— 这就是 Ccoder 的价值主张</p>
            </div>
          </div>
        </GlowCard>
      </SlideIn>
    </section>

    <!-- 核心特性 -->
    <section class="max-w-6xl mx-auto px-6 mb-20">
      <SlideIn direction="up" :delay="100">
        <h2 class="text-3xl md:text-4xl font-bold text-center text-gray-900 dark:text-white mb-4">为什么选择 Ccoder？</h2>
        <p class="text-center text-gray-500 dark:text-dark-400 mb-12 max-w-2xl mx-auto">企业级基础设施 + 开发者友好设计，让 AI 编程回归本质</p>
      </SlideIn>

      <StaggerContainer :stagger-delay="100" :delay="200">
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <GlowCard
            v-for="(feature, index) in features"
            :key="index"
            :glow-color="feature.color"
          >
            <div class="bg-white dark:bg-dark-900 border border-gray-200 dark:border-dark-800 rounded-2xl p-6 hover:shadow-lg transition-all duration-300">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center mb-4" :class="feature.bgClass">
                <Icon :name="feature.icon as any" size="lg" :class="feature.iconClass" />
              </div>
              <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-2">{{ feature.title }}</h3>
              <p class="text-sm text-gray-600 dark:text-dark-400 leading-relaxed">{{ feature.description }}</p>
            </div>
          </GlowCard>
        </div>
      </StaggerContainer>
    </section>

    <!-- 定价套餐 -->
    <section id="pricing" class="max-w-6xl mx-auto px-6 mb-20">
      <SlideIn direction="up" :delay="100">
        <h2 class="text-3xl md:text-4xl font-bold text-center text-gray-900 dark:text-white mb-4">简单透明的定价</h2>
        <p class="text-center text-gray-500 dark:text-dark-400 mb-12">无隐藏费用，按需选择，随时升级</p>
      </SlideIn>

      <StaggerContainer :stagger-delay="100" :delay="200">
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          <GlowCard
            v-for="(plan, index) in pricingPlans"
            :key="index"
            :glow-color="plan.popular ? 'rgb(217, 119, 87)' : 'rgb(99, 102, 241)'"
          >
            <div class="bg-white dark:bg-dark-900 border-2 rounded-2xl p-6 relative hover:shadow-xl transition-all duration-300" :class="plan.popular ? 'border-primary-500' : 'border-gray-200 dark:border-dark-800'">
              <div v-if="plan.popular" class="absolute -top-3 left-1/2 -translate-x-1/2 bg-primary-600 text-white text-xs font-bold px-3 py-1 rounded-full">
                🔥 最受欢迎
              </div>
              <div class="text-center mb-6">
                <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">{{ plan.name }}</h3>
                <div class="flex items-baseline justify-center gap-1">
                  <span class="text-4xl font-extrabold text-gray-900 dark:text-white">{{ plan.price }}</span>
                  <span class="text-gray-500 dark:text-dark-400">/月</span>
                </div>
              </div>
              <ul class="space-y-3 mb-6">
                <li v-for="(feature, idx) in plan.features" :key="idx" class="flex items-start gap-2 text-sm text-gray-600 dark:text-dark-400">
                  <Icon name="check" size="sm" class="text-green-500 flex-shrink-0 mt-0.5" />
                  <span>{{ feature }}</span>
                </li>
              </ul>
              <MagneticButton>
                <button
                  @click="selectPlan(plan)"
                  class="w-full py-3 rounded-xl font-bold transition-all duration-300"
                  :class="plan.popular
                    ? 'bg-primary-600 hover:bg-primary-700 text-white shadow-lg'
                    : 'bg-gray-100 dark:bg-dark-800 hover:bg-gray-200 dark:hover:bg-dark-700 text-gray-900 dark:text-white'"
                >
                  {{ plan.cta }}
                </button>
              </MagneticButton>
            </div>
          </GlowCard>
        </div>
      </StaggerContainer>
    </section>

    <!-- FAQ -->
    <section id="faq" class="max-w-3xl mx-auto px-6 mb-20">
      <SlideIn direction="up" :delay="100">
        <h2 class="text-3xl md:text-4xl font-bold text-center text-gray-900 dark:text-white mb-4">常见问题</h2>
        <p class="text-center text-gray-500 dark:text-dark-400 mb-12">快速了解 Ccoder 的核心功能</p>
      </SlideIn>

      <StaggerContainer :stagger-delay="80" :delay="200">
        <div class="space-y-4">
          <div
            v-for="(faq, index) in faqs"
            :key="index"
            class="bg-white dark:bg-dark-900 border border-gray-200 dark:border-dark-800 rounded-xl overflow-hidden"
          >
            <button
              @click="toggleFaq(index)"
              class="w-full px-6 py-4 flex items-center justify-between text-left hover:bg-gray-50 dark:hover:bg-dark-800 transition-colors"
            >
              <span class="font-semibold text-gray-900 dark:text-white">{{ faq.question }}</span>
              <Icon
                :name="activeFaq === index ? 'chevronUp' : 'chevronDown'"
                size="sm"
                class="text-gray-400 transition-transform"
              />
            </button>
            <div
              v-show="activeFaq === index"
              class="px-6 pb-4 text-sm text-gray-600 dark:text-dark-400 leading-relaxed"
            >
              {{ faq.answer }}
            </div>
          </div>
        </div>
      </StaggerContainer>
    </section>

    <!-- Footer -->
    <footer class="border-t border-gray-200 dark:border-dark-800 py-8">
      <div class="max-w-7xl mx-auto px-6 text-center text-sm text-gray-500 dark:text-dark-400">
        <p>&copy; 2024 Ccoder.me. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Icon from '@/components/icons/Icon.vue'
import { FadeIn, SlideIn, StaggerContainer, GlowCard, MagneticButton } from '@/components/animations'

const router = useRouter()
const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => {
  if (authStore.user?.role === 'admin') return '/admin/dashboard'
  if (authStore.user?.role === 'org_admin') return '/org/dashboard'
  return '/dashboard'
})

const activeFaq = ref<number | null>(null)

const features = [
  {
    icon: 'bolt',
    title: '极速响应',
    description: '企业级专线直连，告别排队等待，秒级响应你的每一次请求',
    color: 'rgb(234, 179, 8)',
    bgClass: 'bg-yellow-100 dark:bg-yellow-900/20',
    iconClass: 'text-yellow-600 dark:text-yellow-400'
  },
  {
    icon: 'shield',
    title: '稳定可靠',
    description: '99.9% 可用性保证，多节点容灾，让你的开发流程永不中断',
    color: 'rgb(34, 197, 94)',
    bgClass: 'bg-green-100 dark:bg-green-900/20',
    iconClass: 'text-green-600 dark:text-green-400'
  },
  {
    icon: 'dollar',
    title: '透明计费',
    description: '实时账单，每笔消费清晰可查，无损计费机制，用多少付多少',
    color: 'rgb(59, 130, 246)',
    bgClass: 'bg-blue-100 dark:bg-blue-900/20',
    iconClass: 'text-blue-600 dark:text-blue-400'
  },
  {
    icon: 'users',
    title: '团队协作',
    description: '免费多用户管理，统一账单，让团队协作更高效',
    color: 'rgb(168, 85, 247)',
    bgClass: 'bg-purple-100 dark:bg-purple-900/20',
    iconClass: 'text-purple-600 dark:text-purple-400'
  },
  {
    icon: 'cloud',
    title: '全球加速',
    description: '智能路由选择，自动选择最优节点，全球范围内都能享受极速体验',
    color: 'rgb(236, 72, 153)',
    bgClass: 'bg-pink-100 dark:bg-pink-900/20',
    iconClass: 'text-pink-600 dark:text-pink-400'
  },
  {
    icon: 'lock',
    title: '安全保障',
    description: '企业级加密传输，数据隐私保护，让你的代码和数据更安全',
    color: 'rgb(239, 68, 68)',
    bgClass: 'bg-red-100 dark:bg-red-900/20',
    iconClass: 'text-red-600 dark:text-red-400'
  }
]

const pricingPlans = [
  {
    name: '体验版',
    price: '¥0',
    features: [
      '$10 纯净体验金',
      '立刻到账，不限速',
      '体验所有核心功能',
      '适合新用户试用'
    ],
    cta: '立即领取',
    popular: false
  },
  {
    name: '个人版',
    price: '¥299',
    features: [
      '获得 $600 调用额度',
      '每日 $30 高速通道',
      '适合个人开发者',
      '无损计费机制'
    ],
    cta: '立即购买',
    popular: false
  },
  {
    name: '专业版',
    price: '¥589',
    features: [
      '获得 $1300 调用额度',
      '每日 $60 高速通道',
      '适合重度用户',
      '优先技术支持'
    ],
    cta: '立即购买',
    popular: true
  },
  {
    name: '团队版',
    price: '¥1189',
    features: [
      '获得 $2500 调用额度',
      '每日 $120 高速通道',
      '适合团队协作',
      '专属客户经理'
    ],
    cta: '立即购买',
    popular: false
  }
]

const faqs = [
  {
    question: '什么是无损计费机制？',
    answer: '我们采用 1:1 透明计费，你在 Ccoder 上消费的每一分钱都会完整转化为 API 调用额度，没有任何中间损耗或隐藏费用。'
  },
  {
    question: '如何保证 99.9% 的可用性？',
    answer: '我们部署了多个全球节点，采用智能路由和自动容灾机制。当某个节点出现问题时，系统会自动切换到备用节点，确保服务不中断。'
  },
  {
    question: '支持哪些 AI 模型？',
    answer: '目前支持 Claude 全系列（Opus、Sonnet、Haiku）、OpenAI Codex、Google Gemini 等主流 AI 模型，并持续增加新模型支持。'
  },
  {
    question: '可以开发票吗？',
    answer: '可以。我们支持开具增值税普通发票和专用发票，购买后在控制台申请即可。'
  },
  {
    question: '如何开始使用？',
    answer: '注册账号后即可获得 $10 体验金，无需绑卡。在控制台创建 API Key，然后在你的开发工具中配置即可开始使用。'
  },
  {
    question: '有技术支持吗？',
    answer: '有。我们提供 7x24 小时技术支持，专业版及以上用户享有优先响应权。同时我们也有详细的文档和社区支持。'
  }
]

function toggleFaq(index: number) {
  activeFaq.value = activeFaq.value === index ? null : index
}

function selectPlan(_plan: any) {
  if (isAuthenticated.value) {
    router.push('/pricing')
  } else {
    router.push('/register')
  }
}

onMounted(() => {
  authStore.checkAuth()
})
</script>

<style scoped>
.text-gradient {
  background: linear-gradient(135deg, #D97757 0%, #6366F1 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.shadow-soft {
  box-shadow: 0 10px 40px -10px rgba(0,0,0,0.08);
}

.shadow-glow {
  box-shadow: 0 0 20px rgba(217, 119, 87, 0.3);
}
</style>

