import restClient from '../../lib/restClient'

export default {
  async list(offset, limit) {
    try {
      const params = { offset, limit }
      const res = await restClient.get('/profiles', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async show(id) {
    try {
      const res = await restClient.get(`/profiles/${id}`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete(id) {
    try {
      await restClient.delete(`/profiles/${id}`)
    } catch (e) {
      throw e
    }
  },
  async listIdentities(id) {
    try {
      const res = await restClient.get(`/profiles/${id}/identities`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async listEvents(id) {
    try {
      const res = await restClient.get(`/profiles/${id}/events`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async refreshComputedTraits(id) {
    try {
      await restClient.post(`/profiles/${id}/refresh-computed-traits`)
    } catch (e) {
      throw e
    }
  },
  async search(id) {
    try {
      const res = await restClient.get(`/profiles/search/${id}`) // TODO: backend
      return res.data
    } catch (e) {
      throw e
    }
  }
}
