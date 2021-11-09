import restClient from '../restClient'

export default {
  async list(params) {
    try {
      const res = await restClient.get('/reactions', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      await restClient.delete(`/reactions/${id}`)
    } catch (e) {
      throw e
    }
  }
}
