import { describe, expect, it } from 'vitest'

import {
  formatRedeemCodesForExport,
  getRedeemCodeExportExtension
} from '@/utils/redeemCodeExport'

describe('redeemCodeExport', () => {
  it('formats codes one per line', () => {
    expect(formatRedeemCodesForExport(['A', 'B', 'C'], 'lines')).toBe('A\nB\nC')
  })

  it('formats codes as comma separated text', () => {
    expect(formatRedeemCodesForExport(['A', 'B', 'C'], 'comma')).toBe('A,B,C')
  })

  it('formats codes as a json array', () => {
    expect(formatRedeemCodesForExport(['A', 'B'], 'json')).toBe('[\n  "A",\n  "B"\n]')
  })

  it('returns the expected download extension', () => {
    expect(getRedeemCodeExportExtension('lines')).toBe('txt')
    expect(getRedeemCodeExportExtension('comma')).toBe('txt')
    expect(getRedeemCodeExportExtension('json')).toBe('json')
  })
})
