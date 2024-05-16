import { VDataTable } from 'vuetify/components'

export const itemsPerPageOptions = [
  {value: 5, title: '5'},
  {value: 10, title: '10'},
  {value: 25, title: '25'},
  {value: 50, title: '50'},
  {value: 100, title: '100'}
]

// milliseconds to delay input before triggering data refresh
export const debounceMs = 400
export const maxDebounceMs = 5000

export function getHeaderListIcon(excludedHeaders: string[], headerKey: string) {
  return excludedHeaders.includes(headerKey) ? 'mdi-close' : 'mdi-check'
}
export function getHeaderListIconColor(excludedHeaders: string[], headerKey: string) {
  return excludedHeaders.includes(headerKey) ? 'error' : 'success'
}

export function getPageTextEstimated(totalItemsEstimated: number) {
  return "~" + totalItemsEstimated.toLocaleString() + " items"
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

    if (options.sortBy[i].order === 'desc') {
      sortParam += '-' + options.sortBy[i].key
    } else {
      sortParam += options.sortBy[i].key
    }
  }
  uri += sortParam

  return uri
}
