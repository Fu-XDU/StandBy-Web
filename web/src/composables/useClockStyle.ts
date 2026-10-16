import { markRaw } from 'vue'
import { useFloatConfig } from '@/composables/useFloatConfig'
import Float, { floatColors, floatNightModeColors } from '@/components/Float.vue'
import Numerical, { numericalColors, numericalNightModeColors } from '@/components/Numerical.vue'
import Analog, { analogColors, analogNightModeColors } from '@/components/Analog.vue'

export interface ClockStyleOption {
  id: string
  name: string
}

export interface ClockStyleDef {
  id: string
  component: ReturnType<typeof markRaw>
  colors: string[][]
  nightModeColors: string[]
}

export interface ClockTypeDef {
  id: string
  styles: ClockStyleDef[]
}

export const CLOCK_TYPES: ClockTypeDef[] = [
  {
    id: 'digital',
    styles: [
      { id: 'float', component: markRaw(Float), colors: floatColors, nightModeColors: floatNightModeColors },
      { id: 'numerical', component: markRaw(Numerical), colors: numericalColors, nightModeColors: numericalNightModeColors },
      { id: 'analog', component: markRaw(Analog), colors: analogColors, nightModeColors: analogNightModeColors },
    ],
  },
]

export const CLOCK_STYLES: ClockStyleOption[] = [
  { id: 'float', name: '浮动' },
  { id: 'numerical', name: '数字' },
  { id: 'analog', name: '指针' },
]

export function useClockStyle() {
  const { clockStyle } = useFloatConfig()
  return { clockStyle, clockStyles: CLOCK_STYLES }
}
