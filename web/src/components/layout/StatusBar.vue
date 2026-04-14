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
    <footer class="h-7 flex items-center px-3 gap-4 text-xs shrink-0 border-t select-none"
        :class="app.isDark ? 'bg-slate-800/60 border-slate-700/50 text-slate-400' : 'bg-slate-50 border-slate-200 text-slate-500'"
    >
        <span class="flex items-center gap-1">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            CPU {{ info.cpu_percent?.toFixed(1) }}%
        </span>
        <span class="flex items-center gap-1">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2 1 3 3 3h10c2 0 3-1 3-3V7c0-2-1-3-3-3H7C5 4 4 5 4 7z" />
            </svg>
            MEM {{ formatMB(info.mem_used) }}MB
        </span>
        <span class="flex-1"></span>
        <span>CCWT v1.0</span>
    </footer>
</template>
