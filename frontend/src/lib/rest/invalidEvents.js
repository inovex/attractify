import restClient from '../restClient'

export default {
  async list(offset, limit) {
    try {
      const params = { offset, limit }
      const res = await restClient.get('/invalid-events', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async show(id) {
    try {
      const res = await restClient.get(`/invalid-events/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      await restClient.delete(`/invalid-events/${id}`)
    } catch (e) {
      throw e
    }
  },
  async listEventNames() {
    try {
      let res = await restClient.get('/invalid-events')
      return res.data.map(item => {
        return { text: item.name, value: item.id }
      })
    } catch (e) {
      throw e
    }
  }
}
