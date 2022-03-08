import restClient from '../../lib/restClient'

export default {
  async events(actionId, range) {
    try {
      const params = { actionId, ...range }
      const res = await restClient.get('/analyze/events', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async rates(actionId, range) {
    try {
      const params = { actionId, ...range }
      const res = await restClient.get('/analyze/rates', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async reach(actionId, range) {
    try {
      const params = { actionId, ...range }
      const res = await restClient.get('/analyze/reach', params)

      return res.data
    } catch (e) {
      throw e
    }
  }
}
