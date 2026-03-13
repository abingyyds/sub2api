<script setup lang="ts">
import { useMotion } from '@vueuse/motion'
import { ref, onMounted, computed } from 'vue'

interface Props {
  direction?: 'up' | 'down' | 'left' | 'right'
  delay?: number
  duration?: number
  distance?: number
}

const props = withDefaults(defineProps<Props>(), {
  direction: 'up',
  delay: 0,
  duration: 600,
  distance: 50
})

const target = ref<HTMLElement>()

const initialTransform = computed(() => {
  switch (props.direction) {
    case 'up':
      return { y: props.distance, opacity: 0 }
    case 'down':
      return { y: -props.distance, opacity: 0 }
    case 'left':
      return { x: props.distance, opacity: 0 }
    case 'right':
      return { x: -props.distance, opacity: 0 }
  }
})

onMounted(() => {
  if (target.value) {
    useMotion(target.value, {
      initial: initialTransform.value,
      enter: {
        x: 0,
        y: 0,
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
