import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

let tabCounter = 0

export const useTerminalStore = defineStore('terminal', () => {
    const tabs = ref([])
    const activeId = ref(null)

    const activeTab = computed(() => tabs.value.find(t => t.id === activeId.value))

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
    }

    function renameTab(id, name) {
        const tab = tabs.value.find(t => t.id === id)
        if (tab) tab.name = name
    }

    function setSession(tabId, sessionId) {
        const tab = tabs.value.find(t => t.id === tabId)
        if (tab) tab.sessionId = sessionId
    }

    return { tabs, activeId, activeTab, addTab, removeTab, renameTab, setSession }
})
