import Vue from 'vue'
import restClient from '../../lib/restClient'

export default {
  namespaced: true,
  state: {
    actions: []
  },
  getters: {
    all(state) {
      return state.actions
    },
    getById(state) {
      return id => {
        for (const action of state.actions) {
          if (action.id === id) {
            return action
          }
        }
        return null
      }
    }
  },
  mutations: {
    clear(state) {
      Vue.set(state, 'actions', [])
    },
    add(state, action) {
      state.actions.push(action)
    },
    update(state, action) {
      const idx = state.actions.findIndex(obj => obj.id == action.id)
      state.actions[idx] = action
    },
    delete(state, id) {
      const idx = state.actions.findIndex(obj => obj.id == id)
      state.actions.splice(idx, 1)
    }
  },
  actions: {
    async list({ commit }) {
      try {
        let res = await restClient.get('/actions')
        commit('clear')

        if (res.data.length > 0) {
          for (const action of res.data) {
            commit('add', action)
          }
        }
      } catch (e) {
        throw e
      }
    },
    async show(_, id) {
      try {
        const res = await restClient.get(`/actions/${id}`)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async create({ commit }, action) {
      try {
        const res = await restClient.post('/actions', action)
        commit('add', res.data)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async update({ commit }, action) {
      try {
        await restClient.put(`/actions/${action.id}`, action)
        commit('update', action)
      } catch (e) {
        throw e
      }
    },
    async updateState({ commit }, req) {
      try {
        await restClient.put(`/actions/${req.action.id}/state`, {
          state: req.state
        })
        req.action.conditions.state = req.state
        commit('update', req.action)
      } catch (e) {
        throw e
      }
    },
    async delete({ commit }, id) {
      try {
        await restClient.delete(`/actions/${id}`)
        commit('delete', id)
      } catch (e) {
        throw e
      }
    },
    async testWebhook(_, data) {
      const req = {
        event: data.event,
        channel: data.channel,
        userId: data.userId,
        properties: data.properties
      }

      try {
        const res = await restClient.post(
          `/actions/${data.actionId}/test-webhook`,
          req
        )

        return res.data
      } catch (e) {
        throw e
      }
    }
  }
}
