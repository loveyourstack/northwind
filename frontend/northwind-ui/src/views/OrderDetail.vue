<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <v-tabs v-model="selectedTab" class="rounded">

                <v-tab class="ml-2" value="details">Details</v-tab>

                <v-tab value="orderItems"
                  v-if="props.id !== 0"
                  @click="setVisited('orderItems')"
                >Order items
                </v-tab>
              </v-tabs>

              <v-window v-model="selectedTab">

                <v-window-item value="details">
                  <OrderForm :id="props.id"
                    @archive="router.push({ name: 'Orders' })"
                    @cancel="router.back"
                    @create="router.push({ name: 'Orders' })"
                    @load=""
                  ></OrderForm>
                </v-window-item>
                
                <v-window-item v-if="props.id !== 0 && visitedTabs.includes('orderItems')" value="orderItems">
                  <v-card>
                    <v-card-text class="pa-0">
                      <OrderItemTable v-if="item" :order_id="props.id" :order_number="item.order_number" :title="'Order #' + item.order_number + ' > Items'" />
                      <v-row class="pt-4 pb-4 pl-4">
                        <v-col>
                          <v-btn icon class="mr-4 mb-1 ml-1" @click="router.back">
                            <v-icon icon="mdi-arrow-left"></v-icon>
                          </v-btn>
                        </v-col>
                      </v-row>
                    </v-card-text>
                  </v-card>
                </v-window-item>

              </v-window>

            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ax from '@/api'
import { useFetch } from '@/composables/fetch'
import { Order } from '@/types/sales'
import OrderForm from '@/components/forms/OrderForm.vue'
import OrderItemTable from '@/components/tables/OrderItemTable.vue';

const props = defineProps<{
  id: number
}>()

const router = useRouter()

const item = ref<Order>()
const baseUrl = '/a/sales/orders'
const itemURL = baseUrl + '/' + props.id

const selectedTab = ref('details')
const visitedTabs = ref<string[]>([]) // allows for lazy loading of tab content
const lsKey = 'order_detail'

function loadItem() {
  useFetch(itemURL, item)
}

function setVisited(tab: string) {
  if (visitedTabs.value.indexOf(tab) === -1) { visitedTabs.value.push(tab) }
}

watch([selectedTab], () => {
  let lsObj = {
    'selectedTab': selectedTab.value
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
})

onMounted(() => {
  // if editing
  if (props.id !== 0) {
    // load order to get order_number
    loadItem()
  }

  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['selectedTab']) { 
    selectedTab.value = lsObj['selectedTab'] 
    visitedTabs.value.push(selectedTab.value)
  }

  // if no ID passed (new item), ensure details tab is selected
  if (props.id === 0 && selectedTab.value !== 'details') {
    selectedTab.value = 'details'
  }
})

</script>