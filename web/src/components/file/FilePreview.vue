<script setup>
import { useAppStore } from '../../stores/app'

const props = defineProps({
    path: String,
    content: String,
})
const emit = defineEmits(['close'])
const app = useAppStore()
</script>

<template>
    <div class="flex flex-col h-full">
        <div class="h-10 flex items-center px-3 gap-2 border-b shrink-0"
            :class="app.isDark ? 'bg-slate-800/60 border-slate-700/50' : 'bg-white border-slate-200'">
            <span class="text-sm truncate flex-1">{{ path }}</span>
            <button @click="emit('close')" class="p-1 rounded hover:bg-slate-700/50">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </button>
        </div>
        <pre class="flex-1 overflow-auto p-4 text-sm font-mono whitespace-pre-wrap"
            :class="app.isDark ? 'text-slate-300' : 'text-slate-700'">{{ content }}</pre>
    </div>
</template>
