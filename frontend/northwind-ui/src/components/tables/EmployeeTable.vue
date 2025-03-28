<template>
  <v-data-table-server
    v-model:items-per-page="itemsPerPage"
    v-model:sortBy="sortBy"
    :headers="selectedHeaders"
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
            <div class="dt-title">{{ props.title ? props.title : 'Employees' }}</div>
          </div>

          <v-btn class="float-end" color="primary" :to="{ name: 'New employee'}">Add</v-btn>

          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn class="float-right mr-5" icon="mdi-dots-vertical" v-bind="props"></v-btn>
            </template>
            <v-list>

              <v-list-item prepend-icon="mdi-selection-ellipse-remove">
                <v-list-item-title class="clickable" @click="resetTable()">Reset table</v-list-item-title>
              </v-list-item>

              <AdjustColsListItem :headers="headers" :excluded-headers="excludedHeaders" @toggle="(key: string) => toggleHeader(key)"></AdjustColsListItem>

              <v-list-item prepend-icon="mdi-download-outline">
                <v-list-item-title class="clickable" @click="fileDownload(excelDlUrl)">Download to Excel</v-list-item-title>
              </v-list-item>

            </v-list>
          </v-menu>

        </v-col>
      </v-row>

      <v-row no-gutters class="mb-1">
        <v-col>
          <v-chip-group column>

            <FilterChipText name="First name" :filterValue="filterFirstName" 
              @closed="filterFirstName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterFirstName = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChipText name="Last name" :filterValue="filterLastName" 
              @closed="filterLastName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterLastName = val; debouncedRefreshItems() }">
            </FilterChipText>

          </v-chip-group>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.hire_date`]="{ item }">
      {{ useDateFormat(item.hire_date, 'DD MMM YYYY').value }}
    </template>

    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" :to="{ name: 'Employee detail', params: { id: item.id }}">
        <v-icon color="primary" icon="mdi-details"></v-icon>
      </v-btn>
    </template>

    <template v-slot:[`bottom`]="{}">
      <DtFooter :totalItemsIsEstimate="totalItemsIsEstimate" :totalItemsEstimated="totalItemsEstimated"></DtFooter>
    </template>

  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onBeforeMount, onMounted } from 'vue'
import { VDataTable } from 'vuetify/components'
import { useFetchDt } from '@/composables/fetch'
import { Employee } from '@/types/hr'
import { getTextFilterUrlParam, itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'
import { useDateFormat, useDebounceFn } from '@vueuse/core'
import AdjustColsListItem from '@/components/AdjustColsListItem.vue'
import DtFooter from '@/components/DtFooter.vue'
import FilterChipText from '@/components/FilterChipText.vue'

const props = defineProps<{
  title?: string
}>()

var headers = [
  { title: 'Title', key: 'title' },  
  { title: 'First name', key: 'first_name' },
  { title: 'Last name', key: 'last_name' },
  { title: 'Job title', key: 'job_title' },
  { title: 'Home phone', key: 'home_phone' },
  { title: 'Reports to', key: 'reports_to' },
  { title: 'Hire date', key: 'hire_date' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const
const excludedHeaders = ref<string[]>([])
const selectedHeaders = ref()

const baseUrl = '/a/hr/employees'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<Employee[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const filterFirstName = ref<string>()
const filterLastName = ref<string>()
const lsKey = 'employees_dt'

function getFilterStr(): string {
  var ret = ''

  // exclude None
  ret += '&id=!-1'

  ret += getTextFilterUrlParam('first_name', filterFirstName.value)
  ret += getTextFilterUrlParam('last_name', filterLastName.value)

  return ret
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

function toggleHeader(key: string) {
  excludedHeaders.value.includes(key) ? excludedHeaders.value = excludedHeaders.value.filter((v) => v != key) : excludedHeaders.value.push(key)
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
}

watch([itemsPerPage, search, sortBy, excludedHeaders], () => {

  let lsObj = {
    'itemsPerPage': itemsPerPage.value,
    'sortBy': sortBy.value,
    'filterFirstName': filterFirstName.value,
    'filterLastName': filterLastName.value,
    'excludedHeaders': excludedHeaders.value,
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
  if (lsObj['filterFirstName']) { filterFirstName.value = lsObj['filterFirstName'] }
  if (lsObj['filterLastName']) { filterLastName.value = lsObj['filterLastName'] }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
