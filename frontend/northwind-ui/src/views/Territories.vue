<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <TerritoryTable :salesman_id="0" />

            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import TerritoryTable from '@/components/tables/TerritoryTable.vue'
</script>
