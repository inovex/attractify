import restClient from '../../lib/restClient'

export default {
  async list(offset, limit) {
    try {
      const params = { offset, limit }
      const res = await restClient.get('/audiences', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async show(id) {
    try {
      const res = await restClient.get(`/audiences/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      await restClient.delete(`/audiences/${id}`)
    } catch (e) {
      throw e
    }
  },
  async create(params) {
    try {
      const res = await restClient.post('/audiences', params)
      return res.data
    } catch (e) {
      throw e
    }
  },
  async update(params) {
    try {
      await restClient.put(`/audiences/${params.id}`, params)
    } catch (e) {
      throw e
    }
  },
  async preview(params) {
    try {
      const res = await restClient.post('/audiences/preview', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async refresh(id) {
    try {
      const res = await restClient.put(`/audiences/${id}/refresh`)
      return res.data
    } catch (e) {
      throw e
    }
  }
}
