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
    class="relative rounded-lg transition-all duration-300"
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
  box-shadow: 0 0 20px var(--glow-color),
              0 0 40px var(--glow-color);
  transform: translateY(-2px);
}

.glow-active::before {
  content: '';
  position: absolute;
  inset: -2px;
  border-radius: inherit;
  background: linear-gradient(
    45deg,
    transparent,
    var(--glow-color),
    transparent
  );
  opacity: var(--glow-intensity);
  z-index: -1;
  animation: glow-pulse 2s ease-in-out infinite;
}

@keyframes glow-pulse {
  0%, 100% {
    opacity: calc(var(--glow-intensity) * 0.5);
  }
  50% {
    opacity: var(--glow-intensity);
  }
}
</style>
