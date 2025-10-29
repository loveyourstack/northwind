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
            <div class="dt-title">{{ props.title ? props.title : 'Products' }}</div>
          </div>

          <v-btn class="float-end" color="primary" :to="{ name: 'New product'}">Add</v-btn>

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

            <FilterChipText name="Name" :filterValue="filterName" 
              @closed="filterName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterName = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChip name="Category" :filterValue="filterCategoryID" :filterText="filterCategoryIDText" @closed="filterCategoryID = undefined; refreshItems()">
              <template #menuContent>
                <v-autocomplete label="Category" v-model="filterCategoryID" autofocus
                  :items="coreStore.categoriesList" item-title="name" item-value="id"
                  @update:model-value="refreshItems"
                ></v-autocomplete>
              </template>
            </FilterChip>

            <FilterChip name="Supplier" :filterValue="filterSupplierID" :filterText="filterSupplierIDText" @closed="filterSupplierID = undefined; refreshItems()">
              <template #menuContent>
                <v-autocomplete label="Supplier" v-model="filterSupplierID" autofocus
                  :items="coreStore.suppliersList" item-title="name" item-value="id"
                  @update:model-value="refreshItems"
                ></v-autocomplete>
              </template>
            </FilterChip>

            <FilterChipBool name="Discontinued" :filterValue="filterDiscontinued" :filterText="filterDiscontinuedText" @closed="filterDiscontinued = undefined; refreshItems()">
              <template #menuContent>
                <v-autocomplete label="Discontinued" v-model="filterDiscontinued" autofocus
                  :items="appStore.booleanOptions"
                  @update:model-value="refreshItems"
                ></v-autocomplete>
              </template>
            </FilterChipBool>

          </v-chip-group>
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
      <router-link :to="{ path: '/suppliers/' + item.supplier_fk }">
        {{item.supplier_company_name}}
      </router-link>
    </template>

    <template v-slot:[`item.unit_price`]="{ item }">
      <span>{{ '$' + item.unit_price.toFixed(2) }}</span>
    </template>

    <template v-slot:[`item.is_discontinued`]="{ item }">
      <v-icon v-if="item.is_discontinued" size="small" icon="mdi-check"></v-icon>
    </template>

    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" :to="{ name: 'Product detail', params: { id: item.id }}">
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
import { useDebounceFn } from '@vueuse/core'
import { VDataTable } from 'vuetify/components'
import { useFetchDt } from '@/composables/fetch'
import { type Product } from '@/types/core'
import { getTextFilterUrlParam, itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'
import { useAppStore } from '@/stores/app'
import { useCoreStore } from '@/stores/core'

const props = defineProps<{
  supplier_id: number // pass 0 rather than null/undefined, easier to handle
  title?: string
}>()

const appStore = useAppStore()
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
const excludedHeaders = ref<string[]>([])
const selectedHeaders = ref()

const baseUrl = '/a/core/products'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<Product[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const filterName = ref<string>()

const filterCategoryID = ref<number>()
const filterCategoryIDText = computed(() => {
  return filterCategoryID.value ? coreStore.categoriesList.find(ele => ele.id === filterCategoryID.value)?.name : ''
})

const filterSupplierID = ref<number>()
const filterSupplierIDText = computed(() => {
  return filterSupplierID.value ? coreStore.suppliersList.find(ele => ele.id === filterSupplierID.value)?.name : ''
})

const filterDiscontinued = ref<boolean>()
const filterDiscontinuedText = computed(() => {
  return filterDiscontinued.value != undefined ? appStore.booleanOptions.find(ele => ele.value === filterDiscontinued.value)?.title : ''
})

const lsKey = 'products_dt'

function getFilterStr(): string {
  var ret = ''

  ret += getTextFilterUrlParam('name', filterName.value)
  
  if (filterCategoryID.value) {
    ret += '&category_fk=' + filterCategoryID.value
  }

  // prop filter overrides regular filter
  if (props.supplier_id > 0) {
    ret += '&supplier_fk=' + props.supplier_id
  } else if (filterSupplierID.value) {
    ret += '&supplier_fk=' + filterSupplierID.value
  }

  // checking bool filter like this so that undefined and null do not count as false
  if (filterDiscontinued.value == false || filterDiscontinued.value == true) {
    ret += '&is_discontinued=' + filterDiscontinued.value
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
  // changing data-table search property to a new value triggers loadItems
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
    'filterName': filterName.value,
    'filterCategoryID': filterCategoryID.value,
    'filterSupplierID': filterSupplierID.value,
    'filterDiscontinued': filterDiscontinued.value,
    'excludedHeaders': excludedHeaders.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
}, { deep: true })

// onBeforeMount, since using onMounted results in the initial data load taking place before this runs
onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['itemsPerPage']) { itemsPerPage.value = lsObj['itemsPerPage'] }
  if (lsObj['sortBy']) { sortBy.value = lsObj['sortBy'] }
  if (lsObj['filterName']) { filterName.value = lsObj['filterName'] }
  if (lsObj['filterCategoryID']) { filterCategoryID.value = lsObj['filterCategoryID'] }
  if (lsObj['filterSupplierID']) { filterSupplierID.value = lsObj['filterSupplierID'] }
  if (lsObj['filterDiscontinued']) { filterDiscontinued.value = lsObj['filterDiscontinued'] }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
