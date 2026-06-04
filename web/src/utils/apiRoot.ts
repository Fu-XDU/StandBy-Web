/** 从当前页面路径推断 API 根前缀（兼容反向代理子路径）。 */
export function resolveApiRoot(): string {
  const path = window.location.pathname.replace(/\/$/, '')
  const marker = '/v1/web'
  const idx = path.indexOf(marker)
  if (idx >= 0) {
    return path.slice(0, idx)
  }
  return ''
}

export function apiUrl(path: string): string {
  const root = resolveApiRoot()
  const normalized = path.startsWith('/') ? path : `/${path}`
  return `${root}${normalized}`
}
