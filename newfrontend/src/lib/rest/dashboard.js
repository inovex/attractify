import restClient from '../../lib/restClient'

export default {
  async load() {
    const res = await restClient.get('/dashboard')

    return res.data
  }
}
