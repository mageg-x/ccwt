import api from './axios'

export const getStatus = () => api.get('/proxy/status')
export const start = (host, port) => api.post('/proxy/start', {
    host: host || undefined,
    port: port || undefined,
})
export const stop = () => api.post('/proxy/stop')
