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
            <div class="dt-title">{{ props.title ? props.title : 'Orders' }}</div>
          </div>
        </v-col>
        <v-col>
          <!--<v-btn class="float-end" color="primary" :to="{ name: 'New order'}">Add</v-btn>-->
          <v-switch hide-details v-model="showFilters" class="float-right mr-7" :color="showFilters ? 'teal' : 'undefined'" density="compact"
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
        <v-icon color="light-blue" icon="mdi-details"></v-icon>
      </v-btn>
    </template>
  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import ax from '@/api'
import { Order } from '@/types/sales'
import { debounceMs, maxDebounceMs, getPageTextEstimated, itemsPerPageOptions, processURIOptions } from '@/composables/datatable'
import { booleanOptions } from '@/composables/form'
import { useCommonStore } from '@/stores/common'
import { useDateFormat, useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  customer_id: number // pass 0 rather than null/undefined, easier to handle
  title?: string
}>()

const commonStore = useCommonStore()

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

const items = ref<Order[]>([])
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

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {

  var myURL = '/a/sales/orders'
  myURL = processURIOptions(myURL, options)

  // filters
  if (filterOrderNumber.value) {
    myURL += '&order_number=~' + filterOrderNumber.value + '~'
  }

  // prop filter overrides regular filter
  if (props.customer_id > 0) {
    myURL += '&customer_fk=' + props.customer_id
  } else if (filterCustomerName.value) {
    myURL += '&customer_company_name=~' + filterCustomerName.value + '~'
  }

  // checking bool filter like this so that undefined and null do not count as false
  if (filterShipped.value == false || filterShipped.value == true) {
    myURL += '&is_shipped=' + filterShipped.value
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
  if (lsObj['filterOrderNumber']) { filterOrderNumber.value = lsObj['filterOrderNumber'] }
  if (lsObj['filterCustomerName']) { filterCustomerName.value = lsObj['filterCustomerName'] }
  if (lsObj['filterShipped']) { filterShipped.value = lsObj['filterShipped'] }
})

</script>
