<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  glowColor?: string
  glowIntensity?: number
}

withDefaults(defineProps<Props>(), {
  glowColor: 'rgb(59, 130, 246)',
  glowIntensity: 0.5
})

const isHovered = ref(false)
</script>

<template>
  <div
    class="relative rounded-2xl transition-all duration-300"
    :class="{ 'glow-active': isHovered }"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false"
    :style="{
      '--glow-color': glowColor,
      '--glow-intensity': glowIntensity
    }"
  >
    <slot />
  </div>
</template>

<style scoped>
.glow-active {
  box-shadow: 0 0 15px color-mix(in srgb, var(--glow-color) 30%, transparent),
              0 0 30px color-mix(in srgb, var(--glow-color) 15%, transparent);
  transform: translateY(-2px);
}
</style>
