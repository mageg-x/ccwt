<script setup>
import { useAppStore } from '../../stores/app'

const app = useAppStore()
const emit = defineEmits(['key'])

const keys = [
    { label: 'Esc', value: '\x1b' },
    { label: 'Tab', value: '\t' },
    { label: 'Ctrl', value: null, modifier: true },
    { label: '↑', value: '\x1b[A' },
    { label: '↓', value: '\x1b[B' },
    { label: '←', value: '\x1b[D' },
    { label: '→', value: '\x1b[C' },
    { label: 'Ctrl+C', value: '\x03' },
    { label: 'Ctrl+D', value: '\x04' },
    { label: 'Ctrl+Z', value: '\x1a' },
    { label: 'Ctrl+L', value: '\x0c' },
    { label: 'Ctrl+R', value: '\x12' },
]

function sendKey(key) {
    if (key.value !== null) {
        emit('key', key.value)
    }
}
</script>

<template>
    <div v-if="app.isMobile" class="flex gap-1 px-2 py-1.5 overflow-x-auto shrink-0 border-t"
        :class="app.isDark ? 'bg-slate-800/80 border-slate-700/50' : 'bg-slate-100 border-slate-200'">
        <button
            v-for="key in keys"
            :key="key.label"
            @click="sendKey(key)"
            class="px-2.5 py-1.5 rounded-lg text-xs font-mono whitespace-nowrap transition-colors active:scale-95"
            :class="app.isDark
                ? 'bg-slate-700/80 text-slate-300 active:bg-indigo-600'
                : 'bg-white text-slate-600 active:bg-indigo-500 active:text-white border border-slate-200'"
        >
            {{ key.label }}
        </button>
    </div>
</template>
