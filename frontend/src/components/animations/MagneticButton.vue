<script setup lang="ts">
import { ref } from 'vue'
import { useMotion } from '@vueuse/motion'

interface Props {
  strength?: number
}

const props = withDefaults(defineProps<Props>(), {
  strength: 0.3
})

const button = ref<HTMLElement>()
const motionInstance = ref()

const handleMouseMove = (e: MouseEvent) => {
  if (!button.value) return

  const rect = button.value.getBoundingClientRect()
  const centerX = rect.left + rect.width / 2
  const centerY = rect.top + rect.height / 2

  const deltaX = (e.clientX - centerX) * props.strength
  const deltaY = (e.clientY - centerY) * props.strength

  if (motionInstance.value) {
    motionInstance.value.apply({
      x: deltaX,
      y: deltaY,
      transition: {
        type: 'spring',
        stiffness: 150,
        damping: 15
      }
    })
  }
}

const handleMouseLeave = () => {
  if (motionInstance.value) {
    motionInstance.value.apply({
      x: 0,
      y: 0,
      transition: {
        type: 'spring',
        stiffness: 150,
        damping: 15
      }
    })
  }
}

const setupMotion = () => {
  if (button.value) {
    motionInstance.value = useMotion(button.value, {
      initial: { x: 0, y: 0 }
    })
  }
}
</script>

<template>
  <div
    ref="button"
    class="inline-block"
    @mousemove="handleMouseMove"
    @mouseleave="handleMouseLeave"
    @vue:mounted="setupMotion"
  >
    <slot />
  </div>
</template>
