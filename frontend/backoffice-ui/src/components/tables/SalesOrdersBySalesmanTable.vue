<template>

  <v-data-table-server
    v-model:items-per-page="itemsPerPage"
    v-model:sortBy="sortBy"
    :headers="headers"
    :hover=true
    :items-length="totalItems"
    :items="items"
    :multi-sort=true
    :search="search"
    :show-current-page=true
    item-value="id"
    :items-per-page-options="itemsPerPageOptions"
    @update:options="loadItems"
    class="dt"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Orders by salesman' }}</div>
          </div>

          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn class="float-right mr-5" icon="mdi-dots-vertical" v-bind="props"></v-btn>
            </template>
            <v-list>

              <v-list-item prepend-icon="mdi-selection-ellipse-remove">
                <v-list-item-title class="clickable" @click="resetTable()">Reset table</v-list-item-title>
              </v-list-item>

              <v-list-item prepend-icon="mdi-download-outline">
                <v-list-item-title class="clickable" @click="fileDownload(excelDlUrl)">Download to Excel</v-list-item-title>
              </v-list-item>

            </v-list>
          </v-menu>

        </v-col>
      </v-row>

      <v-row no-gutters class="mb-1">
        <v-col class="ml-2 mr-2">
          <DateTextField :dateVal="filterFromDate" label="From date"
            @updated="(val: string) => { filterFromDate = val; refreshItems() } "
          ></DateTextField>
        </v-col>

        <v-col class="ml-2 mr-2">
          <DateTextField :dateVal="filterUntilDate" label="Until date"
            @updated="(val: string) => { filterUntilDate = val; refreshItems() } "
          ></DateTextField>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.total_value`]="{ item }">
      <span>{{ Intl.NumberFormat('en-US', { style: "currency", currency: "USD", maximumFractionDigits: 0 }).format(item.total_value) }}</span>
    </template>

    <template v-slot:[`bottom`]="{}">
    </template>
    
  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import { useFetchDt } from '@/composables/fetch'
import { type OrdersBySalesman } from '@/types/sales'
import { itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'

const props = defineProps<{
  title?: string
}>()

var headers = [
  { title: 'Salesman', key: 'salesman' },
  { title: '# orders', key: 'order_count', align: 'end' },
  { title: '# orders shipped', key: 'shipped_count', align: 'end' },
  { title: 'Total value', key: 'total_value', align: 'end' },
] as const

const baseUrl = '/a/sales/orders-by-salesman'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<OrdersBySalesman[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const filterFromDate = ref('2021-04-20') // YYYY-MM-DD

const filterUntilDate = ref('2021-05-06') // YYYY-MM-DD

const lsKey = 'orders_by_salesman_dt'

function getFilterStr(): string {
  var ret = ''

  if (!filterFromDate.value) { return '' }
  ret += '&from_date=' + filterFromDate.value

  if (!filterUntilDate.value) { return '' }
  ret += '&until_date=' + filterUntilDate.value

  return ret
}

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {
  var myURL = processURIOptions(baseUrl, options)

  // from and until date filters are mandatory
  var filterStr = getFilterStr()
  if (!filterStr) { return }

  myURL += filterStr
  useFetchDt(myURL, items, totalItems, totalItemsIsEstimate, totalItemsEstimated)
}

function refreshItems() {
  search.value = String(Date.now())
}

function resetTable() {
  localStorage.removeItem(lsKey)
  window.location.reload()
}

watch([itemsPerPage, search, sortBy], () => {

  let lsObj = {
    'itemsPerPage': itemsPerPage.value,
    'sortBy': sortBy.value,
    'filterFromDate': filterFromDate.value,
    'filterUntilDate': filterUntilDate.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
}, { deep: true })

onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['itemsPerPage']) { itemsPerPage.value = lsObj['itemsPerPage'] }
  if (lsObj['sortBy']) { sortBy.value = lsObj['sortBy'] }
  if (lsObj['filterFromDate']) { filterFromDate.value = lsObj['filterFromDate'] }
  if (lsObj['filterUntilDate']) { filterUntilDate.value = lsObj['filterUntilDate'] }
})

</script>
