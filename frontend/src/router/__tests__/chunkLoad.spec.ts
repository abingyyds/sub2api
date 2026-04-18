import { describe, expect, it } from 'vitest'
import { buildChunkReloadUrl, CHUNK_RELOAD_QUERY_KEY, isChunkLoadError } from '../chunkLoad'

describe('chunkLoad helpers', () => {
  it('identifies common dynamic import failures across browsers', () => {
    expect(isChunkLoadError(new Error('Failed to fetch dynamically imported module'))).toBe(true)
    expect(isChunkLoadError(new Error('Importing a module script failed.'))).toBe(true)
    expect(isChunkLoadError(new Error('error loading dynamically imported module'))).toBe(true)

    const namedError = new Error('Loading chunk 12 failed')
    namedError.name = 'ChunkLoadError'
    expect(isChunkLoadError(namedError)).toBe(true)
  })

  it('ignores unrelated router errors', () => {
    expect(isChunkLoadError(new Error('Navigation aborted'))).toBe(false)
    expect(isChunkLoadError(new Error('Network error'))).toBe(false)
  })

  it('builds a hard reload url for the pending route', () => {
    const url = new URL(buildChunkReloadUrl('/login?redirect=%2Fpricing#top', 123), window.location.origin)

    expect(url.pathname).toBe('/login')
    expect(url.searchParams.get('redirect')).toBe('/pricing')
    expect(url.searchParams.get(CHUNK_RELOAD_QUERY_KEY)).toBe('123')
    expect(url.hash).toBe('#top')
  })
})
