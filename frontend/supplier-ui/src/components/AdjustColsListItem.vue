<template>
  <v-menu :close-on-content-click=false location="start">
    <template v-slot:activator="{ props }">
      <v-list-item v-bind="props" prepend-icon="mdi-table-column">
        <v-list-item-title class="clickable">Adjust columns</v-list-item-title>
      </v-list-item>
    </template>
    <v-list>
      <v-list-item v-for="(header, i) in props.headers" :key="i" :value="header" @click="emit('toggle', header.key)">
        <template v-slot:append>
          <v-icon :icon="getHeaderListIcon(props.excludedHeaders, header.key)" :color="getHeaderListIconColor(props.excludedHeaders, header.key)"></v-icon>
        </template>
        <v-list-item-title class="clickable" v-text="header.title"></v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script lang="ts" setup>

const props = defineProps<{
  headers: readonly any[]
  excludedHeaders: string[]
}>()

const emit = defineEmits<{
  (e: 'toggle', key: string): void
}>()

function getHeaderListIcon(excludedHeaders: string[], headerKey: string) {
  return excludedHeaders.includes(headerKey) ? 'mdi-close' : 'mdi-check'
}
function getHeaderListIconColor(excludedHeaders: string[], headerKey: string) {
  return excludedHeaders.includes(headerKey) ? 'error' : 'success'
}

</script>