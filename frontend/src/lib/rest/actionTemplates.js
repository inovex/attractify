import restClient from '../restClient'

export default {
    async list(offset, limit) {
        try {
            const params = { offset, limit }
            const res = await restClient.get('/actiontemplates', params)

            return res.data
        } catch (e) {
            throw e
        }
    },
    async show(id) {
        try {
            const res = await restClient.get(`/actiontemplates/${id}`)

            return res.data
        } catch (e) {
            throw e
        }
    },
    async delete(id) {
        try {
            await restClient.delete(`/actiontemplates/${id}`)
        } catch (e) {
            throw e
        }
    },
    async create(params) {
        try {
            const res = await restClient.post('/actiontemplates', params)
            return res.data
        } catch (e) {
            throw e
        }
    },
    async update(params) {
        try {
            await restClient.put(`/actiontemplates/${params.id}`, params)
        } catch (e) {
            throw e
        }
    },
    async listNames() {
        try {
            let res = await restClient.get('/actiontemplates')
            return res.data.map(item => {
                return { text: item.name, value: item.id }
            })
        } catch (e) {
            throw e
        }
    }
}
