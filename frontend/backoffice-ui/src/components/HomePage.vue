<template>
  <h1 class="text-h2 font-weight-bold projectFont mb-4">{{ appStore.projectTitle }}</h1>
  <v-card class="mx-auto text-center" max-width="1000">
    <v-card-text>
      <v-sheet color="rgba(0, 0, 0, .12)">
        <v-sparkline 
          :model-value="chartValues"
          :gradient="gradient"
          padding="24"
          stroke-linecap="round"
          smooth
        >
          <template v-slot:label="item">
            {{ item.value }}
          </template>
        </v-sparkline>
      </v-sheet>
    </v-card-text>
  <v-card-text>
    <div class="text-h4 font-weight-thin">
      Order values last 10 weeks ($ k)
    </div>
  </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import ax from '@/api'
import { type OrderValueLatestWeeks } from '@/types/sales'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const chartValues = ref<number[]>([])

const gradient = ['#1feaea', '#ffd200', '#f72047']

function loadChartValues() {
  ax.get('/a/sales/order-value-latest-weeks')
    .then(response => {
      var ov: OrderValueLatestWeeks[]
      ov = response.data.data
      ov.forEach((v) => {
        chartValues.value.push(v.total_value)
      })
    })
    .catch() // handled by interceptor
}

onMounted(() => {
  loadChartValues()
})
</script>
