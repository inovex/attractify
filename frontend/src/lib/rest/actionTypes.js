import restClient from '../restClient'

export default {
    async list(offset, limit) {
        try {
            const params = { offset, limit }
            const res = await restClient.get('/action-types', params)

            return res.data
        } catch (e) {
            throw e
        }
    },
    async show(id) {
        try {
            const res = await restClient.get(`/action-types/${id}`)

            return res.data
        } catch (e) {
            throw e
        }
    },
    async delete(id) {
        try {
            await restClient.delete(`/action-types/${id}`)
        } catch (e) {
            throw e
        }
    },
    async create(params) {
        try {
            const res = await restClient.post('/action-types', params)
            return res.data
        } catch (e) {
            throw e
        }
    },
    async update(params) {
        try {
            await restClient.put(`/action-types/${params.id}`, params)
        } catch (e) {
            throw e
        }
    },
    async nameExists(name) {
        try {
            const res = await restClient.get(`/action-types/exists/${name}`)
            return res.data
        } catch (e) {
            throw e
        }
    }
}
