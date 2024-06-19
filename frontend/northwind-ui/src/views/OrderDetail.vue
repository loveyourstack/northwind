<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <v-tabs v-model="selectedTab" class="rounded">

                <v-tab class="ml-2" value="details">Details</v-tab>

                <v-tab value="orderDetails"
                  v-if="props.id !== 0"
                  @click="setVisited('orderDetails')"
                >Order details
                </v-tab>
              </v-tabs>

              <v-window v-model="selectedTab">

                <v-window-item value="details">
                  <OrderForm :id="props.id"
                    @cancel="router.back"
                    @create="router.push({ name: 'Orders' })"
                    @delete="router.push({ name: 'Orders' })"
                    @load="(name) => { itemName = name }"
                  ></OrderForm>
                </v-window-item>
                
                <v-window-item v-if="props.id !== 0 && visitedTabs.includes('orderDetails')" value="orderDetails">
                  <v-card>
                    <v-card-text class="pa-0">
                      <OrderDetailTable :order_id="props.id" />
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
import OrderForm from '@/components/forms/OrderForm.vue'
import OrderDetailTable from '@/components/tables/OrderDetailTable.vue';

const props = defineProps<{
  id: number
}>()

const router = useRouter()

const selectedTab = ref('details')
const visitedTabs = ref<string[]>([]) // allows for lazy loading of tab content
const lsKey = 'order_detail'

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