import restClient from '../restClient'

export default {
  async list(offset, limit) {

    const params = { offset, limit }
    const res = await restClient.get('/events', params)

    return res.data

  },
  async show(id) {

    const res = await restClient.get(`/events/${id}`)

    return res.data

  },
  async delete(id) {

    await restClient.delete(`/events/${id}`)

  },
  async create(params) {

    const res = await restClient.post('/events', params)

    return res.data

  },
  async update(params) {

    await restClient.put(`/events/${params.id}`, params)

  },
  async listEventNames() {

    let res = await restClient.get('/events')
    return res.data.map(item => {
      return { text: item.name, value: item.id }
    })

  },
  async listProperties(id) {

    let res = await restClient.get(`/events/${id}/properties`)
    return res.data.map(item => {
      return { text: item.key, value: item.key, type: item.type }
    })

  }
}
