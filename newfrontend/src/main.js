import { createApp } from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import notify from './plugins/notify'
import router from './plugins/router'
import { BusPlugin } from './plugins/bus'
import { loadFonts } from './plugins/webfontloader'

loadFonts()

createApp(App)
  .use(vuetify)
  .use(router)
  .use(BusPlugin)
  .use(notify)
  .mount('#app')
