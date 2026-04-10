import { describe, expect, it } from 'vitest'

import { isSidebarPathActive } from '../sidebarPath'

describe('isSidebarPathActive', () => {
  const knownPaths = ['/admin/dashboard', '/admin/ops', '/admin/ops/status', '/admin/users']

  it('matches exact routes', () => {
    expect(isSidebarPathActive('/admin/ops', '/admin/ops', knownPaths)).toBe(true)
    expect(isSidebarPathActive('/admin/ops/status', '/admin/ops/status', knownPaths)).toBe(true)
  })

  it('does not keep parent route active when a more specific sidebar item matches', () => {
    expect(isSidebarPathActive('/admin/ops/status', '/admin/ops', knownPaths)).toBe(false)
  })

  it('keeps parent route active for unmatched descendants', () => {
    expect(isSidebarPathActive('/admin/ops/providers', '/admin/ops', knownPaths)).toBe(true)
  })
})
