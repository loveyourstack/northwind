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
    @update:options="loadItems"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Settings' }}</div>
          </div>

          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn class="float-right mr-5" icon="mdi-dots-vertical" v-bind="props"></v-btn>
            </template>
            <v-list>

              <v-list-item prepend-icon="mdi-selection-ellipse-remove">
                <v-list-item-title class="clickable" @click="resetTable()">Reset table</v-list-item-title>
              </v-list-item>

            </v-list>
          </v-menu>

        </v-col>
      </v-row>

      <v-row no-gutters class="mb-3">
        <v-col>
          <span class="text-body-1">{{ pgVersion }}</span>
        </v-col>
      </v-row>

      <v-row no-gutters class="mb-3">
        <v-col>
          <v-chip-group column>

            <FilterChipText name="Name search" :filterValue="filterName" 
              @closed="filterName = ''; refreshItems()" 
              @updated="(val: string | undefined) => { filterName = val; debouncedRefreshItems() }">
            </FilterChipText>

          </v-chip-group>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.changed`]="{ item }">
      <span v-if="item.changed">{{ item.boot_val ? item.boot_val : '(empty)' }}</span>
    </template>

    <template v-slot:[`bottom`]="{}">
      <DtFooter :totalItemsIsEstimate="totalItemsIsEstimate" :totalItemsEstimated="totalItemsEstimated"></DtFooter>
    </template>

  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount, onMounted } from 'vue'
import { VDataTable } from 'vuetify/components'
import { type PgSetting } from '@/types/tech'
import { useFetch, useFetchDt } from '@/composables/fetch'
import { getTextFilterUrlParam, processURIOptions } from '@/functions/datatable'
import { useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  title?: string
}>()

var headers = [
  { title: 'Name', key: 'name' },
  { title: 'Setting', key: 'setting' },
  { title: 'Unit', key: 'unit' },
  { title: 'Description', key: 'short_desc' },
  { title: 'Boot val, if different', key: 'changed' },
] as const

const baseUrl = '/a/tech/pg-settings'

const items = ref<PgSetting[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const pgVersion = ref<string>('')

const filterName = ref<string>()

const lsKey = 'pg_settings_dt'

function getFilterStr(): string {
  var ret = ''
  ret += getTextFilterUrlParam('name', filterName.value)
  return ret
}

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {
  var myURL = processURIOptions(baseUrl, options)
  myURL += getFilterStr()
  useFetchDt(myURL, items, totalItems, totalItemsIsEstimate, totalItemsEstimated)
}

function loadPgVersion() {
  useFetch('/a/tech/pg-version', pgVersion)
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

watch([itemsPerPage, search, sortBy], () => {

  let lsObj = {
    'itemsPerPage': itemsPerPage.value,
    'sortBy': sortBy.value,
    'filterName': filterName.value,
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
  if (lsObj['filterName']) { filterName.value = lsObj['filterName'] }
})

onMounted(() => {
  loadPgVersion()
})

</script>
