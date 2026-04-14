<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'
import * as systemApi from '../api/system'

const router = useRouter()
const app = useAppStore()
const users = ref([])
const loading = ref(true)

onMounted(async () => {
    await loadUsers()
})

async function loadUsers() {
    loading.value = true
    try {
        const { data } = await systemApi.getUsers()
        users.value = data.users || []
    } catch (e) {
        console.error('加载用户列表失败', e)
    } finally {
        loading.value = false
    }
}

async function toggleRole(user) {
    const newRole = user.role === 'admin' ? 'user' : 'admin'
    if (!confirm(`确定将 ${user.username} 的角色改为 ${newRole}？`)) return
    try {
        await systemApi.updateRole(user.id, newRole)
        await loadUsers()
    } catch (e) {
        alert(e.response?.data?.error || '操作失败')
    }
}

async function delUser(user) {
    if (!confirm(`确定删除用户 ${user.username}？此操作不可恢复。`)) return
    try {
        await systemApi.deleteUser(user.id)
        await loadUsers()
    } catch (e) {
        alert(e.response?.data?.error || '删除失败')
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
            <h1 class="font-semibold">管理面板</h1>
        </header>

        <!-- 用户管理 -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
            <div class="max-w-4xl mx-auto">
                <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
                    <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                    </svg>
                    用户管理
                </h2>

                <div v-if="loading" class="flex items-center justify-center py-12">
                    <svg class="animate-spin w-6 h-6 text-indigo-400" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                </div>

                <div v-else class="rounded-xl border overflow-hidden"
                    :class="app.isDark ? 'border-slate-700/50' : 'border-slate-200'">
                    <table class="w-full">
                        <thead>
                            <tr :class="app.isDark ? 'bg-slate-800/50' : 'bg-slate-50'">
                                <th class="text-left px-4 py-3 text-xs font-medium uppercase tracking-wider text-slate-400">ID</th>
                                <th class="text-left px-4 py-3 text-xs font-medium uppercase tracking-wider text-slate-400">用户名</th>
                                <th class="text-left px-4 py-3 text-xs font-medium uppercase tracking-wider text-slate-400">角色</th>
                                <th class="text-left px-4 py-3 text-xs font-medium uppercase tracking-wider text-slate-400">注册时间</th>
                                <th class="text-right px-4 py-3 text-xs font-medium uppercase tracking-wider text-slate-400">操作</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y" :class="app.isDark ? 'divide-slate-700/30' : 'divide-slate-100'">
                            <tr v-for="user in users" :key="user.id"
                                class="transition-colors"
                                :class="app.isDark ? 'hover:bg-slate-800/30' : 'hover:bg-slate-50'">
                                <td class="px-4 py-3 text-sm">{{ user.id }}</td>
                                <td class="px-4 py-3 text-sm font-medium">
                                    <div class="flex items-center gap-2">
                                        <div class="w-7 h-7 rounded-full bg-indigo-500/20 flex items-center justify-center text-xs font-bold text-indigo-400">
                                            {{ user.username[0]?.toUpperCase() }}
                                        </div>
                                        {{ user.username }}
                                    </div>
                                </td>
                                <td class="px-4 py-3 text-sm">
                                    <span class="px-2 py-0.5 rounded-full text-xs font-medium"
                                        :class="user.role === 'admin'
                                            ? 'bg-amber-500/20 text-amber-400'
                                            : 'bg-slate-500/20 text-slate-400'">
                                        {{ user.role }}
                                    </span>
                                </td>
                                <td class="px-4 py-3 text-sm text-slate-400">{{ user.created_at }}</td>
                                <td class="px-4 py-3 text-sm text-right">
                                    <div class="flex items-center justify-end gap-2">
                                        <button @click="toggleRole(user)"
                                            class="px-2.5 py-1 rounded-lg text-xs transition-colors"
                                            :class="app.isDark ? 'hover:bg-slate-700' : 'hover:bg-slate-200'">
                                            {{ user.role === 'admin' ? '降为用户' : '升为管理员' }}
                                        </button>
                                        <button v-if="user.id !== 1" @click="delUser(user)"
                                            class="px-2.5 py-1 rounded-lg text-xs text-red-400 hover:bg-red-500/20 transition-colors">
                                            删除
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
