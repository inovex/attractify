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
      try {
        const res = await restClient.get('/organization')
        if (res.data) {
          return res.data
        }
      } catch (e) {
        throw e
      }
    },
    async signUp(_, data) {
      try {
        const res = await restClient.post('/organization', data)
        return res.data
      } catch (e) {
        throw e
      }
    },
    async update(_, data) {
      try {
        await restClient.put('/organization', data)
      } catch (e) {
        throw e
      }
    },
    async token() {
      try {
        const res = await restClient.post('/organization/token')
        return res.data
      } catch (e) {
        throw e
      }
    },
    async key(_, password) {
      try {
        const res = await restClient.post('/organization/key', {
          password: password
        })
        return res.data
      } catch (e) {
        throw e
      }
    }
  }
}
