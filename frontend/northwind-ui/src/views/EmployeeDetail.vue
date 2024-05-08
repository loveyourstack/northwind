<template>
  <v-container fluid class="cockpit">
    <v-responsive class="">
      <v-row>
        <v-col>
          <v-tabs style="background: rgb(var(--v-theme-surface));" v-model="selectedTab" class="rounded">
            <v-tab class="ml-2" value="details">Details</v-tab>
          </v-tabs>

          <v-window v-model="selectedTab">

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
