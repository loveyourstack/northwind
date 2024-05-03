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
// https://vuetifyjs.com/en/styles/colors/#material-colors
export default createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          surface: '#FFFDE7', // yellow-lighten-5
          primary: '#F57F17', // yellow-darken-4
          yellow_tint: '#FFFDE7', // yellow-lighten-5
        },
      },
      dark: {
        colors: {
          primary: '#FFF9C4', // yellow-lighten-4
          yellow_tint: '#F57F17', // yellow-darken-4
        }
      }
    },
  },
})
