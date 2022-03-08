import restClient from '../../lib/restClient'

export default {
  async list(offset, limit) {
    try {
      const params = { offset, limit }
      const res = await restClient.get('/computed-traits', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async show(id) {
    try {
      const res = await restClient.get(`/computed-traits/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      await restClient.delete(`/computed-traits/${id}`)
    } catch (e) {
      throw e
    }
  },
  async create(params) {
    try {
      const res = await restClient.post('/computed-traits', params)
      return res.data
    } catch (e) {
      throw e
    }
  },
  async update(params) {
    try {
      await restClient.put(`/computed-traits/${params.id}`, params)
    } catch (e) {
      throw e
    }
  },
  async refresh(id) {
    try {
      const res = await restClient.post(`/computed-traits/${id}/refresh`)
      return res.data
    } catch (e) {
      throw e
    }
  },
  async listTraits() {
    try {
      let res = await restClient.get('/computed-traits')
      if (res.data && res.data.length > 0) {
        return res.data.map(item => {
          return {
            text: item.key,
            value: item.key,
            type: item.type,
            propertyType: item.properties.type
          }
        })
      }
    } catch (e) {
      throw e
    }
  }
}
