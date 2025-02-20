<template>
  <v-row align="center">
    <v-col>
      <v-data-table-footer v-if="props.totalItemsIsEstimate"
        :items-per-page-options="itemsPerPageOptions"
        :page-text="getPageTextEstimated(props.totalItemsEstimated)"
        :show-current-page=true
      ></v-data-table-footer>
      <v-data-table-footer v-else
        :items-per-page-options="itemsPerPageOptions"
        :show-current-page=true
      ></v-data-table-footer>
    </v-col>

    <slot></slot>
  </v-row>
</template>

<script lang="ts" setup>
import { getPageTextEstimated, itemsPerPageOptions } from '@/functions/datatable'

const props = defineProps<{
  totalItemsIsEstimate: boolean
  totalItemsEstimated: number
}>()

</script>