import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersistence from 'vuex-persist'

import actions from './modules/actions'
import organization from './modules/organization'
import user from './modules/user'
import users from './modules/users'

Vue.use(Vuex)

const vuexLocal = new VuexPersistence({
  storage: window.sessionStorage
})

const store = new Vuex.Store({
  modules: {
    actions: actions,
    organization: organization,
    user: user,
    users: users
  },
  plugins: [vuexLocal.plugin]
})

export default store
