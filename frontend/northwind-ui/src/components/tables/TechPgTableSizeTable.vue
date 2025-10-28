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
    @update:options="loadItems"
    class="dt"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Table size' }}</div>
          </div>

          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn class="float-right mr-5" icon="mdi-dots-vertical" v-bind="props"></v-btn>
            </template>
            <v-list>

              <v-list-item prepend-icon="mdi-selection-ellipse-remove">
                <v-list-item-title class="clickable" @click="resetTable()">Reset table</v-list-item-title>
              </v-list-item>

            </v-list>
          </v-menu>

          <span class="float-right mr-4 mt-3 text-body-1">
            Database size: {{ dbSize }}
          </span>

        </v-col>
      </v-row>

      <v-row no-gutters class="mb-3">
        <v-col>
          <v-chip-group column>

            <FilterChipText name="Table search" :filterValue="filterTable" 
              @closed="filterTable = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterTable = val; debouncedRefreshItems() }">
            </FilterChipText>

          </v-chip-group>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.row_estimate`]="{ item }">
      {{ Intl.NumberFormat('en-US').format(item.row_estimate) }}
    </template>

    <template v-slot:[`item.index_bytes`]="{ item }">
      {{ item.index_pretty }} <!-- show pretty, but keep bytes in headers so that sorting works -->
    </template>

    <template v-slot:[`item.table_bytes`]="{ item }">
      {{ item.table_pretty }}
    </template>

    <template v-slot:[`item.toast_bytes`]="{ item }">
      {{ item.toast_pretty }}
    </template>

    <template v-slot:[`item.total_bytes`]="{ item }">
      {{ item.total_pretty }}
    </template>

    <template v-slot:[`item.total_size_share`]="{ item }">
      {{ Intl.NumberFormat('en-US', { style: 'percent', maximumFractionDigits: 2, minimumFractionDigits: 2}).format(item.total_size_share) }}
    </template>

    <template v-slot:[`bottom`]="{}">
      <DtFooter :totalItemsIsEstimate="totalItemsIsEstimate" :totalItemsEstimated="totalItemsEstimated"></DtFooter>
    </template>

  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount, onMounted } from 'vue'
import { VDataTable } from 'vuetify/components'
import { type PgTableSize } from '@/types/tech'
import { useFetch, useFetchDt } from '@/composables/fetch'
import { getTextFilterUrlParam, processURIOptions } from '@/functions/datatable'
import { useDebounceFn } from '@vueuse/core'
import DtFooter from '@/components/DtFooter.vue'
import FilterChipText from '@/components/FilterChipText.vue'

const props = defineProps<{
  title?: string
}>()

var headers = [
  { title: 'Schema', key: 'table_schema' },
  { title: 'Table', key: 'table_name' },
  { title: 'est. # rows', key: 'row_estimate', align: 'end' },
  { title: 'Total', key: 'total_bytes', align: 'end' },
  { title: '%', key: 'total_size_share', align: 'end' },
  { title: 'Table', key: 'table_bytes', align: 'end' },
  { title: 'Indexes', key: 'index_bytes', align: 'end' },
  { title: 'Toast', key: 'toast_bytes', align: 'end' },
] as const

const baseUrl = '/a/tech/pg-table-size'

const items = ref<PgTableSize[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const dbSize = ref<string>('')

const filterTable = ref<string>()

const lsKey = 'pg_tablesize_dt'

function getFilterStr(): string {
  var ret = ''
  ret += getTextFilterUrlParam('table_name', filterTable.value)
  return ret
}

function loadDbSize() {
  useFetch('/a/tech/pg-database-size', dbSize)
}

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {
  var myURL = processURIOptions(baseUrl, options)
  myURL += getFilterStr()
  useFetchDt(myURL, items, totalItems, totalItemsIsEstimate, totalItemsEstimated)
}

const debouncedRefreshItems = useDebounceFn(() => {
  search.value = String(Date.now())
}, import.meta.env.VITE_DEBOUNCE_MS, { maxWait: import.meta.env.VITE_MAX_DEBOUNCE_MS })

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
    'filterTable': filterTable.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
})

onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['itemsPerPage']) { itemsPerPage.value = lsObj['itemsPerPage'] }
  if (lsObj['sortBy']) { sortBy.value = lsObj['sortBy'] }
  if (lsObj['filterTable']) { filterTable.value = lsObj['filterTable'] }
})

onMounted(() => {
  loadDbSize()
})

</script>
