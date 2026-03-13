import { ref, onMounted } from 'vue'
import { useMotion } from '@vueuse/motion'
import { useIntersectionObserver } from '@vueuse/core'

interface ScrollAnimationOptions {
  threshold?: number
  rootMargin?: string
  once?: boolean
}

export function useScrollAnimation(options: ScrollAnimationOptions = {}) {
  const { threshold = 0.1, rootMargin = '0px', once = true } = options
  const target = ref<HTMLElement>()
  const isVisible = ref(false)

  onMounted(() => {
    if (!target.value) return

    const { stop } = useIntersectionObserver(
      target.value,
      ([{ isIntersecting }]) => {
        if (isIntersecting) {
          isVisible.value = true

          useMotion(target.value!, {
            initial: { opacity: 0, y: 50 },
            enter: {
              opacity: 1,
              y: 0,
              transition: {
                duration: 600,
                ease: 'easeOut'
              }
            }
          })

          if (once) {
            stop()
          }
        } else if (!once) {
          isVisible.value = false
        }
      },
      { threshold, rootMargin }
    )
  })

  return {
    target,
    isVisible
  }
}
