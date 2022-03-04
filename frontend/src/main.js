
import store from './store/index'
import notify from './plugins/notify'
import router from './plugins/router'
import { BusPlugin } from './plugins/bus'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'


import { createApp } from 'vue'
import { createVuetify } from 'vuetify'
import App from './components/App.vue'

//const render = render: h => h(App)
const app = createApp(App)
const vuetify = createVuetify({
  components,
  directives,
})

app.use(vuetify)
app.use(router)
app.use(BusPlugin)
app.use(notify)
app.use(store)

app.mount('#app')
/*
Vue.config.productionTip = false
Vue.use(bus)
Vue.use(notify)

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')*/
