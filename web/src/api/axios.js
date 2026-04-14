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
        const cfg = err.config || {}
        const reqUrl = cfg.url || ''
        const isAuthApi = reqUrl.includes('/auth/login') ||
            reqUrl.includes('/auth/register') ||
            reqUrl.includes('/auth/refresh')

        if (err.response?.status === 401 && !cfg._retry && !cfg.__skipRefresh && !isAuthApi) {
            cfg._retry = true
            try {
                await api.post('/auth/refresh', {}, { __skipRefresh: true })
                return api(cfg)
            } catch {
                if (window.location.pathname !== '/login') {
                    window.location.replace('/login')
                }
            }
        }
        return Promise.reject(err)
    }
)

export default api
