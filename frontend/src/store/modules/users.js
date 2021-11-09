import Vue from 'vue'
import restClient from '../../lib/restClient'

export default {
  namespaced: true,
  state: {
    users: []
  },
  getters: {
    all(state) {
      return state.users
    },
    getById(state) {
      return id =>
        state.users.filter(item => {
          return item.id === id
        })
    }
  },
  mutations: {
    clear(state) {
      Vue.set(state, 'users', [])
    },
    add(state, user) {
      state.users.push(user)
    },
    update(state, user) {
      const idx = state.users.findIndex(obj => obj.id == user.id)
      state.users[idx] = user
    },
    delete(state, id) {
      const idx = state.users.findIndex(obj => obj.id == id)
      state.users.splice(idx, 1)
    }
  },
  actions: {
    async list({ commit }) {
      try {
        let res = await restClient.get('/users')
        commit('clear')

        if (res.data && res.data.length > 0) {
          for (const user of res.data) {
            commit('add', user)
          }
        }
      } catch (e) {
        throw e
      }
    },
    async show(_, id) {
      try {
        const res = await restClient.get(`/users/${id}`)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async create({ commit }, user) {
      try {
        const res = await restClient.post('/users', user)
        commit('add', res.data)
        return res.data
      } catch (e) {
        throw e
      }
    },
    resend(_, id) {
      return restClient.post(`/users/${id}/resend`, id)
    },
    async update({ commit }, user) {
      try {
        await restClient.put(`/users/${user.id}`, user)
        commit('update', user)
      } catch (e) {
        throw e
      }
    },
    async delete({ commit }, id) {
      try {
        await restClient.delete(`/users/${id}`)
        commit('delete', id)
      } catch (e) {
        throw e
      }
    }
  }
}
