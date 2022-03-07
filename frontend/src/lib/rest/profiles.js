import restClient from '../../lib/restClient'

export default {
  async list(offset, limit) {

    const params = { offset, limit }
    const res = await restClient.get('/profiles', params)

    return res.data

  },
  async show(id) {

    const res = await restClient.get(`/profiles/${id}`)

    return res.data

  },
  async delete(id) {

    await restClient.delete(`/profiles/${id}`)

  },
  async listIdentities(id) {

    const res = await restClient.get(`/profiles/${id}/identities`)

    return res.data

  },
  async listEvents(id) {

    const res = await restClient.get(`/profiles/${id}/events`)

    return res.data

  },
  async refreshComputedTraits(id) {

    await restClient.post(`/profiles/${id}/refresh-computed-traits`)

  }
}
