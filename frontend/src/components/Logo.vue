<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 100 100"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    :class="className"
  >
    <!-- Background circle with gradient -->
    <defs>
      <linearGradient id="logoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" :style="`stop-color:${gradientStart};stop-opacity:1`" />
        <stop offset="100%" :style="`stop-color:${gradientEnd};stop-opacity:1`" />
      </linearGradient>
    </defs>

    <!-- Outer circle -->
    <circle cx="50" cy="50" r="48" fill="url(#logoGradient)" opacity="0.1" />

    <!-- Code brackets < > -->
    <path
      d="M 30 35 L 20 50 L 30 65"
      :stroke="iconColor"
      stroke-width="4"
      stroke-linecap="round"
      stroke-linejoin="round"
      fill="none"
    />
    <path
      d="M 70 35 L 80 50 L 70 65"
      :stroke="iconColor"
      stroke-width="4"
      stroke-linecap="round"
      stroke-linejoin="round"
      fill="none"
    />

    <!-- Letter C in the center -->
    <path
      d="M 60 40 C 55 35, 45 35, 40 40 C 35 45, 35 55, 40 60 C 45 65, 55 65, 60 60"
      :stroke="iconColor"
      stroke-width="4"
      stroke-linecap="round"
      fill="none"
    />

    <!-- Small dots representing code -->
    <circle cx="50" cy="30" r="2" :fill="iconColor" opacity="0.6" />
    <circle cx="50" cy="70" r="2" :fill="iconColor" opacity="0.6" />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  size?: number | string
  className?: string
  theme?: 'light' | 'dark' | 'auto'
}>(), {
  size: 40,
  className: '',
  theme: 'auto'
})

const gradientStart = computed(() => {
  if (props.theme === 'dark') return '#3b82f6'
  if (props.theme === 'light') return '#2563eb'
  return '#3b82f6' // default
})

const gradientEnd = computed(() => {
  if (props.theme === 'dark') return '#8b5cf6'
  if (props.theme === 'light') return '#7c3aed'
  return '#8b5cf6' // default
})

const iconColor = computed(() => {
  if (props.theme === 'dark') return '#60a5fa'
  if (props.theme === 'light') return '#2563eb'
  return 'currentColor' // auto - follows text color
})
</script>
