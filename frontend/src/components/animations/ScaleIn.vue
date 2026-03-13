<script setup lang="ts">
import { useMotion } from '@vueuse/motion'
import { ref, onMounted } from 'vue'

interface Props {
  delay?: number
  duration?: number
  initialScale?: number
}

const props = withDefaults(defineProps<Props>(), {
  delay: 0,
  duration: 600,
  initialScale: 0.8
})

const target = ref<HTMLElement>()

onMounted(() => {
  if (target.value) {
    useMotion(target.value, {
      initial: {
        scale: props.initialScale,
        opacity: 0
      },
      enter: {
        scale: 1,
        opacity: 1,
        transition: {
          delay: props.delay,
          duration: props.duration,
          ease: 'easeOut'
        }
      }
    })
  }
})
</script>

<template>
  <div ref="target">
    <slot />
  </div>
</template>
