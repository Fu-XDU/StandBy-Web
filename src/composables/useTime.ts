import {ref, onMounted, onUnmounted} from 'vue'

// 定义一个返回当前小时和分钟四个数字的组合逻辑
export function useTime() {
    const digits = ref<string[]>([])

    function updateDigits() {
        const now = new Date()
        const hours = now.getHours().toString().padStart(2, '0')  // "08"
        const minutes = now.getMinutes().toString().padStart(2, '0') // "05"
        digits.value = [...hours, ...minutes]  // ['0', '8', '0', '5']
    }

    let timer: number
    let recalibrationTimer: number

    const startAlignedTimer = () => {
        updateDigits()
        const now = new Date()
        const msUntilNextMinute = (60 - now.getSeconds()) * 1000 - now.getMilliseconds()

        setTimeout(() => {
            updateDigits()
            timer = window.setInterval(updateDigits, 60 * 1000)
        }, msUntilNextMinute)
    }

    const recalibrate = () => {
        if (timer) clearInterval(timer)
        startAlignedTimer()
    }

    onMounted(() => {
        startAlignedTimer()

        // 每 60 分钟重新校准一次定时器
        recalibrationTimer = window.setInterval(() => {
            recalibrate()
        }, 60 * 60 * 1000)
    })

    onUnmounted(() => {
        if (timer) clearInterval(timer)
        if (recalibrationTimer) clearInterval(recalibrationTimer)
    })

    return {digits}
}
