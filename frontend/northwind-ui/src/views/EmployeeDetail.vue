<template>
  <v-container class="">
    <v-responsive class="">
      <v-row>
        <v-col>
          <v-tabs v-model="selectedTab" color="deep-purple-accent-4">
            <v-tab value="details">Details</v-tab>
          </v-tabs>

          <v-window v-model="selectedTab" class="mt-2">

            <v-window-item value="details">
              <EmployeeForm :id="props.id"
                @cancel="router.back"
                @create="router.push({ name: 'Employees' })"
                @delete="router.push({ name: 'Employees' })"
              ></EmployeeForm>
            </v-window-item>
            
          </v-window>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import EmployeeForm from '@/components/forms/EmployeeForm.vue'

const props = defineProps<{
  id: number
}>()

const router = useRouter()

const selectedTab = ref('details')
</script>
