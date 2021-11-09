import restClient from '../restClient'

export default {
  async show() {
    try {
      const res = await restClient.get(`/custom-traits`)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async delete() {
    try {
      await restClient.delete(`/custom-traits`)
    } catch (e) {
      throw e
    }
  },
  async upsert(params) {
    try {
      const res = await restClient.post('/custom-traits', params)

      return res.data
    } catch (e) {
      throw e
    }
  },
  async listProperties() {
    try {
      let res = await restClient.get(`/custom-traits/properties`)
      return res.data.map(item => {
        return { text: item.key, value: item.key, type: item.type }
      })
    } catch (e) {
      throw e
    }
  }
}
