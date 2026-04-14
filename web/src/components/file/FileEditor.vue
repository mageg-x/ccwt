<script setup>
import { ref, onMounted, watch, shallowRef } from 'vue'
import { useAppStore } from '../../stores/app'
import { useFileStore } from '../../stores/file'
import * as monaco from 'monaco-editor'

const app = useAppStore()
const fileStore = useFileStore()
const container = ref(null)
const editor = shallowRef(null)
const saving = ref(false)

onMounted(() => {
    if (!container.value || !fileStore.editingFile) return
    editor.value = monaco.editor.create(container.value, {
        value: fileStore.editingFile.content,
        language: fileStore.editingFile.language,
        theme: app.isDark ? 'vs-dark' : 'vs',
        fontSize: 14,
        fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
        minimap: { enabled: false },
        lineNumbers: 'on',
        wordWrap: 'on',
        scrollBeyondLastLine: false,
        automaticLayout: true,
        padding: { top: 8 },
    })

    // Ctrl+S 保存
    editor.value.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, save)
})

watch(() => app.isDark, (dark) => {
    monaco.editor.setTheme(dark ? 'vs-dark' : 'vs')
})

async function save() {
    if (!editor.value || !fileStore.editingFile) return
    saving.value = true
    try {
        const content = editor.value.getValue()
        await fileStore.saveFile(fileStore.editingFile.path, content)
    } finally {
        saving.value = false
    }
}
</script>

<template>
    <div v-if="fileStore.editingFile" class="flex flex-col h-full">
        <!-- 编辑器标题栏 -->
        <div class="h-10 flex items-center px-3 gap-2 border-b shrink-0"
            :class="app.isDark ? 'bg-slate-800/60 border-slate-700/50' : 'bg-white border-slate-200'">
            <span class="text-sm truncate flex-1">{{ fileStore.editingFile.path }}</span>
            <button @click="save" :disabled="saving"
                class="px-3 py-1 text-xs bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg transition-colors disabled:opacity-50">
                {{ saving ? '保存中...' : '保存' }}
            </button>
            <button @click="fileStore.closeEditor()"
                class="p-1 rounded hover:bg-slate-700/50 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </button>
        </div>
        <div ref="container" class="flex-1"></div>
    </div>
</template>
