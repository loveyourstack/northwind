<template>
  <v-app class="rounded rounded-md">
    <v-app-bar density="compact" elevation="8" class="appbar-bg" image="./../assets/sky.png" theme="dark">

      <v-app-bar-nav-icon variant="text" @click.stop="showNav = !showNav"></v-app-bar-nav-icon>
      <v-img max-height="30px" max-width="30px" src="./../assets/logo.png" class="ml-1"></v-img>
      <v-toolbar-title>
        <span class="text-yellow-darken-2 font-weight-bold projectFont">{{ appStore.projectTitle }}</span>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <div class="text-body-1 mr-2">Username</div>

      <v-btn icon="mdi-theme-light-dark" v-tooltip="'Toggle theme'" @click="toggleTheme"></v-btn>

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

    <v-navigation-drawer v-model="showNav" elevation="8" class="nav-bg-left" floating>
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

    <v-main class="d-flex cockpit">
      <ApiError></ApiError>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { ref, watch, onBeforeMount, onMounted } from 'vue'
import { useTheme } from 'vuetify'
import { useAppStore } from '@/stores/app'
import { useCommonStore } from '@/stores/common'
import { useCoreStore } from '@/stores/core'
import { useHRStore } from '@/stores/hr'
import { useSalesStore } from '@/stores/sales'
import ApiError from '@/components/ApiError.vue'

const theme = useTheme()
const appStore = useAppStore()
const commonStore = useCommonStore()
const coreStore = useCoreStore()
const hrStore = useHRStore()
const salesStore = useSalesStore()

const showNav = ref(true)

const lsKey = 'main'

function toggleTheme () {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

watch([theme.global.name], () => {

  let lsObj = {
    'theme': theme.global.name.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
})

onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['theme']) { theme.global.name.value = lsObj['theme'] }
})

onMounted(() => {
  // refresh stores on page reload
  commonStore.loadCountriesList()
  coreStore.loadCategoriesList()
  coreStore.loadProductsList()
  coreStore.loadSuppliersList()
  hrStore.loadEmployeesList()
  salesStore.loadCustomersList()
  salesStore.loadShippersList()
})
</script>

<style>

.appbar-bg {
  opacity: 0.9;
}

.clickable:hover {
  cursor: pointer;
}

/* using pseudo-element to be able to apply opacity to background image */
.cockpit {
  position: relative; 
  height: 100%;
  display: flex;
}
.cockpit::before {
  content: "";
  background: url('../assets/cockpit.jpg');
  background-size: cover;
  position: absolute;
  top: 0px;
  right: 0px;
  bottom: 0px;
  left: 0px;
  opacity: 0.4;
}
.v-theme--dark .cockpit::before {
  opacity: 0.25;
}

.color-pill {
  border-radius: 0.25em;
  padding: 5px 10px 5px 10px;
}

.dt {
  border-radius: 4px !important;
  min-width: 680px;
  padding: 16px !important;
}

.dt-title {
  font-size: 1.25rem;
  font-weight: 500;
  letter-spacing: 0.0125em;
  line-height: 2rem;
  white-space: nowrap;
}

.dt-title-block {
  float: left;
  padding: 0.3rem 2rem 0.5rem 0;
}

.form-col-s {
  min-width: 350px;
}
.form-col {
  min-width: 450px;
}
.form-col-l {
  min-width: 550px;
}

.fs-std {
  border-color: #FAFAFA80;
}

.nav-bg-left {
  opacity: 0.9;
  background: linear-gradient(90deg, rgb(var(--v-theme-dark_yellow)) 0%, rgb(var(--v-theme-light_yellow)) 50%) !important;
}
.v-theme--dark .nav-bg-left {
  background: linear-gradient(90deg, rgb(var(--v-theme-dark_yellow)) 0%, rgba(33,33,33,1) 50%) !important;
}

.projectFont {
  font-family: "Sriracha", cursive;
}

tbody tr:nth-of-type(even) {
  background-color: rgba(0, 0, 0, .02);
}
.v-theme--dark tbody tr:nth-of-type(even) {
  background-color: rgba(0, 0, 0, .10);
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