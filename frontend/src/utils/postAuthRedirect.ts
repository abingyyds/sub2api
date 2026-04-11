import type { Router } from 'vue-router'

function fallbackPath(path: string | null | undefined, fallback: string): string {
  if (!path) return fallback
  if (!path.startsWith('/')) return fallback
  if (path.startsWith('//')) return fallback
  if (path.includes('://')) return fallback
  if (path.includes('\n') || path.includes('\r')) return fallback
  return path
}

export function sanitizePostAuthRedirect(
  path: string | null | undefined,
  fallback = '/pricing'
): string {
  return fallbackPath(path, fallback)
}

export async function redirectAfterAuth(
  router: Router,
  path: string | null | undefined,
  fallback = '/pricing'
): Promise<void> {
  const safePath = sanitizePostAuthRedirect(path, fallback)
  const resolved = router.resolve(safePath)

  // Pricing currently has an intermittent client-side hydration/transition issue
  // right after authentication. A hard navigation avoids the blank first render.
  if (resolved.path === '/pricing') {
    window.location.assign(resolved.fullPath)
    return
  }

  await router.replace(resolved.fullPath)
}
