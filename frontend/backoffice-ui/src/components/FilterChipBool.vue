<template>
  <v-menu :close-on-content-click="false" transition="scale-transition" offset="5px 0px">
    <template v-slot:activator="{ props }">
      <v-chip v-bind="props" size="large" :prepend-icon="filterValue ? 'mdi-check' : 'mdi-plus'">
        {{ name }}<span v-if="filterText" class="ml-1 mr-1">|</span><span class="ml-1 text-secondary">{{ filterText }}</span>
        <template #close>
          <v-icon v-if="filterValue != undefined" icon="mdi-close-circle" @click.stop="emit('closed')"/>
        </template>
      </v-chip>
    </template>
    <v-card width="300" class="pa-2">
      <slot name="menuContent"></slot>
    </v-card>
  </v-menu>
</template>

<script lang="ts" setup>

const props = defineProps<{
  name: string
  filterValue: boolean | undefined
  filterText: string | undefined
}>()

const emit = defineEmits<{
  (e: 'closed'): void
}>()

</script>