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
          突破算力封锁，<br/>你的 <span class="text-gradient">10x 生产力引擎</span>
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

    <!-- 痛点对比模块 - 两栏布局 -->
    <section class="max-w-5xl mx-auto px-6 mb-20">
      <SlideIn direction="up" :delay="100">
        <GlowCard glow-color="rgb(99, 102, 241)">
          <div class="bg-white dark:bg-dark-900 border border-gray-200 dark:border-dark-800 rounded-3xl p-8 md:p-10 relative overflow-hidden shadow-soft">
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-8 text-center">算一笔账：为什么全职开发者都在用 Ccoder？</h3>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- 左侧：官方渠道痛点 -->
              <div class="bg-red-50/50 dark:bg-red-900/10 border border-red-100 dark:border-red-800/30 rounded-2xl p-6 md:p-8">
                <div class="text-red-500 dark:text-red-400 font-bold mb-4 flex items-center gap-2 text-lg">
                  <Icon name="exclamationTriangle" size="md" /> 使用官方渠道，您需要承担：
                </div>
                <ul class="space-y-4 text-sm text-gray-600 dark:text-dark-400">
                  <li class="flex items-start gap-3">
                    <Icon name="x" size="sm" class="text-red-400 mt-0.5 flex-shrink-0" />
                    <span class="leading-relaxed">高昂的美元汇率 (7.2+)</span>
                  </li>
                  <li class="flex items-start gap-3">
                    <Icon name="x" size="sm" class="text-red-400 mt-0.5 flex-shrink-0" />
                    <span class="leading-relaxed">信用卡跨国手续费及拒付风险</span>
                  </li>
                  <li class="flex items-start gap-3">
                    <Icon name="x" size="sm" class="text-red-400 mt-0.5 flex-shrink-0" />
                    <span class="leading-relaxed">复杂的网络代理与环境配置</span>
                  </li>
                  <li class="flex items-start gap-3">
                    <Icon name="x" size="sm" class="text-red-400 mt-0.5 flex-shrink-0" />
                    <span class="leading-relaxed">官方 API 频繁超限的封号风险</span>
                  </li>
                </ul>
              </div>

              <!-- 右侧：Ccoder 优势 -->
              <div class="bg-orange-50/80 dark:bg-orange-900/10 border border-orange-200 dark:border-orange-800/30 rounded-2xl p-6 md:p-8 relative">
                <div class="absolute top-4 right-4 text-4xl opacity-10">🚀</div>
                <div class="text-orange-700 dark:text-orange-400 font-bold mb-4 flex items-center gap-2 text-lg">
                  <Icon name="checkCircle" size="md" /> 在本站 (以 ¥589 专业版为例)：
                </div>
                <ul class="space-y-4 text-sm text-gray-700 dark:text-dark-300">
                  <li class="flex items-start gap-3">
                    <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                    <span class="leading-relaxed">获得 <strong>$1800</strong> 的官方算力额度</span>
                  </li>
                  <li class="flex items-start gap-3">
                    <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                    <span class="leading-relaxed">享有 <strong>1:1</strong> 的超低汇率标准</span>
                  </li>
                  <li class="flex items-start gap-3">
                    <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                    <span class="leading-relaxed">专属不拥堵的高速调用通道</span>
                  </li>
                  <li class="flex items-start gap-3">
                    <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                    <span class="leading-relaxed">永不断联的自动溢出续航保障</span>
                  </li>
                </ul>
              </div>
            </div>

            <!-- 底部总结 -->
            <div class="mt-8 bg-gray-50 dark:bg-dark-800 rounded-xl p-5 border border-gray-100 dark:border-dark-700 text-center">
              <p class="text-gray-600 dark:text-dark-400 text-[15px] md:text-base leading-relaxed">
                对于任何需要生产力工具的开发者来说，这不仅仅是帮您省钱，<br class="hidden md:block"/>
                <strong class="text-gray-900 dark:text-white font-bold"><Icon name="clock" size="sm" class="text-primary-600 dark:text-primary-400 inline mr-1" /> 更是为您节省了至少 3 小时/天的配置与排错时间。</strong>
              </p>
            </div>
          </div>
        </GlowCard>
      </SlideIn>
    </section>

    <!-- 核心定价卡片 -->
    <section id="pricing" class="pb-24 max-w-7xl mx-auto px-6 relative">
      <SlideIn direction="up" :delay="100">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-extrabold mb-4 text-gray-900 dark:text-white">选择您的生产力阶梯</h2>
          <p class="text-gray-500 dark:text-dark-400">根据您的调用频率与并发需求，选择最匹配的性能方案</p>
        </div>
      </SlideIn>

      <StaggerContainer :stagger-delay="120" :delay="200">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-stretch">

          <!-- 开发者版 ¥299 -->
          <GlowCard glow-color="rgb(234, 179, 8)">
            <div class="bg-white dark:bg-dark-900 p-8 rounded-3xl border-2 border-orange-100 dark:border-orange-800/30 shadow-soft relative flex flex-col hover:shadow-lg transition h-full">
              <div class="absolute -top-3 left-8 bg-orange-50 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400 text-xs font-bold px-3 py-1 rounded-full border border-orange-200 dark:border-orange-700 shadow-sm">
                🔥 个人首选 / 火爆热销
              </div>
              <div class="text-gray-500 dark:text-dark-400 text-sm font-semibold mb-2 mt-2">[适合：独立开发者 / 个人项目]</div>
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">开发者版</h3>
              <div class="text-4xl font-extrabold text-gray-900 dark:text-white mb-8">¥299</div>

              <ul class="space-y-5 mb-8 text-sm text-gray-600 dark:text-dark-400 flex-grow">
                <li class="flex items-start gap-3">
                  <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white">总额 $900：</strong>官方价格体系，无需复杂换算</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white">每日 $30 高速配额：</strong>专属通道，响应丝滑无卡顿</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white">永不掉线保障：</strong>额度用尽自动转按量，业务 7x24 在线</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="check" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white">超高性价比：</strong>个人开发首选，一次付费全月无忧</div>
                </li>
              </ul>
              <MagneticButton>
                <button
                  @click="selectPlan('dev')"
                  class="w-full py-3.5 rounded-xl bg-orange-50 dark:bg-orange-900/20 hover:bg-orange-100 dark:hover:bg-orange-900/30 text-orange-700 dark:text-orange-400 border border-orange-200 dark:border-orange-700 transition font-bold mt-auto"
                >
                  购买开发者版
                </button>
              </MagneticButton>
            </div>
          </GlowCard>

          <!-- 专业生产版 ¥589 (主推) -->
          <GlowCard glow-color="rgb(217, 119, 87)">
            <div class="bg-white dark:bg-dark-900 p-8 rounded-3xl border-[3px] border-primary-500/80 dark:border-primary-400/60 transform lg:scale-105 shadow-2xl shadow-primary-500/20 relative z-10 flex flex-col h-full">
              <div class="absolute -top-4 left-1/2 -translate-x-1/2 bg-gradient-to-r from-orange-500 to-red-500 text-white text-xs font-bold px-5 py-1.5 rounded-full shadow-lg whitespace-nowrap tracking-wide">
                👑 生产环境推荐 / 主推套餐
              </div>

              <div class="text-primary-600 dark:text-primary-400 text-sm font-semibold mb-2 mt-2">[适合：全职开发 / 生产级应用]</div>
              <h3 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">专业生产版</h3>
              <div class="text-5xl font-extrabold text-gray-900 dark:text-white mb-4">¥589</div>

              <div class="bg-primary-500/10 dark:bg-primary-400/10 border border-primary-500/20 dark:border-primary-400/20 rounded-xl p-2.5 mb-6 text-sm text-primary-600 dark:text-primary-400 text-center font-bold">
                折合 $1 = ¥0.33，立省 67% 成本！
              </div>

              <ul class="space-y-5 mb-8 text-sm text-gray-600 dark:text-dark-400 flex-grow">
                <li class="flex items-start gap-3">
                  <Icon name="database" size="sm" class="text-primary-600 dark:text-primary-400 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white text-base">总额 $1800：</strong>充足算力储备，应对高频调用</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="bolt" size="sm" class="text-yellow-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white text-base">每日 $60 高速配额：</strong>生产级优先级，拒绝排队等待</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="dollar" size="sm" class="text-green-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white text-base">无损计费机制：</strong>拒绝倍率陷阱，每一分实打实可用</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="shield" size="sm" class="text-blue-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white text-base">熔断智能防护：</strong>异常流量自动隔离，生产环境稳定如磐石</div>
                </li>
              </ul>
              <MagneticButton>
                <button
                  @click="selectPlan('pro')"
                  class="w-full py-4 rounded-xl bg-primary-600 hover:bg-primary-700 text-white font-bold text-lg shadow-[0_4px_14px_0_rgba(217,119,87,0.39)] transform hover:-translate-y-0.5 transition duration-200 mt-auto"
                >
                  立即购买，释放生产力
                </button>
              </MagneticButton>
            </div>
          </GlowCard>

          <!-- 旗舰团队版 ¥1189 -->
          <GlowCard glow-color="rgb(168, 85, 247)">
            <div class="bg-white dark:bg-dark-900 p-8 rounded-3xl border-2 border-purple-100 dark:border-purple-800/30 shadow-soft relative flex flex-col hover:shadow-lg transition h-full">
              <div class="absolute -top-3 right-8 bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 text-xs font-bold px-3 py-1 rounded-full border border-purple-200 dark:border-purple-700 shadow-sm">
                💎 团队专属 / SLA 保障
              </div>
              <div class="text-gray-500 dark:text-dark-400 text-sm font-semibold mb-2 mt-2">[适合：开发团队 / AI 产品项目]</div>
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">旗舰团队版</h3>
              <div class="text-4xl font-extrabold text-gray-900 dark:text-white mb-8">¥1189</div>

              <ul class="space-y-5 mb-8 text-sm text-gray-600 dark:text-dark-400 flex-grow">
                <li class="flex items-start gap-3">
                  <Icon name="server" size="sm" class="text-gray-700 dark:text-dark-300 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white">总额 $3600：</strong>顶配储备，支撑海量高并发接入</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="badge" size="sm" class="text-gray-700 dark:text-dark-300 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white">每日 $120 高速配额：</strong>团队级 VIP 优先级调度，秒级响应</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="document" size="sm" class="text-gray-400 dark:text-dark-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white">高可用服务保障：</strong>极致架构，满足企业级稳定运营标准</div>
                </li>
                <li class="flex items-start gap-3">
                  <Icon name="refresh" size="sm" class="text-gray-400 dark:text-dark-500 mt-0.5 flex-shrink-0" />
                  <div class="leading-relaxed"><strong class="text-gray-900 dark:text-white">无感溢出续航：</strong>额度海量储备，确保 AI 工作流永不断供</div>
                </li>
              </ul>
              <MagneticButton>
                <button
                  @click="selectPlan('team')"
                  class="w-full py-3.5 rounded-xl bg-purple-50 dark:bg-purple-900/20 hover:bg-purple-100 dark:hover:bg-purple-900/30 text-purple-700 dark:text-purple-400 border border-purple-200 dark:border-purple-700 transition font-bold mt-auto"
                >
                  购买旗舰版
                </button>
              </MagneticButton>
            </div>
          </GlowCard>

        </div>
      </StaggerContainer>
    </section>

    <!-- FAQ -->
    <section id="faq" class="py-24 bg-gray-50 dark:bg-dark-900 border-t border-gray-200 dark:border-dark-800">
      <div class="max-w-4xl mx-auto px-6">
        <SlideIn direction="up" :delay="100">
          <div class="text-center mb-12">
            <h2 class="text-3xl font-extrabold text-gray-900 dark:text-white tracking-tight">常见问题</h2>
            <p class="text-gray-500 dark:text-dark-400 mt-2">关于服务的一些常见疑问</p>
          </div>
        </SlideIn>

        <StaggerContainer :stagger-delay="80" :delay="200">
          <div class="space-y-4">
            <div
              v-for="(faq, index) in faqs"
              :key="index"
              class="bg-white dark:bg-dark-800 border border-gray-200 dark:border-dark-700 rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-shadow"
            >
              <button
                @click="toggleFaq(index)"
                class="w-full flex justify-between items-center p-6 text-left focus:outline-none bg-transparent hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors"
              >
                <span class="text-gray-900 dark:text-white font-bold text-[15px] tracking-wide">{{ faq.question }}</span>
                <Icon
                  :name="activeFaq === index ? 'chevronUp' : 'chevronDown'"
                  size="sm"
                  class="text-blue-500 dark:text-blue-400 transition-transform flex-shrink-0 ml-4"
                />
              </button>
              <div
                v-show="activeFaq === index"
                class="px-6 pb-6 text-sm text-gray-600 dark:text-dark-400 leading-relaxed border-t border-gray-50 dark:border-dark-700"
              >
                {{ faq.answer }}
              </div>
            </div>
          </div>
        </StaggerContainer>
      </div>
    </section>

    <!-- 企业定制 -->
    <section class="max-w-5xl mx-auto px-6 py-20">
      <SlideIn direction="up" :delay="100">
        <GlowCard glow-color="rgb(99, 102, 241)">
          <div class="bg-white dark:bg-dark-900 border border-gray-200 dark:border-dark-800 rounded-3xl p-8 md:p-10 flex flex-col md:flex-row items-center justify-between gap-8 relative overflow-hidden shadow-soft">
            <div class="absolute right-0 top-0 w-64 h-full bg-gradient-to-l from-indigo-50 dark:from-indigo-900/10 to-transparent pointer-events-none"></div>
            <div class="flex items-center gap-6 z-10 w-full md:w-auto">
              <div class="w-16 h-16 rounded-2xl bg-indigo-50 dark:bg-indigo-900/20 flex items-center justify-center text-indigo-600 dark:text-indigo-400 text-3xl shrink-0">
                <Icon name="cube" size="lg" />
              </div>
              <div>
                <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-1">企业定制方案</h3>
                <p class="text-gray-500 dark:text-dark-400 text-sm">私有化部署 / 专属独享池 / SLA 服务保障</p>
              </div>
            </div>
            <div class="z-10 w-full md:w-auto">
              <MagneticButton>
                <button class="w-full md:w-auto bg-white dark:bg-dark-800 border border-gray-200 dark:border-dark-700 hover:border-indigo-500 hover:text-indigo-600 dark:hover:border-indigo-400 dark:hover:text-indigo-400 text-gray-700 dark:text-dark-300 font-bold py-3 px-8 rounded-xl shadow-sm transition-all flex items-center justify-center gap-2">
                  <Icon name="mail" size="sm" /> 联系商务代表
                </button>
              </MagneticButton>
            </div>
          </div>
        </GlowCard>
      </SlideIn>
    </section>

    <!-- Footer -->
    <footer class="bg-white dark:bg-dark-900 border-t border-gray-200 dark:border-dark-800 pt-16 pb-8">
      <div class="max-w-7xl mx-auto px-6 grid grid-cols-1 md:grid-cols-4 gap-12 mb-16">
        <div class="col-span-1">
          <div class="text-2xl font-extrabold tracking-tighter flex items-center gap-2 mb-4">
            <div class="bg-gray-100 dark:bg-dark-800 w-8 h-8 rounded-lg flex items-center justify-center text-gray-900 dark:text-white">
              <Icon name="terminal" size="sm" />
            </div>
            <div>
              <span class="text-gray-900 dark:text-white text-xl">C</span>
              <span class="text-primary-600 dark:text-primary-400 text-xl">coder</span>
              <span class="text-gray-400 dark:text-dark-500 text-xl">.me</span>
            </div>
          </div>
          <p class="text-gray-500 dark:text-dark-400 text-sm mb-4">专业的 AI 生产力接口中转引擎，为开发者提供纯净、稳定、高性价比的算力支持。</p>
        </div>

        <div>
          <h4 class="text-gray-900 dark:text-white font-bold mb-6">产品服务</h4>
          <ul class="space-y-4 text-sm text-gray-500 dark:text-dark-400">
            <li>
              <router-link :to="isAuthenticated ? '/model-plaza' : '/login'" class="hover:text-primary-600 dark:hover:text-primary-400 transition">API 模型列表</router-link>
            </li>
            <li><a href="#pricing" class="hover:text-primary-600 dark:hover:text-primary-400 transition">价格与套餐</a></li>
          </ul>
        </div>

        <div>
          <h4 class="text-gray-900 dark:text-white font-bold mb-6">帮助支持</h4>
          <ul class="space-y-4 text-sm text-gray-500 dark:text-dark-400">
            <li><a href="#faq" class="hover:text-primary-600 dark:hover:text-primary-400 transition">常见问题</a></li>
            <li>
              <router-link :to="isAuthenticated ? '/tutorial' : '/login'" class="hover:text-primary-600 dark:hover:text-primary-400 transition">快速接入</router-link>
            </li>
          </ul>
        </div>

        <div>
          <h4 class="text-gray-900 dark:text-white font-bold mb-6">获得支持</h4>
          <div class="flex items-center gap-3 text-gray-700 dark:text-dark-300 text-sm mb-4 bg-gray-50 dark:bg-dark-800 w-fit px-4 py-2 rounded-lg border border-gray-200 dark:border-dark-700">
            <Icon name="chat" size="sm" class="text-green-500" />
            <span class="font-mono">ccoder-support</span>
          </div>
        </div>
      </div>

      <div class="max-w-7xl mx-auto px-6 text-center text-gray-400 dark:text-dark-500 text-xs border-t border-gray-100 dark:border-dark-800 pt-8">
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

const faqs = [
  {
    question: '如何开始使用 API 服务？',
    answer: '注册账号后，在控制台生成您的专属 API Key。将您代码或应用中的 BaseURL 替换为我们的中转地址即可直接使用，100% 兼容官方格式。'
  },
  {
    question: '如何计费？有最低消费吗？',
    answer: '无最低消费。购买的月卡额度直接存入账户，0 倍率损耗，每一分钱都实打实用于 Token 消耗。'
  },
  {
    question: '服务的稳定性如何？会经常 429 吗？',
    answer: '系统采用企业级多节点负载均衡与智能账号池调度，遭遇官网限流时秒级无感切换。购买专业版/旗舰版更可享有专属高并发通道，告别 429 报错。'
  },
  {
    question: '如何保障数据安全和隐私？',
    answer: '所有请求均采用 256-bit SSL 加密传输，内存转发后即刻销毁，系统绝不落盘、不保存您的任何 Prompt，更不会用于 AI 模型训练。'
  },
  {
    question: '充值的额度会过期吗？',
    answer: '不会。您购买的 API 总额度是永久有效的。套餐内的"每日高速配额"仅用于控制高优先级调度，降级后只会变慢，余额绝对不会清零。'
  }
]

function toggleFaq(index: number) {
  activeFaq.value = activeFaq.value === index ? null : index
}

function selectPlan(_tier: string) {
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
</style>
