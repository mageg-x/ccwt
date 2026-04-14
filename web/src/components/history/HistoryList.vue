<script setup>
import { useAppStore } from '../../stores/app'

const props = defineProps({
    projects: { type: Array, default: () => [] },
})
const emit = defineEmits(['select'])
const app = useAppStore()

function formatTime(t) {
    return new Date(t).toLocaleString('zh-CN')
}
</script>

<template>
    <div class="space-y-3">
        <div v-for="proj in projects" :key="proj.project"
            class="rounded-xl border overflow-hidden"
            :class="app.isDark ? 'border-slate-700/50 bg-slate-800/30' : 'border-slate-200 bg-white'">
            <!-- 项目标题 -->
            <div class="px-4 py-3 border-b flex items-center gap-2"
                :class="app.isDark ? 'border-slate-700/50' : 'border-slate-200'">
                <svg class="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
                <span class="font-medium text-sm">{{ proj.project }}</span>
                <span class="text-xs text-slate-400 ml-auto">{{ proj.sessions.length }} 个会话</span>
            </div>

            <!-- 会话列表 -->
            <div class="divide-y" :class="app.isDark ? 'divide-slate-700/30' : 'divide-slate-100'">
                <button
                    v-for="sess in proj.sessions.slice(0, 10)"
                    :key="sess.file"
                    @click="emit('select', sess.file)"
                    class="w-full text-left px-4 py-2.5 flex items-center gap-3 transition-colors"
                    :class="app.isDark ? 'hover:bg-slate-700/30' : 'hover:bg-slate-50'"
                >
                    <svg class="w-4 h-4 text-slate-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                    </svg>
                    <span class="text-sm truncate flex-1 font-mono text-xs">{{ sess.file }}</span>
                    <span class="text-xs text-slate-400 shrink-0">{{ formatTime(sess.mod_time) }}</span>
                </button>
            </div>
        </div>

        <div v-if="projects.length === 0" class="text-center py-12 text-slate-400">
            <svg class="w-12 h-12 mx-auto mb-3 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <p class="text-sm">暂无会话历史</p>
        </div>
    </div>
</template>
