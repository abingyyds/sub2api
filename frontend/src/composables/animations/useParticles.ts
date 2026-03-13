import { ref, onMounted, onUnmounted } from 'vue'

interface Particle {
  x: number
  y: number
  vx: number
  vy: number
  radius: number
  opacity: number
  color: string
}

interface ParticlesOptions {
  count?: number
  color?: string
  speed?: number
  minRadius?: number
  maxRadius?: number
}

export function useParticles(options: ParticlesOptions = {}) {
  const {
    count = 50,
    color = '#3b82f6',
    speed = 0.5,
    minRadius = 1,
    maxRadius = 3
  } = options

  const canvas = ref<HTMLCanvasElement>()
  const particles = ref<Particle[]>([])
  let ctx: CanvasRenderingContext2D | null = null
  let animationId: number

  const initParticles = (width: number, height: number) => {
    particles.value = Array.from({ length: count }, () => ({
      x: Math.random() * width,
      y: Math.random() * height,
      vx: (Math.random() - 0.5) * speed,
      vy: (Math.random() - 0.5) * speed,
      radius: Math.random() * (maxRadius - minRadius) + minRadius,
      opacity: Math.random() * 0.5 + 0.2,
      color
    }))
  }

  const animate = () => {
    if (!ctx || !canvas.value) return

    const width = canvas.value.width
    const height = canvas.value.height

    ctx.clearRect(0, 0, width, height)

    particles.value.forEach(particle => {
      particle.x += particle.vx
      particle.y += particle.vy

      if (particle.x < 0 || particle.x > width) particle.vx *= -1
      if (particle.y < 0 || particle.y > height) particle.vy *= -1

      if (ctx) {
        ctx.beginPath()
        ctx.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2)
        ctx.fillStyle = `${particle.color}${Math.floor(particle.opacity * 255).toString(16).padStart(2, '0')}`
        ctx.fill()
      }
    })

    animationId = requestAnimationFrame(animate)
  }

  const resize = () => {
    if (!canvas.value) return
    canvas.value.width = canvas.value.offsetWidth
    canvas.value.height = canvas.value.offsetHeight
    initParticles(canvas.value.width, canvas.value.height)
  }

  const start = () => {
    if (!canvas.value) return
    ctx = canvas.value.getContext('2d')
    resize()
    animate()
  }

  const stop = () => {
    if (animationId) {
      cancelAnimationFrame(animationId)
    }
  }

  onMounted(() => {
    if (canvas.value) {
      start()
      window.addEventListener('resize', resize)
    }
  })

  onUnmounted(() => {
    stop()
    window.removeEventListener('resize', resize)
  })

  return {
    canvas,
    particles,
    start,
    stop,
    resize
  }
}
