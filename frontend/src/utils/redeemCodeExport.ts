export type RedeemCodeExportFormat = 'lines' | 'comma' | 'json'

export function formatRedeemCodesForExport(
  codes: string[],
  format: RedeemCodeExportFormat
): string {
  switch (format) {
    case 'comma':
      return codes.join(',')
    case 'json':
      return JSON.stringify(codes, null, 2)
    case 'lines':
    default:
      return codes.join('\n')
  }
}

export function getRedeemCodeExportExtension(format: RedeemCodeExportFormat): 'txt' | 'json' {
  return format === 'json' ? 'json' : 'txt'
}
