<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const isRegister = ref(false)
const username = ref('')
const password = ref('')
const inviteCode = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
    error.value = ''
    loading.value = true
    try {
        if (isRegister.value) {
            await auth.register(username.value, password.value, inviteCode.value)
        } else {
            await auth.login(username.value, password.value)
        }
        router.push('/')
    } catch (e) {
        error.value = e.response?.data?.error || '操作失败'
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-900 p-4">
        <!-- 背景装饰 -->
        <div class="fixed inset-0 overflow-hidden pointer-events-none">
            <div class="absolute top-1/4 -left-20 w-72 h-72 bg-indigo-500/10 rounded-full blur-3xl"></div>
            <div class="absolute bottom-1/4 -right-20 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl"></div>
        </div>

        <div class="w-full max-w-md relative z-10">
            <!-- Logo -->
            <div class="text-center mb-8">
                <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-indigo-500/20 border border-indigo-500/30 mb-4">
                    <svg class="w-8 h-8 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                </div>
                <h1 class="text-3xl font-bold text-white tracking-tight">CCWT</h1>
                <p class="text-slate-400 mt-1">Claude Code Web Terminal</p>
            </div>

            <!-- 表单卡片 -->
            <div class="bg-slate-800/50 backdrop-blur-xl border border-slate-700/50 rounded-2xl p-8 shadow-2xl">
                <h2 class="text-xl font-semibold text-white mb-6">
                    {{ isRegister ? '创建账户' : '登录' }}
                </h2>

                <form @submit.prevent="submit" class="space-y-5">
                    <div>
                        <label class="block text-sm font-medium text-slate-300 mb-2">用户名</label>
                        <input
                            v-model="username"
                            type="text"
                            required
                            autocomplete="username"
                            class="w-full px-4 py-3 bg-slate-900/50 border border-slate-600/50 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500/50 transition-all"
                            placeholder="输入用户名"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-slate-300 mb-2">密码</label>
                        <input
                            v-model="password"
                            type="password"
                            required
                            autocomplete="current-password"
                            class="w-full px-4 py-3 bg-slate-900/50 border border-slate-600/50 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500/50 transition-all"
                            placeholder="输入密码（至少6位）"
                        />
                    </div>
                    <div v-if="isRegister">
                        <label class="block text-sm font-medium text-slate-300 mb-2">邀请码</label>
                        <input
                            v-model="inviteCode"
                            type="text"
                            autocomplete="off"
                            class="w-full px-4 py-3 bg-slate-900/50 border border-slate-600/50 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500/50 transition-all"
                            placeholder="输入邀请码（如需要）"
                        />
                    </div>

                    <div v-if="error" class="p-3 bg-red-500/10 border border-red-500/30 rounded-xl text-red-400 text-sm">
                        {{ error }}
                    </div>

                    <button
                        type="submit"
                        :disabled="loading"
                        class="w-full py-3 px-4 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-xl transition-all duration-200 shadow-lg shadow-indigo-500/25 hover:shadow-indigo-500/40"
                    >
                        <span v-if="loading" class="inline-flex items-center gap-2">
                            <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                            </svg>
                            处理中...
                        </span>
                        <span v-else>{{ isRegister ? '注册' : '登录' }}</span>
                    </button>
                </form>

                <div class="mt-6 text-center">
                    <button
                        @click="isRegister = !isRegister; error = ''"
                        class="text-indigo-400 hover:text-indigo-300 text-sm transition-colors"
                    >
                        {{ isRegister ? '已有账户？去登录' : '没有账户？去注册' }}
                    </button>
                </div>
            </div>

            <p class="text-center text-slate-500 text-xs mt-6">
                首个注册的用户将自动成为管理员
            </p>
        </div>
    </div>
</template>
