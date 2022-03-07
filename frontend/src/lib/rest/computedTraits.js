import restClient from '../../lib/restClient'

export default {
  async list(offset, limit) {

    const params = { offset, limit }
    const res = await restClient.get('/computed-traits', params)

    return res.data

  },
  async show(id) {

    const res = await restClient.get(`/computed-traits/${id}`)

    return res.data

  },
  async delete(id) {

    await restClient.delete(`/computed-traits/${id}`)

  },
  async create(params) {

    const res = await restClient.post('/computed-traits', params)
    return res.data

  },
  async update(params) {

    await restClient.put(`/computed-traits/${params.id}`, params)

  },
  async refresh(id) {

    const res = await restClient.post(`/computed-traits/${id}/refresh`)
    return res.data

  },
  async listTraits() {

    let res = await restClient.get('/computed-traits')
    if (res.data && res.data.length > 0) {
      return res.data.map(item => {
        return {
          text: item.key,
          value: item.key,
          type: item.type,
          propertyType: item.properties.type
        }
      })
    }

  }
}
