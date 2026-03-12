<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950">
    <!-- Navigation -->
    <nav class="sticky top-0 z-50 border-b border-gray-200/50 bg-white/80 backdrop-blur-lg dark:border-dark-800/50 dark:bg-dark-900/80">
      <div class="mx-auto max-w-4xl px-6 py-4">
        <div class="flex items-center justify-between">
          <router-link to="/home" class="flex items-center gap-3">
            <Logo :size="32" theme="auto" />
            <span class="text-lg font-bold text-gray-900 dark:text-white">ccoder</span>
          </router-link>
          <router-link to="/home" class="text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">
            &larr; 返回首页
          </router-link>
        </div>
      </div>
    </nav>

    <!-- Content -->
    <main class="mx-auto max-w-4xl px-6 py-16">
      <h1 class="mb-2 text-3xl font-bold text-gray-900 dark:text-white">用户协议</h1>
      <p class="mb-12 text-sm text-gray-500 dark:text-dark-400">更新日期：2026年1月1日</p>

      <div class="space-y-8">
        <!-- 前言 -->
        <section>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <p class="text-gray-700 dark:text-dark-300">
              欢迎使用 ccoder 平台（以下简称"本平台"）。在使用本平台提供的各项服务之前，请您仔细阅读并充分理解本用户协议（以下简称"本协议"）的全部内容。您注册、登录、使用本平台服务的行为，即视为您已阅读并同意接受本协议的全部条款。如您不同意本协议的任何内容，请立即停止使用本平台服务。
            </p>
          </div>
        </section>

        <!-- 第一条：服务内容 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第一条 服务内容</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>ccoder 是一个 AI API 聚合与转发服务平台，为用户提供统一的接口来访问多种人工智能模型及相关服务。</li>
              <li>本平台通过聚合第三方 AI 服务提供商的 API，为用户提供便捷、高效的 API 调用体验，包括但不限于大语言模型、图像生成、语音处理等 AI 能力。</li>
              <li>本平台有权根据业务发展需要，随时调整、新增或终止部分或全部服务内容，并在平台上发布公告通知用户。</li>
              <li>本平台提供的服务依赖于第三方 AI 服务提供商，因第三方原因导致的服务变更或中断，本平台将尽合理努力通知用户，但不承担因此造成的损失。</li>
            </ol>
          </div>
        </section>

        <!-- 第二条：账号管理 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第二条 账号管理</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>您在注册本平台账号时，应提供真实、准确、完整的信息，并在信息发生变更时及时更新。</li>
              <li>您应妥善保管账号及密码，因账号密码保管不当导致的任何损失由您自行承担。如发现账号存在异常使用情况，请立即联系本平台。</li>
              <li>每位用户仅可注册一个账号，不得通过任何方式转让、借用、出租或出售账号。</li>
              <li>如本平台发现您注册信息存在虚假、违规等情况，有权暂停或终止您的账号使用权限，由此产生的后果由您自行承担。</li>
              <li>您的 API Key 属于敏感凭证信息，请妥善保管。因 API Key 泄露产生的一切费用及法律责任，由您自行承担。</li>
            </ol>
          </div>
        </section>

        <!-- 第三条：使用规范 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第三条 使用规范</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>您在使用本平台服务时，应遵守中华人民共和国及您所在地区的相关法律法规，不得利用本平台从事任何违法违规活动。</li>
              <li>
                您不得利用本平台服务进行以下行为：
                <ul class="mt-2 list-disc space-y-1 pl-5 text-gray-600 dark:text-dark-400">
                  <li>生成、传播违反法律法规或公序良俗的内容；</li>
                  <li>侵犯他人知识产权、隐私权、名誉权等合法权益；</li>
                  <li>发送垃圾信息、恶意代码或进行网络攻击；</li>
                  <li>干扰、破坏本平台或第三方服务的正常运行；</li>
                  <li>对本平台进行反向工程、反编译或以其他方式试图获取本平台源代码；</li>
                  <li>以任何方式规避本平台的技术保护措施或使用限制；</li>
                  <li>其他违反本协议或本平台规则的行为。</li>
                </ul>
              </li>
              <li>如您违反上述规范，本平台有权立即暂停或终止您的服务，并保留追究法律责任的权利。</li>
            </ol>
          </div>
        </section>

        <!-- 第四条：API使用与调用 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第四条 API 使用与调用</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>您通过本平台调用 API 服务时，应遵守本平台的调用规则及各 AI 模型提供方的使用条款。</li>
              <li>本平台采用预付费或后付费的计费方式，具体费率以平台页面展示为准。本平台有权根据运营成本变化合理调整价格，并提前通知用户。</li>
              <li>您应合理使用 API 调用额度，不得通过自动化程序恶意刷取额度或滥用服务资源。</li>
              <li>本平台将尽力保障 API 服务的稳定性和可用性，但不保证服务完全不中断或无错误。因系统维护、升级或不可抗力等原因导致的服务中断，本平台将尽快恢复并通知用户。</li>
              <li>您通过 API 传输和处理的数据内容由您自行负责。本平台不会主动访问或使用您的业务数据，但为保障服务安全，本平台有权对异常调用行为进行监控和处理。</li>
            </ol>
          </div>
        </section>

        <!-- 第五条：费用与支付 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第五条 费用与支付</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>本平台各项服务的具体收费标准以平台页面公示为准。</li>
              <li>您应按照本平台规定的方式及时支付相关费用。如因欠费导致服务中断，本平台不承担相关责任。</li>
              <li>由于数字服务的特殊性，已消耗的 API 调用额度原则上不予退款。如因本平台原因导致的服务异常，您可联系客服协商处理。</li>
              <li>本平台有权根据市场情况及运营成本对服务价格进行调整，调整前将在平台上提前公告。</li>
            </ol>
          </div>
        </section>

        <!-- 第六条：知识产权 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第六条 知识产权</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>本平台的商标、标识、界面设计、软件程序、技术文档等所有知识产权均归 ccoder 所有，未经书面许可，任何人不得擅自使用。</li>
              <li>您通过本平台 API 生成的内容，其知识产权归属遵循相关法律法规及各 AI 模型提供方的约定。</li>
              <li>您应保证通过本平台提交的输入内容不侵犯任何第三方的知识产权。因您的输入内容引发的知识产权纠纷，由您自行承担全部责任。</li>
            </ol>
          </div>
        </section>

        <!-- 第七条：隐私保护 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第七条 隐私保护</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>本平台重视用户隐私保护，将依据相关法律法规收集、使用和存储您的个人信息。</li>
              <li>本平台仅收集提供服务所必要的信息，不会向第三方出售您的个人信息。</li>
              <li>为保障服务质量和安全，本平台可能会记录 API 调用日志（不包含具体传输内容），相关日志将在合理期限后自动清除。</li>
            </ol>
          </div>
        </section>

        <!-- 第八条：免责声明 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第八条 免责声明</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>本平台作为 AI API 聚合与转发服务平台，仅提供技术中转服务，不对 AI 模型生成的内容的准确性、完整性、合法性及适用性做出任何保证。</li>
              <li>您基于 AI 生成内容作出的任何判断或决策，由您自行承担风险和责任。</li>
              <li>因不可抗力（包括但不限于自然灾害、政策变化、网络攻击、第三方服务中断等）导致的服务中断或数据损失，本平台不承担责任。</li>
              <li>本平台不对因您违反本协议或相关法律法规而导致的任何直接或间接损失承担责任。</li>
              <li>在法律允许的最大范围内，本平台对因使用或无法使用本服务而产生的任何间接、附带、特殊、惩罚性损害赔偿不承担责任。</li>
            </ol>
          </div>
        </section>

        <!-- 第九条：协议修改与终止 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第九条 协议修改与终止</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>本平台有权根据业务发展及法律法规的要求修订本协议。修订后的协议将在平台上公布，自公布之日起生效。</li>
              <li>如您在协议修订后继续使用本平台服务，视为您同意修订后的协议内容。如您不同意修订内容，请停止使用本平台服务。</li>
              <li>您可随时注销账号并停止使用本平台服务。账号注销后，本平台将在合理期限内删除您的个人信息，但法律法规要求保留的除外。</li>
              <li>本平台有权在您违反本协议或相关法律法规时，暂停或终止向您提供服务，且无需承担任何责任。</li>
            </ol>
          </div>
        </section>

        <!-- 第十条：其他条款 -->
        <section>
          <h2 class="mb-4 text-lg font-semibold text-gray-700 dark:text-dark-200">第十条 其他条款</h2>
          <div class="rounded-xl border border-gray-200 bg-white p-6 dark:border-dark-700 dark:bg-dark-800">
            <ol class="list-decimal space-y-3 pl-5 text-gray-700 dark:text-dark-300">
              <li>本协议的成立、生效、履行、解释及纠纷解决，适用中华人民共和国法律。</li>
              <li>因本协议引起的或与本协议有关的争议，双方应首先友好协商解决；协商不成的，任何一方有权向本平台所在地有管辖权的人民法院提起诉讼。</li>
              <li>本协议中的任何条款被认定为无效或不可执行的，不影响其他条款的效力。</li>
              <li>本平台未行使或延迟行使本协议项下的任何权利，不构成对该权利的放弃。</li>
              <li>如您对本协议有任何疑问，欢迎通过平台内的客服渠道与我们联系。</li>
            </ol>
          </div>
        </section>
      </div>
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-200 px-6 py-8 dark:border-dark-800">
      <div class="mx-auto max-w-4xl">
        <div class="flex flex-col items-center gap-4">
          <p class="text-xs text-gray-400 dark:text-dark-500">&copy; {{ currentYear }} ccoder. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Logo from '@/components/Logo.vue'

const currentYear = computed(() => new Date().getFullYear())
</script>
