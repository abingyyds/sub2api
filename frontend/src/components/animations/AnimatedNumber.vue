<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  value: number
  duration?: number
  decimals?: number
}

const props = withDefaults(defineProps<Props>(), {
  duration: 1000,
  decimals: 0
})

const displayValue = ref(0)
const target = ref<HTMLElement>()

watch(() => props.value, (newValue) => {
  const startValue = displayValue.value
  const endValue = newValue
  const startTime = Date.now()

  const animate = () => {
    const elapsed = Date.now() - startTime
    const progress = Math.min(elapsed / props.duration, 1)

    const easeOutQuart = 1 - Math.pow(1 - progress, 4)
    displayValue.value = startValue + (endValue - startValue) * easeOutQuart

    if (progress < 1) {
      requestAnimationFrame(animate)
    } else {
      displayValue.value = endValue
    }
  }

  animate()
}, { immediate: true })
</script>

<template>
  <span ref="target">
    {{ displayValue.toFixed(decimals) }}
  </span>
</template>
