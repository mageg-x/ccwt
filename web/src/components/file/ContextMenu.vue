<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../../stores/app'
import { useDialogStore } from '../../stores/dialog'
import { useFileStore } from '../../stores/file'
import * as fileApi from '../../api/file'

const props = defineProps({
    x: Number,
    y: Number,
    node: Object,
})
const emit = defineEmits(['close'])
const { t } = useI18n()
const app = useAppStore()
const dialog = useDialogStore()
const fileStore = useFileStore()

function copyPath() {
    navigator.clipboard.writeText(props.node.path)
    emit('close')
}

async function newFile() {
    const name = await dialog.prompt(t('contextMenu.enterFileName'), {
        title: t('contextMenu.newFile'),
        placeholder: t('contextMenu.fileNameExample'),
        okText: t('contextMenu.create'),
    })
    if (!name || !name.trim()) return
    const cleanName = name.trim()
    const path = props.node.is_dir ? `${props.node.path}/${cleanName}` : `${props.node.path.replace(/\/[^/]+$/, '')}/${cleanName}`
    try {
        await fileApi.writeFile(path, '')
        fileStore.loadTree()
        emit('close')
    } catch (e) {
        await dialog.alert(e.response?.data?.error || t('contextMenu.createFileFailed'), { title: t('contextMenu.operationFailed') })
    }
}

async function newFolder() {
    const name = await dialog.prompt(t('contextMenu.enterFolderName'), {
        title: t('contextMenu.newFolder'),
        placeholder: t('contextMenu.folderNameExample'),
        okText: t('contextMenu.create'),
    })
    if (!name || !name.trim()) return
    const cleanName = name.trim()
    const path = props.node.is_dir ? `${props.node.path}/${cleanName}` : `${props.node.path.replace(/\/[^/]+$/, '')}/${cleanName}`
    try {
        await fileApi.mkdir(path)
        fileStore.loadTree()
        emit('close')
    } catch (e) {
        await dialog.alert(e.response?.data?.error || t('contextMenu.createFolderFailed'), { title: t('contextMenu.operationFailed') })
    }
}

async function rename() {
    const name = await dialog.prompt(t('contextMenu.enterNewName'), {
        title: t('contextMenu.rename'),
        defaultValue: props.node.name,
        placeholder: t('contextMenu.newName'),
        okText: t('contextMenu.save'),
    })
    if (!name || !name.trim()) return
    const cleanName = name.trim()
    if (cleanName === props.node.name) return
    const parts = props.node.path.split('/')
    parts[parts.length - 1] = cleanName
    try {
        await fileApi.renameFile(props.node.path, parts.join('/'))
        fileStore.loadTree()
        emit('close')
    } catch (e) {
        await dialog.alert(e.response?.data?.error || t('contextMenu.renameFailed'), { title: t('contextMenu.operationFailed') })
    }
}

async function del() {
    const ok = await dialog.confirm(t('contextMenu.confirmDelete', { name: props.node.name }), {
        title: t('contextMenu.confirmDeleteTitle'),
        okText: t('contextMenu.delete'),
    })
    if (!ok) return
    try {
        await fileApi.deleteFile(props.node.path)
        fileStore.loadTree()
        emit('close')
    } catch (e) {
        await dialog.alert(e.response?.data?.error || t('contextMenu.deleteFailed'), { title: t('contextMenu.operationFailed') })
    }
}

function download() {
    const link = document.createElement('a')
    link.href = fileApi.downloadUrl(props.node.path)
    link.style.display = 'none'
    // 通过 download 属性触发下载行为，避免打开新窗口预览
    link.download = props.node.is_dir ? `${props.node.name}.zip` : props.node.name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    emit('close')
}

function handleClickOutside(e) {
    emit('close')
}

onMounted(() => {
    setTimeout(() => document.addEventListener('click', handleClickOutside), 10)
})
onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
    <div
        class="fixed z-50 py-1 rounded-xl shadow-2xl border min-w-[160px]"
        :class="app.isDark ? 'bg-slate-800 border-slate-700' : 'bg-white border-slate-200'"
        :style="{ left: x + 'px', top: y + 'px' }"
    >
        <button @click="newFile" class="w-full text-left px-4 py-2 text-sm hover:bg-indigo-500/20 transition-colors flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
            {{ t('contextMenu.newFile') }}
        </button>
        <button @click="newFolder" class="w-full text-left px-4 py-2 text-sm hover:bg-indigo-500/20 transition-colors flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" /></svg>
            {{ t('contextMenu.newFolder') }}
        </button>
        <div class="my-1 border-t" :class="app.isDark ? 'border-slate-700' : 'border-slate-200'"></div>
        <button @click="rename" class="w-full text-left px-4 py-2 text-sm hover:bg-indigo-500/20 transition-colors flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
            {{ t('contextMenu.rename') }}
        </button>
        <button @click="copyPath" class="w-full text-left px-4 py-2 text-sm hover:bg-indigo-500/20 transition-colors flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
            {{ t('contextMenu.copyPath') }}
        </button>
        <button @click="download" class="w-full text-left px-4 py-2 text-sm hover:bg-indigo-500/20 transition-colors flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
            {{ t('contextMenu.download') }}
        </button>
        <div class="my-1 border-t" :class="app.isDark ? 'border-slate-700' : 'border-slate-200'"></div>
        <button @click="del" class="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-red-500/20 transition-colors flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
            {{ t('contextMenu.delete') }}
        </button>
    </div>
</template>
