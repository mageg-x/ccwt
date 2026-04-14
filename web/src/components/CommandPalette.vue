<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'
import { useTerminalStore } from '../stores/terminal'

const emit = defineEmits(['close', 'action'])
const router = useRouter()
const app = useAppStore()
const termStore = useTerminalStore()
const query = ref('')
const selectedIdx = ref(0)

const commands = [
    { id: 'new-term', label: '新建终端', icon: '🖥️', action: () => emit('action', 'newTab') },
    { id: 'close-term', label: '关闭当前终端', icon: '✖️', action: () => { if (termStore.activeId) termStore.removeTab(termStore.activeId) } },
    { id: 'toggle-sidebar', label: '切换侧边栏', icon: '📁', action: () => app.toggleSidebar() },
    { id: 'toggle-theme', label: '切换主题', icon: '🎨', action: () => app.toggleTheme() },
    { id: 'history', label: '会话历史', icon: '📜', action: () => router.push('/history') },
    { id: 'admin', label: '管理面板', icon: '⚙️', action: () => router.push('/admin') },
    { id: 'refresh-tree', label: '刷新文件树', icon: '🔄', action: () => emit('action', 'refreshTree') },
]

const filtered = computed(() => {
    if (!query.value) return commands
    const q = query.value.toLowerCase()
    return commands.filter(c => c.label.toLowerCase().includes(q))
})

function execute(cmd) {
    cmd.action()
    emit('close')
}

function onKeydown(e) {
    if (e.key === 'Escape') {
        emit('close')
    } else if (e.key === 'ArrowDown') {
        e.preventDefault()
        selectedIdx.value = Math.min(selectedIdx.value + 1, filtered.value.length - 1)
    } else if (e.key === 'ArrowUp') {
        e.preventDefault()
        selectedIdx.value = Math.max(selectedIdx.value - 1, 0)
    } else if (e.key === 'Enter') {
        e.preventDefault()
        if (filtered.value[selectedIdx.value]) {
            execute(filtered.value[selectedIdx.value])
        }
    }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))
</script>

<template>
    <div class="fixed inset-0 z-50 flex items-start justify-center pt-[15vh]" @click.self="emit('close')">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>
        <div class="relative w-full max-w-lg rounded-2xl shadow-2xl border overflow-hidden"
            :class="app.isDark ? 'bg-slate-800 border-slate-700' : 'bg-white border-slate-200'">
            <!-- 搜索框 -->
            <div class="flex items-center gap-3 px-4 py-3 border-b"
                :class="app.isDark ? 'border-slate-700' : 'border-slate-200'">
                <svg class="w-5 h-5 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                <input
                    v-model="query"
                    ref="input"
                    autofocus
                    placeholder="输入命令..."
                    class="flex-1 bg-transparent outline-none text-sm"
                    :class="app.isDark ? 'text-white placeholder-slate-400' : 'text-slate-800 placeholder-slate-500'"
                />
                <kbd class="px-1.5 py-0.5 text-xs rounded border"
                    :class="app.isDark ? 'bg-slate-700 border-slate-600 text-slate-400' : 'bg-slate-100 border-slate-200 text-slate-500'">
                    ESC
                </kbd>
            </div>

            <!-- 命令列表 -->
            <div class="max-h-72 overflow-y-auto py-1">
                <button
                    v-for="(cmd, idx) in filtered"
                    :key="cmd.id"
                    @click="execute(cmd)"
                    @mouseenter="selectedIdx = idx"
                    class="w-full flex items-center gap-3 px-4 py-2.5 text-sm transition-colors"
                    :class="idx === selectedIdx
                        ? (app.isDark ? 'bg-indigo-500/20 text-white' : 'bg-indigo-50 text-indigo-700')
                        : (app.isDark ? 'text-slate-300' : 'text-slate-700')"
                >
                    <span class="text-base">{{ cmd.icon }}</span>
                    <span>{{ cmd.label }}</span>
                </button>
                <div v-if="filtered.length === 0" class="px-4 py-6 text-center text-sm text-slate-400">
                    没有匹配的命令
                </div>
            </div>
        </div>
    </div>
</template>
