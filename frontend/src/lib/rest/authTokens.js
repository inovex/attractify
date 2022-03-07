import restClient from '../../lib/restClient'

export default {
  async list() {

    const res = await restClient.get('/auth-tokens')

    return res.data

  },
  async show(id) {

    const res = await restClient.get(`/auth-tokens/${id}`)

    return res.data

  },
  async delete(id) {

    const res = await restClient.delete(`/auth-tokens/${id}`)

    return res.data

  },
  async create(params) {

    const res = await restClient.post('/auth-tokens', params)

    return res.data

  }
}
