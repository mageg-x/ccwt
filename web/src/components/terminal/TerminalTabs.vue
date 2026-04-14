<script setup>
import { useAppStore } from '../../stores/app'
import { useTerminalStore } from '../../stores/terminal'

const app = useAppStore()
const termStore = useTerminalStore()
const emit = defineEmits(['newTab'])

function closeTab(e, id) {
    e.stopPropagation()
    termStore.removeTab(id)
}

function dblClick(tab) {
    const name = prompt('终端名称:', tab.name)
    if (name) termStore.renameTab(tab.id, name)
}
</script>

<template>
    <div class="h-9 flex items-end gap-0.5 px-2 overflow-x-auto shrink-0 scrollbar-none"
        :class="app.isDark ? 'bg-slate-900/50' : 'bg-slate-100'">
        <div
            v-for="tab in termStore.tabs"
            :key="tab.id"
            @click="termStore.activeId = tab.id"
            @dblclick="dblClick(tab)"
            @keydown.enter.prevent="termStore.activeId = tab.id"
            @keydown.space.prevent="termStore.activeId = tab.id"
            role="button"
            tabindex="0"
            class="flex items-center gap-1.5 px-3 py-1.5 text-xs rounded-t-lg transition-colors max-w-[160px] group"
            :class="tab.id === termStore.activeId
                ? (app.isDark ? 'bg-slate-800 text-white' : 'bg-white text-slate-800')
                : (app.isDark ? 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50' : 'text-slate-500 hover:text-slate-700 hover:bg-slate-200/50')"
        >
            <svg class="w-3 h-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span class="truncate">{{ tab.name }}</span>
            <button @click="closeTab($event, tab.id)"
                class="p-0.5 rounded opacity-0 group-hover:opacity-100 hover:bg-red-500/30 transition-all">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </button>
        </div>

        <!-- 新建标签 -->
        <button @click="emit('newTab')"
            class="p-1.5 rounded-t-lg transition-colors"
            :class="app.isDark ? 'text-slate-400 hover:text-white hover:bg-slate-800/50' : 'text-slate-400 hover:text-slate-700 hover:bg-slate-200/50'"
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
