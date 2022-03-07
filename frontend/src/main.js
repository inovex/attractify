import { createApp } from 'vue'
import store from './store/index'
import notify from './plugins/notify'
import router from './plugins/router'
import vuetify from './plugins/vuetify'
import { BusPlugin } from './plugins/bus'



import App from './components/App.vue'

const app = createApp(App)

app.use(vuetify)
app.use(router)
app.use(BusPlugin)
app.use(notify)
app.use(store)

app.mount('#app')

