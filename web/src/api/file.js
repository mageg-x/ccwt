import api from './axios'

export const getTree = (path = '.') => api.get('/files/tree', { params: { path } })
export const listDir = (path = '.') => api.get('/files/list', { params: { path } })
export const readFile = (path) => api.get('/files/read', { params: { path } })
export const writeFile = (path, content) => api.post('/files/write', { path, content })
export const mkdir = (path) => api.post('/files/mkdir', { path })
export const deleteFile = (path) => api.delete('/files', { params: { path } })
export const renameFile = (oldPath, newPath) => api.post('/files/rename', { old_path: oldPath, new_path: newPath })
export const moveFile = (srcPath, dstDir) => api.post('/files/move', { src_path: srcPath, dst_dir: dstDir })
export const uploadFile = (path, file) => {
    const fd = new FormData()
    fd.append('file', file)
    fd.append('path', path)
    return api.post('/files/upload', fd)
}
export const downloadUrl = (path) => `/api/files/download?path=${encodeURIComponent(path)}`
