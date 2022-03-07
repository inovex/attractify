import restClient from '../restClient'

export default {
  async show() {

    const res = await restClient.get(`/custom-traits`)

    return res.data

  },
  async delete() {

    await restClient.delete(`/custom-traits`)

  },
  async upsert(params) {

    const res = await restClient.post('/custom-traits', params)

    return res.data

  },
  async listProperties() {

    let res = await restClient.get(`/custom-traits/properties`)
    return res.data.map(item => {
      return { text: item.key, value: item.key, type: item.type }
    })

  }
}
