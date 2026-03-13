<script setup lang="ts">
import { useMotion } from '@vueuse/motion'
import { ref, onMounted } from 'vue'

interface Props {
  delay?: number
  duration?: number
}

const props = withDefaults(defineProps<Props>(), {
  delay: 0,
  duration: 600
})

const target = ref<HTMLElement>()

onMounted(() => {
  if (target.value) {
    useMotion(target.value, {
      initial: { opacity: 0 },
      enter: {
        opacity: 1,
        transition: {
          delay: props.delay,
          duration: props.duration
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
