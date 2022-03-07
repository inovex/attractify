import Vue from 'vue'
import restClient from '../../lib/restClient'

export default {
  namespaced: true,
  state: {
    user: null
  },
  getters: {
    get(state) {
      return state.user
    }
  },
  mutations: {
    clear(state) {
      Vue.set(state, 'user', null)
    },
    set(state, user) {
      state.user = user
    },
    setProperty(state, data) {
      Vue.set(state.user, data.property, data.value)
    }
  },
  actions: {
    async signIn({ commit }, credentials) {
      let res = await restClient.post('/user/session', credentials)
      if (res.data) {
        commit('set', res.data)
      }
    },
    async signOut({ commit }) {
      await restClient.delete('/user/session')
      commit('clear')
    },
    async show() {
      const res = await restClient.get('/user')
      if (res.data) {
        return res.data
      }
    },
    async activate({ commit }, data) {
      const params = {
        name: data.name,
        password: data.password
      }
      const res = await restClient.post('/user', params, data.token)
      if (res.data) {
        commit('set', res.data)
      }
      return res.data
    },
    async reset(_, data) {
      await restClient.post('/user/reset-password', { email: data })
    },
    async resetPassword(_, data) {
      const params = {
        password: data.password
      }
      await restClient.put('/user/reset-password', params, data.token)
    },
    async updatePassword(_, data) {
      await restClient.put('/user/password', data)
    },
    async update({ commit }, data) {
      await restClient.put('/user', data)
      if (data.email) {
        commit('setProperty', { property: 'email', value: data.email })
      }

      if (data.name) {
        commit('setProperty', { property: 'name', value: data.name })
      }
    },
    async token(_, type) {
      const res = await restClient.post('/user/token', { type: type })
      return res.data
    },
    async key(_, password) {
      const res = await restClient.post('/user/key', {
        password: password
      })
      return res.data
    }
  }
}
