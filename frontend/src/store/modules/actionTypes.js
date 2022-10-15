import Vue from 'vue'
import restClient from '../../lib/restClient'

export default {
  namespaced: true,
  state: {
    actionTypes: []
  },
  getters: {
    all(state) {
      return state.actionTypes
    },
    getById(state) {
      return id => {
        for (const actionType of state.actionTypes) {
          if (actionType.id === id) {
            return actionType
          }
        }
        return null
      }
    }
  },
  mutations: {
    clear(state) {
      Vue.set(state, 'actionTypes', [])
    },
    add(state, actionType) {
      state.actionTypes.push(actionType)
    },
    update(state, actionType) {
      const idx = state.actionTypes.findIndex(obj => obj.id == actionType.id)
      state.actionTypes[idx] = actionType
    },
    delete(state, id) {
      const idx = state.actionTypes.findIndex(obj => obj.id == id)
      state.actionTypes.splice(idx, 1)
    }
  },
  actions: {
    async list({ commit }) {
      try {
        let res = await restClient.get('/action-types')
        commit('clear')

        if (res.data.length > 0) {
          for (const actionType of res.data) {
            commit('add', actionType)
          }
        }
      } catch (e) {
        throw e
      }
    },
    async show(_, id) {
      try {
        const res = await restClient.get(`/action-types/${id}`)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async create({ commit }, actionType) {
      try {
        const res = await restClient.post('/action-types', actionType)
        commit('add', res.data)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async update({ commit }, actionType) {
      try {
        await restClient.put(`/action-types/${actionType.id}`, actionType)
        commit('update', actionType)
      } catch (e) {
        throw e
      }
    },
    async delete({ commit }, id) {
      try {
        await restClient.delete(`/action-types/${id}`)
        commit('delete', id)
      } catch (e) {
        throw e
      }
    }
  }
}
