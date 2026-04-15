import api from './axios'

export const getSettings = () => api.get('/settings')
export const getSetting = (key) => api.get('/settings/get', { params: { key } })
export const updateSetting = (key, value) => api.put('/settings', { key, value })