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

      let res = await restClient.get('/users')
      commit('clear')

      if (res.data && res.data.length > 0) {
        for (const user of res.data) {
          commit('add', user)
        }
      }

    },
    async show(_, id) {

      const res = await restClient.get(`/users/${id}`)
      return res.data

    },
    async create({ commit }, user) {

      const res = await restClient.post('/users', user)
      commit('add', res.data)
      return res.data

    },
    resend(_, id) {
      return restClient.post(`/users/${id}/resend`, id)
    },
    async update({ commit }, user) {

      await restClient.put(`/users/${user.id}`, user)
      commit('update', user)

    },
    async delete({ commit }, id) {

      await restClient.delete(`/users/${id}`)
      commit('delete', id)

    }
  }
}
