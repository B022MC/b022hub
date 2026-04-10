import { describe, expect, it, vi } from 'vitest'
import { defineComponent } from 'vue'
import { mount } from '@vue/test-utils'

import OpsStatusMatrixTable from '../OpsStatusMatrixTable.vue'
import type { OpsStatusMatrixRow } from '@/api/admin/ops'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string | number>) => {
        if (key === 'admin.ops.statusMatrix.tooltips.excludedErrors' && params) {
          return `已排除失败：${params.count}`
        }
        if (key === 'admin.ops.statusMatrix.table.excludedShort') {
          return '排除'
        }
        return key
      }
    })
  }
})

const DataTableStub = defineComponent({
  name: 'DataTable',
  props: {
    data: {
      type: Array,
      default: () => []
    }
  },
  template: `
    <div class="data-table-stub">
      <div
        v-for="row in data"
        :key="row.model"
        class="data-table-row"
      >
        <slot name="cell-cache_hit_rate" :row="row" :value="row.cache_hit_rate" />
        <slot name="cell-availability" :row="row" :value="row.availability" />
        <slot name="cell-timeline" :row="row" :value="row.buckets" />
      </div>
    </div>
  `
})

const TimelineStub = defineComponent({
  name: 'OpsStatusTimeline',
  template: '<div class="timeline-stub" />'
})

const rows: OpsStatusMatrixRow[] = [
  {
    platform: 'openai',
    group_id: 1,
    group_name: '主通道',
    model: 'gpt-4.1',
    success_count: 0,
    error_count: 0,
    excluded_error_count: 2,
    availability: null,
    cache_hit_rate: null,
    last_checked_at: null,
    last_success_at: null,
    last_latency_ms: null,
    buckets: []
  }
]

describe('OpsStatusMatrixTable', () => {
  it('shows "-" when availability is null', () => {
    const wrapper = mount(OpsStatusMatrixTable, {
      props: {
        rows,
        timeRange: '90m'
      },
      global: {
        stubs: {
          DataTable: DataTableStub,
          OpsStatusTimeline: TimelineStub,
          PlatformIcon: true
        }
      }
    })

    expect(wrapper.get('[data-testid="availability-value"]').text()).toBe('-')
  })

  it('renders excluded-error tooltip text without affecting availability value', () => {
    const wrapper = mount(OpsStatusMatrixTable, {
      props: {
        rows,
        timeRange: '24h'
      },
      global: {
        stubs: {
          DataTable: DataTableStub,
          OpsStatusTimeline: TimelineStub,
          PlatformIcon: true
        }
      }
    })

    expect(wrapper.get('[data-testid="availability-tooltip"]').attributes('title')).toBe('已排除失败：2')
    expect(wrapper.get('[data-testid="availability-value"]').text()).toBe('-')
  })

  it('renders cache hit rate when present', () => {
    const wrapper = mount(OpsStatusMatrixTable, {
      props: {
        rows: [{ ...rows[0], cache_hit_rate: 0.2 }],
        timeRange: '90m'
      },
      global: {
        stubs: {
          DataTable: DataTableStub,
          OpsStatusTimeline: TimelineStub,
          PlatformIcon: true
        }
      }
    })

    expect(wrapper.get('[data-testid="cache-hit-rate-value"]').text()).toBe('20.00%')
  })
})
