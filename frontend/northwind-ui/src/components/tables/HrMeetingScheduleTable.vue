<template>
  <v-dialog v-model="showDialog" persistent width="auto">
    <HrMeetingScheduleForm :id="editID"
      @cancel="showDialog = false"
      @create="showDialog = false; refreshItems()"
      @delete="showDialog = false; refreshItems()"
      @update="showDialog = false; refreshItems()"
    ></HrMeetingScheduleForm>
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
      <v-row align="center">
        <v-col>
          <div class="dt-title-block">
            <div class="dt-title">{{ props.title ? props.title : 'Meeting schedule' }}</div>
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

            <FilterChip name="Sched. time >" :filterValue="filterSchedTime" :filterText="filterSchedTime" @closed="filterSchedTime = undefined; refreshItems()">
              <template #menuContent>
                <TimeTextField :timeVal="filterSchedTime" label="Sched. time >"
                  @updated="(val: string | undefined) => { filterSchedTime = val; refreshItems() } "
                ></TimeTextField>
              </template>
            </FilterChip>

          </v-chip-group>
        </v-col>
      </v-row>
    </template>

    <template v-slot:[`item.day`]="{ item }">
      <span v-if="item.day !== 'None'">{{ item.day }}</span>
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
import { ref, computed, watch, onBeforeMount } from 'vue'
import { VDataTable } from 'vuetify/components'
import { useFetchDt } from '@/composables/fetch'
import { type MeetingSchedule } from '@/types/hr'
import { itemsPerPageOptions, processURIOptions } from '@/functions/datatable'
import { fileDownload } from '@/functions/file'

const props = defineProps<{
  title?: string
}>()

// declaring headers like this to avoid typescript issue
var headers = [
  { title: 'Name', key: 'name' },
  { title: 'Frequency', key: 'frequency' },
  { title: 'Day', key: 'day' },
  { title: 'Scheduled time', key: 'scheduled_time', align: 'center' },
  { title: 'Actions', key: 'actions', sortable: false },
] as const

const baseUrl = '/a/hr/meeting-schedule'
const excelDlUrl = computed(() => {
  return baseUrl + '?xformat=excel' + getFilterStr()
}) 

const items = ref<MeetingSchedule[]>([])
const itemsPerPage = ref(10)
const sortBy = ref<any>()
const search = ref('')
const totalItems = ref(0)
const totalItemsIsEstimate = ref(false)
const totalItemsEstimated = ref(0)

const editID = ref(0)
const showDialog = ref(false)

const filterSchedTime = ref<string>() // hh24:mm

const lsKey = 'meeting_schedule_dt'

function getFilterStr(): string {
  var ret = ''

  if (filterSchedTime.value) {
    ret += '&scheduled_time=>eq' + filterSchedTime.value
  }

  return ret
}

function loadItems(options: { page: number, itemsPerPage: number, sortBy: VDataTable['sortBy'] }) {
  var myURL = processURIOptions(baseUrl, options)
  myURL += getFilterStr()

  useFetchDt(myURL, items, totalItems, totalItemsIsEstimate, totalItemsEstimated)
}

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
  'filterSchedTime': filterSchedTime.value,
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
  if (lsObj['filterSchedTime']) { filterSchedTime.value = lsObj['filterSchedTime'] }
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