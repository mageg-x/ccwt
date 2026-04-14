<script setup>
import { onMounted } from 'vue'
import { useAppStore } from '../../stores/app'
import { useFileStore } from '../../stores/file'
import FileTree from '../file/FileTree.vue'

const app = useAppStore()
const fileStore = useFileStore()

onMounted(() => {
    fileStore.loadTree()
})

const emit = defineEmits(['cd', 'openFile'])
</script>

<template>
    <Transition name="slide">
        <aside v-show="app.sidebarOpen"
            class="flex flex-col border-r overflow-hidden shrink-0"
            :class="[
                app.isDark ? 'bg-slate-800/40 border-slate-700/50' : 'bg-slate-50 border-slate-200',
                app.isMobile ? 'fixed inset-y-12 left-0 z-40 w-72 shadow-2xl' : 'w-40'
            ]"
        >
            <!-- 标题栏 -->
            <div class="h-10 flex items-center px-3 gap-2 border-b shrink-0"
                :class="app.isDark ? 'border-slate-700/50' : 'border-slate-200'">
                <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
                <span class="text-sm font-medium flex-1">文件</span>
                <button @click="fileStore.loadTree()" class="p-1 rounded hover:bg-slate-700/50 transition-colors" title="刷新">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                </button>
            </div>

            <!-- 文件树 -->
            <div class="flex-1 overflow-y-auto p-1">
                <div v-if="fileStore.loading" class="flex items-center justify-center py-8">
                    <svg class="animate-spin w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                </div>
                <FileTree
                    v-else-if="fileStore.tree"
                    :node="fileStore.tree"
                    :depth="0"
                    @cd="(p) => emit('cd', p)"
                    @open-file="(p) => emit('openFile', p)"
                />
                <div v-else class="text-sm text-slate-500 text-center py-8">
                    工作区为空
                </div>
            </div>
        </aside>
    </Transition>

    <!-- 移动端遮罩 -->
    <div v-if="app.isMobile && app.sidebarOpen"
        @click="app.sidebarOpen = false"
        class="fixed inset-0 bg-black/50 z-30"
    ></div>
</template>
