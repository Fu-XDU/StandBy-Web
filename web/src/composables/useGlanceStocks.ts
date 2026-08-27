export const GLANCE_API_URL = 'https://dev.flxdu.cn/glance/api/menu'

export interface GlanceMenuItem {
  title: string
  action?: string
  value?: string
  status_title?: string
  children?: GlanceMenuItem[]
}

export interface GlanceMenuResponse {
  title?: string
  refresh_after_seconds?: number
  menu?: GlanceMenuItem[]
}

export interface StockItem {
  value: string
  name: string
  price: string
  fullTitle: string
}

/**
 * 递归从 glance menu 结构中提取出所有可选项
 */
export function extractStockItems(menu?: GlanceMenuItem[]): StockItem[] {
  if (!menu || !Array.isArray(menu)) return []

  const result: StockItem[] = []

  function traverse(items: GlanceMenuItem[]) {
    for (const item of items) {
      if (item.action === 'select' && item.value) {
        let price = item.status_title?.trim() ?? ''
        let name = ''

        if (price && item.title.includes(price)) {
          name = item.title.replace(price, '').trim()
        } else {
          // 备用分割逻辑：通过空格拆分出前缀名称与末尾价格
          const parts = item.title.trim().split(/\s{2,}|\t+/)
          if (parts.length >= 2) {
            name = parts[0].trim()
            if (!price) price = parts.slice(1).join(' ').trim()
          } else {
            name = item.title.trim()
          }
        }

        result.push({
          value: item.value,
          name: name || item.value,
          price: price || '--',
          fullTitle: item.title,
        })
      }

      if (item.children && Array.isArray(item.children)) {
        traverse(item.children)
      }
    }
  }

  traverse(menu)
  return result
}

/**
 * 统一的网络请求方法，带超时保护
 */
export async function fetchGlanceMenu(timeoutMs = 8000): Promise<GlanceMenuResponse> {
  const controller = new AbortController()
  const timer = setTimeout(() => controller.abort(), timeoutMs)

  try {
    const res = await fetch(GLANCE_API_URL, {
      signal: controller.signal,
      headers: {
        Accept: 'application/json',
      },
    })
    if (!res.ok) {
      throw new Error(`HTTP ${res.status} ${res.statusText}`)
    }
    return (await res.json()) as GlanceMenuResponse
  } finally {
    clearTimeout(timer)
  }
}
