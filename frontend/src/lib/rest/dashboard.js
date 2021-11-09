import restClient from '../../lib/restClient'

export default {
  async load() {
    try {
      const res = await restClient.get('/dashboard')

      return res.data
    } catch (e) {
      throw e
    }
  }
}
