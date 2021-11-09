import Vue from 'vue'

import vuetify from './plugins/vuetify'
import store from './store/index'
import router from './plugins/router'
import bus from './plugins/bus'
import notify from './plugins/notify'

import App from './components/App.vue'

Vue.config.productionTip = false
Vue.use(bus)
Vue.use(notify)

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
