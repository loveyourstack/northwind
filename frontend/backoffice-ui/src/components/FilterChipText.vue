<template>
  <v-menu :close-on-content-click="false" transition="scale-transition" offset="5px 0px">
    <template v-slot:activator="{ props }">
      <v-chip v-bind="props" size="large" :prepend-icon="filterValue ? 'mdi-check' : 'mdi-plus'">
        {{ name }}<span v-if="filterValue" class="ml-1 mr-1">|</span><span class="ml-1 text-secondary">{{ filterValue }}</span>
        <template #close>
          <v-icon v-if="filterValue" icon="mdi-close-circle" @click.stop="emit('closed')"/>
        </template>
      </v-chip>
    </template>
    <v-card width="300" class="pa-2">
      <v-text-field :label="name" v-model="localVal" autofocus
        @update:model-value="emit('updated', localVal)"
      ></v-text-field>
    </v-card>
  </v-menu>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'

const props = defineProps<{
  name: string
  filterValue: string | undefined
}>()

const emit = defineEmits<{
  (e: 'closed'): void
  (e: 'updated', val: string | undefined): void
}>()

// need local val for v-text-field v-model: can't use readonly prop filterValue
const localVal = ref<string>()

watch(() => props.filterValue, () => {
  localVal.value = props.filterValue
})

</script>