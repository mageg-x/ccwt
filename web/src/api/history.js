import api from './axios'

export const getProjects = () => api.get('/history/projects')
export const getSession = (file) => api.get('/history/session', { params: { file } })
