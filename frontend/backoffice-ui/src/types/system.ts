
export interface ApiError {
  method: string
  url: string
  errMsg: string
}

export interface GetMetadata {
  count: number
  total_count: number
  total_count_is_estimated: boolean
}

export interface SelectionItem {
  id: number
  name: string
}
