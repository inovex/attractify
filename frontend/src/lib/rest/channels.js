import restClient from '../../lib/restClient'

export default {
  async list() {

    const res = await restClient.get('/channels')

    return res.data

  },
  async select() {

    const res = await restClient.get('/channels')

    return res.data.map(c => {
      return { text: c.name, value: c.key }
    })

  },
  async show(id) {

    const res = await restClient.get(`/channels/${id}`)

    return res.data

  },
  async delete(id) {

    const res = await restClient.delete(`/channels/${id}`)

    return res.data

  },
  async create(params) {

    const res = await restClient.post('/channels', params)

    return res.data

  },
  async update(params) {

    const res = await restClient.put(`/channels/${params.id}`, params)

    return res.data

  }
}
