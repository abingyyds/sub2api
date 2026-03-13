<script setup lang="ts">
interface Props {
  duration?: number
  delay?: number
}

withDefaults(defineProps<Props>(), {
  duration: 2000,
  delay: 0
})
</script>

<template>
  <div class="relative inline-block overflow-hidden">
    <slot />
    <div
      class="absolute inset-0 -translate-x-full animate-shimmer bg-gradient-to-r from-transparent via-white/20 to-transparent"
      :style="{
        animationDuration: `${duration}ms`,
        animationDelay: `${delay}ms`
      }"
    />
  </div>
</template>

<style scoped>
@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.animate-shimmer {
  animation: shimmer 2s infinite;
}
</style>
