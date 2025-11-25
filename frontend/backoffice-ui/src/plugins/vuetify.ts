/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  defaults: {
    VDataTableServer: {
      style: 'border-radius: 4px !important; min-width: 680px; padding: 16px !important;'
    }
  },
  theme: {
    defaultTheme: 'system',
  },
})
