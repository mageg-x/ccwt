<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick, computed } from 'vue'
import { useTerminalStore } from '../../stores/terminal'
import { useAppStore } from '../../stores/app'
import { useFileStore } from '../../stores/file'
import { useTerminal } from '../../composables/useTerminal'

const props = defineProps({
    tabId: { type: String, required: true },
})

const termStore = useTerminalStore()
const app = useAppStore()
const fileStore = useFileStore()
const containerRef = ref(null)
const isActive = computed(() => termStore.activeId === props.tabId)

let ws = null
let closed = false
let refreshTimer = null
let activateRaf = null
let activateFallbackTimer = null
let openResizeTimer = null

function scheduleTreeRefresh() {
    if (refreshTimer) clearTimeout(refreshTimer)
    refreshTimer = setTimeout(() => {
        fileStore.loadTree()
    }, 280)
}

const { term, mount, write, fit, focus } = useTerminal(containerRef, {
    onData: (data) => {
        if (ws?.readyState === WebSocket.OPEN) {
            ws.send(data)
            if (data.includes('\r')) {
                scheduleTreeRefresh()
            }
        }
    },
    onResize: (rows, cols) => {
        if (ws?.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: 'resize', rows, cols }))
        }
    },
    // 终端优先 dark/shell，避免浅色主题下 CLI 配色可读性问题
    theme: computed(() => app.theme === 'light' ? 'dark' : app.theme),
})

function syncPtySize() {
    fit()
    if (ws?.readyState !== WebSocket.OPEN) return
    const rawRows = term.value?.rows
    const rawCols = term.value?.cols
    // 防止异常测量值导致后端列数过大/过小，出现输入不换行
    const rows = Math.min(300, Math.max(8, Number(rawRows) || 0))
    const cols = Math.min(400, Math.max(40, Number(rawCols) || 0))
    if (!rows || !cols) return
    ws.send(JSON.stringify({ type: 'resize', rows, cols }))
}

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
        // 连接建立后强制同步一次尺寸，避免后端仍停留在默认 80 列
        if (openResizeTimer) clearTimeout(openResizeTimer)
        openResizeTimer = setTimeout(() => {
            syncPtySize()
            openResizeTimer = null
        }, 40)
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

function sendInput(data) {
    if (!data) return
    if (ws?.readyState === WebSocket.OPEN) {
        ws.send(data)
        if (data.includes('\r')) {
            scheduleTreeRefresh()
        }
    }
}

onMounted(async () => {
    await nextTick()
    mount()
    connect()
    if (isActive.value) focus()
})

onUnmounted(() => {
    closed = true
    if (refreshTimer) clearTimeout(refreshTimer)
    if (openResizeTimer) clearTimeout(openResizeTimer)
    if (activateRaf) cancelAnimationFrame(activateRaf)
    if (activateFallbackTimer) clearTimeout(activateFallbackTimer)
    ws?.close()
})

function syncActiveTerminal() {
    if (!isActive.value) return
    syncPtySize()
    focus()
}

watch(isActive, (active) => {
    if (!active) return
    if (activateRaf) cancelAnimationFrame(activateRaf)
    if (activateFallbackTimer) clearTimeout(activateFallbackTimer)

    nextTick(() => {
        activateRaf = requestAnimationFrame(() => {
            syncActiveTerminal()
            activateRaf = null
        })
        // 某些浏览器在标签切换时 RAF 可能延后，补一次轻量兜底
        activateFallbackTimer = setTimeout(() => {
            syncActiveTerminal()
            activateFallbackTimer = null
        }, 80)
    })
})

defineExpose({ sendInput, focus, fit })
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
