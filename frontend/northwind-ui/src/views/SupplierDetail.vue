<template>
  <v-container class="">
    <v-responsive class="">
      <v-row>
        <v-col>

          <v-tabs v-model="selectedTab" color="deep-purple-accent-4">

            <v-tab value="details">Details</v-tab>

            <v-tab value="products"
              v-if="props.id !== 0"
              @click="setVisited('products')"
            >Products
            </v-tab>

          </v-tabs>

          <v-window v-model="selectedTab" class="mt-2">

            <!-- details tab not lazy loaded, others are -->
            <v-window-item value="details">
              <SupplierForm :id="props.id"
                @cancel="router.back"
                @create="router.push({ name: 'Suppliers' })"
                @delete="router.push({ name: 'Suppliers' })"
                @load="(name) => { itemName = name }"
              ></SupplierForm>
            </v-window-item>

            <v-window-item v-if="props.id !== 0 && visitedTabs.includes('products')" value="products">
              <v-card>
                <v-card-text class="pt-0">
                  <ProductTable :supplier_id="props.id" :title="itemName + ' > Products'" />
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
import SupplierForm from '@/components/forms/SupplierForm.vue'
import ProductTable from '@/components/tables/ProductTable.vue'

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
