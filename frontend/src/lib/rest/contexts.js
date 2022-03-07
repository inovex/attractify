import restClient from '../restClient'

export default {
  async list(offset, limit) {

    const params = { offset, limit }
    const res = await restClient.get('/contexts', params)

    return res.data

  },
  async show(id) {

    const res = await restClient.get(`/contexts/${id}`)

    return res.data

  },
  async delete(id) {

    await restClient.delete(`/contexts/${id}`)

  },
  async create(params) {

    const res = await restClient.post('/contexts', params)

    return res.data

  },
  async update(params) {

    await restClient.put(`/contexts/${params.id}`, params)

  },
  async listEventNames() {

    let res = await restClient.get('/contexts')
    return res.data.map(item => {
      return { text: item.name, value: item.id }
    })

  },
  async listProperties(channel) {

    let res = await restClient.get(`/contexts/${channel}/properties`)
    return res.data.map(item => {
      return { text: item.key, value: item.key, type: item.type }
    })

  }
}
