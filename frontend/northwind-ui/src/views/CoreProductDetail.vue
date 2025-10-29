<template>
  <v-container fluid>
    <v-responsive>
      <v-row>
        <v-col cols="auto">
          <v-card>
            <v-card-text class="pa-0">

              <CoreProductForm :id="props.id"
                @cancel="router.back"
                @create="router.push({ name: 'Products' })"
                @delete="router.push({ name: 'Products' })"
              ></CoreProductForm>

            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'

const props = defineProps<{
  id: number
}>()

const router = useRouter()

</script>
