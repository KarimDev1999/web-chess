export interface TimeControlPreset {
  label: string
  base: number
  increment: number
}

export interface TimeCategory {
  id: string
  label: string
  icon: string
  presets: TimeControlPreset[]
}

export const TIME_CATEGORIES: TimeCategory[] = [
  {
    id: 'bullet',
    label: 'Bullet',
    icon: '🏃',
    presets: [
      { label: '1+0', base: 60, increment: 0 },
      { label: '1|1', base: 60, increment: 1 },
      { label: '2|1', base: 120, increment: 1 },
    ],
  },
  {
    id: 'blitz',
    label: 'Blitz',
    icon: '⚡',
    presets: [
      { label: '3+0', base: 180, increment: 0 },
      { label: '3|2', base: 180, increment: 2 },
      { label: '5+0', base: 300, increment: 0 },
      { label: '5|3', base: 300, increment: 3 },
      { label: '10+0', base: 600, increment: 0 },
    ],
  },
  {
    id: 'rapid',
    label: 'Rapid',
    icon: '🕐',
    presets: [
      { label: '10+0', base: 600, increment: 0 },
      { label: '10|5', base: 600, increment: 5 },
      { label: '15|10', base: 900, increment: 10 },
      { label: '25+0', base: 1500, increment: 0 },
    ],
  },
  {
    id: 'classical',
    label: 'Classical',
    icon: '♟️',
    presets: [
      { label: '30+0', base: 1800, increment: 0 },
      { label: '30|20', base: 1800, increment: 20 },
      { label: '60|30', base: 3600, increment: 30 },
    ],
  },
]

const UNTIMED_CATEGORY: TimeCategory = {
  id: 'unlimited',
  label: 'Unlimited',
  icon: '♾️',
  presets: [{ label: 'Unlimited', base: 0, increment: 0 }],
}

export const ALL_CATEGORIES = [UNTIMED_CATEGORY, ...TIME_CATEGORIES]
