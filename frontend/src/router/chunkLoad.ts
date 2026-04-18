export const CHUNK_RELOAD_QUERY_KEY = '__chunk_reload'

export function isChunkLoadError(error: unknown): boolean {
  if (error instanceof Error && error.name === 'ChunkLoadError') {
    return true
  }

  const message = error instanceof Error
    ? `${error.name}: ${error.message}`
    : String(error ?? '')
  const normalizedMessage = message.toLowerCase()

  return [
    'failed to fetch dynamically imported module',
    'error loading dynamically imported module',
    'importing a module script failed',
    'loading chunk',
    'loading css chunk',
    'unable to preload css',
    'unable to preload dynamically imported module'
  ].some((pattern) => normalizedMessage.includes(pattern))
}

export function buildChunkReloadUrl(path: string, reloadMarker: string | number): string {
  const url = new URL(path, window.location.origin)
  url.searchParams.set(CHUNK_RELOAD_QUERY_KEY, String(reloadMarker))
  return url.toString()
}
