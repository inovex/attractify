import restClient from '../../lib/restClient'

export default {
  async list(offset, limit) {

    const params = { offset, limit }
    const res = await restClient.get('/actions', params)

    return res.data

  },
  async show(id) {

    const res = await restClient.get(`/actions/${id}`)

    return res.data

  },
  async delete(id) {

    await restClient.delete(`/actions/${id}`)

  },
  async duplicate(id) {

    await restClient.post(`/actions/${id}/duplicate`)

  },
  async create(params) {

    const res = await restClient.post('/actions', params)
    return res.data

  },
  async update(params) {

    await restClient.put(`/actions/${params.id}`, params)

  },
  async updateState(params) {

    await restClient.put(`/actions/${params.id}/state`, params)

  },
  async listNames() {

    let res = await restClient.get('/actions')
    return res.data.map(item => {
      return { text: item.name, value: item.id }
    })

  }
}
