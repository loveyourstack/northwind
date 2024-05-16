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
    class="pa-4 rounded"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center" class="pb-2">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Employees' }}</div>
          </div>
        </v-col>

        <v-col>
          <v-btn class="float-end" color="primary" :to="{ name: 'New employee'}">Add</v-btn>

          <v-menu :close-on-content-click=false>
            <template v-slot:activator="{ props }">
              <v-btn density="comfortable" flat class="float-right mr-7" icon="mdi-table-column" v-bind="props"></v-btn>
            </template>
            <v-list>
              <v-list-item v-for="(header, i) in headers" :key="i" :value="header" @click="toggleHeader(header.key)">
                <template v-slot:append>
                  <v-icon :icon="getHeaderListIcon(excludedHeaders, header.key)" :color="getHeaderListIconColor(excludedHeaders, header.key)"></v-icon>
                </template>
                <v-list-item-title class="clickable" v-text="header.title"></v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>

          <v-switch hide-details v-model="showFilters" class="float-right mr-7" :color="showFilters ? 'teal' : 'undefined'" density="compact"
          ></v-switch>
          <v-icon icon="mdi-filter-variant" class="float-right mt-2 mr-3"></v-icon>
        </v-col>
      </v-row>

      <v-row v-show="showFilters" class="mt-0">
        <v-col cols="12" sm="6" lg="4">
          <v-text-field label="First name search" v-model="filterFirstName" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="4">
          <v-text-field label="Last name search" v-model="filterLastName" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
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

    <template v-if="totalItemsIsEstimate" v-slot:[`bottom`]="{}">
      <v-data-table-footer
        :items-per-page-options="itemsPerPageOptions"
        :page-text="getPageTextEstimated(totalItemsEstimated)"
        :show-current-page=true
      ></v-data-table-footer>
    </template>

  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount, onMounted } from 'vue'
import { VDataTable } from 'vuetify/components'
import ax from '@/api'
import { Employee } from '@/types/hr'
import { debounceMs, maxDebounceMs, getHeaderListIcon, getHeaderListIconColor, getPageTextEstimated, itemsPerPageOptions, processURIOptions } from '@/composables/datatable'
import { useCommonStore } from '@/stores/common'
import { useDateFormat, useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  title?: string
}>()

const commonStore = useCommonStore()

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

const items = ref<Employee[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const showFilters = ref(false)
const filterFirstName = ref<string>()
const filterLastName = ref<string>()
const lsKey = 'employees_dt'

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {

  var myURL = '/a/hr/employees'
  myURL = processURIOptions(myURL, options)

  // filters
  if (filterFirstName.value) {
    myURL += '&first_name=~' + filterFirstName.value + '~'
  }
  if (filterLastName.value) {
    myURL += '&last_name=~' + filterLastName.value + '~'
  }

  ax.get(myURL)
  .then(resp => {
      items.value = resp.data.data
      totalItemsIsEstimate.value = resp.headers['x-total-count-estimated'] === 'true' ? true : false
      if (totalItemsIsEstimate.value) {
        totalItemsEstimated.value = Number(resp.headers['x-total-count'])
        totalItems.value = 101 // workaround so that next page button is enabled even if estimate is too low
      } else {
        totalItems.value = Number(resp.headers['x-total-count'])
      }
    })
    .catch() // handled by interceptor
}

const debouncedRefreshItems = useDebounceFn(() => {
  search.value = String(Date.now())
}, debounceMs, { maxWait: maxDebounceMs })

function refreshItems() {
  search.value = String(Date.now())
}

function toggleHeader(key: string) {
  excludedHeaders.value.includes(key) ? excludedHeaders.value = excludedHeaders.value.filter((v) => v != key) : excludedHeaders.value.push(key)
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
}

watch([itemsPerPage, showFilters, search, sortBy, excludedHeaders], () => {

  if (!showFilters.value) {
    filterFirstName.value = undefined
    filterLastName.value = undefined
    refreshItems()
  }

  let lsObj = {
    'itemsPerPage': itemsPerPage.value,
    'sortBy': sortBy.value,
    'showFilters': showFilters.value,
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
  if (lsObj['showFilters']) { showFilters.value = lsObj['showFilters'] }
  if (lsObj['filterFirstName']) { filterFirstName.value = lsObj['filterFirstName'] }
  if (lsObj['filterLastName']) { filterLastName.value = lsObj['filterLastName'] }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
