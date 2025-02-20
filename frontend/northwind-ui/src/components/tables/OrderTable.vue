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
    :row-props="(item: Order) => getRowClass(item)"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Orders' }}</div>
            <div>Unshipped orders are marked red.</div>
          </div>

          <v-btn class="float-end" color="primary" :to="{ name: 'New order'}">Add</v-btn>

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
        <v-col cols="12" sm="6" lg="4">
          <v-text-field label="Order # search" v-model="filterOrderNumber" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="4" v-if="props.customer_id == 0">
          <v-text-field label="Customer search" v-model="filterCustomerName" clearable
            @update:model-value="debouncedRefreshItems"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" lg="4">
          <v-autocomplete label="Shipped" v-model="filterShipped" clearable
            :items="coreStore.booleanOptions"
            @update:model-value="refreshItems"
          ></v-autocomplete>
        </v-col>
      </v-row>

      <v-row class="mt-0">
        <v-col cols="12" sm="6" lg="4">
          <v-menu v-model="showOrderDateDp" :close-on-content-click="false">
            <template #activator="{ props }">
              <v-text-field id="order-date-text" v-bind="props" label="Order date >" prepend-icon="mdi-calendar" readonly clearable
                :model-value="filterOrderDate ? useDateFormat(filterOrderDate, 'DD MMM YYYY').value : undefined"
                @click:clear="filterOrderDate = undefined; refreshItems()"
              ></v-text-field>
            </template>
            <template #default>
              <v-date-picker color="primary" v-model="filterOrderDate" @update:model-value="showOrderDateDp = false; refreshItems()"></v-date-picker>
            </template>
          </v-menu>
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
      <span>{{ '$' + item.order_value.toFixed(2) }}</span>
    </template>

    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" :to="{ name: 'Order detail', params: { id: item.id }}">
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
import { useDateFormat, useDebounceFn } from '@vueuse/core'
import { useTheme } from 'vuetify'
import { VDataTable } from 'vuetify/components'
import { useFetchDt } from '@/composables/fetch'
import { useCoreStore } from '@/stores/core'
import { Order } from '@/types/sales'
import { GetMetadata } from '@/types/system'
import { getHeaderListIcon, getHeaderListIconColor, itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'
import DtFooter from '@/components/DtFooter.vue'

const props = defineProps<{
  customer_id: number // pass 0 rather than null/undefined, easier to handle
  title?: string
}>()

const coreStore = useCoreStore()

const theme = useTheme()

var headers = [
  { title: 'Order #', key: 'order_number' },  
  { title: 'Customer', key: 'customer_company_name' },
  { title: '# items', key: 'order_item_count', align: 'end' },
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

const filterOrderNumber = ref<string>()
const filterCustomerName = ref<string>()
const filterShipped = ref<boolean>()
/*
  filter type Date:
  - use menu/text/date picker with a "show" bool variable to control menu visibility
  - useDateFormat in getFilterStr
  - save as YYYY-MM-DD string in LS
  - default empty string when writing to LS in watch (ideally undefined, but this causes type error), useDateFormat if truthy
  - new Date() when reading from LS in onBeforeMount
*/
const filterOrderDate = ref<Date>()
const lsKey = 'orders_dt'

const showOrderDateDp = ref(false)

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

  if (filterOrderDate.value) {
    ret += '&order_date=>' + useDateFormat(filterOrderDate, 'YYYY-MM-DD').value
  }

  return ret
}

function getRowClass(item: Order) {
  if (!item.is_shipped) {
    return theme.global.current.value.dark ? { style: 'background-color: #480505;' } : { class: 'bg-red-lighten-5' }
  }
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
    'filterOrderNumber': filterOrderNumber.value,
    'filterCustomerName': filterCustomerName.value,
    'filterShipped': filterShipped.value,
    'filterOrderDate': '',
    'excludedHeaders': excludedHeaders.value,
  }
  if (filterOrderDate.value) {
    lsObj.filterOrderDate = useDateFormat(filterOrderDate, 'YYYY-MM-DD').value
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
  if (lsObj['filterOrderNumber']) { filterOrderNumber.value = lsObj['filterOrderNumber'] }
  if (lsObj['filterCustomerName']) { filterCustomerName.value = lsObj['filterCustomerName'] }
  if (lsObj['filterShipped']) { filterShipped.value = lsObj['filterShipped'] }
  if (lsObj['filterOrderDate']) { filterOrderDate.value = new Date(lsObj['filterOrderDate']) }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
