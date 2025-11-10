<template>
  <v-dialog v-model="showDialog" persistent width="auto">
    <HrEmployeeForm :id="editID"
      @cancel="showDialog = false"
      @create="showDialog = false; refreshItems()"
      @delete="showDialog = false; refreshItems()"
      @update="showDialog = false; refreshItems()"
    ></HrEmployeeForm>
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
      <v-row align="center">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Employees' }}</div>
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

            <FilterChipText name="Given name" :filterValue="filterGivenName" 
              @closed="filterGivenName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterGivenName = val; debouncedRefreshItems() }">
            </FilterChipText>

            <FilterChipText name="Family name" :filterValue="filterFamilyName" 
              @closed="filterFamilyName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterFamilyName = val; debouncedRefreshItems() }">
            </FilterChipText>

          </v-chip-group>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.hire_date`]="{ item }">
      {{ useDateFormat(item.hire_date, 'DD MMM YYYY').value }}
    </template>

    <template v-slot:[`item.actions`]="{ item }">
      <v-btn icon flat size="small" @click="editID = item.id; showDialog = true">
        <v-icon color="primary" icon="mdi-pencil"></v-icon>
      </v-btn>
      <v-btn icon flat size="small" v-tooltip="'Territories'" :to="{ name: 'Employee detail', params: { id: item.id }}">
        <v-icon color="primary" icon="mdi-land-plots"></v-icon>
      </v-btn>
    </template>

    <template v-slot:[`bottom`]="{}">
      <DtFooter :totalItemsIsEstimate="totalItemsIsEstimate" :totalItemsEstimated="totalItemsEstimated"></DtFooter>
    </template>

  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onBeforeMount, onMounted } from 'vue'
import { VDataTable } from 'vuetify/components'
import { useFetchDt } from '@/composables/fetch'
import { type Employee } from '@/types/hr'
import { getTextFilterUrlParam, itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'
import { useDateFormat, useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  title?: string
}>()

var headers = [
  { title: 'Title', key: 'title' },  
  { title: 'Given name', key: 'given_name' },
  { title: 'Family name', key: 'family_name' },
  { title: 'Job title', key: 'job_title' },
  { title: 'Home phone', key: 'home_phone' },
  { title: 'Reports to', key: 'reports_to' },
  { title: 'Hire date', key: 'hire_date' },
  { title: '# territories', key: 'territory_count', align: 'end' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const
const excludedHeaders = ref<string[]>([])
const selectedHeaders = ref()

const baseUrl = '/a/hr/employees'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<Employee[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const editID = ref(0)
const showDialog = ref(false)

const filterFamilyName = ref<string>()
const filterGivenName = ref<string>()
const lsKey = 'employees_dt'

function getFilterStr(): string {
  var ret = ''

  // exclude None
  ret += '&id=!-1'

  ret += getTextFilterUrlParam('family_name', filterFamilyName.value)
  ret += getTextFilterUrlParam('given_name', filterGivenName.value)

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
    'filterFamilyName': filterFamilyName.value,
    'filterGivenName': filterGivenName.value,
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
  if (lsObj['filterFamilyName']) { filterFamilyName.value = lsObj['filterFamilyName'] }
  if (lsObj['filterGivenName']) { filterGivenName.value = lsObj['filterGivenName'] }
  if (lsObj['excludedHeaders']) { excludedHeaders.value = lsObj['excludedHeaders'] }
})

onMounted(() => {
  selectedHeaders.value = headers.filter((v) => !excludedHeaders.value.includes(v.key))
})

</script>
