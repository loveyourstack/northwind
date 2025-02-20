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

              <v-list-item prepend-icon="mdi-download-outline">
                <v-list-item-title class="clickable" @click="fileDownload(excelDlUrl)">Download to Excel</v-list-item-title>
              </v-list-item>

              <v-menu :close-on-content-click=false location="start">
                <template v-slot:activator="{ props }">
                  <v-list-item v-bind="props" prepend-icon="mdi-table-column">
                    <v-list-item-title class="clickable">Adjust columns</v-list-item-title>
                  </v-list-item>
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

            </v-list>
          </v-menu>

        </v-col>
      </v-row>

      <v-row class="mt-0">
        <v-col cols="12" sm="6" lg="3">
          <v-text-field label="Code search" v-model="filterCode" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="3">
          <v-text-field label="Company name search" v-model="filterCompanyName" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="3">
          <v-text-field label="Contact name search" v-model="filterContactName" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="3">
          <v-autocomplete label="Country" v-model="filterCountryID" clearable
            :items="commonStore.activeCountriesList" item-title="name" item-value="id"
            @update:model-value="refreshItems"
          ></v-autocomplete>
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
import { Customer } from '@/types/sales'
import { GetMetadata } from '@/types/system'
import { getHeaderListIcon, getHeaderListIconColor, itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'
import { useCommonStore } from '@/stores/common'
import { useDebounceFn } from '@vueuse/core'
import DtFooter from '@/components/DtFooter.vue'

const props = defineProps<{
  title?: string
}>()

const commonStore = useCommonStore()

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
const metadata = ref<GetMetadata>()
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
const lsKey = 'customers_dt'

function getFilterStr(): string {
  var ret = ''

  if (filterCode.value) {
    ret += '&code=~' + filterCode.value + '~'
  }
  if (filterCompanyName.value) {
    ret += '&company_name=~' + filterCompanyName.value + '~'
  }
  if (filterContactName.value) {
    ret += '&contact_name=~' + filterContactName.value + '~'
  }
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
