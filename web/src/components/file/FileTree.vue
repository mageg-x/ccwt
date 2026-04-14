<script setup>
import { ref } from 'vue'
import { useAppStore } from '../../stores/app'
import ContextMenu from './ContextMenu.vue'

const props = defineProps({
    node: { type: Object, required: true },
    depth: { type: Number, default: 0 },
})

const emit = defineEmits(['cd', 'openFile'])
const app = useAppStore()
const expanded = ref(props.depth < 2) // 前2层默认展开
const contextMenu = ref(null) // { x, y, node }

function toggle() {
    if (props.node.is_dir) {
        expanded.value = !expanded.value
    }
}

function handleClick() {
    if (props.node.is_dir) {
        toggle()
        emit('cd', props.node.path)
    } else {
        emit('openFile', props.node.path)
    }
}

function onContextMenu(e) {
    e.preventDefault()
    contextMenu.value = { x: e.clientX, y: e.clientY, node: props.node }
}

// 文件图标
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
            class="flex items-center gap-1.5 py-1 px-2 rounded-md cursor-pointer text-sm select-none group"
            :class="app.isDark ? 'hover:bg-slate-700/50 text-slate-300' : 'hover:bg-slate-200/80 text-slate-700'"
            :style="{ paddingLeft: (depth * 16 + 8) + 'px' }"
        >
            <!-- 展开箭头 / 文件图标 -->
            <span v-if="node.is_dir" class="w-4 text-center transition-transform" :class="{ 'rotate-90': expanded }">
                <svg class="w-3 h-3 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
            </span>
            <span v-if="node.is_dir" class="text-xs">{{ expanded ? '📂' : '📁' }}</span>
            <span v-else class="w-4"></span>
            <span v-if="!node.is_dir" class="text-xs">{{ fileIcon(node.name) }}</span>

            <span class="truncate flex-1" :title="node.name">{{ node.name }}</span>
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
