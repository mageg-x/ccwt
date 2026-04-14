import api from './axios'

export const login = (data) => api.post('/auth/login', data)
export const register = (data) => api.post('/auth/register', data)
export const logout = () => api.post('/auth/logout')
export const getMe = () => api.get('/auth/me')
export const refresh = () => api.post('/auth/refresh')
