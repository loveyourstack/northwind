<template>
  <v-dialog v-model="showDialog" persistent width="auto">
    <SalesTerritoryForm :id="editID"
      @cancel="showDialog = false"
      @create="showDialog = false; refreshItems()"
      @delete="showDialog = false; refreshItems()"
      @update="showDialog = false; refreshItems()"
    ></SalesTerritoryForm>
  </v-dialog>

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
      <v-row align="center" class="pb-2" style="min-width: 800px;">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Territories' }}</div>
          </div>

          <v-btn class="float-end" color="primary" @click="editID = 0; showDialog = true">Add</v-btn>

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

            <FilterChipText name="Code" :filterValue="filterCode" 
              @closed="filterCode = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterCode = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChip name="Region" :filterValue="filterRegion" :filterText="filterRegion" @closed="filterRegion = ''; refreshItems()">
              <template #menuContent>
                <v-autocomplete label="Region" v-model="filterRegion" autofocus
                  :items="salesStore.regions"
                  @update:model-value="refreshItems"
                ></v-autocomplete>
              </template>
            </FilterChip>

            <FilterChip name="Salesman" :filterValue="filterSalesmanID" :filterText="filterSalesmanIDText" @closed="filterSalesmanID = undefined; refreshItems()">
              <template #menuContent>
                <v-autocomplete label="Salesman" v-model="filterSalesmanID" autofocus
                  :items="hrStore.employeesList" item-title="name" item-value="id"
                  @update:model-value="refreshItems"
                ></v-autocomplete>
              </template>
            </FilterChip>

          </v-chip-group>
        </v-col>

      </v-row>
    </template>

    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" @click="editID = item.id; showDialog = true">
        <v-icon color="primary" icon="mdi-pencil"></v-icon>
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
import { type Territory } from '@/types/sales'
import { getTextFilterUrlParam, itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'
import { useHRStore } from '@/stores/hr'
import { useSalesStore } from '@/stores/sales'
import AdjustColsListItem from '@/components/AdjustColsListItem.vue'
import DtFooter from '@/components/DtFooter.vue'
import FilterChip from '@/components/FilterChip.vue'
import FilterChipText from '@/components/FilterChipText.vue'
import SalesTerritoryForm from '@/components/forms/SalesTerritoryForm.vue'

const props = defineProps<{
  salesman_id: number
  title?: string
}>()

const hrStore = useHRStore()
const salesStore = useSalesStore()

var headers = [
  { title: 'Name', key: 'name' },  
  { title: 'Code', key: 'code' },
  { title: 'Region', key: 'region' },
  { title: 'Salesman', key: 'salesman' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const
const excludedHeaders = ref<string[]>([])
const selectedHeaders = ref()

const baseUrl = '/a/sales/territories'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<Territory[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const editID = ref(0)
const showDialog = ref(false)

const filterCode = ref<string>()
const filterName = ref<string>()
const filterRegion = ref<string>()

const filterSalesmanID = ref<number>()
const filterSalesmanIDText = computed(() => {
  return filterSalesmanID.value ? hrStore.employeesList.find(ele => ele.id === filterSalesmanID.value)?.name : ''
})

const lsKey = 'territories_dt'

function getFilterStr(): string {
  var ret = ''

  // prop filter overrides regular filter
  if (props.salesman_id > 0) {
    ret += '&salesman_fk=' + props.salesman_id
  } else if (filterSalesmanID.value) {
    ret += '&salesman_fk=' + filterSalesmanID.value
  }

  ret += getTextFilterUrlParam('code', filterCode.value)
  ret += getTextFilterUrlParam('name', filterName.value)
  ret += getTextFilterUrlParam('region', filterRegion.value)

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
    'filterName': filterName.value,
    'filterRegion': filterRegion.value,
    'filterSalesmanID': filterSalesmanID.value,
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
  if (lsObj['filterCode']) { filterCode.value = lsObj['filterCode'] }
  if (lsObj['filterName']) { filterName.value = lsObj['filterName'] }
  if (lsObj['filterRegion']) { filterRegion.value = lsObj['filterRegion'] }
  if (lsObj['filterSalesmanID']) { filterSalesmanID.value = lsObj['filterSalesmanID'] }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
