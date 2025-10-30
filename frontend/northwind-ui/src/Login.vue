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
                  <span class="ml-4">{{ appStore.projectTitle }}</span>
                </v-card-title>
                <v-card-text class="pa-5">
                  <v-form ref="loginForm">

                    <v-autocomplete label="Employee" v-model="employeeName"
                      :items="hrStore.mandEmployeesList" item-title="name" item-value="name"
                      :rules="[(v: number) => !!v || 'Employee is required']"
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
import { useAppStore } from '@/stores/app'
import { useHRStore } from '@/stores/hr'

const appStore = useAppStore()
const hrStore = useHRStore()

const showLoginForm = ref(false)
const loginForm = ref()
const employeeName = ref('')

const route = useRoute()
const router = useRouter()

function login() {
  localStorage.setItem('token', employeeName.value)
  loginSuccess(employeeName.value)
}

function tokenLogin(employeeName: string) {
  loginSuccess(employeeName)
}

function loginSuccess(employeeName: string) {

  // flag user as authenticated, write user props
  auth.user.authenticated = true
  auth.user.name = employeeName

  // redirect to intended path
  var destPath = '/home'
  if (route.query.to) {
    destPath = route.query.to.toString()
  }
  router.push({ path: destPath })
}

onMounted(() => {
  hrStore.loadEmployeesList()

  if (localStorage.getItem('token') && route.query.to) {
    tokenLogin(localStorage.getItem('token')!)
  } else {
    showLoginForm.value = true
  }
})
</script>

<style scoped>
</style>