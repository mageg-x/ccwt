import api from './axios'

export const getStatus = () => api.get('/proxy/status')
export const start = () => api.post('/proxy/start')
export const stop = () => api.post('/proxy/stop')
