<template>
  <v-container class="">
    <v-responsive class="">
      <v-row class="d-flex pt-1">
        <v-col cols="auto">
          <ProductTable :supplier_id="0" />
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import ProductTable from '@/components/tables/ProductTable.vue'
</script>
