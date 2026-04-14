import api from './axios'

export const recognize = (audioBlob) => {
    const fd = new FormData()
    fd.append('audio', audioBlob, 'recording.wav')
    return api.post('/voice/recognize', fd)
}
export const getStatus = () => api.get('/voice/status')
