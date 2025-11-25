<template>
  <v-app class="rounded rounded-md">
    <v-app-bar density="compact" elevation="8">

      <v-app-bar-nav-icon variant="text" v-tooltip="'Toggle left menu'" @click.stop="showLeftNav = !showLeftNav"></v-app-bar-nav-icon>
      <v-img max-height="30px" max-width="30px" src="./../assets/logo.png" class="ml-1"></v-img>
      <v-toolbar-title>
        <span class="font-weight-bold projectFont">{{ appStore.company }} - {{ appStore.projectTitle }}</span>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <div class="text-body-1 mr-2">{{ auth.user.name }}</div>

      <v-btn icon="mdi-theme-light-dark" v-tooltip="'Toggle theme'" @click="theme.toggle()"></v-btn>

      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn icon="mdi-dots-vertical" v-bind="props"></v-btn>
        </template>
        
        <v-list>
          <v-list-item prepend-icon="mdi-logout">
            <v-list-item-title class="clickable" @click="logout()">Logout</v-list-item-title>
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
import { useRouter } from 'vue-router'
import auth from '@/auth'
import { useAppStore } from '@/stores/app'
import { useCoreStore } from '@/stores/core'
import { useHRStore } from '@/stores/hr'
import { useSalesStore } from '@/stores/sales'

const theme = useTheme()
const router = useRouter()
const appStore = useAppStore()
const coreStore = useCoreStore()
const hrStore = useHRStore()
const salesStore = useSalesStore()

const showLeftNav = ref(true)
const showRightNav = ref(true)

const lsKey = 'main'

function logout() {
  auth.logout()
  router.push({ path: '/login' })
}

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
  //hrStore.loadEmployeesList() // loaded by Login.vue
  salesStore.loadCustomersList()
  salesStore.loadShippersList()
})
</script>

<style scoped>

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

</style>