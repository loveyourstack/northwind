<template>
  <v-dialog v-model="showDialog" persistent width="auto">
    <OrderItemForm :order_id="order_id" :order_number="props.order_number" :id="editID"
      @archive="showDialog = false; refreshItems()"
      @cancel="showDialog = false"
      @create="showDialog = false; refreshItems()"
      @update="showDialog = false; refreshItems()"
    ></OrderItemForm>
  </v-dialog>

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
    class="dt"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center" class="pb-2">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Order items' }}</div>
          </div>

          <v-btn class="float-end" color="primary" @click="editID = 0; showDialog = true">Add</v-btn>

          <v-btn icon flat size="small" class="float-right mr-7" v-tooltip="'Download to Excel'" @click="fileDownload(excelDlUrl)">
            <v-icon icon="mdi-file-download-outline"></v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.product_name`]="{ item }">
      <router-link :to="{ name: 'Product detail', params: {id: item.product_fk }}">
        {{item.product_name}}
      </router-link>
    </template>

    <template v-slot:[`item.unit_price`]="{ item }">
      <span>{{ '$' + item.unit_price.toFixed(2) }}</span>
    </template>

    <template v-slot:[`item.discount`]="{ item }">
      <span v-if="item.discount">{{ Intl.NumberFormat('en-US', { style: 'percent', maximumFractionDigits: 0, minimumFractionDigits: 0}).format(item.discount) }}</span>
    </template>

    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" @click="editID = item.id; showDialog = true">
        <v-icon color="primary" icon="mdi-pencil"></v-icon>
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
import { ref, computed, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import ax from '@/api'
import { OrderItem } from '@/types/sales'
import { GetMetadata } from '@/types/system'
import { getPageTextEstimated, itemsPerPageOptions, processURIOptions } from '@/composables/datatable'
import { fileDownload } from '@/composables/file'
import OrderItemForm from '@/components/forms/OrderItemForm.vue'

const props = defineProps<{
  order_id: number // pass 0 rather than null/undefined, easier to handle
  order_number: number
  title?: string
}>()

var headers = [
  { title: 'Product', key: 'product_name' },
  { title: 'Quantity', key: 'quantity', align: 'end' },
  { title: 'Unit price', key: 'unit_price', align: 'end' },
  { title: 'Discount', key: 'discount', align: 'end' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const

const baseUrl = '/a/sales/order-items'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<OrderItem[]>([])
const metadata = ref<GetMetadata>()
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const editID = ref(0)
const showDialog = ref(false)

const lsKey = 'order_items_dt'

function getFilterStr(): string {
  var ret = ''

  if (props.order_id > 0) {
    ret += '&order_fk=' + props.order_id
  }

  return ret
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

function refreshItems() {
  search.value = String(Date.now())
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
