<template>
  <v-card v-if="item" variant="flat">
    <v-card-title class="pt-6">
      {{ cardTitle }}
    </v-card-title>
    <v-card-text class="pa-5">
      <v-form ref="itemForm">

        <v-row>

          <v-col class="form-col">
            <v-autocomplete label="Product" v-model="item.product_fk" @update:model-value="loadProduct"
              :items="coreStore.productsList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Product is required']"
            ></v-autocomplete>

            <v-text-field label="Quantity" v-model.number="item.quantity" type="number"
              :rules="[(v: number) => v > 0 || 'Quantity is required']"
            ></v-text-field>

            <v-text-field label="Unit price" v-model.number="item.unit_price" type="number" prefix="$"
              :rules="[(v: number) => v > 0 || 'Unit price is required']"
            ></v-text-field>

            <v-text-field label="Discount" v-model.number="item.discount" type="number"
              :rules="[(v: number) => (v >= 0 && v <= 1) || 'Discount must be between 0 and 1']"
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

            <v-btn v-if="props.id !== 0" color="error" class="mt-2" style="float: right;" @click="archiveItem">Archive</v-btn>

          </v-col>
        </v-row>

      </v-form>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import ax from '@/api'
import { useFetch } from '@/composables/fetch'
import { type Product } from '@/types/core'
import { type OrderItem, type OrderItemInput, NewOrderItem, GetOrderItemInputFromItem } from '@/types/sales'
import { useCoreStore } from '@/stores/core'

/*
Composition relationship: order_id is mandatory. Order items are only added in the context of a selected order
*/
const props = defineProps<{
  order_id: number
  order_number: number
  id: number
}>()

const emit = defineEmits<{
  (e: 'archive'): void
  (e: 'cancel'): void
  (e: 'create', newID: number): void
  (e: 'load'): void
  (e: 'update'): void
}>()

const coreStore = useCoreStore()

const saving = ref(false)

const item = ref<OrderItem>()
const baseURL = '/a/sales/order-items'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const cardTitle = computed(() => {
  var ret = 'Order #' + props.order_number + ': '
  props.id !== 0 ? ret += 'detail id ' + props.id : ret += 'new item'
  return ret
})

function archiveItem() {
  if (!confirm('Are you sure?')) {
    return
  }

  ax.delete(itemURL + '/archive')
    .then(() => {
      coreStore.loadCategoriesList()
      emit('archive')
    })
    .catch() // handled by interceptor
}

function loadItem() {
  useFetch(itemURL, item, () => { emit('load') })
}

function loadProduct() {
  if (!item.value!.product_fk) {
    return
  }

  ax.get('/a/core/products/' + item.value!.product_fk)
    .then(response => {
      var p: Product
      p = response.data.data
      item.value!.unit_price = p.unit_price
    })
    .catch() // handled by interceptor
}

async function saveItem() {
  const {valid} = await itemForm.value?.validate()
  if (!valid) {
    return
  }

  saving.value = true

  var saveItem: OrderItemInput = GetOrderItemInputFromItem(item.value!)

  if (props.id !== 0) {
    await ax.put(itemURL, saveItem)
      .then(() => {
        showSaved.value = true
        setTimeout(() => { showSaved.value = false }, import.meta.env.VITE_FADE_MS)
        loadItem()
        emit('update')
      })
      .catch() // handled by interceptor
      .finally(() => saving.value = false)
    return
  }

  await ax.post(baseURL, saveItem)
    .then(response => {
      saveBtnLabel.value = 'Save'
      showSaved.value = true
      setTimeout(() => { showSaved.value = false }, import.meta.env.VITE_FADE_MS)
      emit('create', response.data.data)
    })
    .catch() // handled by interceptor
    .finally(() => saving.value = false)
}

onMounted(() => {
  if (props.id !== 0) {
    loadItem()
  } else {
    saveBtnLabel.value = 'Create'
    item.value = NewOrderItem(props.order_id)
  }
})
</script>
