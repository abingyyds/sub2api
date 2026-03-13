declare module '@vueuse/motion' {
  import type { Directive } from 'vue'

  export interface MotionVariants {
    initial?: Record<string, any>
    enter?: Record<string, any>
    leave?: Record<string, any>
    visible?: Record<string, any>
    visibleOnce?: Record<string, any>
  }

  export interface MotionInstance {
    apply: (variant: Record<string, any>) => void
    stop: () => void
  }

  export function useMotion(
    target: any,
    variants: MotionVariants
  ): MotionInstance

  export const MotionPlugin: {
    install: (app: any, options?: any) => void
  }

  export const vMotion: Directive
}
