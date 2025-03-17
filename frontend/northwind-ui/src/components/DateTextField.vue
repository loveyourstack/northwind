<template>
  <v-menu v-model="showDateDp" :close-on-content-click="false">
    <template #activator="{ props }">
      <v-text-field v-bind="props" :class="tFClass" :label="label ? label : 'Date'" prepend-inner-icon="mdi-calendar" readonly :clearable="clearable" :hide-details="hideDetails"
        :model-value="dateVal ? useDateFormat(dateVal, 'DD MMM YYYY').value : undefined" :rules="rules"
        @click:clear="emit('cleared')"
      ></v-text-field>
    </template>
    <template #default>
      <v-date-picker color="primary" v-model="localDateVal" :max="max" @update:model-value="showDateDp = false; emit('updated', localDateVal)"></v-date-picker>
    </template>
  </v-menu>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { useDateFormat } from '@vueuse/core'

const props = defineProps<{
  dateVal: Date | undefined
  label?: string
  tFClass?: string // tF = text field
  clearable?: boolean
  hideDetails?: boolean
  max?: Date
  rules?: readonly any[]
}>()

const emit = defineEmits<{
  (e: 'cleared'): void
  (e: 'updated', val: Date | undefined): void
}>()

const showDateDp = ref(false)

// need local val for v-date-picker v-model: can't use readonly prop dateVal
const localDateVal = ref<Date>()

watch(() => props.dateVal, () => {
  localDateVal.value = props.dateVal
}, { immediate: true }) // so that props.dateVal is read immediately and the initial date in the date picker is highlighted

</script>