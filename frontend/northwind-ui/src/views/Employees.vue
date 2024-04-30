<template>
  <v-container fluid class="cockpit">
    <v-responsive class="">
      <v-row class="d-flex pt-1">
        <v-col cols="auto">
          <EmployeeTable />
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import EmployeeTable from '@/components/tables/EmployeeTable.vue'
</script>
