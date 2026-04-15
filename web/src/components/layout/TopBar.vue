<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../../stores/app'
import { useAuthStore } from '../../stores/auth'
import { useTerminalStore } from '../../stores/terminal'
import VoiceInput from '../VoiceInput.vue'

const router = useRouter()
const { t, locale } = useI18n()
const app = useAppStore()
const auth = useAuthStore()
const termStore = useTerminalStore()

const showVoice = ref(false)
const showUserMenu = ref(false)
const showLangMenu = ref(false)

const languages = [
    { code: 'zh', label: '中文' },
    { code: 'en', label: 'English' }
]

function setLanguage(code) {
    locale.value = code
    localStorage.setItem('locale', code)
    showLangMenu.value = false
}

const emit = defineEmits(['voiceResult'])

async function doLogout() {
    termStore.reset('anonymous')
    await auth.logout()
    router.push('/login')
}

// 点击其他地方关闭下拉菜单
function handleClickOutside(event) {
    // 检查点击是否在语言菜单或用户菜单内部
    const langMenu = document.querySelector('.lang-menu-container')
    const userMenu = document.querySelector('.user-menu-container')
    
    if (langMenu && !langMenu.contains(event.target)) {
        showLangMenu.value = false
    }
    
    if (userMenu && !userMenu.contains(event.target)) {
        showUserMenu.value = false
    }
}

onMounted(() => {
    document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
    <header class="h-14 flex items-center px-3 md:px-4 gap-2 border-b shrink-0 select-none backdrop-blur-xl topbar-shell relative z-20"
        :class="app.isDark ? 'bg-slate-900/55 border-slate-700/40 text-slate-200' : 'bg-white/80 border-slate-200/90 text-slate-800'"
    >
        <!-- 汉堡菜单 -->
        <button @click="app.toggleSidebar" class="p-2 rounded-xl hover:bg-slate-700/40 transition-colors" :title="t('nav.toggleSidebar')">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
        </button>

        <!-- Logo -->
        <div class="hidden sm:flex items-center gap-2">
            <img src="/src/assets/imgs/logo.png" alt="CCWT Logo" class="w-7 h-7 object-contain" />
            <span class="font-semibold tracking-[0.08em]"
                :class="app.isDark ? 'text-slate-100/95' : 'text-slate-800'">CCWT</span>
        </div>

        <div class="flex-1"></div>

        <!-- 语音输入 -->
        <button @click="showVoice = !showVoice"
            class="p-2 rounded-xl hover:bg-slate-700/40 transition-colors" :title="t('nav.voiceInput')">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
            </svg>
        </button>

        <!-- 主题切换 -->
        <button @click="app.toggleTheme"
            class="p-2 rounded-xl hover:bg-slate-700/40 transition-colors" :title="t('nav.toggleTheme')">
            <svg v-if="app.isDark" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
        </button>

        <!-- 语言切换 -->
        <div class="relative lang-menu-container">
            <button @click.stop="showLangMenu = !showLangMenu"
                class="p-2 rounded-xl hover:bg-slate-700/40 transition-colors" :title="t('nav.language')">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129" />
                </svg>
            </button>

            <Transition name="fade">
                <div v-if="showLangMenu" @click.stop
                    class="absolute right-0 top-full mt-1 w-32 py-1 rounded-xl shadow-xl z-50 border"
                    :class="app.isDark ? 'bg-slate-800 border-slate-700' : 'bg-white border-slate-200'"
                >
                    <button v-for="lang in languages" :key="lang.code"
                        @click="setLanguage(lang.code)"
                        class="w-full text-left px-4 py-2 text-sm hover:bg-slate-700/50 transition-colors flex items-center justify-between"
                        :class="locale === lang.code ? 'text-indigo-400' : ''"
                    >
                        {{ lang.label }}
                        <span v-if="locale === lang.code" class="text-indigo-400">✓</span>
                    </button>
                </div>
            </Transition>
        </div>

        <!-- 历史 -->
        <button @click="router.push('/history')"
            class="p-2 rounded-xl hover:bg-slate-700/40 transition-colors" :title="t('nav.sessionHistory')">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
        </button>

        <!-- 设置 -->
        <button @click="router.push('/settings')"
            class="p-2 rounded-xl hover:bg-slate-700/40 transition-colors" :title="t('nav.settings')">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
        </button>

        <!-- 用户菜单 -->
        <div class="relative user-menu-container">
            <button @click.stop="showUserMenu = !showUserMenu"
                class="flex items-center gap-2 p-2 rounded-xl hover:bg-slate-700/40 transition-colors">
                <div class="w-7 h-7 rounded-full bg-indigo-500/30 flex items-center justify-center text-xs font-bold text-indigo-300">
                    {{ auth.user?.username?.[0]?.toUpperCase() }}
                </div>
                <span class="text-sm hidden sm:inline">{{ auth.user?.username }}</span>
            </button>

            <Transition name="fade">
                <div v-if="showUserMenu" @click.stop
                    class="absolute right-0 top-full mt-1 w-48 py-1 rounded-xl shadow-xl z-50 border"
                    :class="app.isDark ? 'bg-slate-800 border-slate-700' : 'bg-white border-slate-200'"
                >
                    <div class="px-4 py-2 text-xs text-slate-400 border-b" :class="app.isDark ? 'border-slate-700' : 'border-slate-200'">
                        {{ auth.user?.role === 'admin' ? t('user.admin') : t('user.user') }}
                    </div>
                    <button v-if="auth.isAdmin" @click="router.push('/admin')"
                        class="w-full text-left px-4 py-2 text-sm hover:bg-slate-700/50 transition-colors">
                        {{ t('user.adminPanel') }}
                    </button>
                    <button @click="doLogout"
                        class="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-slate-700/50 transition-colors">
                        {{ t('nav.logout') }}
                    </button>
                </div>
            </Transition>
        </div>
    </header>

    <!-- 语音输入弹窗 -->
    <VoiceInput v-if="showVoice" @close="showVoice = false" @result="(t) => { emit('voiceResult', t); showVoice = false }" />
</template>
