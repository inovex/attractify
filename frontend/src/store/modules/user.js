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
      try {
        let res = await restClient.post('/user/session', credentials)
        if (res.data) {
          commit('set', res.data)
        }
      } catch (e) {
        throw e
      }
    },
    async signOut({ commit }) {
      try {
        await restClient.delete('/user/session')
      } catch (e) {
        throw e
      }
      commit('clear')
    },
    async show() {
      try {
        const res = await restClient.get('/user')
        if (res.data) {
          return res.data
        }
      } catch (e) {
        throw e
      }
    },
    async activate({ commit }, data) {
      try {
        const params = {
          name: data.name,
          password: data.password
        }
        const res = await restClient.post('/user', params, data.token)
        if (res.data) {
          commit('set', res.data)
        }
        return res.data
      } catch (e) {
        throw e
      }
    },
    async reset(_, data) {
      try {
        await restClient.post('/user/reset-password', { email: data })
      } catch (e) {
        throw e
      }
    },
    async resetPassword(_, data) {
      try {
        const params = {
          password: data.password
        }
        await restClient.put('/user/reset-password', params, data.token)
      } catch (e) {
        throw e
      }
    },
    async updatePassword(_, data) {
      try {
        await restClient.put('/user/password', data)
      } catch (e) {
        throw e
      }
    },
    async update({ commit }, data) {
      try {
        await restClient.put('/user', data)
        if (data.email) {
          commit('setProperty', { property: 'email', value: data.email })
        }

        if (data.name) {
          commit('setProperty', { property: 'name', value: data.name })
        }
      } catch (e) {
        throw e
      }
    },
    async token(_, type) {
      try {
        const res = await restClient.post('/user/token', { type: type })
        return res.data
      } catch (e) {
        throw e
      }
    },
    async key(_, password) {
      try {
        const res = await restClient.post('/user/key', {
          password: password
        })
        return res.data
      } catch (e) {
        throw e
      }
    }
  }
}
