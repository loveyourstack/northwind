<template>
  <v-card v-if="item" variant="flat">
    <v-card-title class="pt-6">
      {{ cardTitle }}
    </v-card-title>
    <v-card-text class="pa-5">
      <v-form ref="itemForm">

        <v-row>

          <v-col class="form-col">
            <v-text-field label="Name" v-model="item.name"
              :rules="[(v: string) => !!v || 'Name is required']"
            ></v-text-field>

            <v-text-field label="Code" v-model="item.code"
              :rules="[(v: string) => !!v || 'Code is required']"
            ></v-text-field>

            <v-autocomplete label="Region" v-model="item.region"
              :items="salesStore.regions"
              :rules="[(v: string) => !!v || 'Region is required']"
            ></v-autocomplete>
            
            <v-autocomplete label="Salesman" v-model="item.salesman_fk"
              :items="hrStore.employeesList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Salesman is required']"
            ></v-autocomplete>
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
import { Territory, TerritoryInput, NewTerritory, GetTerritoryInputFromItem } from '@/types/sales'
import { fadeMs } from '@/composables/form'
import { useHRStore } from '@/stores/hr'
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

const hrStore = useHRStore()
const salesStore = useSalesStore()

const saving = ref(false)

const item = ref<Territory>()
const baseURL = '/a/sales/territories'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const cardTitle = computed(() => {
  return props.id !== 0 ? item.value!.name : 'New territory'
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

  var saveItem: TerritoryInput = GetTerritoryInputFromItem(item.value!)

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
    item.value = NewTerritory()
  }
})
</script>
