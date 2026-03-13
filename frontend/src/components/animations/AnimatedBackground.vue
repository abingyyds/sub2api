<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface Props {
  type?: 'particles' | 'gradient-mesh' | 'waves'
  particleCount?: number
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'particles',
  particleCount: 50,
  color: '#3b82f6'
})

const canvas = ref<HTMLCanvasElement>()
let ctx: CanvasRenderingContext2D | null = null
let animationId: number
let particles: Particle[] = []

interface Particle {
  x: number
  y: number
  vx: number
  vy: number
  radius: number
  opacity: number
}

const initParticles = () => {
  if (!canvas.value) return

  const width = canvas.value.width
  const height = canvas.value.height

  particles = Array.from({ length: props.particleCount }, () => ({
    x: Math.random() * width,
    y: Math.random() * height,
    vx: (Math.random() - 0.5) * 0.5,
    vy: (Math.random() - 0.5) * 0.5,
    radius: Math.random() * 2 + 1,
    opacity: Math.random() * 0.5 + 0.2
  }))
}

const drawParticles = () => {
  if (!ctx || !canvas.value) return

  const width = canvas.value.width
  const height = canvas.value.height

  ctx.clearRect(0, 0, width, height)

  particles.forEach(particle => {
    particle.x += particle.vx
    particle.y += particle.vy

    if (particle.x < 0 || particle.x > width) particle.vx *= -1
    if (particle.y < 0 || particle.y > height) particle.vy *= -1

    if (ctx) {
      ctx.beginPath()
      ctx.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2)
      ctx.fillStyle = `${props.color}${Math.floor(particle.opacity * 255).toString(16).padStart(2, '0')}`
      ctx.fill()
    }
  })

  animationId = requestAnimationFrame(drawParticles)
}

const drawGradientMesh = () => {
  if (!ctx || !canvas.value) return

  const width = canvas.value.width
  const height = canvas.value.height
  const time = Date.now() * 0.001

  ctx.clearRect(0, 0, width, height)

  const gradient = ctx.createLinearGradient(
    0,
    0,
    width * Math.cos(time * 0.5),
    height * Math.sin(time * 0.5)
  )

  gradient.addColorStop(0, `${props.color}10`)
  gradient.addColorStop(0.5, `${props.color}30`)
  gradient.addColorStop(1, `${props.color}10`)

  if (ctx) {
    ctx.fillStyle = gradient
    ctx.fillRect(0, 0, width, height)
  }

  animationId = requestAnimationFrame(drawGradientMesh)
}

const resizeCanvas = () => {
  if (!canvas.value) return

  canvas.value.width = canvas.value.offsetWidth
  canvas.value.height = canvas.value.offsetHeight

  if (props.type === 'particles') {
    initParticles()
  }
}

onMounted(() => {
  if (!canvas.value) return

  ctx = canvas.value.getContext('2d')
  resizeCanvas()

  if (props.type === 'particles') {
    drawParticles()
  } else if (props.type === 'gradient-mesh') {
    drawGradientMesh()
  }

  window.addEventListener('resize', resizeCanvas)
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
  window.removeEventListener('resize', resizeCanvas)
})
</script>

<template>
  <canvas
    ref="canvas"
    class="absolute inset-0 w-full h-full pointer-events-none"
    :class="{ 'opacity-50': type === 'particles' }"
  />
</template>
