import { defineAsyncComponent, type AsyncComponentLoader, type Component } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import ChunkLoadError from '@/components/common/ChunkLoadError.vue'
import { getCurrentLocationPath, isChunkLoadError, recoverFromChunkError } from './chunkLoad'

export type ComponentImportFn = AsyncComponentLoader<Component>

export const chunkResilientLoaderKey = Symbol('chunkResilientLoader')

type ChunkResilientComponent = Component & {
  [chunkResilientLoaderKey]?: ComponentImportFn
}

export function defineChunkResilientAsyncComponent(
  loader: ComponentImportFn,
  getTargetPath: () => string = getCurrentLocationPath
): Component {
  const component = defineAsyncComponent({
    loader,
    errorComponent: ChunkLoadError,
    delay: 150,
    timeout: 30_000,
    onError(error, _retry, fail) {
      if (isChunkLoadError(error) && recoverFromChunkError(getTargetPath(), error)) {
        return
      }
      fail()
    }
  }) as ChunkResilientComponent

  Object.defineProperty(component, chunkResilientLoaderKey, {
    value: loader,
    enumerable: false
  })

  return component
}

export function getChunkResilientLoader(component: unknown): ComponentImportFn | null {
  const candidate = component as Partial<ChunkResilientComponent> | null
  return candidate?.[chunkResilientLoaderKey] || null
}

function wrapLazyComponent(component: unknown): unknown {
  if (typeof component !== 'function') {
    return component
  }

  return defineChunkResilientAsyncComponent(component as ComponentImportFn)
}

export function wrapRouteLazyComponents(routeRecords: RouteRecordRaw[]): RouteRecordRaw[] {
  return routeRecords.map((route) => {
    const wrapped: RouteRecordRaw = { ...route }

    if ('component' in wrapped && wrapped.component) {
      wrapped.component = wrapLazyComponent(wrapped.component) as RouteRecordRaw['component']
    }

    if ('components' in wrapped && wrapped.components) {
      wrapped.components = Object.fromEntries(
        Object.entries(wrapped.components).map(([name, component]) => [
          name,
          wrapLazyComponent(component)
        ])
      ) as RouteRecordRaw['components']
    }

    if (wrapped.children) {
      wrapped.children = wrapRouteLazyComponents(wrapped.children)
    }

    return wrapped
  })
}
