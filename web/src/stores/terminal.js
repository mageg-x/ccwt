import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

let tabCounter = 0
const STORAGE_KEY_PREFIX = 'ccwt-terminal-state'
const LEGACY_STORAGE_KEY = 'ccwt-terminal-state'

export const useTerminalStore = defineStore('terminal', () => {
    const tabs = ref([])
    const activeId = ref(null)
    const storageScope = ref('anonymous')

    const activeTab = computed(() => tabs.value.find(t => t.id === activeId.value))
    const hasInited = ref(false)

    function getStorageKey() {
        return `${STORAGE_KEY_PREFIX}:${storageScope.value}`
    }

    function clearState() {
        tabs.value = []
        activeId.value = null
    }

    function reset(scope = storageScope.value || 'anonymous') {
        storageScope.value = scope || 'anonymous'
        clearState()
        persist()
    }

    function persist() {
        localStorage.setItem(
            getStorageKey(),
            JSON.stringify({
                tabs: tabs.value,
                activeId: activeId.value,
            })
        )
    }

    function init(scope = 'anonymous') {
        const nextScope = scope || 'anonymous'
        const scopeChanged = storageScope.value !== nextScope
        if (hasInited.value && !scopeChanged) return

        storageScope.value = nextScope
        hasInited.value = true
        clearState()

        // 清理历史全局 key，避免 A/B 用户串会话
        localStorage.removeItem(LEGACY_STORAGE_KEY)

        try {
            const raw = localStorage.getItem(getStorageKey())
            if (!raw) return
            const state = JSON.parse(raw)
            if (!Array.isArray(state?.tabs)) return

            const restored = state.tabs
                .filter(t => t && typeof t.id === 'string')
                .map((t, i) => ({
                    id: t.id,
                    name: typeof t.name === 'string' && t.name.trim() ? t.name : `终端 ${i + 1}`,
                    sessionId: typeof t.sessionId === 'string' && t.sessionId ? t.sessionId : null,
                }))

            tabs.value = restored
            activeId.value = restored.some(t => t.id === state.activeId)
                ? state.activeId
                : (restored[0]?.id || null)

            for (const t of restored) {
                const n = Number((t.id.match(/^term-(\d+)$/) || [])[1] || 0)
                if (n > tabCounter) tabCounter = n
            }
        } catch {
            // ignore corrupted local storage
        }
    }

    function addTab(name) {
        tabCounter++
        const id = `term-${tabCounter}`
        const tab = {
            id,
            name: name || `终端 ${tabCounter}`,
            sessionId: null, // 由 WebSocket 连接后填充
        }
        tabs.value.push(tab)
        activeId.value = id
        persist()
        return tab
    }

    function removeTab(id) {
        const idx = tabs.value.findIndex(t => t.id === id)
        if (idx === -1) return
        tabs.value.splice(idx, 1)
        if (activeId.value === id) {
            // 切换到相邻标签
            if (tabs.value.length > 0) {
                activeId.value = tabs.value[Math.min(idx, tabs.value.length - 1)].id
            } else {
                activeId.value = null
            }
        }
        persist()
    }

    function renameTab(id, name) {
        const tab = tabs.value.find(t => t.id === id)
        if (tab) tab.name = name
        persist()
    }

    function setSession(tabId, sessionId) {
        const tab = tabs.value.find(t => t.id === tabId)
        if (tab) tab.sessionId = sessionId
        persist()
    }

    return { tabs, activeId, activeTab, init, reset, addTab, removeTab, renameTab, setSession }
})
