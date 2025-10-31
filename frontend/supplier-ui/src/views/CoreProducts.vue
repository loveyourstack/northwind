<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <CoreProductTable :supplier_id=supplierId />

            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'

const supplierId = ref(0)

onMounted(() => {
  supplierId.value = Number(localStorage.getItem('token'))
})

</script>
