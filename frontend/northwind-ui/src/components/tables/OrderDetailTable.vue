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
      <v-row align="center" class="pb-2" style="min-width: 800px;">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Order details' }}</div>
          </div>
        </v-col>
      </v-row>
    </template>
    <template v-slot:[`item.product_name`]="{ item }">
      <router-link :to="{ name: 'Product detail', params: {id: item.product_fk }}">
        {{item.product_name}}
      </router-link>
    </template>
    <template v-slot:[`item.discount`]="{ item }">
      <span v-if="item.discount">{{ Intl.NumberFormat('en-US', { style: 'percent', maximumFractionDigits: 0, minimumFractionDigits: 0}).format(item.discount) }}</span>
    </template>
  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import ax from '@/api'
import { OrderDetail } from '@/types/sales'
import { getPageTextEstimated, itemsPerPageOptions, processURIOptions } from '@/composables/datatable'
import { useCoreStore } from '@/stores/core'

const props = defineProps<{
  order_id: number // pass 0 rather than null/undefined, easier to handle
  title?: string
}>()

const coreStore = useCoreStore()

var headers = [
  { title: 'Order number', key: 'order_number' },
  { title: 'Product', key: 'product_name' },
  { title: 'Quantity', key: 'quantity', align: 'end' },
  { title: 'Unit price', key: 'unit_price', align: 'end' },
  { title: 'Discount', key: 'discount', align: 'end' },
] as const

const items = ref<OrderDetail[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const lsKey = 'order_details_dt'

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {

  var myURL = '/a/sales/order-details'
  myURL = processURIOptions(myURL, options)

  // filters
  if (props.order_id > 0) {
    myURL += '&order_fk=' + props.order_id
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

watch([itemsPerPage, search, sortBy], () => {

  let lsObj = {
    'itemsPerPage': itemsPerPage.value,
    'sortBy': sortBy.value,
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
})

</script>
