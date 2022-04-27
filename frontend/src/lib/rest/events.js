import restClient from '../restClient'
import invalidEventClient from '../rest/invalidEvents'

export default {
  async list(offset, limit) {
    try {
      const params = { offset, limit }
      const res = await restClient.get('/events', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async show(id) {
    try {
      const res = await restClient.get(`/events/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      await restClient.delete(`/events/${id}`)
      const res = await restClient.get('/invalid-events', {})

      for (let i = res.data.length - 1; i >= 0; i--) {
        if (res.data[i].eventId === id) {
          invalidEventClient.delete(res.data[i].id)
        }
      }

    } catch (e) {
      throw e
    }
  },
  async create(params) {
    try {
      const res = await restClient.post('/events', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async update(params) {
    try {
      await restClient.put(`/events/${params.id}`, params)
    } catch (e) {
      throw e
    }
  },
  async listEventNames() {
    try {
      let res = await restClient.get('/events')
      return res.data.map(item => {
        return { text: item.name, value: item.id }
      })
    } catch (e) {
      throw e
    }
  },
  async listProperties(id) {
    try {
      let res = await restClient.get(`/events/${id}/properties`)
      return res.data.map(item => {
        return { text: item.key, value: item.key, type: item.type }
      })
    } catch (e) {
      throw e
    }
  }
}
