<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <OrderDetailTable :order_id="props.id" />
              
              <v-row class="pt-4 pb-4 pl-4">
                <v-col>
                  <v-btn icon class="mr-4 mb-1 ml-1" @click="router.back">
                    <v-icon icon="mdi-arrow-left"></v-icon>
                  </v-btn>
                </v-col>
              </v-row>

            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import OrderDetailTable from '@/components/tables/OrderDetailTable.vue';

const props = defineProps<{
  id: number
}>()

const router = useRouter()

const selectedTab = ref('orderDetails')

</script>