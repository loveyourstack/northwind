<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <EmployeeForm :id="props.id"
                @cancel="router.back"
                @create="router.push({ name: 'Employees' })"
                @delete="router.push({ name: 'Employees' })"
              ></EmployeeForm>

            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'
import EmployeeForm from '@/components/forms/EmployeeForm.vue'

const props = defineProps<{
  id: number
}>()

const router = useRouter()

</script>
