<template>
  <v-list>

    <v-list-item link title="Countries" to="/countries" prepend-icon="mdi-earth" class="mt-2"></v-list-item>

    <v-divider></v-divider>
    <v-list-subheader title="Tech" class="text-uppercase mt-2 clickable" @click="showTechItems = !showTechItems"></v-list-subheader>
    <div v-if="showTechItems">
      <v-list-item link title="Database" to="/database" prepend-icon="mdi-database-search"></v-list-item>
    </div>

  </v-list>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'

const showTechItems = ref(true)

const lsKey = 'right_nav_list'

watch([showTechItems], () => {

  let lsObj = {
    'showTechItems': showTechItems.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
})

onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['showTechItems'] !== undefined) { showTechItems.value = lsObj['showTechItems'] }
})

</script>
