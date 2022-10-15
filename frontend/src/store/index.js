import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersistence from 'vuex-persist'

import actions from './modules/actions'
import actionTypes from './modules/actionTypes'
import organization from './modules/organization'
import user from './modules/user'
import users from './modules/users'

Vue.use(Vuex)

const vuexLocal = new VuexPersistence({
  storage: window.localStorage
})

const store = new Vuex.Store({
  modules: {
    actions: actions,
    actionTypes: actionTypes,
    organization: organization,
    user: user,
    users: users
  },
  plugins: [vuexLocal.plugin]
})

export default store
