<template>
  <v-container fluid class="cockpit">
    <v-responsive class="">
      <v-row class="d-flex pt-1">
        <v-col cols="auto">
          <OrderTable :customer_id="0" />
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import OrderTable from '@/components/tables/OrderTable.vue'
</script>
