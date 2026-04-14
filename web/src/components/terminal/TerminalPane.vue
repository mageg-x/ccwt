<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick, computed } from 'vue'
import { useTerminalStore } from '../../stores/terminal'
import { useAppStore } from '../../stores/app'
import { useTerminal } from '../../composables/useTerminal'

const props = defineProps({
    tabId: { type: String, required: true },
})

const termStore = useTerminalStore()
const app = useAppStore()
const containerRef = ref(null)
const isActive = computed(() => termStore.activeId === props.tabId)

let ws = null
let closed = false

const { term, mount, write, fit, focus } = useTerminal(containerRef, {
    onData: (data) => {
        if (ws?.readyState === WebSocket.OPEN) {
            ws.send(data)
        }
    },
    onResize: (rows, cols) => {
        if (ws?.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: 'resize', rows, cols }))
        }
    },
    theme: computed(() => app.theme),
})

function connect() {
    const tab = termStore.tabs.find(t => t.id === props.tabId)
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
    const token = document.cookie.split('; ').find(c => c.startsWith('token='))?.split('=')[1] || ''
    let url = `${proto}//${location.host}/ws/terminal?token=${token}`
    if (tab?.sessionId) {
        url += `&session_id=${tab.sessionId}`
    }

    ws = new WebSocket(url)
    ws.binaryType = 'arraybuffer'

    ws.onopen = () => {
        // 连接成功
    }

    ws.onmessage = (e) => {
        if (e.data instanceof ArrayBuffer) {
            write(e.data)
        } else {
            try {
                const msg = JSON.parse(e.data)
                if (msg.type === 'session') {
                    termStore.setSession(props.tabId, msg.data)
                } else if (msg.type === 'exit') {
                    write('\r\n\x1b[33m[会话已结束]\x1b[0m\r\n')
                }
            } catch {
                write(e.data)
            }
        }
    }

    ws.onclose = () => {
        if (!closed) {
            write('\r\n\x1b[31m[连接断开，正在重连...]\x1b[0m\r\n')
            setTimeout(connect, 2000)
        }
    }

    ws.onerror = () => {}
}

onMounted(async () => {
    await nextTick()
    mount()
    connect()
    if (isActive.value) focus()
})

onUnmounted(() => {
    closed = true
    ws?.close()
})

watch(isActive, (active) => {
    if (active) {
        nextTick(() => {
            fit()
            focus()
        })
    }
})

defineExpose({ write, focus, fit })
</script>

<template>
    <div
        v-show="isActive"
        ref="containerRef"
        class="terminal-pane w-full h-full"
    ></div>
</template>

<style scoped>
.terminal-pane {
    padding: 6px;
    box-sizing: border-box;
    border: 0;
    background: transparent;
    box-shadow: none;
}
</style>
