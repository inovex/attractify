import restClient from '../restClient'

export default {
  async list(params) {

    const res = await restClient.get('/reactions', params)

    return res.data

  },
  async delete(id) {

    await restClient.delete(`/reactions/${id}`)

  }
}
