<template>
  <v-card v-if="item" variant="flat">
    <v-card-title class="pt-6">
      {{ cardTitle }}
    </v-card-title>
    <v-card-subtitle v-if="props.id !== 0">
      {{ item.order_item_count }} item(s) with total value: ${{ item.order_value }}
    </v-card-subtitle>
    <v-card-text class="pa-5">
      <v-form ref="itemForm">

        <v-row>

          <v-col cols="12" md="6" class="form-col">

            <v-autocomplete label="Customer" v-model="item.customer_fk" :disabled="props.id === 0 && props.customer_id > 0"
              :items="salesStore.customersList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Customer is required']"
            ></v-autocomplete>

            <DateTextField :dateVal="item.order_date" label="Order date" clearable @cleared="item.order_date = undefined"
              @updated="(val: string | undefined) => { item!.order_date = val }"
              :rules="[(v: string) => !!v || 'Order date is required']"
            ></DateTextField>

            <v-autocomplete label="Salesman" v-model="item.salesman_fk"
              :items="hrStore.employeesList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Salesman is required']"
            ></v-autocomplete>

            <DateTextField :dateVal="item.required_date" label="Required date" clearable @cleared="item.required_date = undefined"
              @updated="(val: string | undefined) => { item!.required_date = val }"
              :rules="[(v: string) => !!v || 'Required date is required']"
            ></DateTextField>

            <v-autocomplete label="Shipper" v-model="item.shipper_fk"
              :items="salesStore.shippersList" item-title="company_name" item-value="id"
              :rules="[(v: number) => !!v || 'Shipper is required']"
            ></v-autocomplete>

            <v-btn v-if="item.order_item_count" class="mb-4" density="compact" @click="calculateFreightCost">Calculate freight cost</v-btn>

            <v-text-field v-if="item.order_item_count" label="Freight cost" v-model.number="item.freight_cost" type="number" prefix="$" disabled
            ></v-text-field>

            <v-checkbox v-if="props.id !== 0" label="Shipped" v-model="item.is_shipped" @update:model-value="item.shipped_date = useDateFormat(useNow().value, 'YYYY-MM-DD').value"></v-checkbox>

            <DateTextField v-if="item.is_shipped" :dateVal="item.shipped_date" label="Shipped date" clearable @cleared="item.shipped_date = undefined"
              @updated="(val: string | undefined) => { item!.shipped_date = val }"
              :rules="[(v: string) => !!v || 'Shipped date is required']"
            ></DateTextField>

          </v-col>
          <v-col cols="12" md="6" class="form-col">

            <fieldset class="pa-4 fs-std mb-4">

              <legend class="pl-2 pr-2">Delivery address</legend>

              <v-btn :disabled="!item.customer_fk" class="mb-4" density="compact" @click="fillFromCustomer">Fill from customer</v-btn>

              <v-text-field label="Company" v-model="item.dest_company_name"
              ></v-text-field>

              <v-text-field label="Address" v-model="item.dest_address"
              ></v-text-field>

              <v-text-field label="City" v-model="item.dest_city"
              ></v-text-field>

              <v-text-field label="State" v-model="item.dest_state"
              ></v-text-field>

              <v-text-field label="Postal code" v-model="item.dest_postal_code"
              ></v-text-field>

              <v-autocomplete label="Country" v-model="item.dest_country_fk"
                :items="coreStore.activeCountriesList" item-title="name" item-value="id"
                :rules="[(v: number) => !!v || 'Country is required']"
              ></v-autocomplete>

            </fieldset>

          </v-col>
        </v-row>

        <v-row align="center">
          <v-col cols="12">

            <v-btn icon class="mr-4 mb-1 ml-1" @click="$emit('cancel')">
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
import { useNow, useDateFormat } from '@vueuse/core'
import ax from '@/api'
import { useFetch } from '@/composables/fetch'
import { type Customer, type Order, type OrderInput, NewOrder, GetOrderInputFromItem } from '@/types/sales'
import { useCoreStore } from '@/stores/core'
import { useHRStore } from '@/stores/hr'
import { useSalesStore } from '@/stores/sales'

const props = defineProps<{
  id: number
  customer_id: number
}>()

const emit = defineEmits<{
  (e: 'archive'): void
  (e: 'cancel'): void
  (e: 'create', newID: number): void
  (e: 'load', order_number: number): void
  (e: 'update'): void
}>()

const coreStore = useCoreStore()
const hrStore = useHRStore()
const salesStore = useSalesStore()

const saving = ref(false)

const item = ref<Order>()
const baseURL = '/a/sales/orders'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const cardTitle = computed(() => {
  return props.id !== 0 ? 'Order #' + item.value!.order_number : 'New Order'
})

function archiveItem() {
  if (!confirm('Are you sure?')) {
    return
  }

  ax.delete(itemURL + '/archive')
    .then(() => {
      emit('archive')
    })
    .catch() // handled by interceptor
}

function calculateFreightCost() {
  if (!item.value?.order_value) {
    return
  }

  // enter dummy cost based on order value, round to 2 dp
  item.value!.freight_cost = Math.round(item.value!.order_value * (Math.random() + 0.1) * 100) / 100
}

function fillFromCustomer() {
  if (!item.value?.customer_fk) {
    return
  }

  ax.get('/a/sales/customers/' + item.value!.customer_fk)
    .then((response: any) => {
      var cust:Customer = response.data.data
      item.value!.dest_company_name = cust.company_name
      item.value!.dest_address = cust.address ? cust.address : ''
      item.value!.dest_city = cust.city ? cust.city : ''
      item.value!.dest_state = cust.state
      item.value!.dest_postal_code = cust.postal_code
      item.value!.dest_country_fk = cust.country_fk
    })
    .catch() // handled by interceptor
}

function loadItem() {
  useFetch(itemURL, item, () => emit('load', item.value!.order_number) )
}

async function saveItem() {
  const {valid} = await itemForm.value?.validate()
  if (!valid) {
    return
  }

  saving.value = true

  var saveItem: OrderInput = GetOrderInputFromItem(item.value!)

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
    .then((response: any) => {
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
    item.value = NewOrder()
    
    if (props.customer_id) { item.value.customer_fk = props.customer_id }
  }
})
</script>
