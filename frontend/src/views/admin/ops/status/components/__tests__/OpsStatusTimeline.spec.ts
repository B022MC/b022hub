import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import OpsStatusTimeline from '../OpsStatusTimeline.vue'
import type { OpsStatusMatrixBucket } from '@/api/admin/ops'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string | number>) => {
        if (key === 'admin.ops.statusMatrix.tooltips.timeline' && params) {
          return `${params.status}|${params.success}|${params.error}|${params.excluded}`
        }
        return key
      }
    })
  }
})

const makeBuckets = (count: number, status: OpsStatusMatrixBucket['status']): OpsStatusMatrixBucket[] => {
  const start = new Date('2026-04-10T00:00:00Z').getTime()
  return Array.from({ length: count }, (_, index) => ({
    bucket_start: new Date(start + index * 300000).toISOString(),
    bucket_end: new Date(start + (index + 1) * 300000).toISOString(),
    success_count: status === 'ok' ? 1 : 0,
    error_count: status === 'down' ? 1 : 0,
    excluded_error_count: 0,
    status
  }))
}

describe('OpsStatusTimeline', () => {
  it('renders 18 cells for a 90 minute window', () => {
    const wrapper = mount(OpsStatusTimeline, {
      props: {
        buckets: makeBuckets(18, 'ok')
      }
    })

    expect(wrapper.findAll('[data-testid="status-cell"]')).toHaveLength(18)
    expect(wrapper.attributes('data-bucket-count')).toBe('18')
  })

  it('renders 24 cells and applies status tone classes', () => {
    const buckets = makeBuckets(24, 'nodata')
    buckets[2] = {
      ...buckets[2],
      status: 'warn',
      success_count: 2,
      error_count: 1
    }

    const wrapper = mount(OpsStatusTimeline, {
      props: { buckets }
    })

    const cells = wrapper.findAll('[data-testid="status-cell"]')
    expect(cells).toHaveLength(24)
    expect(cells[2].classes()).toContain('status-timeline-cell-warn')
    expect(cells[2].attributes('data-status')).toBe('warn')
  })
})
