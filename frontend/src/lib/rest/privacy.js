import restClient from '../restClient'

export default {
  async export(params) {

    await restClient.post('/privacy/export', params)

  },
  async deletion(params) {

    await restClient.post('/privacy/deletion', params)

  },
  async lockedIdentities() {

    const res = await restClient.get('/privacy/locked-identities')

    return res.data

  },
  async createLockedIdentities(params) {

    await restClient.post('/privacy/locked-identities', params)

  },
  async deleteLockedIdentities(id) {

    await restClient.delete(`/privacy/locked-identities/${id}`)

  }
}
