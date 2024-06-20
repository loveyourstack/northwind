<template>
  <v-card v-if="item" variant="flat">
    <v-card-title class="pt-6">
      {{ cardTitle }}
    </v-card-title>
    <v-card-text class="pa-5">
      <v-form ref="itemForm">

        <v-row>

          <v-col class="form-col">
            <v-autocomplete label="Product" v-model="item.product_fk"
              :items="coreStore.productsList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Product is required']"
            ></v-autocomplete>

            <v-text-field label="Quantity" v-model.number="item.quantity"
              :rules="[(v: string) => !!v || 'Quantity is required']"
            ></v-text-field>

            <v-text-field label="Unit price" v-model.number="item.unit_price"
              :rules="[(v: string) => !!v || 'Description is required']"
            ></v-text-field>

            <v-text-field label="Discount" v-model.number="item.discount"
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
import { OrderDetail, OrderDetailInput, NewOrderDetail, GetOrderDetailInputFromItem } from '@/types/sales'
import { fadeMs } from '@/composables/form'
import { useCoreStore } from '@/stores/core'

const props = defineProps<{
  order_id: number
  id: number
}>()

const emit = defineEmits<{
  (e: 'cancel'): void
  (e: 'create', newID: number): void
  (e: 'delete'): void
  (e: 'load'): void
  (e: 'update'): void
}>()

const coreStore = useCoreStore()

const saving = ref(false)

const item = ref<OrderDetail>()
const baseURL = '/a/sales/order-details'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const cardTitle = computed(() => {
  var ret = 'Order id: ' + props.order_id + ': '
  props.id !== 0 ? ret += 'detail id ' + props.id : ret += 'new detail'
  return ret
})

function deleteItem() {
  if (!confirm('Are you sure?')) {
    return
  }

  ax.delete(itemURL)
    .then(() => {
      coreStore.loadCategoriesList()
      emit('delete')
    })
    .catch() // handled by interceptor
}

function loadItem() {
  ax.get(itemURL)
    .then(response => {
      item.value = response.data.data
      emit('load')
    })
    .catch() // handled by interceptor
}

async function saveItem() {
  const {valid} = await itemForm.value?.validate()
  if (!valid) {
    return
  }

  saving.value = true

  var saveItem: OrderDetailInput = GetOrderDetailInputFromItem(item.value!)

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
      saveBtnLabel.value = 'Save'
      showSaved.value = true
      setTimeout(() => { showSaved.value = false }, fadeMs)
      var newItem: OrderDetail = response.data.data
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
    item.value = NewOrderDetail(props.order_id)
  }
})
</script>
