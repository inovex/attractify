import Vue from 'vue'
import Vuetify from 'vuetify/lib'

Vue.use(Vuetify)

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: '#29326B',
        buttontext: '#FFFFFF'
      },
      dark: {
        primary: '#28DECE',
        buttontext: '#000000'
      }
    },
    options: { customProperties: true },
  }
})
