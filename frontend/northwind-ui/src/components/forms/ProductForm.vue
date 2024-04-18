<template>
  <v-card v-if="item" variant="flat">
    <v-card-title class="pt-6">
      {{ cardTitle }}
    </v-card-title>
    <v-card-text class="pa-5">
      <v-form ref="itemForm">

        <v-row>

          <v-col cols="12" md="6">

            <v-text-field label="Name" v-model="item.name"
              :rules="[(v: string) => !!v || 'Name is required']"
            ></v-text-field>

            <v-autocomplete label="Category" v-model="item.category_fk"
              :items="coreStore.categoriesList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Category is required']"
            ></v-autocomplete>

            <v-autocomplete label="Supplier" v-model="item.supplier_fk"
              :items="coreStore.suppliersList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Supplier is required']"
            ></v-autocomplete>

            <v-text-field label="Qty / unit" v-model="item.quantity_per_unit"
              :rules="[(v: number) => !!v || 'Qty / unit is required']"
            ></v-text-field>

            <v-autocomplete label="Discontinued" v-model="item.is_discontinued"
              :items="booleanOptions"
            ></v-autocomplete>

          </v-col>
          <v-col cols="12" md="6">

            <v-text-field label="Unit price" v-model.number="item.unit_price"
              :rules="[(v: number) => !!v || 'Unit price is required']"
            ></v-text-field>

            <v-text-field label="Units in stock" v-model.number="item.units_in_stock"
            ></v-text-field>

            <v-text-field label="Units on order" v-model.number="item.units_on_order"
            ></v-text-field>

            <v-text-field label="Reorder level" v-model.number="item.reorder_level"
            ></v-text-field>

          </v-col>
        </v-row>

        <v-row align="center">
          <v-col cols="12">

            <v-btn icon class="mr-4 mb-1 ml-1" @click="emit('cancel')">
              <v-icon icon="mdi-arrow-left"></v-icon>
            </v-btn>

            <v-btn color="primary" :loading="saving" @click="saveItem">{{ saveBtnLabel }}</v-btn>

            <v-fade-transition mode="out-in">
              <v-btn color="green darken-1" variant="text" v-show="showSaved">Saved</v-btn>
            </v-fade-transition>

            <v-btn v-if="props.id !== 0" color="error" class="mt-2" style="float: right;" @click="deleteItem">Delete</v-btn>

          </v-col>
        </v-row>

      </v-form>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import ax from '@/api'
import { Product, ProductInput, NewProduct, GetProductInputFromItem } from '@/types/core'
import { booleanOptions, fadeMs } from '@/composables/form'
import { useCoreStore } from '@/stores/core'

const props = defineProps<{
  id: number
}>()

const emit = defineEmits<{
  (e: 'cancel'): void
  (e: 'create', newID: number): void
  (e: 'delete'): void
  (e: 'load', name: string): void
  (e: 'update'): void
}>()

const coreStore = useCoreStore()

const saving = ref(false)

const item = ref<Product>()
const baseURL = '/a/core/products'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const cardTitle = computed(() => {
  return props.id !== 0 ? item.value!.name : 'New product'
})

function deleteItem() {
  if (!confirm('Are you sure?')) {
    return
  }

  ax.delete(itemURL)
    .then(() => {
      emit('delete')
    })
    .catch() // handled by interceptor
}

function loadItem() {
  ax.get(itemURL)
    .then(response => {
      item.value = response.data.data
      emit('load', item.value!.name)
    })
    .catch() // handled by interceptor
}

async function saveItem() {
  const {valid} = await itemForm.value?.validate()
  if (!valid) {
    return
  }

  saving.value = true

  var saveItem: ProductInput = GetProductInputFromItem(item.value!)

  if (props.id !== 0) {
    await ax.put(itemURL, saveItem)
      .then(() => {
        showSaved.value = true
        setTimeout(() => { showSaved.value = false }, fadeMs)
        loadItem()
        emit('update')
      })
      .catch() // handled by interceptor
      .finally(() => saving.value = false)
    return
  }

  await ax.post(baseURL, saveItem)
    .then(response => {
      // component does not get remounted, need to make Save button changes here
      saveBtnLabel.value = 'Save'
      showSaved.value = true
      setTimeout(() => { showSaved.value = false }, fadeMs)
      var newItem: Product = response.data.data
      emit('create', newItem.id)
    })
    .catch() // handled by interceptor
    .finally(() => saving.value = false)
}

onMounted(() => {
  if (props.id !== 0) {
    loadItem()
  } else {
    saveBtnLabel.value = 'Create'
    item.value = NewProduct()
  }
})
</script>
