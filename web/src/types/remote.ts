export type RemotePageId = 'float'

export type RemoteSyncAction = 'apply_server' | 'store_client'

export interface RemoteSyncResponse<TConfig = unknown> {
  action: RemoteSyncAction
  updatedAt: number
  config?: TConfig
}

export interface RemoteSyncRequest<TConfig = unknown> {
  deviceId: string
  pageId: RemotePageId
  updatedAt: number
  config: TConfig
}
