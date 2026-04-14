<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'
import { useTerminalStore } from '../stores/terminal'
import { useFileStore } from '../stores/file'
import TopBar from '../components/layout/TopBar.vue'
import StatusBar from '../components/layout/StatusBar.vue'
import Sidebar from '../components/layout/Sidebar.vue'
import TerminalTabs from '../components/terminal/TerminalTabs.vue'
import TerminalPane from '../components/terminal/TerminalPane.vue'
import VirtualKeys from '../components/terminal/VirtualKeys.vue'
import FileEditor from '../components/file/FileEditor.vue'
import CommandPalette from '../components/CommandPalette.vue'

const auth = useAuthStore()
const app = useAppStore()
const termStore = useTerminalStore()
const fileStore = useFileStore()
const termRefs = ref({})

onMounted(async () => {
    await auth.fetchMe()
    // 自动创建第一个终端
    if (termStore.tabs.length === 0) {
        newTab()
    }
})

function newTab() {
    termStore.addTab()
}

// 文件树点击目录 → 向活跃终端发送 cd 命令
function handleCd(path) {
    const active = termStore.activeTab
    if (active) {
        const termPane = termRefs.value[active.id]
        if (termPane) {
            termPane.write(`cd ${path}\r`)
        }
    }
    // 移动端自动关闭侧边栏
    if (app.isMobile) {
        app.sidebarOpen = false
    }
}

// 文件树点击文件 → 打开编辑器
function handleOpenFile(path) {
    fileStore.openFile(path)
    if (app.isMobile) {
        app.sidebarOpen = false
    }
}

// 语音识别结果 → 输入到活跃终端
function handleVoiceResult(text) {
    const active = termStore.activeTab
    if (active) {
        const termPane = termRefs.value[active.id]
        if (termPane) {
            termPane.write(text)
        }
    }
}

// 虚拟按键 → 输入到活跃终端
function handleVirtualKey(key) {
    const active = termStore.activeTab
    if (active) {
        const termPane = termRefs.value[active.id]
        if (termPane) {
            termPane.write(key)
            termPane.focus()
        }
    }
}

// 命令面板操作
function handleCmdAction(action) {
    if (action === 'newTab') newTab()
    else if (action === 'refreshTree') fileStore.loadTree()
}

// Cmd+K 快捷键
function onKeydown(e) {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault()
        app.toggleCmdPalette()
    }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))
</script>

<template>
    <div class="h-full flex flex-col" :class="app.isDark ? 'bg-slate-900 text-slate-200' : 'bg-white text-slate-800'">
        <TopBar @voiceResult="handleVoiceResult" />

        <div class="flex flex-1 overflow-hidden">
            <Sidebar @cd="handleCd" @openFile="handleOpenFile" />

            <!-- 主内容区 -->
            <main class="flex-1 flex flex-col overflow-hidden min-w-0">
                <!-- 文件编辑器（悬浮覆盖在终端上方） -->
                <div v-if="fileStore.editingFile" class="flex-1 flex flex-col overflow-hidden">
                    <FileEditor />
                </div>

                <!-- 终端区域 -->
                <div v-show="!fileStore.editingFile" class="flex-1 flex flex-col overflow-hidden">
                    <TerminalTabs @newTab="newTab" />

                    <!-- 终端面板 -->
                    <div class="flex-1 relative overflow-hidden" :class="app.isDark ? 'bg-slate-900' : 'bg-white'">
                        <TerminalPane
                            v-for="tab in termStore.tabs"
                            :key="tab.id"
                            :tabId="tab.id"
                            :ref="el => { if (el) termRefs[tab.id] = el }"
                        />

                        <!-- 空状态 -->
                        <div v-if="termStore.tabs.length === 0"
                            class="absolute inset-0 flex flex-col items-center justify-center gap-4 text-slate-400">
                            <svg class="w-16 h-16 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                            </svg>
                            <p class="text-sm">没有打开的终端</p>
                            <button @click="newTab"
                                class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white text-sm rounded-xl transition-colors">
                                新建终端
                            </button>
                        </div>
                    </div>

                    <!-- 移动端虚拟按键 -->
                    <VirtualKeys @key="handleVirtualKey" />
                </div>
            </main>
        </div>

        <StatusBar />

        <!-- 快捷指令面板 -->
        <CommandPalette
            v-if="app.cmdPaletteOpen"
            @close="app.cmdPaletteOpen = false"
            @action="handleCmdAction"
        />
    </div>
</template>
