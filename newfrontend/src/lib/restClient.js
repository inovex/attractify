import axios from 'axios'
import store from '../store/index.js'

const baseURL = process.env.VUE_APP_PLATFORM_ENDPOINT

const restClient = (token = null) => {
  let headers = {}
  if (token) {
    headers.Authorization = 'Bearer ' + token
  } else {
    const user = store.getters['user/get']

    if (user && user.token) {
      headers.Authorization = 'Bearer ' + user.token
    }
  }

  return axios.create({
    baseURL: baseURL,
    headers: headers
  })
}

export default {
  get: async (path, params, token = null) => {
    let cli = await restClient(token)
    return cli.get(path, { params: params })
  },
  post: async (path, data, token = null) => {
    let cli = await restClient(token)
    return cli.post(path, data)
  },
  put: async (path, data, token = null) => {
    let cli = await restClient(token)
    return cli.put(path, data)
  },
  delete: async (path, token = null) => {
    let cli = await restClient(token)
    return cli.delete(path)
  }
}
