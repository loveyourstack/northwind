<template>
  <v-menu :close-on-content-click="false" transition="scale-transition" offset="5px 0px">
    <template v-slot:activator="{ props }">
      <v-chip v-bind="props" size="large" :prepend-icon="filterValue.operator ? 'mdi-check' : 'mdi-plus'">
        {{ name }}<span v-if="filterText" class="ml-1 mr-1">|</span><span class="ml-1 text-secondary">{{ filterText }}</span>
        <template #close>
          <v-icon v-if="filterValue.operator" icon="mdi-close-circle" @click.stop="emit('closed')"/>
        </template>
      </v-chip>
    </template>
    <v-card width="300" class="pa-2">

      <v-autocomplete label="Operator" v-model="filterValue.operator" autofocus
        :items="coreStore.operatorsList" @update:model-value="emit('updated')"
      ></v-autocomplete>

      <v-text-field type="number" label="Value" v-model.number="filterValue.value"
        @update:model-value="emit('updatedDebounce')"
      ></v-text-field>

      <v-text-field v-if="filterValue.operator === '<=>'" type="number" label="Value" v-model.number="filterValue.value_upper"
        @update:model-value="emit('updatedDebounce')"
      ></v-text-field>

    </v-card>
  </v-menu>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { type NumericFilter } from '@/types/core'
import { useCoreStore } from '@/stores/core'
import { getNumericFilterDisplayText } from '@/functions/datatable'

const props = defineProps<{
  name: string
  filterValue: NumericFilter
  isPercent?: boolean
}>()

const emit = defineEmits<{
  (e: 'closed'): void
  (e: 'updated'): void
  (e: 'updatedDebounce'): void
}>()

const coreStore = useCoreStore()

const filterText = computed(() => {
  return getNumericFilterDisplayText(props.filterValue, props.isPercent)
})

</script>