<script setup>
import { computed, ref } from 'vue'
import { useAppStore } from '../../stores/app'
import { useDialogStore } from '../../stores/dialog'
import { useFileStore } from '../../stores/file'
import * as fileApi from '../../api/file'
import ContextMenu from './ContextMenu.vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
    node: { type: Object, required: true },
    depth: { type: Number, default: 0 },
})

const emit = defineEmits(['cd', 'openFile'])
const app = useAppStore()
const dialog = useDialogStore()
const fileStore = useFileStore()
const { t } = useI18n()
const expanded = computed(() => fileStore.isExpanded(props.node.path, props.depth))
const contextMenu = ref(null) // { x, y, node }
const dragOver = ref(false)

function toggle() {
    if (props.node.is_dir) {
        fileStore.toggleExpanded(props.node.path, props.depth)
    }
}

function handleClick() {
    if (!props.node.is_dir) {
        emit('openFile', props.node.path)
    }
}

function handleDirCd() {
    if (props.node.is_dir) {
        emit('cd', props.node.path)
    }
}

function onContextMenu(e) {
    e.preventDefault()
    contextMenu.value = { x: e.clientX, y: e.clientY, node: props.node }
}

function onDragStart(e) {
    if (!e.dataTransfer) return
    const payload = {
        path: props.node.path,
        name: props.node.name,
        isDir: props.node.is_dir,
    }
    fileStore.setDraggingNode(payload)
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', props.node.path)
    e.dataTransfer.setData('application/x-ccwt-node', JSON.stringify(payload))
    // 隐藏默认拖拽预览，避免出现巨大方块
    const ghost = document.createElement('canvas')
    ghost.width = 1
    ghost.height = 1
    e.dataTransfer.setDragImage(ghost, 0, 0)
}

function onDragEnd() {
    dragOver.value = false
    fileStore.setDraggingNode(null)
}

function canDropOnCurrentNode(sourcePath, sourceIsDir) {
    if (!props.node.is_dir || !sourcePath) return false
    if (sourcePath === props.node.path) return false
    if (sourceIsDir && props.node.path.startsWith(`${sourcePath}/`)) return false
    return true
}

function onDragOver(e) {
    if (!props.node.is_dir) return
    const sourceMeta = fileStore.draggingNode
    // 外部拖拽（桌面文件）：统一允许 copy drop，兼容各浏览器类型差异
    if (!sourceMeta) {
        e.preventDefault()
        e.dataTransfer.dropEffect = 'copy'
        dragOver.value = true
        return
    }

    const sourcePath = sourceMeta?.path || ''
    if (!canDropOnCurrentNode(sourcePath, !!sourceMeta?.isDir)) return
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
    dragOver.value = true
}

function onDragLeave() {
    dragOver.value = false
}

async function onDrop(e) {
    dragOver.value = false
    if (!props.node.is_dir) return
    e.preventDefault()

    const { files: droppedFiles, hasDirectory } = collectDroppedFiles(e.dataTransfer)
    if (droppedFiles.length > 0) {
        await uploadDroppedFiles(droppedFiles)
        if (hasDirectory) {
            await dialog.alert(t('fileTree.dirUploadIgnored'), { title: t('fileTree.uploadNoticeTitle') })
        }
        return
    }
    if (hasDirectory) {
        await dialog.alert(t('fileTree.dirUploadNotSupported'), { title: t('fileTree.uploadNoticeTitle') })
        return
    }

    const sourceMeta = fileStore.draggingNode
    const sourcePath = sourceMeta?.path || e.dataTransfer?.getData('text/plain') || ''
    if (!sourcePath || !sourceMeta) return

    if (!canDropOnCurrentNode(sourcePath, !!sourceMeta?.isDir)) {
        await dialog.alert(t('fileTree.invalidMoveTarget'), { title: t('fileTree.moveFailedTitle') })
        return
    }

    try {
        await fileApi.moveFile(sourcePath, props.node.path)
        await fileStore.loadTree()
        fileStore.setExpanded(props.node.path, true)
    } catch (err) {
        await dialog.alert(err?.response?.data?.error || t('fileTree.moveFailedMessage'), { title: t('fileTree.moveFailedTitle') })
    }
}

function collectDroppedFiles(dataTransfer) {
    if (!dataTransfer) return { files: [], hasDirectory: false }
    const directFiles = Array.from(dataTransfer.files || [])
    if (directFiles.length > 0) {
        return { files: directFiles, hasDirectory: false }
    }

    const out = []
    let hasDirectory = false
    const items = Array.from(dataTransfer.items || [])
    if (items.length > 0) {
        for (const item of items) {
            if (item.kind !== 'file') continue
            const entry = item.webkitGetAsEntry ? item.webkitGetAsEntry() : null
            if (entry && entry.isDirectory) {
                hasDirectory = true
                continue
            }
            const file = item.getAsFile ? item.getAsFile() : null
            if (file) out.push(file)
        }
    }
    return { files: out, hasDirectory }
}

async function uploadDroppedFiles(files) {
    const tasks = files.map((file) => fileApi.uploadFile(props.node.path, file))
    const results = await Promise.allSettled(tasks)
    const failed = results.filter((r) => r.status === 'rejected').length
    await fileStore.loadTree()
    fileStore.setExpanded(props.node.path, true)
    if (failed > 0) {
        await dialog.alert(t('fileTree.uploadPartialFailed', { count: failed }), { title: t('fileTree.uploadResultTitle') })
    }
}

function fileIcon(name) {
    const ext = name.split('.').pop()?.toLowerCase()
    const icons = {
        js: '📄', ts: '📄', vue: '💚', py: '🐍', go: '🔵',
        json: '📋', yaml: '📋', yml: '📋', md: '📝', html: '🌐',
        css: '🎨', sh: '⚙️', sql: '🗄️', txt: '📄',
    }
    return icons[ext] || '📄'
}

</script>

<template>
    <div>
        <div
            @click="handleClick"
            @contextmenu="onContextMenu"
            @touchstart.passive=""
            :draggable="true"
            @dragstart="onDragStart"
            @dragend="onDragEnd"
            @dragover="onDragOver"
            @dragleave="onDragLeave"
            @drop="onDrop"
            class="flex items-center gap-1.5 py-1 px-2 rounded-md cursor-pointer text-sm select-none group"
            :class="[
                app.isDark ? 'hover:bg-slate-700/50 text-slate-300' : 'hover:bg-slate-200/80 text-slate-700',
                dragOver ? (app.isDark ? 'bg-indigo-500/20 ring-1 ring-indigo-400/50' : 'bg-indigo-100 ring-1 ring-indigo-300') : '',
            ]"
            :style="{ paddingLeft: (depth * 16 + 8) + 'px' }"
        >
            <!-- 同一列：目录箭头 / 文件图标（左对齐） -->
            <button
                v-if="node.is_dir"
                type="button"
                @click.stop="toggle"
                class="w-4 flex items-center justify-start transition-transform"
                :class="{ 'rotate-90': expanded }"
                :title="expanded ? '收起目录' : '展开目录'"
            >
                <svg class="w-3 h-3 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
            </button>
            <span v-else class="w-4 text-xs leading-none flex items-center justify-start">{{ fileIcon(node.name) }}</span>

            <span
                v-if="node.is_dir"
                @click.stop="handleDirCd"
                class="truncate flex-1 hover:underline font-semibold"
                :title="node.name"
            >{{ node.name }}</span>
            <span
                v-else
                class="truncate flex-1 font-normal"
                :title="node.name"
            >{{ node.name }}</span>
        </div>

        <!-- 子节点 -->
        <template v-if="node.is_dir && expanded && node.children">
            <FileTree
                v-for="child in node.children"
                :key="child.path"
                :node="child"
                :depth="depth + 1"
                @cd="(p) => emit('cd', p)"
                @open-file="(p) => emit('openFile', p)"
            />
        </template>

        <!-- 右键菜单 -->
        <ContextMenu
            v-if="contextMenu"
            :x="contextMenu.x"
            :y="contextMenu.y"
            :node="contextMenu.node"
            @close="contextMenu = null"
        />
    </div>
</template>
