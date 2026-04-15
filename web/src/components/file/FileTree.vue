<script setup>
import { computed, ref } from 'vue'
import { useAppStore } from '../../stores/app'
import { useFileStore } from '../../stores/file'
import ContextMenu from './ContextMenu.vue'

const props = defineProps({
    node: { type: Object, required: true },
    depth: { type: Number, default: 0 },
})

const emit = defineEmits(['cd', 'openFile'])
const app = useAppStore()
const fileStore = useFileStore()
const expanded = computed(() => fileStore.isExpanded(props.node.path, props.depth))
const contextMenu = ref(null) // { x, y, node }

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
