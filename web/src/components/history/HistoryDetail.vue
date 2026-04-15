<script setup>
import { ref } from 'vue'
import { useAppStore } from '../../stores/app'

const props = defineProps({
    entries: { type: Array, default: () => [] },
    file: String,
})
const emit = defineEmits(['back'])
const app = useAppStore()
const search = ref('')

// 提取消息文本
function getText(entry) {
    if (!entry) return ''
    if (typeof entry.message === 'string') return entry.message
    if (typeof entry.content === 'string') return entry.content
    if (typeof entry.text === 'string') return entry.text
    if (entry.message?.content) {
        if (typeof entry.message.content === 'string') return entry.message.content
        if (Array.isArray(entry.message.content)) {
            return entry.message.content
                .filter(c => c.type === 'text')
                .map(c => c.text)
                .join('\n')
        }
    }
    const fallback = entry.message ?? entry.content ?? entry
    const asString = JSON.stringify(fallback, null, 2)
    return typeof asString === 'string' ? asString : ''
}

function isUser(entry) {
    return entry.type === 'human' || entry.type === 'user'
}

function isCode(text) {
    return typeof text === 'string' && text.includes('```')
}

// 提取代码块
function extractBlocks(text) {
    const safeText = typeof text === 'string' ? text : ''
    const parts = []
    const regex = /```(\w*)\n([\s\S]*?)```/g
    let last = 0
    let match
    while ((match = regex.exec(safeText)) !== null) {
        if (match.index > last) {
            parts.push({ type: 'text', content: safeText.slice(last, match.index) })
        }
        parts.push({ type: 'code', lang: match[1], content: match[2] })
        last = match.index + match[0].length
    }
    if (last < safeText.length) {
        parts.push({ type: 'text', content: safeText.slice(last) })
    }
    return parts.length ? parts : [{ type: 'text', content: safeText }]
}

function copyCode(code) {
    navigator.clipboard.writeText(code)
}

const filtered = ref(null)
function doSearch() {
    if (!search.value) {
        filtered.value = null
        return
    }
    const q = search.value.toLowerCase()
    filtered.value = props.entries.filter((e) => {
        const text = getText(e)
        return text.toLowerCase().includes(q)
    })
}

const displayEntries = ref(null)
</script>

<template>
    <div class="flex flex-col h-full">
        <!-- 标题栏 -->
        <div class="flex items-center gap-3 px-4 py-3 border-b shrink-0"
            :class="app.isDark ? 'border-slate-700/50' : 'border-slate-200'">
            <button @click="emit('back')" class="p-1.5 rounded-lg hover:bg-slate-700/50 transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
            </button>
            <span class="text-sm font-mono truncate flex-1">{{ file }}</span>
        </div>

        <!-- 搜索 -->
        <div class="px-4 py-2 border-b shrink-0" :class="app.isDark ? 'border-slate-700/50' : 'border-slate-200'">
            <input v-model="search" @input="doSearch" placeholder="搜索会话内容..."
                class="w-full px-3 py-2 rounded-lg text-sm bg-transparent border outline-none"
                :class="app.isDark ? 'border-slate-600 focus:border-indigo-500' : 'border-slate-300 focus:border-indigo-500'" />
        </div>

        <!-- 消息列表 -->
        <div class="flex-1 overflow-y-auto p-4 space-y-4">
            <div v-for="(entry, idx) in (filtered || entries)" :key="idx"
                class="flex" :class="isUser(entry) ? 'justify-end' : 'justify-start'">
                <div class="max-w-[85%] rounded-2xl px-4 py-3 text-sm"
                    :class="isUser(entry)
                        ? 'bg-indigo-600 text-white rounded-br-md'
                        : (app.isDark ? 'bg-slate-800 border border-slate-700/50 rounded-bl-md' : 'bg-slate-100 border border-slate-200 rounded-bl-md')">

                    <!-- 渲染内容，支持代码块 -->
                    <template v-for="(block, bi) in extractBlocks(getText(entry))" :key="bi">
                        <p v-if="block.type === 'text'" class="whitespace-pre-wrap break-words">{{ block.content }}</p>
                        <div v-else class="my-2 rounded-lg overflow-hidden"
                            :class="app.isDark ? 'bg-slate-900' : 'bg-slate-800'">
                            <div class="flex items-center justify-between px-3 py-1.5 text-xs text-slate-400">
                                <span>{{ block.lang || 'code' }}</span>
                                <button @click="copyCode(block.content)"
                                    class="hover:text-white transition-colors flex items-center gap-1">
                                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                                    </svg>
                                    复制
                                </button>
                            </div>
                            <pre class="px-3 pb-3 text-xs overflow-x-auto text-green-300"><code>{{ block.content }}</code></pre>
                        </div>
                    </template>
                </div>
            </div>

            <div v-if="entries.length === 0" class="text-center py-12 text-slate-400 text-sm">
                会话为空
            </div>
        </div>
    </div>
</template>
