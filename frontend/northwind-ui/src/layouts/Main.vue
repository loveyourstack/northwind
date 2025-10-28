<template>
  <v-app class="rounded rounded-md">
    <v-app-bar density="compact" elevation="8">

      <v-app-bar-nav-icon variant="text" v-tooltip="'Toggle left menu'" @click.stop="showLeftNav = !showLeftNav"></v-app-bar-nav-icon>
      <v-img max-height="30px" max-width="30px" src="./../assets/logo.png" class="ml-1"></v-img>
      <v-toolbar-title>
        <span class="font-weight-bold projectFont">{{ appStore.projectTitle }}</span>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <div class="text-body-1 mr-2">Username</div>

      <v-btn icon="mdi-theme-light-dark" v-tooltip="'Toggle theme'" @click="theme.toggle()"></v-btn>

      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn icon="mdi-dots-vertical" v-bind="props"></v-btn>
        </template>
        
        <v-list>
          <v-list-item prepend-icon="mdi-logout">
            <v-list-item-title class="clickable">Logout</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>

      <v-app-bar-nav-icon variant="text" v-tooltip="'Toggle right menu'" @click.stop="showRightNav = !showRightNav"></v-app-bar-nav-icon>
    </v-app-bar>

    <v-navigation-drawer v-model="showLeftNav" elevation="8" floating>
      <LeftNavList />
    </v-navigation-drawer>

    <v-navigation-drawer location="right" v-model="showRightNav" elevation="8" floating>
      <RightNavList />
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
import { useCoreStore } from '@/stores/core'
import { useHRStore } from '@/stores/hr'
import { useSalesStore } from '@/stores/sales'
import ApiError from '@/components/ApiError.vue'
import LeftNavList from '@/components/LeftNavList.vue'
import RightNavList from '@/components/RightNavList.vue'

const theme = useTheme()
const appStore = useAppStore()
const coreStore = useCoreStore()
const hrStore = useHRStore()
const salesStore = useSalesStore()

const showLeftNav = ref(true)
const showRightNav = ref(true)

const lsKey = 'main'

watch([showRightNav, theme.global.name], () => {

  let lsObj = {
    'showRightNav': showRightNav.value,
    'theme': theme.name.value,
  }
  localStorage.setItem(lsKey, JSON.stringify(lsObj))
})

onBeforeMount(() => {
  var lsJSON = localStorage.getItem(lsKey)
  if (!lsJSON) {
    return
  }

  let lsObj = JSON.parse(lsJSON)
  if (lsObj['showRightNav'] !== undefined) { showRightNav.value = lsObj['showRightNav'] }
  if (lsObj['theme']) { theme.change(lsObj['theme']) }
})

onMounted(() => {
  // refresh stores on page reload
  coreStore.loadCategoriesList()
  coreStore.loadCountriesList()
  coreStore.loadProductsList()
  coreStore.loadSuppliersList()
  hrStore.loadEmployeesList()
  salesStore.loadCustomersList()
  salesStore.loadShippersList()
})
</script>

<style>

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
  float: left !important;
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