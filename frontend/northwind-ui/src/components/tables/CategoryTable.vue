<template>
  <v-dialog v-model="showDialog" persistent width="auto">
    <CategoryForm :id="editID"
      @cancel="showDialog = false"
      @create="showDialog = false; refreshItems()"
      @delete="showDialog = false; refreshItems()"
      @update="showDialog = false; refreshItems()"
    ></CategoryForm>
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
    class="pa-4 rounded"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center" class="pb-2">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Categories' }}</div>
          </div>
        </v-col>

        <v-col>
          <v-btn class="float-end" color="primary" @click="editID = 0; showDialog = true">Add</v-btn>

          <v-btn icon flat size="small" class="float-right mr-7" :href="excelDlUrl" download>
            <v-icon icon="mdi-file-download-outline"></v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.name`]="{ item }">
      <span v-if="item.color_hex" 
        class="color-pill" 
        :class="item.color_is_light ? 'text-black' : 'text-white'"
        :style="'background-color: ' + item.color_hex + ';'"
      >{{ item.name }}</span>
      <span v-else>{{ item.name }}</span>
    </template>

    <template v-slot:[`item.color_hex`]="{ item }">
      <span v-if="item.color_hex" class="dot" :style="'background-color: ' + item.color_hex + ';'"></span>
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
import { ref, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import ax from '@/api'
import { Category } from '@/types/core'
import { getPageTextEstimated, itemsPerPageOptions, processURIOptions } from '@/composables/datatable'
import CategoryForm from '@/components/forms/CategoryForm.vue'

const props = defineProps<{
  title?: string
}>()

// declaring headers like this to avoid typescript issue
var headers = [
  { title: 'Name', key: 'name' },
  { title: 'Description', key: 'description' },
  { title: 'Color', key: 'color_hex' },
  { title: '# products', key: 'active_product_count', align: 'end' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const

const baseUrl = '/a/core/categories'
const excelDlUrl = import.meta.env.VITE_API_URL +  baseUrl + '?xformat=excel'

const items = ref<Category[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const editID = ref(0)
const showDialog = ref(false)

const lsKey = 'categories_dt'

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {

  var myURL = baseUrl
  myURL = processURIOptions(myURL, options)

  // exclude None
  myURL += '&id=!-1'

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

<style scoped>
.dot {
  height: 25px;
  width: 25px;
  border-radius: 50%;
  display: inline-block;
  margin-top: 5px;
}
</style>