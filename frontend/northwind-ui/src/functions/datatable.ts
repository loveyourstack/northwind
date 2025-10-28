import { VDataTable } from 'vuetify/components'
import { type NumericFilter } from '@/types/core'

export const itemsPerPageOptions = [
  {value: 5, title: '5'},
  {value: 10, title: '10'},
  {value: 25, title: '25'},
  {value: 50, title: '50'},
  {value: 100, title: '100'}
]

export function getHeaderListIcon(excludedHeaders: string[], headerKey: string) {
  return excludedHeaders.includes(headerKey) ? 'mdi-close' : 'mdi-check'
}
export function getHeaderListIconColor(excludedHeaders: string[], headerKey: string) {
  return excludedHeaders.includes(headerKey) ? 'error' : 'success'
}

export function getNumericFilterDisplayText(nf: NumericFilter, isPercent?: boolean): string {

  if (!nf.operator) { 
    return ''
  }

  var val: number = nf.value ? nf.value : 0 // ensure that 0 rather than '' is shown when user deletes value in textbox
  var val_upper: number = nf.value_upper ? nf.value_upper : 0

  var val_s: string = val.toString()
  if (isPercent) {
    val_s += '%'
  }
  var val_upper_s: string = val_upper.toString()
  if (isPercent) {
    val_upper_s += '%'
  }

  if (nf.operator === '<=>') {
    if (val_upper < val) {
      return 'invalid'
    }

    return nf.operator + ' ' + val_s + ' and ' + val_upper_s
  }

  return nf.operator + ' ' + val_s
}

export function getNumericFilterUrlParams(param: string, nf: NumericFilter, isPercent?: boolean): string {

  if (!nf.operator) { 
    return ''
  }

  var val: number = nf.value ? nf.value : 0 // ensure that 0 rather than '' is passed to API when user deletes value in textbox
  var val_upper: number = nf.value_upper ? nf.value_upper : 0

  if (isPercent) {
    val /= 100
    val_upper /= 100
  }

  if (nf.operator === '<=>') {

    if (val_upper < val) {
      return ''
    }

    return '&' + param + '=>eq' + val + '&' + param + '=<eq' + val_upper
  }
  
  return '&' + param + '=' + getOperatorParam(nf.operator) + val
}


export function getOperatorParam(operator: string): string {
  switch (operator) {
    case '<':
      return operator
    case '<=':
      return '<eq'
    case '=':
      return '' // empty so that param becomes x=val
    case '>=':
      return '>eq'
    case '>':
      return operator
    case '!=':
      return '!'
    default:
      return 'unknown'
  }
}

export function getPageTextEstimated(totalItemsEstimated: number) {
  return "~" + totalItemsEstimated.toLocaleString() + " items"
}

export function getTextFilterUrlParam(param: string, filterStr: string | undefined): string {

  if (!filterStr) { 
    return ''
  }

  // if filterStr starts with !
  if (filterStr.startsWith('!')) {

    // return a "not contains" filter
    return '&' + param + '=!~' + filterStr.substring(1) + '~'
  }

  // return a "contains any" filter, which allows use of "|" for OR conditions
  return '&' + param + '=~[' + filterStr + ']~'
}

// processURIOptions translates the datatable pagination and sorting options into URL params for the API
// uri = the API uri before pagination and sorting are added
// options = the options object from v-data-table-server
// NB: default sorting is handled by API and not needed here
export function processURIOptions (uri: string, options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {

  if (options.page > 0) {
    uri += '?xpage=' + options.page
  } else {
    uri += '?xpage=1'
  }

  if (options.itemsPerPage > 0) {
    uri += '&xper_page=' + options.itemsPerPage
  } else {
    uri += '&xper_page=10'
  }

  var sortParam = ''
  for (var i = 0; i < options.sortBy.length; i++) {

    if (i === 0) {
      sortParam = '&xsort='
    } else {
      sortParam += ','
    }

    if (options.sortBy[i]!.order === 'desc') {
      sortParam += '-' + options.sortBy[i]!.key
    } else {
      sortParam += options.sortBy[i]!.key
    }
  }
  uri += sortParam

  return uri
}
