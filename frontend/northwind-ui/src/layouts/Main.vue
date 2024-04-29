<template>
  <v-app class="rounded rounded-md">
    <v-app-bar density="compact">

      <v-app-bar-nav-icon variant="text" @click.stop="showNav = !showNav"></v-app-bar-nav-icon>
      <v-img max-height="30px" max-width="30px" src="./../assets/logo.png" class="ml-1"></v-img>
      <v-toolbar-title>{{ appStore.projectTitle }}</v-toolbar-title>

      <v-spacer></v-spacer>

      <div class="text-body-1 mr-2">Username</div>

      <v-btn icon="mdi-theme-light-dark" @click="toggleTheme"></v-btn>

      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn icon="mdi-dots-vertical" v-bind="props"></v-btn>
        </template>
        
        <v-list>
          <v-list-item>
            <v-list-item-title class="clickable">Logout</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>

    </v-app-bar>

    <v-navigation-drawer v-model="showNav" theme="dark" image="./../assets/sidebar.jpg">
      <v-list density="compact">
        <v-list-item link title="Home" to="/home"></v-list-item>

        <v-divider></v-divider>
        <v-list-subheader title="Import" class="text-uppercase mt-2"></v-list-subheader>
        <v-list-item link title="Categories" to="/categories" prepend-icon="mdi-food-variant"></v-list-item>
        <v-list-item link title="Products" to="/products" prepend-icon="mdi-food-apple"></v-list-item>
        <v-list-item link title="Suppliers" to="/suppliers" prepend-icon="mdi-chef-hat"></v-list-item>

        <v-divider></v-divider>
        <v-list-subheader title="Sales" class="text-uppercase mt-2"></v-list-subheader>
        <v-list-item link title="Customers" to="/customers" prepend-icon="mdi-account-box-outline"></v-list-item>
        <v-list-item link title="Orders" to="/orders" prepend-icon="mdi-hand-extended"></v-list-item>

        <v-divider></v-divider>
        <v-list-subheader title="HR" class="text-uppercase mt-2"></v-list-subheader>
        <v-list-item link title="Employees" to="/employees" prepend-icon="mdi-account-circle"></v-list-item>

      </v-list>
    </v-navigation-drawer>

    <v-main class="d-flex" style="min-height: 300px;">
      <ApiError></ApiError>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useTheme } from 'vuetify'
import { useAppStore } from '@/stores/app'
import { useCommonStore } from '@/stores/common'
import { useCoreStore } from '@/stores/core'
import { useHRStore } from '@/stores/hr'
import ApiError from '@/components/ApiError.vue'

const theme = useTheme()
const appStore = useAppStore()
const commonStore = useCommonStore()
const coreStore = useCoreStore()
const hrStore = useHRStore()

const showNav = ref(true)

function toggleTheme () {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

onMounted(() => {
  // refresh stores on page reload
  commonStore.loadCountriesList()
  coreStore.loadCategoriesList()
  coreStore.loadSuppliersList()
  hrStore.loadEmployeesList()
})
</script>

<style>
.clickable:hover {
  cursor: pointer;
}

.dt-title {
  font-size: 1.25rem;
  font-weight: 500;
  letter-spacing: 0.0125em;
  line-height: 2rem;
}

.dt-title-block {
  padding: 0.5rem 1rem 0.5rem 0;
}

tbody tr:nth-of-type(even) {
  background-color: rgba(0, 0, 0, .01);
}

.color-pill {
  /* light background-color needs to be set */
  border-radius: 0.25em;
  padding: 5px 10px 5px 10px;
}
.v-theme--dark .color-pill {
  /* in dark theme, retain contrast to light background-color */
  color: black;
}

.v-data-table-footer {
  justify-content: flex-start;
  margin-top: 1rem;
}

.v-data-table__td a {
  color: #1E88E5;
  text-decoration: none;
}

.v-field--disabled {
  opacity: 0.7 !important;
}
</style>