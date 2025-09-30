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
            <div class="dt-title">{{ props.title ? props.title : 'Table and index bloat' }}</div>
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

    <template v-slot:[`item.index_bloat`]="{ item }">
      {{ item.index_bloat.toFixed(1) }}
    </template>

    <template v-slot:[`item.index_waste`]="{ item }">
      {{ item.index_waste_pretty }}
    </template>

    <template v-slot:[`item.table_bloat`]="{ item }">
      {{ item.table_bloat.toFixed(1) }}
    </template>

    <template v-slot:[`item.table_waste`]="{ item }">
      {{ item.table_waste_pretty }}
    </template>

    <template v-slot:[`bottom`]="{}">
      <DtFooter :totalItemsIsEstimate="totalItemsIsEstimate" :totalItemsEstimated="totalItemsEstimated"></DtFooter>
    </template>

  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import { PgBloat } from '@/types/tech'
import { useFetchDt } from '@/composables/fetch'
import { processURIOptions } from '@/functions/datatable'
import DtFooter from '@/components/DtFooter.vue'

const props = defineProps<{
  title?: string
}>()

var headers = [
  { title: 'Schema', key: 'table_schema' },
  { title: 'Table', key: 'table_name' },
  { title: 'Tbl bloat', key: 'table_bloat', align: 'end' },
  { title: 'Tbl waste', key: 'table_waste', align: 'end' },
  { title: 'Index', key: 'index_name' },
  { title: 'Idx bloat', key: 'index_bloat', align: 'end' },
  { title: 'Idx waste', key: 'index_waste', align: 'end' },
] as const

const baseUrl = '/a/tech/pg-bloat'

const items = ref<PgBloat[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const lsKey = 'pg_bloat_dt'

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
