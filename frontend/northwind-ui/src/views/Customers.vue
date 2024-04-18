<template>
  <v-container class="">
    <v-responsive class="">
      <v-row class="d-flex pt-1">
        <v-col cols="auto">
          <CustomerTable />
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import CustomerTable from '@/components/tables/CustomerTable.vue'
</script>
