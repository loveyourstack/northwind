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

            <FilterChipText name="Order #" :filterValue="filterOrderNumber" 
              @closed="filterOrderNumber = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterOrderNumber = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChipText name="Customer name" :filterValue="filterCustomerName" 
              @closed="filterCustomerName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterCustomerName = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChipBool name="Shipped" :filterValue="filterShipped" :filterText="filterShippedText" @closed="filterShipped = undefined; refreshItems()">
              <template #menuContent>
                <v-autocomplete label="Shipped" v-model="filterShipped" autofocus
                  :items="appStore.booleanOptions"
                  @update:model-value="refreshItems"
                ></v-autocomplete>
              </template>
            </FilterChipBool>

            <FilterChip name="Order date >" :filterValue="filterOrderDate" :filterText="filterOrderDateText" @closed="filterOrderDate = undefined; refreshItems()">
              <template #menuContent>
                <DateTextField :dateVal="filterOrderDate" label="Order date >"
                  @updated="(val: string | undefined) => { filterOrderDate = val; refreshItems() } "
                ></DateTextField>
              </template>
            </FilterChip>

          </v-chip-group>
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
import { useAppStore } from '@/stores/app'
import { Order } from '@/types/sales'
import { getTextFilterUrlParam, itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'
import AdjustColsListItem from '@/components/AdjustColsListItem.vue'
import DateTextField from '@/components/DateTextField.vue'
import DtFooter from '@/components/DtFooter.vue'
import FilterChip from '@/components/FilterChip.vue'
import FilterChipBool from '@/components/FilterChipBool.vue'
import FilterChipText from '@/components/FilterChipText.vue'

const props = defineProps<{
  customer_id: number // pass 0 rather than null/undefined, easier to handle
  title?: string
}>()

const appStore = useAppStore()

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
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const filterCustomerName = ref<string>()

const filterOrderDate = ref<string>() // YYYY-MM-DD
const filterOrderDateText = computed(() => {
  return filterOrderDate.value != undefined ? useDateFormat(filterOrderDate, 'DD MMM YYYY').value : ''
})

const filterOrderNumber = ref<string>()

const filterShipped = ref<boolean>()
const filterShippedText = computed(() => {
  return filterShipped.value != undefined ? appStore.booleanOptions.find(ele => ele.value === filterShipped.value)?.title : ''
})

const lsKey = 'orders_dt'

function getFilterStr(): string {
  var ret = ''

  // prop filter overrides regular filter
  if (props.customer_id > 0) {
    ret += '&customer_fk=' + props.customer_id
  } else if (filterCustomerName.value) {
    ret += getTextFilterUrlParam('customer_company_name', filterCustomerName.value)
  }

  if (filterOrderDate.value) {
    ret += '&order_date=>eq' + filterOrderDate.value
  }

  if (filterOrderNumber.value) {
    ret += '&order_number=~' + filterOrderNumber.value + '~'
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
    'filterOrderDate': filterOrderDate.value,
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
  if (lsObj['filterOrderNumber']) { filterOrderNumber.value = lsObj['filterOrderNumber'] }
  if (lsObj['filterCustomerName']) { filterCustomerName.value = lsObj['filterCustomerName'] }
  if (lsObj['filterShipped']) { filterShipped.value = lsObj['filterShipped'] }
  if (lsObj['filterOrderDate']) { filterOrderDate.value = lsObj['filterOrderDate'] }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
