<template>
  <v-card v-if="item" variant="flat">
    <v-card-title class="pt-6">
      {{ cardTitle }}
    </v-card-title>
    <v-card-text class="pa-5">
      <v-form ref="itemForm">

        <v-row>

          <v-col cols="12" md="6" class="form-col">

            <v-text-field label="Code" v-model="item.code"
              :rules=codeRules
            ></v-text-field>

            <v-text-field label="Company name" v-model="item.company_name"
              :rules="[(v: string) => !!v || 'Company name is required']"
            ></v-text-field>

            <v-text-field label="Contact name" v-model="item.contact_name"
              :rules="[(v: string) => !!v || 'Contact name is required']"
            ></v-text-field>

            <v-text-field label="Contact title" v-model="item.contact_title"
            ></v-text-field>

            <v-text-field label="Phone" v-model="item.phone"
            ></v-text-field>

          </v-col>
          <v-col cols="12" md="6" class="form-col">

            <v-text-field label="Address" v-model="item.address"
            ></v-text-field>

            <v-text-field label="City" v-model="item.city"
            ></v-text-field>

            <v-text-field label="State" v-model="item.state"
            ></v-text-field>

            <v-text-field label="Postal code" v-model="item.postal_code"
            ></v-text-field>

            <v-autocomplete label="Country" v-model="item.country_fk"
              :items="commonStore.activeCountriesList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Country is required']"
            ></v-autocomplete>

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
import ax from '@/api'
import { Customer, CustomerInput, NewCustomer, GetCustomerInputFromItem } from '@/types/sales'
import { fadeMs } from '@/composables/form'
import { useCommonStore } from '@/stores/common'

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

const commonStore = useCommonStore()

const saving = ref(false)

const item = ref<Customer>()
const baseURL = '/a/sales/customers'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const cardTitle = computed(() => {
  return props.id !== 0 ? item.value!.name : 'New Customer'
})

const codeRules = [
  (v: any) => !!v || 'Code is required',
  (v: any) => (v.length === 5) || 'Code must be 5 letters',
  (v: any) => (v === v.toUpperCase()) || 'Code must be uppercase',
]

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

  var saveItem: CustomerInput = GetCustomerInputFromItem(item.value!)

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
      var newItem: Customer = response.data.data
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
    item.value = NewCustomer()
  }
})
</script>
