
export interface PgBloat {
  index_bloat: number
  index_name: string
  index_waste: number
  index_waste_pretty: string
  table_bloat: number
  table_name: string
  table_schema: string
  table_waste: number
  table_waste_pretty: string
}

export interface PgQuery {
  application_name: string
  client_addr: string
  pid: number
  query: string
  query_start: Date
  state: string
  usename: string
}

export interface PgSetting {
  boot_val: number
  changed: boolean
  context: Date
  extra_desc: string
  name: string
  setting: string
  short_desc: string
  unit: string
}

export interface PgTableSize {
  index_bytes: number
  index_pretty: string
  row_estimate: number
  table_bytes: number
  table_pretty: string
  table_name: string
  table_schema: string
  toast_bytes: number
  toast_pretty: string
  total_bytes: number
  total_pretty: string
  total_size_share: number
}

export interface PgUnusedIdx {
  index_name: string
  index_size: number
  index_size_pretty: string
  index_scans: number
  last_idx_scan: Date
  table_name: string
  table_schema: string
}
