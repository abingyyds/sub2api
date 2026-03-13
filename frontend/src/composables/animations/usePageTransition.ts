import { ref } from 'vue'

export function usePageTransition() {
  const isTransitioning = ref(false)

  const beforeEnter = () => {
    isTransitioning.value = true
  }

  const afterEnter = () => {
    isTransitioning.value = false
  }

  return {
    isTransitioning,
    beforeEnter,
    afterEnter
  }
}
