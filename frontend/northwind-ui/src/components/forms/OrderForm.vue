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

            <v-autocomplete label="Customer" v-model="item.customer_fk"
              :items="salesStore.customersList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Customer is required']"
            ></v-autocomplete>

            <v-dialog v-model="showOrderDateDp">
              <template #activator="{ props }">
                <v-text-field v-bind="props" label="Order date" prepend-icon="mdi-calendar" readonly 
                  :model-value="item.order_date_d ? useDateFormat(item.order_date_d, 'DD MMM YYYY').value : undefined"
                  :rules="[(v: string) => !!v || 'Order date is required']"
                ></v-text-field>
              </template>
              <template #default>
                <v-date-picker color="primary" v-model="item.order_date_d" @update:model-value="showOrderDateDp = false"></v-date-picker>
              </template>
            </v-dialog>

            <v-autocomplete label="Salesman" v-model="item.salesman_fk"
              :items="hrStore.employeesList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Salesman is required']"
            ></v-autocomplete>

            <v-dialog v-model="showRequiredDateDp">
              <template #activator="{ props }">
                <v-text-field v-bind="props" label="Required date" prepend-icon="mdi-calendar" readonly 
                  :model-value="item.required_date_d ? useDateFormat(item.required_date_d, 'DD MMM YYYY').value : undefined"
                  :rules="[(v: string) => !!v || 'Required date is required']"
                ></v-text-field>
              </template>
              <template #default>
                <v-date-picker color="primary" v-model="item.required_date_d" @update:model-value="showRequiredDateDp = false"></v-date-picker>
              </template>
            </v-dialog>

            <v-autocomplete label="Shipper" v-model="item.shipper_fk"
              :items="salesStore.shippersList" item-title="company_name" item-value="id"
              :rules="[(v: number) => !!v || 'Shipper is required']"
            ></v-autocomplete>

            <v-checkbox v-if="props.id !== 0" label="Shipped" v-model="item.is_shipped" @update:model-value="item.shipped_date_d = useNow().value"></v-checkbox>

            <v-dialog v-if="item.is_shipped" v-model="showShippedDateDp">
              <template #activator="{ props }">
                <v-text-field v-bind="props" label="Shipped date" prepend-icon="mdi-calendar" readonly 
                  :model-value="item.shipped_date_d ? useDateFormat(item.shipped_date_d, 'DD MMM YYYY').value : undefined"
                  :rules="[(v: string) => !!v || 'Shipped date is required']"
                ></v-text-field>
              </template>
              <template #default>
                <v-date-picker color="primary" v-model="item.shipped_date_d" @update:model-value="showShippedDateDp = false"></v-date-picker>
              </template>
            </v-dialog>

            <v-text-field label="Freight cost" v-model.number="item.freight_cost" type="number" prefix="$" disabled
            ></v-text-field>

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
                :items="commonStore.activeCountriesList" item-title="name" item-value="id"
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

            <v-btn v-if="props.id !== 0" color="error" class="mt-2" style="float: right;" @click="deleteItem">Delete</v-btn>

          </v-col>
        </v-row>

      </v-form>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useDateFormat, useNow } from '@vueuse/core'
import ax from '@/api'
import { Customer, Order, OrderInput, NewOrder, GetOrderInputFromItem } from '@/types/sales'
import { fadeMs } from '@/composables/form'
import { useCommonStore } from '@/stores/common'
import { useHRStore } from '@/stores/hr'
import { useSalesStore } from '@/stores/sales'

const props = defineProps<{
  id: number
}>()

const emit = defineEmits<{
  (e: 'cancel'): void
  (e: 'create', newID: number): void
  (e: 'delete'): void
  (e: 'load', order_number: number): void
  (e: 'update'): void
}>()

const commonStore = useCommonStore()
const hrStore = useHRStore()
const salesStore = useSalesStore()

const saving = ref(false)

const item = ref<Order>()
const baseURL = '/a/sales/orders'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const showOrderDateDp = ref(false)
const showRequiredDateDp = ref(false)
const showShippedDateDp = ref(false)

const cardTitle = computed(() => {
  return props.id !== 0 ? 'Order #' + item.value!.order_number : 'New Order'
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

function fillFromCustomer() {
  if (!item.value?.customer_fk) {
    return
  }

  ax.get('/a/sales/customers/' + item.value!.customer_fk)
    .then(response => {
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
  ax.get(itemURL)
    .then(response => {
      item.value = response.data.data

      // assign date objects from YYYY-MM-DD strings
      item.value!.order_date_d = new Date(item.value!.order_date!)
      item.value!.required_date_d = new Date(item.value!.required_date!)
      if (item.value!.shipped_date) {
        item.value!.shipped_date_d = new Date(item.value!.shipped_date!)
      }

      emit('load', item.value!.order_number)
    })
    .catch() // handled by interceptor
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
      var newItem: Order = response.data.data
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
    item.value = NewOrder()
  }
})
</script>
