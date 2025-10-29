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
            <div class="dt-title">{{ props.title ? props.title : 'Customers' }}</div>
          </div>

          <v-btn class="float-end" color="primary" :to="{ name: 'New customer'}">Add</v-btn>

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

            <FilterChipText name="Code" :filterValue="filterCode" 
              @closed="filterCode = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterCode = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChipText name="Company name" :filterValue="filterCompanyName" 
              @closed="filterCompanyName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterCompanyName = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChipText name="Contact name" :filterValue="filterContactName" 
              @closed="filterContactName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterContactName = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChip name="Country" :filterValue="filterCountryID" :filterText="filterCountryIDText" @closed="filterCountryID = undefined; refreshItems()">
              <template #menuContent>
                <v-autocomplete label="Country" v-model="filterCountryID" autofocus
                  :items="coreStore.activeCountriesList" item-title="name" item-value="id"
                  @update:model-value="refreshItems"
                ></v-autocomplete>
              </template>
            </FilterChip>

          </v-chip-group>
        </v-col>
       </v-row>
    </template>

    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" :to="{ name: 'Customer detail', params: { id: item.id }}">
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
import { type Customer } from '@/types/sales'
import { getTextFilterUrlParam, itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'
import { useCoreStore } from '@/stores/core'
import { useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  title?: string
}>()

const coreStore = useCoreStore()

var headers = [
  { title: 'Code', key: 'code' },
  { title: 'Company', key: 'company_name' },
  { title: 'Contact', key: 'contact_name' },
  { title: 'Job title', key: 'contact_title' },
  { title: 'Phone', key: 'phone' },
  { title: 'Country', key: 'country' },
  { title: '# orders', key: 'order_count', align: 'end' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const
const excludedHeaders = ref<string[]>([])
const selectedHeaders = ref()

const baseUrl = '/a/sales/customers'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<Customer[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const filterCode = ref<string>()
const filterCompanyName = ref<string>()
const filterContactName = ref<string>()

const filterCountryID = ref<number>()
const filterCountryIDText = computed(() => {
  return filterCountryID.value ? coreStore.activeCountriesList.find(ele => ele.id === filterCountryID.value)?.name : ''
})

const lsKey = 'customers_dt'

function getFilterStr(): string {
  var ret = ''

  ret += getTextFilterUrlParam('code', filterCode.value)
  ret += getTextFilterUrlParam('company_name', filterCompanyName.value)
  ret += getTextFilterUrlParam('contact_name', filterContactName.value)

  if (filterCountryID.value) {
    ret += '&country_fk=' + filterCountryID.value
  }

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
    'filterCode': filterCode.value,
    'filterCompanyName': filterCompanyName.value,
    'filterContactName': filterContactName.value,
    'filterCountryID': filterCountryID.value,
    'excludedHeaders': excludedHeaders.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
}, { deep: true }) // needed for excludedHeaders array

onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['itemsPerPage']) { itemsPerPage.value = lsObj['itemsPerPage'] }
  if (lsObj['sortBy']) { sortBy.value = lsObj['sortBy'] }
  if (lsObj['filterCode']) { filterCode.value = lsObj['filterCode'] }
  if (lsObj['filterCompanyName']) { filterCompanyName.value = lsObj['filterCompanyName'] }
  if (lsObj['filterContactName']) { filterContactName.value = lsObj['filterContactName'] }
  if (lsObj['filterCountryID']) { filterCountryID.value = lsObj['filterCountryID'] }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
