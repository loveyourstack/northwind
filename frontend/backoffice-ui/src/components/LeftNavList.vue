<template>
  <v-list>
    <v-list-item link title="Home" to="/home"></v-list-item>

    <v-list-item link title="Categories" to="/categories" prepend-icon="mdi-food-variant"></v-list-item>
    <v-list-item link title="Products" to="/products" prepend-icon="mdi-food-apple"></v-list-item>
    <v-list-item link title="Suppliers" to="/suppliers" prepend-icon="mdi-chef-hat"></v-list-item>

    <v-divider></v-divider>
    <v-list-subheader title="Sales" class="text-uppercase mt-2 clickable" @click="showSalesItems = !showSalesItems"></v-list-subheader>
    <div v-if="showSalesItems">
      <v-list-item link title="Customers" to="/customers" prepend-icon="mdi-account-box-outline"></v-list-item>
      <v-list-item link title="Orders" to="/orders" prepend-icon="mdi-hand-extended"></v-list-item>
      <v-list-item link title="Territories" to="/territories" prepend-icon="mdi-land-plots"></v-list-item>
    </div>

    <v-divider></v-divider>
    <v-list-subheader title="HR" class="text-uppercase mt-2 clickable" @click="showHrItems = !showHrItems"></v-list-subheader>
    <div v-if="showHrItems">
      <v-list-item link title="Calendar" to="/calendar" prepend-icon="mdi-calendar"></v-list-item>
      <v-list-item link title="Employees" to="/employees" prepend-icon="mdi-account-circle"></v-list-item>
      <v-list-item link title="Meeting schedule" to="/meeting-schedule" prepend-icon="mdi-account-group"></v-list-item>
    </div>

  </v-list>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount } from 'vue'

const showHrItems = ref(true)
const showSalesItems = ref(true)

const lsKey = 'left_nav_list'

watch([showHrItems, showSalesItems], () => {

  let lsObj = {
    'showHrItems': showHrItems.value,
    'showSalesItems': showSalesItems.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
})

onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['showHrItems'] !== undefined) { showHrItems.value = lsObj['showHrItems'] }
  if (lsObj['showSalesItems'] !== undefined) { showSalesItems.value = lsObj['showSalesItems'] }
})

</script>
