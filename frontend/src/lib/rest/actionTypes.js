import restClient from '../restClient'

export default {
    async list(offset, limit) {
        try {
            const params = { offset, limit }
            const res = await restClient.get('/actiontypes', params)

            return res.data
        } catch (e) {
            throw e
        }
    },
    async show(id) {
        try {
            const res = await restClient.get(`/actiontypes/${id}`)

            return res.data
        } catch (e) {
            throw e
        }
    },
    async inUse(id) {
        try {
            const res = await restClient.get(`/actiontypes/${id}/used`)
            console.log(res.data)
            return res.data.inUse
        } catch (e) {
            throw e
        }
    },
    async delete(id) {
        try {
            await restClient.delete(`/actiontypes/${id}`)
        } catch (e) {
            throw e
        }
    },
    async create(params) {
        try {
            const res = await restClient.post('/actiontypes', params)
            return res.data
        } catch (e) {
            throw e
        }
    },
    async update(params) {
        try {
            await restClient.put(`/actiontypes/${params.id}`, params)
        } catch (e) {
            throw e
        }
    },
    async listNames() {
        try {
            let res = await restClient.get('/actiontypes')
            return res.data.map(item => {
                return { text: item.name, value: item.id }
            })
        } catch (e) {
            throw e
        }
    }
}
