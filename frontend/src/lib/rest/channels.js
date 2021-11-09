import restClient from '../../lib/restClient'

export default {
  async list() {
    try {
      const res = await restClient.get('/channels')

      return res.data
    } catch (e) {
      throw e
    }
  },
  async select() {
    try {
      const res = await restClient.get('/channels')

      return res.data.map(c => {
        return { text: c.name, value: c.key }
      })
    } catch (e) {
      throw e
    }
  },
  async show(id) {
    try {
      const res = await restClient.get(`/channels/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      const res = await restClient.delete(`/channels/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async create(params) {
    try {
      const res = await restClient.post('/channels', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async update(params) {
    try {
      const res = await restClient.put(`/channels/${params.id}`, params)

      return res.data
    } catch (e) {
      throw e
    }
  }
}
