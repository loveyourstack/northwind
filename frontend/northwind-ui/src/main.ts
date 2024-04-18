/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'

// Plugins
import vuetify from '@/plugins/vuetify'
import pinia from '@/stores'
import router from '@/router'

const app = createApp(App)

app.use(vuetify)
app.use(router)
app.use(pinia)

app.mount('#app')
