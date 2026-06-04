import type { RemotePageId } from '@/types/remote'
import { FLOAT_PAGE_ID } from '@/config/floatConfig'

/** 已接入远程同步的页面 ID；新增钟表样式时在此扩展并注册对应 composable。 */
export const REMOTE_PAGE_IDS = [FLOAT_PAGE_ID] as const satisfies readonly RemotePageId[]
