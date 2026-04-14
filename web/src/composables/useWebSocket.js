import { ref, onUnmounted } from 'vue'

export function useWebSocket(url, { onMessage, onOpen, onClose, onError, autoReconnect = true } = {}) {
    const ws = ref(null)
    const connected = ref(false)
    let reconnTimer = null
    let closed = false

    function connect() {
        if (closed) return
        const socket = new WebSocket(url)
        socket.binaryType = 'arraybuffer'
        ws.value = socket

        socket.onopen = () => {
            connected.value = true
            onOpen?.()
        }
        socket.onmessage = (e) => {
            onMessage?.(e)
        }
        socket.onclose = (e) => {
            connected.value = false
            onClose?.(e)
            if (autoReconnect && !closed) {
                reconnTimer = setTimeout(connect, 2000)
            }
        }
        socket.onerror = (e) => {
            onError?.(e)
        }
    }

    function send(data) {
        if (ws.value?.readyState === WebSocket.OPEN) {
            ws.value.send(data)
        }
    }

    function sendJSON(obj) {
        send(JSON.stringify(obj))
    }

    function close() {
        closed = true
        clearTimeout(reconnTimer)
        ws.value?.close()
    }

    connect()

    onUnmounted(() => {
        close()
    })

    return { ws, connected, send, sendJSON, close }
}
