import { ref, watch } from 'vue'

interface CounterAnimationOptions {
  duration?: number
  decimals?: number
  easing?: (t: number) => number
}

export function useCounterAnimation(
  targetValue: number | (() => number),
  options: CounterAnimationOptions = {}
) {
  const {
    duration = 1000,
    decimals = 0,
    easing = (t: number) => 1 - Math.pow(1 - t, 4) // easeOutQuart
  } = options

  const current = ref(0)
  const isAnimating = ref(false)

  const animate = (from: number, to: number) => {
    isAnimating.value = true
    const startTime = Date.now()

    const step = () => {
      const elapsed = Date.now() - startTime
      const progress = Math.min(elapsed / duration, 1)
      const easedProgress = easing(progress)

      current.value = from + (to - from) * easedProgress

      if (progress < 1) {
        requestAnimationFrame(step)
      } else {
        current.value = to
        isAnimating.value = false
      }
    }

    step()
  }

  const getValue = () => {
    return typeof targetValue === 'function' ? targetValue() : targetValue
  }

  // 初始化
  const initialValue = getValue()
  animate(0, initialValue)

  // 监听变化
  if (typeof targetValue !== 'function') {
    watch(() => targetValue, (newValue) => {
      animate(current.value, newValue)
    })
  }

  return {
    current,
    isAnimating,
    animate,
    formatted: () => current.value.toFixed(decimals)
  }
}
