import restClient from '../../lib/restClient'

export default {
  async list(offset, limit) {
    try {
      const params = { offset, limit }
      const res = await restClient.get('/actions', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async show(id) {
    try {
      const res = await restClient.get(`/actions/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      await restClient.delete(`/actions/${id}`)
    } catch (e) {
      throw e
    }
  },
  async duplicate(id) {
    try {
      await restClient.post(`/actions/${id}/duplicate`)
    } catch (e) {
      throw e
    }
  },
  async create(params) {
    try {
      const res = await restClient.post('/actions', params)
      return res.data
    } catch (e) {
      throw e
    }
  },
  async update(params) {
    try {
      await restClient.put(`/actions/${params.id}`, params)
    } catch (e) {
      throw e
    }
  },
  async updateState(params) {
    try {
      await restClient.put(`/actions/${params.id}/state`, params)
    } catch (e) {
      throw e
    }
  },
  async listNames() {
    try {
      let res = await restClient.get('/actions')
      return res.data.map(item => {
        return { text: item.name, value: item.id }
      })
    } catch (e) {
      throw e
    }
  },
  async simulate(params) {
    try {
      params.time = Date.now()
      const res = await restClient.post(`/action-simulation`, params)
      return res.data
    } catch (e) {
      throw e
    }
  }
}
