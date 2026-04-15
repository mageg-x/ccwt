<script setup>
import { useAppStore } from '../../stores/app'
import { useTerminalStore } from '../../stores/terminal'
import { useDialogStore } from '../../stores/dialog'

const app = useAppStore()
const termStore = useTerminalStore()
const dialog = useDialogStore()
const emit = defineEmits(['newTab'])

function closeTab(e, id) {
    e.stopPropagation()
    termStore.removeTab(id)
}

async function dblClick(tab) {
    const name = await dialog.prompt('请输入终端名称', {
        title: '重命名终端',
        defaultValue: tab.name,
        placeholder: '终端名称',
        okText: '保存',
    })
    if (name && name.trim()) termStore.renameTab(tab.id, name.trim())
}
</script>

<template>
    <div class="h-11 flex items-center gap-1 px-2 overflow-x-auto shrink-0 scrollbar-none border-b"
        :class="app.isDark ? 'bg-slate-900/40 border-slate-700/40' : 'bg-white/70 border-slate-200/90'">
        <div
            v-for="tab in termStore.tabs"
            :key="tab.id"
            @click="termStore.activeId = tab.id"
            @dblclick="dblClick(tab)"
            @keydown.enter.prevent="termStore.activeId = tab.id"
            @keydown.space.prevent="termStore.activeId = tab.id"
            role="button"
            tabindex="0"
            class="flex items-center gap-1.5 px-3 py-1.5 text-xs rounded-xl transition-all max-w-[180px] group border"
            :class="tab.id === termStore.activeId
                ? (app.isDark ? 'bg-cyan-400/10 text-cyan-100 border-cyan-300/40 shadow-[0_0_24px_rgba(34,211,238,0.15)]' : 'bg-cyan-50 text-cyan-800 border-cyan-300/60')
                : (app.isDark ? 'text-slate-400 border-white/10 hover:text-slate-100 hover:bg-white/5' : 'text-slate-500 border-slate-200 hover:text-slate-700 hover:bg-slate-100')"
        >
            <svg class="w-3 h-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span class="truncate">{{ tab.name }}</span>
            <button @click="closeTab($event, tab.id)"
                class="p-0.5 rounded-md opacity-0 group-hover:opacity-100 hover:bg-red-500/30 transition-all">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </button>
        </div>

        <!-- 新建标签 -->
        <button @click="emit('newTab')"
            class="p-2 rounded-xl transition-colors border"
            :class="app.isDark ? 'text-slate-300 border-white/10 hover:text-white hover:bg-white/5' : 'text-slate-500 border-slate-200 hover:text-slate-700 hover:bg-slate-100'"
            title="新建终端">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
        </button>
    </div>
</template>

<style scoped>
.scrollbar-none::-webkit-scrollbar { display: none; }
.scrollbar-none { -ms-overflow-style: none; scrollbar-width: none; }
</style>
