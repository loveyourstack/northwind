<template>
  <v-container fluid class="cockpit">
    <v-responsive class="">
      <v-row>
        <v-col>
          <v-tabs style="background: rgb(var(--v-theme-surface));" v-model="selectedTab" class="rounded">

            <v-tab class="ml-2" value="details">Details</v-tab>

            <v-tab value="orders"
              v-if="props.id !== 0"
              @click="setVisited('orders')"
            >Orders
            </v-tab>
          </v-tabs>

          <v-window v-model="selectedTab">

            <v-window-item value="details">
              <CustomerForm :id="props.id"
                @cancel="router.back"
                @create="router.push({ name: 'Customers' })"
                @delete="router.push({ name: 'Customers' })"
                @load="(name) => { itemName = name }"
              ></CustomerForm>
            </v-window-item>
            
            <v-window-item v-if="props.id !== 0 && visitedTabs.includes('orders')" value="orders">
              <v-card>
                <v-card-text class="pt-0">
                  <OrderTable :customer_id="props.id" :title="itemName + ' > Orders'" />
                  <v-row class="pt-4">
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
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import CustomerForm from '@/components/forms/CustomerForm.vue'
import OrderTable from '@/components/tables/OrderTable.vue'

const props = defineProps<{
  id: number
}>()

const router = useRouter()

const selectedTab = ref('details')
const visitedTabs = ref<string[]>([]) // allows for lazy loading of tab content
const itemName = ref('')
const lsKey = 'customer_detail'

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
