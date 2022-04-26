import restClient from '../restClient'
import eventClient from '../rest/events'

export default {
  async list(offset, limit) {
    try {
      const params = { offset, limit }
      const res = await restClient.get('/invalid-events', params)

      for (let i = res.data.length - 1; i >= 0; i--) {
        eventClient.show(res.data[i].eventId).then(eventRes => {
          res.data[i].name = eventRes.name
        }).catch(() => {
          // remove last element
          res.data.pop()
          this.delete(res.data[i].id)
        })
      }

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
