import restClient from '../restClient'

export default {
  async list(params) {
    try {
      const res = await restClient.get('/event-log', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      await restClient.delete(`/event-log/${id}`)
    } catch (e) {
      throw e
    }
  },
  async listEventNames() {
    try {
      let res = await restClient.get('/events')
      return res.data.map(item => {
        return { text: item.name, value: item.name, id: item.id }
      })
    } catch (e) {
      throw e
    }
  }
}
