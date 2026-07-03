import {useFloatConfig} from '@/composables/useFloatConfig'

export interface ClockStyleOption {
  id: string
  name: string
}

export const CLOCK_STYLES: ClockStyleOption[] = [
  {id: 'float', name: '浮动'},
  {id: 'numerical', name: '数字'},
]

export function useClockStyle() {
  const {clockStyle} = useFloatConfig()
  return {clockStyle, clockStyles: CLOCK_STYLES}
}
