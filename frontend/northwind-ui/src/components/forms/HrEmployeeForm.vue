<template>
  <v-card v-if="item" variant="flat">
    <v-card-title class="pt-6">
      {{ cardTitle }}
    </v-card-title>
    <v-card-text class="pa-5">
      <v-form ref="itemForm">

        <v-row>

          <v-col cols="12" md="6" class="form-col">

            <v-text-field label="Name" v-model="item.name" hint="If possible, first name initial and full last name"
              :rules="[(v: string) => !!v || 'Name is required']"
            ></v-text-field>

            <v-text-field label="Title" v-model="item.title"
              :rules="[(v: string) => !!v || 'Title is required']"
            ></v-text-field>

            <v-text-field label="First name" v-model="item.first_name"
              :rules="[(v: string) => !!v || 'First name is required']"
            ></v-text-field>

            <v-text-field label="Last name" v-model="item.last_name"
              :rules="[(v: string) => !!v || 'Last name is required']"
            ></v-text-field>

            <v-text-field label="Date of birth" type="date" v-model="item.date_of_birth"
              :rules="[(v: string) => !!v || 'Date of birth is required']"
            ></v-text-field>

            <v-autocomplete label="Reports to" v-model="item.reports_to_fk"
              :items="hrStore.employeesList" item-title="name" item-value="id"
              :rules="[(v: number) => !!v || 'Reports to is required']"
            ></v-autocomplete>

            <v-text-field label="Job title" v-model="item.job_title"
              :rules="[(v: string) => !!v || 'Job title is required']"
            ></v-text-field>

            <v-text-field label="Hire date" type="date" v-model="item.hire_date"
              :rules="[(v: Date) => !!v || 'Hire date is required']"
            ></v-text-field>

          </v-col>
          <v-col cols="12" md="6" class="form-col">

            <fieldset class="pa-4 fs-std mb-4">

              <legend class="pl-2 pr-2">Address</legend>

              <v-text-field label="Address" v-model="item.address"
                :rules="[(v: string) => !!v || 'Address is required']"
              ></v-text-field>

              <v-text-field label="City" v-model="item.city"
                :rules="[(v: string) => !!v || 'City is required']"
              ></v-text-field>

              <v-text-field label="State" v-model="item.state"
              ></v-text-field>

              <v-text-field label="Postal code" v-model="item.postal_code"
                :rules="[(v: string) => !!v || 'Postal code is required']"
              ></v-text-field>

              <v-autocomplete label="Country" v-model="item.country_fk"
                :items="coreStore.countriesList" item-title="name" item-value="id"
                :rules="[(v: number) => !!v || 'Country is required']"
              ></v-autocomplete>

            </fieldset>

            <v-text-field label="Home phone" v-model="item.home_phone"
            ></v-text-field>

          </v-col>
        </v-row>

        <v-row>
          <v-col>
            <v-textarea label="Notes" v-model="item.notes"
            ></v-textarea>
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
import { useFetch } from '@/composables/fetch'
import { Employee, EmployeeInput, NewEmployee, GetEmployeeInputFromItem } from '@/types/hr'
import { useCoreStore } from '@/stores/core'
import { useHRStore } from '@/stores/hr'

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
const hrStore = useHRStore()

const saving = ref(false)

const item = ref<Employee>()
const baseURL = '/a/hr/employees'
const itemURL = baseURL + '/' + props.id
const itemForm = ref()
const saveBtnLabel = ref('Save')
const showSaved = ref(false)

const showDateDialog = ref(false)

const cardTitle = computed(() => {
  return props.id !== 0 ? item.value!.name : 'New Employee'
})

function deleteItem() {
  if (!confirm('Are you sure?')) {
    return
  }

  ax.delete(itemURL)
    .then(() => {
      hrStore.loadEmployeesList()
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

  var saveItem: EmployeeInput = GetEmployeeInputFromItem(item.value!)

  if (props.id !== 0) {
    await ax.put(itemURL, saveItem)
      .then(() => {
        hrStore.loadEmployeesList()
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
      hrStore.loadEmployeesList()
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
    item.value = NewEmployee()
  }
})
</script>
