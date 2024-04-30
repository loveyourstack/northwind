<template>
  <v-container fluid class="cockpit">
    <v-responsive class="">
      <v-row>
        <v-col>

          <v-tabs color="deep-purple-accent-4">

            <v-tab value="orderDetails">Order details</v-tab>

          </v-tabs>

          <v-window class="mt-2">

            <v-window-item value="orderDetails">
              <v-card>
                <v-card-text class="pt-0">
                  <OrderDetailTable :order_id="props.id" />
                  <v-row class="pt-4">
                    <v-col>
                      <v-btn icon class="mr-4 mb-1 ml-1" @click="router.back">
                        <v-icon icon="mdi-arrow-left"></v-icon>
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-card-text>
              </v-card>
            </v-window-item>

          </v-window>

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