import { describe, it, expect, vi, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import ImportDataModal from '@/components/admin/account/ImportDataModal.vue'

const { showError, showSuccess, importData } = vi.hoisted(() => ({
  showError: vi.fn(),
  showSuccess: vi.fn(),
  importData: vi.fn()
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError,
    showSuccess
  })
}))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    accounts: {
      importData
    }
  }
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  })
}))

describe('ImportDataModal', () => {
  beforeEach(() => {
    showError.mockReset()
    showSuccess.mockReset()
    importData.mockReset()
  })

  it('未选择文件时提示错误', async () => {
    const wrapper = mount(ImportDataModal, {
      props: { show: true, groups: [] },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' }
        }
      }
    })

    await wrapper.find('form').trigger('submit')
    expect(showError).toHaveBeenCalledWith('admin.accounts.dataImportSelectFile')
  })

  it('无效 JSON 时提示解析失败', async () => {
    const wrapper = mount(ImportDataModal, {
      props: { show: true, groups: [] },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' }
        }
      }
    })

    const input = wrapper.find('input[type="file"]')
    const file = new File(['invalid json'], 'data.json', { type: 'application/json' })
    Object.defineProperty(file, 'text', {
      value: () => Promise.resolve('invalid json')
    })
    Object.defineProperty(input.element, 'files', {
      value: [file]
    })

    await input.trigger('change')
    await wrapper.find('form').trigger('submit')
    await Promise.resolve()

    expect(showError).toHaveBeenCalledWith('admin.accounts.dataImportParseFailed')
  })

  it('选择分组后会将 group_ids 一起提交', async () => {
    importData.mockResolvedValue({
      proxy_created: 0,
      proxy_reused: 0,
      proxy_failed: 0,
      account_created: 1,
      account_failed: 0,
      errors: []
    })

    const wrapper = mount(ImportDataModal, {
      props: {
        show: true,
        groups: [
          {
            id: 2,
            name: 'codex',
            platform: 'openai'
          } as any
        ]
      },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' },
          GroupBadge: { template: '<span><slot /></span>' }
        }
      }
    })

    const input = wrapper.find('input[type="file"]')
    const file = new File(['{"type":"sub2api-data","version":1,"proxies":[],"accounts":[]}'], 'data.json', {
      type: 'application/json'
    })
    Object.defineProperty(file, 'text', {
      value: () => Promise.resolve('{"type":"sub2api-data","version":1,"proxies":[],"accounts":[]}')
    })
    Object.defineProperty(input.element, 'files', {
      value: [file]
    })

    await input.trigger('change')
    await wrapper.find('input[type="checkbox"]').setValue(true)
    await wrapper.find('form').trigger('submit')
    await Promise.resolve()

    expect(importData).toHaveBeenCalledWith({
      data: {
        type: 'sub2api-data',
        version: 1,
        proxies: [],
        accounts: []
      },
      group_ids: [2],
      skip_default_group_bind: true
    })
  })

  it('支持一次选择多个文件并逐个导入', async () => {
    importData
      .mockResolvedValueOnce({
        proxy_created: 1,
        proxy_reused: 0,
        proxy_failed: 0,
        account_created: 1,
        account_failed: 0,
        errors: []
      })
      .mockResolvedValueOnce({
        proxy_created: 0,
        proxy_reused: 1,
        proxy_failed: 0,
        account_created: 2,
        account_failed: 0,
        errors: []
      })

    const wrapper = mount(ImportDataModal, {
      props: { show: true, groups: [] },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' }
        }
      }
    })

    const input = wrapper.find('input[type="file"]')
    const firstFile = new File(
      ['{"type":"sub2api-data","version":1,"proxies":[{"name":"p1"}],"accounts":[{"name":"a1"}]}'],
      'first.json',
      {
        type: 'application/json'
      }
    )
    const secondFile = new File(
      ['{"type":"sub2api-data","version":1,"proxies":[{"name":"p2"}],"accounts":[{"name":"a2"},{"name":"a3"}]}'],
      'second.json',
      {
        type: 'application/json'
      }
    )

    Object.defineProperty(firstFile, 'text', {
      value: () =>
        Promise.resolve('{"type":"sub2api-data","version":1,"proxies":[{"name":"p1"}],"accounts":[{"name":"a1"}]}')
    })
    Object.defineProperty(secondFile, 'text', {
      value: () =>
        Promise.resolve(
          '{"type":"sub2api-data","version":1,"proxies":[{"name":"p2"}],"accounts":[{"name":"a2"},{"name":"a3"}]}'
        )
    })
    Object.defineProperty(input.element, 'files', {
      value: [firstFile, secondFile]
    })

    await input.trigger('change')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(input.attributes('multiple')).toBeDefined()
    expect(importData).toHaveBeenCalledTimes(2)
    expect(importData).toHaveBeenNthCalledWith(1, {
      data: {
        type: 'sub2api-data',
        version: 1,
        proxies: [{ name: 'p1' }],
        accounts: [{ name: 'a1' }]
      },
      group_ids: undefined,
      skip_default_group_bind: true
    })
    expect(importData).toHaveBeenNthCalledWith(2, {
      data: {
        type: 'sub2api-data',
        version: 1,
        proxies: [{ name: 'p2' }],
        accounts: [{ name: 'a2' }, { name: 'a3' }]
      },
      group_ids: undefined,
      skip_default_group_bind: true
    })
    expect(showSuccess).toHaveBeenCalledWith('admin.accounts.dataImportSuccess')
    expect(wrapper.emitted('imported')).toBeTruthy()
  })
})
