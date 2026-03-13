<script setup lang="ts">
import { useMotion } from '@vueuse/motion'
import { ref, onMounted } from 'vue'

interface Props {
  delay?: number
  staggerDelay?: number
  duration?: number
}

const props = withDefaults(defineProps<Props>(), {
  delay: 0,
  staggerDelay: 100,
  duration: 600
})

const container = ref<HTMLElement>()

onMounted(() => {
  if (container.value) {
    const children = Array.from(container.value.children) as HTMLElement[]

    children.forEach((child, index) => {
      useMotion(child, {
        initial: {
          opacity: 0,
          y: 20
        },
        enter: {
          opacity: 1,
          y: 0,
          transition: {
            delay: props.delay + (index * props.staggerDelay),
            duration: props.duration,
            ease: 'easeOut'
          }
        }
      })
    })
  }
})
</script>

<template>
  <div ref="container">
    <slot />
  </div>
</template>
