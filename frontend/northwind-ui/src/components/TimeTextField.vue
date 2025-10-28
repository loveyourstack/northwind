<template>
  <v-menu v-model="showTimeP" :close-on-content-click="false">
    <template #activator="{ props }">
      <v-text-field v-bind="props" :class="tFClass" :label="label ? label : 'Time'" prepend-inner-icon="mdi-clock-outline" readonly :clearable="clearable" :hide-details="hideDetails"
        :model-value="timeVal ? timeVal : undefined" :rules="rules"
        @click:clear="emit('cleared')"
      ></v-text-field>
    </template>
    <template #default>
      <v-time-picker format="24hr" v-model="localTimeVal" :allowed-minutes="m => m % 5 === 0"
        @update:model-value="showTimeP = false; emit('updated', localTimeVal)">
      </v-time-picker>
    </template>
  </v-menu>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'

const props = defineProps<{
  timeVal: string | undefined // hh24:mm
  label?: string
  tFClass?: string // tF = text field
  clearable?: boolean
  hideDetails?: boolean
  rules?: readonly any[]
}>()

const emit = defineEmits<{
  (e: 'cleared'): void
  (e: 'updated', val: string | undefined): void // hh24:mm
}>()

const showTimeP = ref(false)

// need local val for v-time-picker v-model: can't use readonly prop timeVal
const localTimeVal = ref<string>()

watch(() => props.timeVal, () => {
  localTimeVal.value = (props.timeVal as string)
}, { immediate: true }) // so that props.timeVal is read immediately and the initial time in the time picker is highlighted

</script>