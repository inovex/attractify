import restClient from '../../lib/restClient'

export default {
  async events(actionId, range) {

    const params = { actionId, ...range }
    const res = await restClient.get('/analyze/events', params)

    return res.data

  },
  async rates(actionId, range) {

    const params = { actionId, ...range }
    const res = await restClient.get('/analyze/rates', params)

    return res.data

  },
  async reach(actionId, range) {

    const params = { actionId, ...range }
    const res = await restClient.get('/analyze/reach', params)

    return res.data

  }
}
