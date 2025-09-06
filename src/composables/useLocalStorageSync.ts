import {type Ref, watch} from "vue"

export function useLocalStorageSync<T>(
    key: string,
    target: Ref<T>,
    options?: {
        parser?: (val: any) => T
        serializer?: (val: T) => any
        deep?: boolean
    }
): void

export function useLocalStorageSync(
    vars: Record<string, Ref<any>>,
    options?: {
        parser?: (val: any, key: string) => any
        serializer?: (val: any, key: string) => any
        deep?: boolean
    }
): void

export function useLocalStorageSync(
    keyOrVars: string | Record<string, Ref<any>>,
    targetOrOptions?: Ref<any> | any,
    maybeOptions: any = {}
) {
    if (typeof keyOrVars === "string") {
        // 单个变量模式
        const key = keyOrVars
        const target = targetOrOptions as Ref<any>
        const { parser, serializer, deep } = maybeOptions

        const raw = localStorage.getItem(key)
        if (raw !== null) {
            try {
                const parsed = JSON.parse(raw)
                target.value = parser ? parser(parsed, key) : parsed
            } catch (e) {
                console.warn(`Failed to parse localStorage for ${key}`, e)
            }
        }

        watch(
            target,
            (newVal) => {
                try {
                    const val = serializer ? serializer(newVal, key) : newVal
                    localStorage.setItem(key, JSON.stringify(val))
                } catch (e) {
                    console.warn(`Failed to stringify localStorage for ${key}`, e)
                }
            },
            { deep }
        )
    } else {
        // 批量模式
        const vars = keyOrVars as Record<string, Ref<any>>
        const { parser, serializer, deep } = targetOrOptions || {}

        for (const [key, ref] of Object.entries(vars)) {
            useLocalStorageSync(key, ref, {
                parser: parser ? (val: any) => parser(val, key) : undefined,
                serializer: serializer ? (val: any) => serializer(val, key) : undefined,
                deep,
            })
        }
    }
}