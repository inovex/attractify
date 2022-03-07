import restClient from '../restClient'

export default {
  async list(params) {

    const res = await restClient.get('/event-log', params)

    return res.data

  },
  async delete(id) {

    await restClient.delete(`/event-log/${id}`)

  },
  async listEventNames() {

    let res = await restClient.get('/events')
    return res.data.map(item => {
      return { text: item.name, value: item.name, id: item.id }
    })

  }
}
