import { ref } from 'vue'
import { useMotion } from '@vueuse/motion'

interface HoverAnimationOptions {
  scale?: number
  duration?: number
}

export function useHoverAnimation(options: HoverAnimationOptions = {}) {
  const { scale = 1.05, duration = 200 } = options
  const target = ref<HTMLElement>()
  const isHovered = ref(false)

  const setupHover = (element: HTMLElement) => {
    target.value = element

    const motionInstance = useMotion(element, {
      initial: { scale: 1 },
      enter: { scale: 1 }
    })

    element.addEventListener('mouseenter', () => {
      isHovered.value = true
      motionInstance.apply({
        scale,
        transition: { duration }
      })
    })

    element.addEventListener('mouseleave', () => {
      isHovered.value = false
      motionInstance.apply({
        scale: 1,
        transition: { duration }
      })
    })
  }

  return {
    target,
    isHovered,
    setupHover
  }
}
