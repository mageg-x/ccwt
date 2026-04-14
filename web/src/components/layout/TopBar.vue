<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../../stores/app'
import { useAuthStore } from '../../stores/auth'
import VoiceInput from '../VoiceInput.vue'
import ProxyPanel from '../ProxyPanel.vue'

const router = useRouter()
const app = useAppStore()
const auth = useAuthStore()

const showVoice = ref(false)
const showProxy = ref(false)
const showUserMenu = ref(false)

const emit = defineEmits(['voiceResult'])

async function doLogout() {
    await auth.logout()
    router.push('/login')
}
</script>

<template>
    <header class="h-12 flex items-center px-3 gap-2 border-b shrink-0 select-none"
        :class="app.isDark ? 'bg-slate-800/80 border-slate-700/50 text-slate-200' : 'bg-white border-slate-200 text-slate-800'"
    >
        <!-- 汉堡菜单 -->
        <button @click="app.toggleSidebar" class="p-2 rounded-lg hover:bg-slate-700/50 transition-colors" title="切换侧边栏">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
        </button>

        <!-- Logo -->
        <span class="font-bold text-indigo-400 tracking-tight hidden sm:inline">CCWT</span>

        <div class="flex-1"></div>

        <!-- 代理开关 -->
        <button v-if="auth.isAdmin" @click="showProxy = !showProxy"
            class="p-2 rounded-lg hover:bg-slate-700/50 transition-colors relative" title="SOCKS5 代理">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
            </svg>
        </button>

        <!-- 语音输入 -->
        <button @click="showVoice = !showVoice"
            class="p-2 rounded-lg hover:bg-slate-700/50 transition-colors" title="语音输入">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
            </svg>
        </button>

        <!-- 主题切换 -->
        <button @click="app.toggleTheme"
            class="p-2 rounded-lg hover:bg-slate-700/50 transition-colors" title="切换主题">
            <svg v-if="app.isDark" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
        </button>

        <!-- 历史 -->
        <button @click="router.push('/history')"
            class="p-2 rounded-lg hover:bg-slate-700/50 transition-colors" title="会话历史">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
        </button>

        <!-- 用户菜单 -->
        <div class="relative">
            <button @click="showUserMenu = !showUserMenu"
                class="flex items-center gap-2 p-2 rounded-lg hover:bg-slate-700/50 transition-colors">
                <div class="w-7 h-7 rounded-full bg-indigo-500/30 flex items-center justify-center text-xs font-bold text-indigo-300">
                    {{ auth.user?.username?.[0]?.toUpperCase() }}
                </div>
                <span class="text-sm hidden sm:inline">{{ auth.user?.username }}</span>
            </button>

            <Transition name="fade">
                <div v-if="showUserMenu" @click="showUserMenu = false"
                    class="absolute right-0 top-full mt-1 w-48 py-1 rounded-xl shadow-xl z-50 border"
                    :class="app.isDark ? 'bg-slate-800 border-slate-700' : 'bg-white border-slate-200'"
                >
                    <div class="px-4 py-2 text-xs text-slate-400 border-b" :class="app.isDark ? 'border-slate-700' : 'border-slate-200'">
                        {{ auth.user?.role === 'admin' ? '管理员' : '用户' }}
                    </div>
                    <button v-if="auth.isAdmin" @click="router.push('/admin')"
                        class="w-full text-left px-4 py-2 text-sm hover:bg-slate-700/50 transition-colors">
                        管理面板
                    </button>
                    <button @click="doLogout"
                        class="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-slate-700/50 transition-colors">
                        退出登录
                    </button>
                </div>
            </Transition>
        </div>
    </header>

    <!-- 语音输入弹窗 -->
    <VoiceInput v-if="showVoice" @close="showVoice = false" @result="(t) => { emit('voiceResult', t); showVoice = false }" />

    <!-- 代理面板 -->
    <ProxyPanel v-if="showProxy" @close="showProxy = false" />
</template>
