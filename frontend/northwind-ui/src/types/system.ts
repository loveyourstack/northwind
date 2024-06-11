
export interface ApiError {
  method: string
  url: string
  errMsg: string
}

export interface GetMetadata {
  count: number
  total_count: number
  total_count_is_estimated: boolean
  last_sync_at: Date
}

export interface SelectionItem {
  id: number
  name: string
}
