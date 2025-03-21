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

            <v-text-field label="Description" v-model="item.description"
              :rules="[(v: string) => !!v || 'Description is required']"
            ></v-text-field>

            <v-color-picker label="Color" v-model="item.color_hex" :mode="'hex'" hide-inputs show-swatches
            ></v-color-picker>

            <v-btn size="small" icon class="mt-2" @click="item.color_hex = ''">
              <v-icon class="float-right" icon="mdi-close"></v-icon>
            </v-btn>
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
import { Category, CategoryInput, NewCategory, GetCategoryInputFromItem } from '@/types/core'
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

const item = ref<Category>()
const baseURL = '/a/core/categories'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const cardTitle = computed(() => {
  return props.id !== 0 ? item.value!.name : 'New category'
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
  useFetch(itemURL, item, () => { emit('load', item.value!.name) })
}

async function saveItem() {
  const {valid} = await itemForm.value?.validate()
  if (!valid) {
    return
  }

  saving.value = true

  var saveItem: CategoryInput = GetCategoryInputFromItem(item.value!)

  if (props.id !== 0) {
    await ax.put(itemURL, saveItem)
      .then(() => {
        coreStore.loadCategoriesList()
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
      coreStore.loadCategoriesList()
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
    item.value = NewCategory()
  }
})
</script>
