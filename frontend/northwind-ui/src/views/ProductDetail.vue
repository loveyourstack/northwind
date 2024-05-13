<template>
  <v-container fluid class="cockpit">
    <v-responsive class="">
      <v-row>
        <v-col cols="auto">
          <ProductForm :id="props.id"
            @cancel="router.back"
            @create="router.push({ name: 'Products' })"
            @delete="router.push({ name: 'Products' })"
          ></ProductForm>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'
import ProductForm from '@/components/forms/ProductForm.vue'

const props = defineProps<{
  id: number
}>()

const router = useRouter()

</script>
