import restClient from '../../lib/restClient'

export default {
  async list() {
    try {
      const res = await restClient.get('/auth-tokens')

      return res.data
    } catch (e) {
      throw e
    }
  },
  async show(id) {
    try {
      const res = await restClient.get(`/auth-tokens/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      const res = await restClient.delete(`/auth-tokens/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async create(params) {
    try {
      const res = await restClient.post('/auth-tokens', params)

      return res.data
    } catch (e) {
      throw e
    }
  }
}
