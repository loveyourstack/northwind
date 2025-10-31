<template>
  <v-app v-if="showLoginForm" class="rounded rounded-md">
    <v-main class="d-flex" style="min-height: 300px;">
      <v-container class="fill-height">
        <v-responsive class="align-center text-center fill-height">
          <v-row>
            <v-col>

              <ApiError></ApiError>

              <v-card max-width="500" class="mx-auto">
                <v-card-title class="d-flex flex-wrap justify-center mt-4">
                  <span class="ml-4">{{ appStore.company }} - {{ appStore.projectTitle }}</span>
                </v-card-title>
                <v-card-text class="pa-5">
                  <v-form ref="loginForm">

                    <v-autocomplete label="Supplier" v-model="supplier" return-object
                      :items="coreStore.suppliersList" item-title="name" item-value="id"
                      :rules="[(v: number) => !!v || 'Supplier is required']"
                    ></v-autocomplete>

                    <v-btn color="secondary" block class="mt-2" @click="login">Login</v-btn>

                  </v-form>
                  
                </v-card-text>
              </v-card>

            </v-col>
          </v-row>
        </v-responsive>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import auth from '@/auth'
import ax from '@/api'
import { useAppStore } from '@/stores/app'
import { useCoreStore } from '@/stores/core'
import { type Supplier } from '@/types/core'

const appStore = useAppStore()
const coreStore = useCoreStore()

const showLoginForm = ref(false)
const loginForm = ref()
const supplier = ref<Supplier>()

const route = useRoute()
const router = useRouter()

function login() {
  if (!supplier.value) {
    return
  }

  localStorage.setItem('token', String(supplier.value.id))
  localStorage.setItem('name', supplier.value.name)
  loginSuccess(supplier.value.id, supplier.value.name)
}

function tokenLogin(supplierId: number, supplierName: string) {
  loginSuccess(supplierId, supplierName)
}

function loginSuccess(supplierId: number, supplierName: string) {

  // flag user as authenticated, write user props
  auth.user.authenticated = true

  if (supplierName) {
    auth.user.name = supplierName
  }

  ax.defaults.headers.common.Authorization = supplierId

  // redirect to intended path
  var destPath = '/home'
  if (route.query.to) {
    destPath = route.query.to.toString()
  }
  router.push({ path: destPath })
}

onMounted(() => {

  // set dummy supplier so that supplier list can be loaded before auth
  ax.defaults.headers.common.Authorization = '1'
  coreStore.loadSuppliersList()

  if (localStorage.getItem('token') && route.query.to) {
    tokenLogin(Number(localStorage.getItem('token')), localStorage.getItem('name')!)
  } else {
    showLoginForm.value = true
  }
})
</script>

<style scoped>
</style>