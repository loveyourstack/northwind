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
    class="pa-4 rounded"
  >
    <template v-if="totalItemsIsEstimate" v-slot:[`bottom`]="{}">
      <v-data-table-footer
        :items-per-page-options="itemsPerPageOptions"
        :page-text="getPageTextEstimated(totalItemsEstimated)"
        :show-current-page=true
      ></v-data-table-footer>
    </template>
    <template v-slot:[`top`]="{}">
      <v-row align="center" class="pb-2">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Products' }}</div>
          </div>
        </v-col>
        <v-col>
          <v-btn class="float-end" color="primary" :to="{ name: 'New product'}">Add</v-btn>
          <v-switch hide-details v-model="showFilters" class="float-right mr-7" :color="showFilters ? 'teal' : 'undefined'" density="compact"
          ></v-switch>
          <v-icon icon="mdi-filter-variant" class="float-right mt-2 mr-3"></v-icon>
        </v-col>
      </v-row>
      <v-row v-show="showFilters" class="mt-0">
        <v-col cols="12" sm="6" lg="3">
          <v-text-field label="Name search" v-model="filterName" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="3">
          <v-autocomplete label="Category" v-model="filterCategoryID" clearable
            :items="coreStore.categoriesList" item-title="name" item-value="id"
            @update:model-value="refreshItems"
          ></v-autocomplete>
        </v-col>
        <v-col cols="12" sm="6" lg="3">
          <v-autocomplete label="Supplier" v-model="filterSupplierID" clearable
            :items="coreStore.suppliersList" item-title="name" item-value="id"
            @update:model-value="refreshItems"
            :disabled="props.supplier_id > 0"
          ></v-autocomplete>
        </v-col>
        <v-col cols="12" sm="6" lg="3">
          <v-autocomplete label="Discontinued" v-model="filterDiscontinued"
            :items="booleanOptions"
            @update:model-value="refreshItems"
          ></v-autocomplete>
        </v-col>
      </v-row>
    </template>
    <template v-slot:[`item.category`]="{ item }">
      <span v-if="item.category_color_hex" 
        class="color-pill" 
        :class="item.category_color_is_light ? 'text-black' : 'text-white'"
        :style="'background-color: ' + item.category_color_hex + ';'"
      >{{ item.category }}</span>
      <span v-else>{{ item.category }}</span>
    </template>
    <template v-slot:[`item.supplier_company_name`]="{ item }">
      <router-link :to="{ name: 'Supplier detail', params: {id: item.supplier_fk }}">
        {{item.supplier_company_name}}
      </router-link>
    </template>
    <template v-slot:[`item.unit_price`]="{ item }">
      <span>{{ item.unit_price.toFixed(2) }}</span>
    </template>
    <template v-slot:[`item.is_discontinued`]="{ item }">
      <v-icon v-if="item.is_discontinued" size="small" icon="mdi-check"></v-icon>
    </template>
    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" :to="{ name: 'Product detail', params: { id: item.id }}">
        <v-icon color="primary" icon="mdi-details"></v-icon>
      </v-btn>
    </template>
  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import ax from '@/api'
import { Product } from '@/types/core'
import { debounceMs, maxDebounceMs, getPageTextEstimated, itemsPerPageOptions, processURIOptions } from '@/composables/datatable'
import { booleanOptions } from '@/composables/form'
import { useCoreStore } from '@/stores/core'
import { useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  supplier_id: number // pass 0 rather than null/undefined, easier to handle
  title?: string
}>()

const coreStore = useCoreStore()

var headers = [
  { title: 'Name', key: 'name' },
  { title: 'Category', key: 'category' },
  { title: 'Supplier', key: 'supplier_company_name' },
  { title: 'Qty / unit', key: 'quantity_per_unit' },
  { title: 'Unit price', key: 'unit_price', align: 'end' },
  { title: 'Units in stock', key: 'units_in_stock', align: 'end' },
  { title: 'Units on order', key: 'units_on_order', align: 'end' },
  { title: 'Reorder level', key: 'reorder_level', align: 'end' },
  { title: 'Discontinued', key: 'is_discontinued', align: 'center' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const

const items = ref<Product[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const showFilters = ref(false)
const filterName = ref<string>()
const filterCategoryID = ref<number>()
const filterSupplierID = ref<number>()
const filterDiscontinued = ref(false)
const lsKey = 'products_dt'

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {

  var myURL = '/a/core/products'
  myURL = processURIOptions(myURL, options)

  // filters
  if (filterName.value) {
    myURL += '&name=~' + filterName.value + '~'
  }
  if (filterCategoryID.value) {
    myURL += '&category_fk=' + filterCategoryID.value
  }

  // prop filter overrides regular filter
  if (props.supplier_id > 0) {
    myURL += '&supplier_fk=' + props.supplier_id
  } else if (filterSupplierID.value) {
    myURL += '&supplier_fk=' + filterSupplierID.value
  }

  myURL += '&is_discontinued=' + filterDiscontinued.value

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
  // changing data-table search property to a new value triggers loadItems
  search.value = String(Date.now())
}

watch([itemsPerPage, showFilters, search, sortBy], () => {

  if (!showFilters.value) {
    filterName.value = undefined
    filterCategoryID.value = undefined
    filterSupplierID.value = undefined
    filterDiscontinued.value = false
    refreshItems()
  }

  let lsObj = {
    'itemsPerPage': itemsPerPage.value,
    'sortBy': sortBy.value,
    'showFilters': showFilters.value,
    'filterName': filterName.value,
    'filterCategoryID': filterCategoryID.value,
    'filterSupplierID': filterSupplierID.value,
    'filterDiscontinued': filterDiscontinued.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
})

// onBeforeMount, since using onMounted results in the initial data load taking place before this runs
onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['itemsPerPage']) { itemsPerPage.value = lsObj['itemsPerPage'] }
  if (lsObj['sortBy']) { sortBy.value = lsObj['sortBy'] }
  if (lsObj['showFilters']) { showFilters.value = lsObj['showFilters'] }
  if (lsObj['filterName']) { filterName.value = lsObj['filterName'] }
  if (lsObj['filterCategoryID']) { filterCategoryID.value = lsObj['filterCategoryID'] }
  if (lsObj['filterSupplierID']) { filterSupplierID.value = lsObj['filterSupplierID'] }
  if (lsObj['filterDiscontinued']) { filterDiscontinued.value = lsObj['filterDiscontinued'] }
})

</script>
