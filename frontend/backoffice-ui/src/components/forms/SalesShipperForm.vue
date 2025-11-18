<template>
  <v-card v-if="item" variant="flat">
    <v-card-title class="pt-6">
      {{ cardTitle }}
    </v-card-title>
    <v-card-text class="pa-5">
      <v-form ref="itemForm">

        <v-row>
          <v-col class="form-col">

            <v-text-field label="Company name" v-model="item.company_name"
              :rules="[(v: string) => !!v || 'Company name is required']"
            ></v-text-field>

            <v-text-field label="Phone" v-model="item.phone"
              :rules="[(v: string) => !!v || 'Phone is required']"
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
import { useFetch } from '@/composables/fetch'
import { callDelete } from '@/functions/form'
import { type Shipper, type ShipperInput, NewShipper, GetShipperInputFromItem } from '@/types/sales'
import { useSalesStore } from '@/stores/sales'

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

const salesStore = useSalesStore()

const saving = ref(false)

const item = ref<Shipper>()
const baseURL = '/a/sales/shippers'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const cardTitle = computed(() => {
  return props.id !== 0 ? item.value!.company_name : 'New shipper'
})

function deleteItem() {
  callDelete(itemURL, () => { salesStore.loadShippersList(); emit('delete') })
}

function loadItem() {
  useFetch(itemURL, item, () => { emit('load', item.value!.company_name) })
}

async function saveItem() {
  const {valid} = await itemForm.value?.validate()
  if (!valid) {
    return
  }

  saving.value = true

  var saveItem: ShipperInput = GetShipperInputFromItem(item.value!)

  if (props.id !== 0) {
    await ax.put(itemURL, saveItem)
      .then(() => {
        salesStore.loadShippersList()
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
      salesStore.loadShippersList()
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
    item.value = NewShipper()
  }
})
</script>
