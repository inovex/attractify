import Vue from 'vue'
import restClient from '../../lib/restClient'

export default {
  namespaced: true,
  state: {
    actiontemplates: []
  },
  getters: {
    all(state) {
      return state.actiontemplates
    },
    getById(state) {
      return id => {
        for (const actiontemplate of state.actiontemplates) {
          if (actiontemplate.id === id) {
            return actiontemplate
          }
        }
        return null
      }
    }
  },
  mutations: {
    clear(state) {
      Vue.set(state, 'actiontemplates', [])
    },
    add(state, actiontemplate) {
      state.actiontemplates.push(actiontemplate)
    },
    update(state, actiontemplate) {
      const idx = state.actiontemplates.findIndex(obj => obj.id == actiontemplate.id)
      state.actiontemplates[idx] = actiontemplate
    },
    delete(state, id) {
      const idx = state.actiontemplates.findIndex(obj => obj.id == id)
      state.actiontemplates.splice(idx, 1)
    }
  },
  actiontemplates: {
    async list({ commit }) {
      try {
        let res = await restClient.get('/actiontemplates')
        commit('clear')

        if (res.data.length > 0) {
          for (const actiontemplate of res.data) {
            commit('add', actiontemplate)
          }
        }
      } catch (e) {
        throw e
      }
    },
    async show(_, id) {
      try {
        const res = await restClient.get(`/actiontemplates/${id}`)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async create({ commit }, actiontemplate) {
      try {
        const res = await restClient.post('/actiontemplates', actiontemplate)
        commit('add', res.data)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async update({ commit }, actiontemplate) {
      try {
        await restClient.put(`/actiontemplates/${actiontemplate.id}`, actiontemplate)
        commit('update', actiontemplate)
      } catch (e) {
        throw e
      }
    },
    async delete({ commit }, id) {
      try {
        await restClient.delete(`/actiontemplates/${id}`)
        commit('delete', id)
      } catch (e) {
        throw e
      }
    }
  }
}
