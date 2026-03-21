export type HomeIconName =
  | 'arrowRight'
  | 'badge'
  | 'bolt'
  | 'book'
  | 'brain'
  | 'chart'
  | 'checkCircle'
  | 'cloud'
  | 'cpu'
  | 'creditCard'
  | 'cube'
  | 'database'
  | 'dollar'
  | 'fire'
  | 'globe'
  | 'key'
  | 'server'
  | 'shield'
  | 'sparkles'
  | 'swap'
  | 'terminal'
  | 'users'

export type HomeTone = 'primary' | 'sky' | 'amber' | 'rose' | 'violet' | 'slate'

export interface HomeTagItem {
  icon: HomeIconName
  label: string
}

export interface HomeInfoCard {
  icon: HomeIconName
  title: string
  description: string
  tone: HomeTone
}

export interface HomeStepItem extends HomeInfoCard {
  index: string
}

export interface HomeScenarioItem extends HomeInfoCard {
  eyebrow: string
  points: string[]
}

export interface HomeProviderItem {
  name: string
  shortLabel: string
  detail: string
  status: 'supported' | 'soon'
  tone: HomeTone
}

export interface HomeComparisonRow {
  feature: string
  official: string
  platform: string
}
