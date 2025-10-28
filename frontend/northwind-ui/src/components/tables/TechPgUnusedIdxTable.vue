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
    class="dt"
  >
    <template v-slot:[`top`]="{}">
      <v-row align="center">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Unused indexes' }}</div>
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

    </template>

    <template v-slot:[`item.index_size`]="{ item }">
      {{ item.index_size_pretty }}
    </template>

    <template v-slot:[`item.last_idx_scan`]="{ item }">
      <span v-if="new Date(item.last_idx_scan).getFullYear() >= 2002">{{ useDateFormat(item.last_idx_scan, 'DD MMM YYYY HH:mm:ss').value }}</span>
      <span v-else>Never</span>
    </template>

    <template v-slot:[`bottom`]="{}">
      <DtFooter :totalItemsIsEstimate="totalItemsIsEstimate" :totalItemsEstimated="totalItemsEstimated"></DtFooter>
    </template>

  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'
import { useDateFormat, useNow } from '@vueuse/core'
import { VDataTable } from 'vuetify/components'
import { type PgUnusedIdx } from '@/types/tech'
import { useFetchDt } from '@/composables/fetch'
import { processURIOptions } from '@/functions/datatable'
import DtFooter from '@/components/DtFooter.vue'

const props = defineProps<{
  title?: string
}>()

var headers = [
  { title: 'Schema', key: 'table_schema' },
  { title: 'Table', key: 'table_name' },
  { title: 'Index', key: 'index_name' },
  { title: 'Size', key: 'index_size' },
  { title: 'Scans', key: 'index_scans' },
  { title: 'Last scan', key: 'last_idx_scan' },
] as const

const baseUrl = '/a/tech/pg-unused-indexes'

const items = ref<PgUnusedIdx[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const lsKey = 'pg_unused_idxs_dt'

function getFilterStr(): string {
  var ret = ''
  return ret
}

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {
  var myURL = processURIOptions(baseUrl, options)
  myURL += getFilterStr()
  useFetchDt(myURL, items, totalItems, totalItemsIsEstimate, totalItemsEstimated)
}

function resetTable() {
  localStorage.removeItem(lsKey)
  window.location.reload()
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
