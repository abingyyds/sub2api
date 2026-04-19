export const CHUNK_RELOAD_QUERY_KEY = '__chunk_reload'
export const CHUNK_RELOAD_STORAGE_KEY = 'chunk_reload_attempted'
const CHUNK_RELOAD_COOLDOWN_MS = 10_000
const CHUNK_ERROR_PATTERNS = [
  'failed to fetch dynamically imported module',
  'error loading dynamically imported module',
  'importing a module script failed',
  'loading chunk',
  'loading css chunk',
  'unable to preload css',
  'unable to preload dynamically imported module',
  'vite preload',
  'preload error'
]

type ChunkErrorTarget = EventTarget | null | undefined

function extractTargetAssetUrl(target: ChunkErrorTarget): string {
  if (!target) {
    return ''
  }

  if (target instanceof HTMLScriptElement) {
    return target.src || ''
  }

  if (target instanceof HTMLLinkElement) {
    return target.href || ''
  }

  return ''
}

function collectErrorFragments(error: unknown, visited = new Set<unknown>()): string[] {
  if (error == null || visited.has(error)) {
    return []
  }

  if (typeof error === 'object' || typeof error === 'function') {
    visited.add(error)
  }

  if (typeof error === 'string') {
    return [error]
  }

  if (typeof ErrorEvent !== 'undefined' && error instanceof ErrorEvent) {
    return [
      error.message,
      extractTargetAssetUrl(error.target),
      ...collectErrorFragments(error.error, visited)
    ].filter(Boolean)
  }

  if (typeof PromiseRejectionEvent !== 'undefined' && error instanceof PromiseRejectionEvent) {
    return collectErrorFragments(error.reason, visited)
  }

  if (error instanceof Error) {
    const cause = 'cause' in error ? (error as Error & { cause?: unknown }).cause : undefined
    return [
      error.name,
      error.message,
      ...collectErrorFragments(cause, visited)
    ].filter(Boolean)
  }

  if (typeof error === 'object') {
    const candidate = error as {
      message?: unknown
      reason?: unknown
      cause?: unknown
      payload?: unknown
      target?: EventTarget | null
      srcElement?: EventTarget | null
      filename?: unknown
      href?: unknown
      src?: unknown
    }

    return [
      typeof candidate.message === 'string' ? candidate.message : '',
      typeof candidate.filename === 'string' ? candidate.filename : '',
      typeof candidate.href === 'string' ? candidate.href : '',
      typeof candidate.src === 'string' ? candidate.src : '',
      extractTargetAssetUrl(candidate.target),
      extractTargetAssetUrl(candidate.srcElement),
      ...collectErrorFragments(candidate.reason, visited),
      ...collectErrorFragments(candidate.cause, visited),
      ...collectErrorFragments(candidate.payload, visited)
    ].filter(Boolean)
  }

  return [String(error)]
}

function looksLikeAssetLoadFailure(error: unknown): boolean {
  const fragments = collectErrorFragments(error)

  return fragments.some((fragment) => {
    const normalized = fragment.toLowerCase()
    return normalized.includes('/assets/') && (normalized.includes('.js') || normalized.includes('.css'))
  })
}

export function isChunkLoadError(error: unknown): boolean {
  if (error instanceof Error && error.name === 'ChunkLoadError') {
    return true
  }

  const normalizedFragments = collectErrorFragments(error)
    .map((fragment) => fragment.toLowerCase())
    .filter(Boolean)

  return normalizedFragments.some((fragment) => {
    return CHUNK_ERROR_PATTERNS.some((pattern) => fragment.includes(pattern))
  }) || looksLikeAssetLoadFailure(error)
}

export function buildChunkReloadUrl(path: string, reloadMarker: string | number): string {
  const url = new URL(path, window.location.origin)
  url.searchParams.set(CHUNK_RELOAD_QUERY_KEY, String(reloadMarker))
  return url.toString()
}

export function recoverFromChunkError(path: string, error: unknown): boolean {
  if (!isChunkLoadError(error)) {
    return false
  }

  const lastReload = sessionStorage.getItem(CHUNK_RELOAD_STORAGE_KEY)
  const now = Date.now()

  if (lastReload && now - Number(lastReload) <= CHUNK_RELOAD_COOLDOWN_MS) {
    return false
  }

  sessionStorage.setItem(CHUNK_RELOAD_STORAGE_KEY, String(now))
  window.location.assign(buildChunkReloadUrl(path, now))
  return true
}

export function installChunkRecoveryListeners(getTargetPath: () => string): () => void {
  const handleRecovery = (error: unknown, preventDefault?: () => void): void => {
    if (recoverFromChunkError(getTargetPath(), error)) {
      preventDefault?.()
    }
  }

  const handlePreloadError = (event: Event): void => {
    const customEvent = event as Event & { payload?: unknown; preventDefault?: () => void }
    handleRecovery(customEvent.payload ?? customEvent, () => customEvent.preventDefault?.())
  }

  const handleUnhandledRejection = (event: PromiseRejectionEvent): void => {
    handleRecovery(event, () => event.preventDefault())
  }

  const handleWindowError = (event: ErrorEvent): void => {
    handleRecovery(event, () => event.preventDefault())
  }

  window.addEventListener('vite:preloadError', handlePreloadError as EventListener)
  window.addEventListener('unhandledrejection', handleUnhandledRejection)
  window.addEventListener('error', handleWindowError)

  return () => {
    window.removeEventListener('vite:preloadError', handlePreloadError as EventListener)
    window.removeEventListener('unhandledrejection', handleUnhandledRejection)
    window.removeEventListener('error', handleWindowError)
  }
}
