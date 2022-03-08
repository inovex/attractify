import restClient from '../../lib/restClient'

export default {
  namespaced: true,
  state: {
    organization: {}
  },
  getters: {
    get(state) {
      return state.organization
    }
  },
  actions: {
    async show() {
      const res = await restClient.get('/organization')
      if (res.data) {
        return res.data
      }
    },
    async signUp(_, data) {
      const res = await restClient.post('/organization', data)
      return res.data
    },
    async update(_, data) {
      await restClient.put('/organization', data)
    },
    async token() {
      const res = await restClient.post('/organization/token')
      return res.data
    },
    async key(_, password) {
      const res = await restClient.post('/organization/key', {
        password: password
      })
      return res.data
    }
  }
}
