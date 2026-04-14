import api from './axios'

export const getInfo = () => api.get('/system/info')
export const getUsers = () => api.get('/admin/users')
export const deleteUser = (id) => api.delete(`/admin/users/${id}`)
export const updateRole = (id, role) => api.put(`/admin/users/${id}/role`, { role })
