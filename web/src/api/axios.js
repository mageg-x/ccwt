import axios from 'axios'

const api = axios.create({
    baseURL: '/api',
    timeout: 30000,
    withCredentials: true,
})

// 响应拦截器：token 过期自动刷新
api.interceptors.response.use(
    res => res,
    async err => {
        if (err.response?.status === 401 && !err.config._retry) {
            err.config._retry = true
            try {
                await axios.post('/api/auth/refresh', {}, { withCredentials: true })
                return api(err.config)
            } catch {
                window.location.href = '/login'
            }
        }
        return Promise.reject(err)
    }
)

export default api
