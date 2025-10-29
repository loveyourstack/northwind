<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <v-tabs v-model="selectedTab" class="rounded">

                <v-tab class="ml-2" value="details">Details</v-tab>

                <v-tab value="products"
                  v-if="props.id !== 0"
                  @click="setVisited('products')"
                >Products
                </v-tab>

              </v-tabs>

              <v-window v-model="selectedTab">

                <!-- details tab not lazy loaded, others are -->
                <v-window-item value="details">
                  <CoreSupplierForm :id="props.id"
                    @cancel="router.back"
                    @create="router.push({ name: 'Suppliers' })"
                    @delete="router.push({ name: 'Suppliers' })"
                    @load="(name: string) => { itemName = name }"
                  ></CoreSupplierForm>
                </v-window-item>

                <v-window-item v-if="props.id !== 0 && visitedTabs.includes('products')" value="products">
                  <v-card>
                    <v-card-text class="pa-0">
                      <CoreProductTable :supplier_id="props.id" :title="itemName + ' > Products'" />
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

const props = defineProps<{
  id: number
}>()

const router = useRouter()

const selectedTab = ref('details')
const visitedTabs = ref<string[]>([]) // allows for lazy loading of tab content
const itemName = ref('')
const lsKey = 'supplier_detail'

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
