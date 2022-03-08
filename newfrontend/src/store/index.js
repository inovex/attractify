import { createStore } from 'vuex'
import VuexPersistence from 'vuex-persist'

import actions from './modules/actions'
import organization from './modules/organization'
import user from './modules/user'
import users from './modules/users'

const vuexLocal = new VuexPersistence({
  storage: window.sessionStorage
})

const store = createStore({
  modules: {
    actions: actions,
    organization: organization,
    user: user,
    users: users
  },
  plugins: [vuexLocal.plugin],
  state() {
    return {
      count: 1
    }
  }
})

export default store
