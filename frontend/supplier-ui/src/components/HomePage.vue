<template>
  <h1 class="text-h2 font-weight-bold projectFont mb-4">{{ appStore.projectTitle }}</h1>
</template>

<script lang="ts" setup>
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

</script>
