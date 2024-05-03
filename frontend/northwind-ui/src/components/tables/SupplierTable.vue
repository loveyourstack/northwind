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
    class="pa-4"
  >
    <template v-if="totalItemsIsEstimate" v-slot:[`bottom`]="{}">
      <v-data-table-footer
        :items-per-page-options="itemsPerPageOptions"
        :page-text="getPageTextEstimated(totalItemsEstimated)"
        :show-current-page=true
      ></v-data-table-footer>
    </template>
    <template v-slot:[`top`]="{}">
      <v-row align="center" class="pb-2" style="min-width: 800px;">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Suppliers' }}</div>
          </div>
        </v-col>
        <v-col>
          <v-btn class="float-end" color="primary" :to="{ name: 'New supplier'}">Add</v-btn>
          <v-switch hide-details v-model="showFilters" class="float-right mr-7" :color="showFilters ? 'teal' : 'undefined'" density="compact"
          ></v-switch>
          <v-icon icon="mdi-filter-variant" class="float-right mt-2 mr-3"></v-icon>
        </v-col>
      </v-row>
      <v-row v-show="showFilters" class="mt-0">
        <v-col cols="12" sm="6" lg="4">
          <v-text-field label="Company name search" v-model="filterCompanyName" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="4">
          <v-text-field label="Contact name search" v-model="filterContactName" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="4">
          <v-autocomplete label="Country" v-model="filterCountryID" clearable
            :items="commonStore.activeCountriesList" item-title="name" item-value="id"
            @update:model-value="refreshItems"
          ></v-autocomplete>
        </v-col>
      </v-row>
    </template>
    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" :to="{ name: 'Supplier detail', params: { id: item.id }}">
        <v-icon color="primary" icon="mdi-details"></v-icon>
      </v-btn>
    </template>
  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import ax from '@/api'
import { Supplier } from '@/types/core'
import { debounceMs, maxDebounceMs, getPageTextEstimated, itemsPerPageOptions, processURIOptions } from '@/composables/datatable'
import { useCommonStore } from '@/stores/common'
import { useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  title?: string
}>()

const commonStore = useCommonStore()

var headers = [
  { title: 'Company', key: 'company_name' },
  { title: 'Contact', key: 'contact_name' },
  { title: 'Job title', key: 'contact_title' },
  { title: 'Phone', key: 'phone' },
  { title: 'Country', key: 'country' },
  { title: '# products', key: 'active_product_count', align: 'end' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const

const items = ref<Supplier[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const showFilters = ref(false)
const filterCompanyName = ref<string>()
const filterContactName = ref<string>()
const filterCountryID = ref<number>()
const lsKey = 'suppliers_dt'

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {

  var myURL = '/a/core/suppliers'
  myURL = processURIOptions(myURL, options)

  // filters
  if (filterCompanyName.value) {
    myURL += '&company_name=~' + filterCompanyName.value + '~'
  }
  if (filterContactName.value) {
    myURL += '&contact_name=~' + filterContactName.value + '~'
  }
  if (filterCountryID.value) {
    myURL += '&country_fk=' + filterCountryID.value
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

watch([itemsPerPage, showFilters, search, sortBy], () => {

  if (!showFilters.value) {
    filterCompanyName.value = undefined
    filterContactName.value = undefined
    filterCountryID.value = undefined
    refreshItems()
  }

  let lsObj = {
    'itemsPerPage': itemsPerPage.value,
    'sortBy': sortBy.value,
    'showFilters': showFilters.value,
    'filterCompanyName': filterCompanyName.value,
    'filterContactName': filterContactName.value,
    'filterCountryID': filterCountryID.value,
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
  if (lsObj['showFilters']) { showFilters.value = lsObj['showFilters'] }
  if (lsObj['filterCompanyName']) { filterCompanyName.value = lsObj['filterCompanyName'] }
  if (lsObj['filterContactName']) { filterContactName.value = lsObj['filterContactName'] }
  if (lsObj['filterCountryID']) { filterCountryID.value = lsObj['filterCountryID'] }
})

</script>
