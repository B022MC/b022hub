import { sanitizeUrl } from './url'

const ORIGINAL_PROJECT_HOSTS = new Set([
  'github.com',
  'www.github.com',
  'raw.githubusercontent.com',
])

const ORIGINAL_PROJECT_PATH_SEGMENT = '/wei-shaw/sub2api'

export function isOriginalProjectUrl(value: string): boolean {
  const normalized = sanitizeUrl(value)
  if (!normalized) {
    return false
  }

  try {
    const parsed = new URL(normalized)
    const hostname = parsed.hostname.toLowerCase()
    const pathname = parsed.pathname.toLowerCase()

    return ORIGINAL_PROJECT_HOSTS.has(hostname) && pathname.includes(ORIGINAL_PROJECT_PATH_SEGMENT)
  } catch {
    return false
  }
}

export function sanitizeProjectExternalUrl(value?: string | null): string {
  if (typeof value !== 'string') {
    return ''
  }

  const normalized = sanitizeUrl(value)
  if (!normalized || isOriginalProjectUrl(normalized)) {
    return ''
  }

  return normalized
}
