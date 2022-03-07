import restClient from '../../lib/restClient'

export default {
  async list(offset, limit) {

    const params = { offset, limit }
    const res = await restClient.get('/audiences', params)

    return res.data

  },
  async show(id) {

    const res = await restClient.get(`/audiences/${id}`)

    return res.data

  },
  async delete(id) {

    await restClient.delete(`/audiences/${id}`)

  },
  async create(params) {

    const res = await restClient.post('/audiences', params)
    return res.data

  },
  async update(params) {

    await restClient.put(`/audiences/${params.id}`, params)

  },
  async preview(params) {

    const res = await restClient.post('/audiences/preview', params)

    return res.data

  },
  async refresh(id) {

    const res = await restClient.put(`/audiences/${id}/refresh`)
    return res.data

  }
}
