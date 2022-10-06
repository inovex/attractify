import Vue from 'vue'
import restClient from '../../lib/restClient'

export default {
  namespaced: true,
  state: {
    actiontypes: []
  },
  getters: {
    all(state) {
      return state.actiontypes
    },
    getById(state) {
      return id => {
        for (const actiontype of state.actiontypes) {
          if (actiontype.id === id) {
            return actiontype
          }
        }
        return null
      }
    }
  },
  mutations: {
    clear(state) {
      Vue.set(state, 'actiontypes', [])
    },
    add(state, actiontype) {
      state.actiontypes.push(actiontype)
    },
    update(state, actiontype) {
      const idx = state.actiontypes.findIndex(obj => obj.id == actiontype.id)
      state.actiontypes[idx] = actiontype
    },
    delete(state, id) {
      const idx = state.actiontypes.findIndex(obj => obj.id == id)
      state.actiontypes.splice(idx, 1)
    }
  },
  actions: {
    async list({ commit }) {
      try {
        let res = await restClient.get('/actiontypes')
        commit('clear')

        if (res.data.length > 0) {
          for (const actiontype of res.data) {
            commit('add', actiontype)
          }
        }
      } catch (e) {
        throw e
      }
    },
    async show(_, id) {
      try {
        const res = await restClient.get(`/actiontypes/${id}`)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async create({ commit }, actiontype) {
      try {
        const res = await restClient.post('/actiontypes', actiontype)
        commit('add', res.data)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async update({ commit }, actiontype) {
      try {
        await restClient.put(`/actiontypes/${actiontype.id}`, actiontype)
        commit('update', actiontype)
      } catch (e) {
        throw e
      }
    },
    async delete({ commit }, id) {
      try {
        await restClient.delete(`/actiontypes/${id}`)
        commit('delete', id)
      } catch (e) {
        throw e
      }
    }
  }
}
