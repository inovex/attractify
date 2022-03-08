import restClient from '../restClient'

export default {
  async export(params) {
    try {
      await restClient.post('/privacy/export', params)
    } catch (e) {
      throw e
    }
  },
  async deletion(params) {
    try {
      await restClient.post('/privacy/deletion', params)
    } catch (e) {
      throw e
    }
  },
  async lockedIdentities() {
    try {
      const res = await restClient.get('/privacy/locked-identities')

      return res.data
    } catch (e) {
      throw e
    }
  },
  async createLockedIdentities(params) {
    try {
      await restClient.post('/privacy/locked-identities', params)
    } catch (e) {
      throw e
    }
  },
  async deleteLockedIdentities(id) {
    try {
      await restClient.delete(`/privacy/locked-identities/${id}`)
    } catch (e) {
      throw e
    }
  }
}
