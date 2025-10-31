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
            <div class="dt-title">{{ props.title ? props.title : 'Queries' }}</div>
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

          <v-btn class="float-right mt-1" :loading="running" @click="startLongRunningQuery">Start long-running query</v-btn>

        </v-col>
      </v-row>

    </template>

    <template v-slot:[`item.query_start`]="{ item }">
      {{ useDateFormat(item.query_start, 'DD MMM YYYY HH:mm:ss').value }}
    </template>

    <template v-slot:[`bottom`]="{}">
      <DtFooter :totalItemsIsEstimate="totalItemsIsEstimate" :totalItemsEstimated="totalItemsEstimated"></DtFooter>
    </template>

  </v-data-table-server>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'
import ax from '@/api'
import { useDateFormat } from '@vueuse/core'
import { VDataTable } from 'vuetify/components'
import { type PgQuery } from '@/types/tech'
import { useFetchDt } from '@/composables/fetch'
import { processURIOptions } from '@/functions/datatable'

const props = defineProps<{
  title?: string
}>()

var headers = [
  { title: 'Pid', key: 'pid' },  
  { title: 'State', key: 'state' },
  { title: 'Started at', key: 'query_start' },
  { title: 'Application', key: 'application_name' },
  { title: 'User', key: 'usename' },
  { title: 'IP', key: 'client_addr' },
  { title: 'Query', key: 'query' },
] as const

const baseUrl = '/a/tech/pg-queries'

const items = ref<PgQuery[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const running = ref(false)

const lsKey = 'pg_queries_dt'

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

function startLongRunningQuery() {
  var myURL = '/a/tech/long-running-query'
  running.value = true
  ax.get(myURL)
    .then(resp => {
      alert('done')
    })
    .catch() // handled by interceptor
    .finally(() => { running.value = false })
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
