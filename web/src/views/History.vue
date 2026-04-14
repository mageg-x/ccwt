<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'
import * as historyApi from '../api/history'
import HistoryList from '../components/history/HistoryList.vue'
import HistoryDetail from '../components/history/HistoryDetail.vue'

const router = useRouter()
const app = useAppStore()
const projects = ref([])
const loading = ref(true)
const selectedFile = ref(null)
const sessionEntries = ref([])
const detailLoading = ref(false)

onMounted(async () => {
    try {
        const { data } = await historyApi.getProjects()
        projects.value = data.projects || []
    } catch (e) {
        console.error('加载历史失败', e)
    } finally {
        loading.value = false
    }
})

async function selectSession(file) {
    selectedFile.value = file
    detailLoading.value = true
    try {
        const { data } = await historyApi.getSession(file)
        sessionEntries.value = data.entries || []
    } catch {
        sessionEntries.value = []
    } finally {
        detailLoading.value = false
    }
}
</script>

<template>
    <div class="h-full flex flex-col" :class="app.isDark ? 'bg-slate-900 text-slate-200' : 'bg-white text-slate-800'">
        <!-- 顶部栏 -->
        <header class="h-12 flex items-center px-4 gap-3 border-b shrink-0"
            :class="app.isDark ? 'bg-slate-800/80 border-slate-700/50' : 'bg-white border-slate-200'">
            <button @click="router.push('/')" class="p-2 rounded-lg hover:bg-slate-700/50 transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
            </button>
            <h1 class="font-semibold">会话历史</h1>
        </header>

        <!-- 内容 -->
        <div class="flex-1 overflow-hidden">
            <!-- 详情视图 -->
            <HistoryDetail
                v-if="selectedFile"
                :entries="sessionEntries"
                :file="selectedFile"
                @back="selectedFile = null"
            />

            <!-- 列表视图 -->
            <div v-else class="h-full overflow-y-auto p-4">
                <div v-if="loading" class="flex items-center justify-center py-12">
                    <svg class="animate-spin w-6 h-6 text-indigo-400" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                </div>
                <HistoryList v-else :projects="projects" @select="selectSession" />
            </div>
        </div>
    </div>
</template>
