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
  async delete(id) {
    try {
      await restClient.delete(`/invalid-events/${id}`)
    } catch (e) {
      throw e
    }
  }
}
