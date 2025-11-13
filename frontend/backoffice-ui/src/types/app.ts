
export interface NumericFilter {
  operator: string
  value: number
  value_upper: number // treated as upper limit when operator is 'between'
}