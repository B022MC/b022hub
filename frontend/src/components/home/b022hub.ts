export type B022HubIconName = 'cpu' | 'sparkles' | 'swap' | 'terminal'
export type B022HubVendorId = 'claude' | 'openai' | 'gemini'

export interface B022HubInstallOption {
  key: string
  label: string
  hint: string
  command: string
}

export interface B022HubConfigFile {
  path: string
  languageLabel: string
  content: string
}

export interface B022HubFeatureCard {
  icon: B022HubIconName
  title: string
  description: string
  status: 'completed' | 'building'
}
