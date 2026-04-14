<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAppStore } from '../../stores/app'
import * as systemApi from '../../api/system'

const app = useAppStore()
const info = ref({ cpu_percent: 0, mem_used: 0, mem_total: 0 })
let timer = null

function formatMB(bytes) {
    return (bytes / 1024 / 1024).toFixed(0)
}

async function fetchInfo() {
    try {
        const { data } = await systemApi.getInfo()
        info.value = data.system
    } catch { /* 忽略 */ }
}

onMounted(() => {
    fetchInfo()
    timer = setInterval(fetchInfo, 10000) // 每10秒更新
})
onUnmounted(() => clearInterval(timer))
</script>

<template>
    <footer class="h-8 flex items-center px-3 md:px-4 gap-4 text-xs shrink-0 border-t select-none backdrop-blur-xl"
        :class="app.isDark ? 'bg-slate-900/45 border-slate-700/30 text-slate-300' : 'bg-white/80 border-slate-200/90 text-slate-600'"
    >
        <span class="flex items-center gap-1.5 px-2 py-0.5 rounded-md"
            :class="app.isDark ? 'bg-white/5 text-cyan-200/90' : 'bg-slate-100 text-slate-700'">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            CPU {{ info.cpu_percent?.toFixed(1) }}%
        </span>
        <span class="flex items-center gap-1.5 px-2 py-0.5 rounded-md"
            :class="app.isDark ? 'bg-white/5 text-emerald-200/90' : 'bg-slate-100 text-slate-700'">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2 1 3 3 3h10c2 0 3-1 3-3V7c0-2-1-3-3-3H7C5 4 4 5 4 7z" />
            </svg>
            MEM {{ formatMB(info.mem_used) }}MB
        </span>
        <span class="flex-1"></span>
        <span class="uppercase tracking-[0.08em] text-[10px] opacity-80">CCWT v1.0</span>
    </footer>
</template>
