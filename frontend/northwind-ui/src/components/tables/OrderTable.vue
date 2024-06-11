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
    :row-props="({ item }) => getRowClass(item)"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center" class="pb-2">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Orders' }}</div>
            <div>Unshipped orders are marked red.</div>
          </div>

          <!--<v-btn class="float-end" color="primary" :to="{ name: 'New order'}">Add</v-btn>-->

          <v-btn icon flat size="small" class="float-right mr-7" v-tooltip="'Download to Excel'" @click="fileDownload(excelDlUrl)">
            <v-icon icon="mdi-file-download-outline"></v-icon>
          </v-btn>

          <v-menu :close-on-content-click=false>
            <template v-slot:activator="{ props }">
              <v-btn density="comfortable" v-tooltip="'Adjust columns'" flat class="float-right mr-5" icon="mdi-table-column" v-bind="props"></v-btn>
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

          <v-switch hide-details v-model="showFilters" v-tooltip="'Show filters'" class="float-right mr-7" :color="showFilters ? 'teal' : 'undefined'" density="compact"
          ></v-switch>
          <v-icon icon="mdi-filter-variant" class="float-right mt-2 mr-3"></v-icon>
        </v-col>
      </v-row>

      <v-row v-show="showFilters" class="mt-0">
        <v-col cols="12" sm="6" lg="4">
          <v-text-field label="Order # search" v-model="filterOrderNumber" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="4">
          <v-text-field label="Customer search" v-model="filterCustomerName" clearable
            @update:model-value="debouncedRefreshItems"
            :disabled="props.customer_id > 0"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="4">
          <v-autocomplete label="Shipped" v-model="filterShipped" clearable
            :items="booleanOptions"
            @update:model-value="refreshItems"
          ></v-autocomplete>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.customer_company_name`]="{ item }">
      <router-link :to="{ name: 'Customer detail', params: {id: item.customer_fk }}">
        {{item.customer_company_name}}
      </router-link>
    </template>

    <template v-slot:[`item.is_shipped`]="{ item }">
      <v-icon v-if="item.is_shipped" size="small" icon="mdi-check"></v-icon>
    </template>

    <template v-slot:[`item.order_date`]="{ item }">
      {{ useDateFormat(item.order_date, 'DD MMM YYYY').value }}
    </template>

    <template v-slot:[`item.required_date`]="{ item }">
      {{ useDateFormat(item.required_date, 'DD MMM YYYY').value }}
    </template>

    <template v-slot:[`item.shipped_date`]="{ item }">
      <span v-if="item.is_shipped">{{ useDateFormat(item.shipped_date, 'DD MMM YYYY').value }}</span>
    </template>

    <template v-slot:[`item.order_value`]="{ item }">
      <span>{{ item.order_value.toFixed(2) }}</span>
    </template>

    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" :to="{ name: 'Order detail', params: { id: item.id }}">
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
import { ref, computed, watch, onBeforeMount, onMounted } from 'vue'
import { useTheme } from 'vuetify'
import { VDataTable } from 'vuetify/components'
import ax from '@/api'
import { Order } from '@/types/sales'
import { GetMetadata } from '@/types/system'
import { debounceMs, maxDebounceMs, getHeaderListIcon, getHeaderListIconColor, getPageTextEstimated, itemsPerPageOptions, processURIOptions } from '@/composables/datatable'
import { fileDownload } from '@/composables/file'
import { booleanOptions } from '@/composables/form'
import { useDateFormat, useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  customer_id: number // pass 0 rather than null/undefined, easier to handle
  title?: string
}>()

const theme = useTheme()

var headers = [
  { title: 'Order #', key: 'order_number' },  
  { title: 'Customer', key: 'customer_company_name' },
  { title: '# details', key: 'order_detail_count', align: 'end' },
  { title: 'Value', key: 'order_value', align: 'end' },
  { title: 'Order date', key: 'order_date' },
  { title: 'Salesman', key: 'salesman' },
  { title: 'Req. date', key: 'required_date' },
  { title: 'Shipper', key: 'shipper_company_name' },
  { title: 'Shipped', key: 'is_shipped' },
  { title: 'Shipped date', key: 'shipped_date' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const
const excludedHeaders = ref<string[]>([])
const selectedHeaders = ref()

const baseUrl = '/a/sales/orders'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<Order[]>([])
const metadata = ref<GetMetadata>()
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const showFilters = ref(false)
const filterOrderNumber = ref<string>()
const filterCustomerName = ref<string>()
const filterShipped = ref<boolean>()
const lsKey = 'orders_dt'

function getFilterStr(): string {
  var ret = ''

  if (filterOrderNumber.value) {
    ret += '&order_number=~' + filterOrderNumber.value + '~'
  }

  // prop filter overrides regular filter
  if (props.customer_id > 0) {
    ret += '&customer_fk=' + props.customer_id
  } else if (filterCustomerName.value) {
    ret += '&customer_company_name=~' + filterCustomerName.value + '~'
  }

  // checking bool filter like this so that undefined and null do not count as false
  if (filterShipped.value == false || filterShipped.value == true) {
    ret += '&is_shipped=' + filterShipped.value
  }

  return ret
}

function getRowClass(item: Order) {
  if (!item.is_shipped) {
    return theme.global.current.value.dark ? { style: 'background-color: #480505;' } : { class: 'bg-red-lighten-5' }
  }
}

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {

  var myURL = baseUrl
  myURL = processURIOptions(myURL, options)
  myURL += getFilterStr()

  ax.get(myURL)
  .then(resp => {
      items.value = resp.data.data
      metadata.value = resp.data.metadata

      totalItemsIsEstimate.value = metadata.value!.total_count_is_estimated
      if (totalItemsIsEstimate.value) {
        totalItemsEstimated.value = metadata.value!.total_count
        totalItems.value = 101 // workaround so that next page button is enabled even if estimate is too low
      } else {
        totalItems.value = metadata.value!.total_count
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
    filterOrderNumber.value = undefined
    filterCustomerName.value = undefined
    filterShipped.value = undefined
    refreshItems()
  }

  let lsObj = {
    'itemsPerPage': itemsPerPage.value,
    'sortBy': sortBy.value,
    'showFilters': showFilters.value,
    'filterOrderNumber': filterOrderNumber.value,
    'filterCustomerName': filterCustomerName.value,
    'filterShipped': filterShipped.value,
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
  if (lsObj['filterOrderNumber']) { filterOrderNumber.value = lsObj['filterOrderNumber'] }
  if (lsObj['filterCustomerName']) { filterCustomerName.value = lsObj['filterCustomerName'] }
  if (lsObj['filterShipped']) { filterShipped.value = lsObj['filterShipped'] }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
