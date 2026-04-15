<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'
import { useDialogStore } from '../stores/dialog'
import * as settingsApi from '../api/settings'
import * as proxyApi from '../api/proxy'

const router = useRouter()
const app = useAppStore()
const dialog = useDialogStore()

const activeTab = ref('voice')
const loading = ref(false)
const settings = ref({})
const proxyRunning = ref(false)
const proxyAddress = ref('')
const proxyIP = ref('0.0.0.0')
const proxyPort = ref(1080)

const tabs = [
    { key: 'voice', label: '语音识别', icon: 'M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z' },
    { key: 'proxy', label: 'SOCKS5 代理', icon: 'M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9' },
    { key: 'about', label: '关于', icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z' },
]

async function loadSettings() {
    loading.value = true
    try {
        const { data } = await settingsApi.getSettings()
        const map = {}
        for (const s of data.settings || []) {
            map[s.key] = s.value
        }
        settings.value = map
    } catch (e) {
        console.error('加载设置失败', e)
    } finally {
        loading.value = false
    }
}

async function loadProxyStatus() {
    try {
        const { data } = await proxyApi.getStatus()
        proxyRunning.value = data.running
        proxyAddress.value = data.address || ''
        proxyIP.value = data.bind_host || settings.value['proxy.ip'] || '0.0.0.0'
        proxyPort.value = data.port || parseInt(settings.value['proxy.port']) || 1080
    } catch (e) {
        console.error('获取代理状态失败', e)
    }
}

async function saveSetting(key, value) {
    try {
        await settingsApi.updateSetting(key, String(value))
    } catch (e) {
        await dialog.alert(e.response?.data?.error || '保存失败', { title: '保存失败' })
    }
}

async function toggleVoice() {
    const newVal = settings.value['voice.enabled'] !== 'true' && settings.value['voice.enabled'] !== '1'
    await saveSetting('voice.enabled', newVal)
    settings.value['voice.enabled'] = String(newVal)
}

async function toggleProxy() {
    if (proxyRunning.value) {
        try {
            await proxyApi.stop()
        } catch (e) {
            await dialog.alert(e.response?.data?.error || '停止失败', { title: '停止失败' })
        }
    } else {
        try {
            await proxyApi.start(settings.value['proxy.ip'] || '0.0.0.0', proxyPort.value)
        } catch (e) {
            await dialog.alert(e.response?.data?.error || '启动失败', { title: '启动失败' })
        }
    }
    await loadProxyStatus()
}

async function saveProxyIP() {
    await saveSetting('proxy.ip', proxyIP.value)
}

async function saveVoiceConfig() {
    await saveSetting('voice.app_id', settings.value['voice.app_id'])
    await saveSetting('voice.api_key', settings.value['voice.api_key'])
    await saveSetting('voice.secret', settings.value['voice.secret'])
}

async function saveProxyPort() {
    await saveSetting('proxy.port', proxyPort.value)
}

async function copyProxyAddr() {
    const displayIP = proxyIP.value === '0.0.0.0' ? '127.0.0.1' : proxyIP.value
    navigator.clipboard.writeText(`socks5://${displayIP}:${proxyPort.value}`)
}

onMounted(() => {
    loadSettings()
    loadProxyStatus()
})
</script>

<template>
    <div class="h-full flex flex-col" :class="app.isDark ? 'bg-slate-900 text-slate-200' : 'bg-white text-slate-800'">
        <header class="h-12 flex items-center px-4 gap-3 border-b shrink-0"
            :class="app.isDark ? 'bg-slate-800/80 border-slate-700/50' : 'bg-white border-slate-200'">
            <button @click="router.push('/')" class="p-2 rounded-lg hover:bg-slate-700/50 transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
            </button>
            <h1 class="font-semibold">设置</h1>
        </header>

        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
            <div class="max-w-2xl mx-auto">
                <div class="flex gap-1 p-1 rounded-xl mb-6"
                    :class="app.isDark ? 'bg-slate-800/50' : 'bg-slate-100'">
                    <button v-for="tab in tabs" :key="tab.key" @click="activeTab = tab.key"
                        class="flex-1 flex items-center justify-center gap-2 py-2.5 px-3 rounded-lg text-sm font-medium transition-all"
                        :class="activeTab === tab.key
                            ? app.isDark ? 'bg-slate-700 text-white shadow' : 'bg-white text-slate-800 shadow-sm'
                            : app.isDark ? 'text-slate-400 hover:text-slate-200' : 'text-slate-500 hover:text-slate-700'">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="tab.icon" />
                        </svg>
                        <span class="hidden sm:inline">{{ tab.label }}</span>
                    </button>
                </div>

                <div v-if="loading" class="flex items-center justify-center py-12">
                    <svg class="animate-spin w-6 h-6 text-indigo-400" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                </div>

                <div v-else>
                    <div v-show="activeTab === 'voice'" class="space-y-6">
                        <div class="rounded-xl border p-4" :class="app.isDark ? 'border-slate-700/50 bg-slate-800/30' : 'border-slate-200 bg-slate-50'">
                            <div class="flex items-center justify-between mb-4">
                                <div>
                                    <h3 class="font-medium">语音识别开关</h3>
                                    <p class="text-xs mt-1" :class="app.isDark ? 'text-slate-400' : 'text-slate-500'">启用后可用语音输入命令</p>
                                </div>
                                <button @click="toggleVoice"
                                    class="relative w-12 h-6 rounded-full transition-colors"
                                    :class="settings['voice.enabled'] === 'true' || settings['voice.enabled'] === '1' ? 'bg-indigo-500' : 'bg-slate-600'">
                                    <span class="absolute top-1 w-4 h-4 rounded-full bg-white shadow transition-transform"
                                        :class="settings['voice.enabled'] === 'true' || settings['voice.enabled'] === '1' ? 'left-7' : 'left-1'" />
                                </button>
                            </div>
                        </div>

                        <div class="rounded-xl border p-4" :class="app.isDark ? 'border-slate-700/50 bg-slate-800/30' : 'border-slate-200 bg-slate-50'">
                            <h3 class="font-medium mb-4">百度语音 API 配置</h3>
                            <p class="text-xs mb-4" :class="app.isDark ? 'text-slate-400' : 'text-slate-500'">
                                请前往 <a href="https://console.bce.baidu.com/" target="_blank" class="text-indigo-400 hover:underline">百度智能云控制台</a> 创建语音应用，获取 App ID、API Key 和 Secret Key
                            </p>
                            <div class="space-y-3">
                                <div>
                                    <label class="block text-xs mb-1.5" :class="app.isDark ? 'text-slate-300' : 'text-slate-600'">App ID</label>
                                    <input type="text" v-model="settings['voice.app_id']"
                                        class="w-full px-3 py-2 rounded-lg border text-sm font-mono outline-none"
                                        :class="app.isDark ? 'bg-slate-900/60 border-slate-600 text-slate-200 focus:border-indigo-500' : 'bg-white border-slate-300 text-slate-700 focus:border-indigo-500'" />
                                </div>
                                <div>
                                    <label class="block text-xs mb-1.5" :class="app.isDark ? 'text-slate-300' : 'text-slate-600'">API Key</label>
                                    <input type="text" v-model="settings['voice.api_key']"
                                        class="w-full px-3 py-2 rounded-lg border text-sm font-mono outline-none"
                                        :class="app.isDark ? 'bg-slate-900/60 border-slate-600 text-slate-200 focus:border-indigo-500' : 'bg-white border-slate-300 text-slate-700 focus:border-indigo-500'" />
                                </div>
                                <div>
                                    <label class="block text-xs mb-1.5" :class="app.isDark ? 'text-slate-300' : 'text-slate-600'">Secret Key</label>
                                    <input type="password" v-model="settings['voice.secret']"
                                        class="w-full px-3 py-2 rounded-lg border text-sm font-mono outline-none"
                                        :class="app.isDark ? 'bg-slate-900/60 border-slate-600 text-slate-200 focus:border-indigo-500' : 'bg-white border-slate-300 text-slate-700 focus:border-indigo-500'" />
                                </div>
                            </div>
                            <div class="flex justify-end mt-4">
                                <button @click="saveVoiceConfig"
                                    class="px-4 py-2 rounded-lg text-sm bg-indigo-600 text-white hover:bg-indigo-500 transition-colors">
                                    保存
                                </button>
                            </div>
                        </div>
                    </div>

                    <div v-show="activeTab === 'proxy'" class="space-y-6">
                        <div class="rounded-xl border p-4" :class="app.isDark ? 'border-slate-700/50 bg-slate-800/30' : 'border-slate-200 bg-slate-50'">
                            <div class="flex items-center justify-between mb-4">
                                <div>
                                    <h3 class="font-medium">SOCKS5 代理服务</h3>
                                    <p class="text-xs mt-1" :class="app.isDark ? 'text-slate-400' : 'text-slate-500'">
                                        {{ proxyRunning ? '运行中' : '已停止' }}
                                        <span v-if="proxyRunning" class="ml-2 font-mono text-indigo-400">{{ proxyAddress }}</span>
                                    </p>
                                </div>
                                <button @click="toggleProxy"
                                    class="relative w-12 h-6 rounded-full transition-colors"
                                    :class="proxyRunning ? 'bg-indigo-500' : 'bg-slate-600'">
                                    <span class="absolute top-1 w-4 h-4 rounded-full bg-white shadow transition-transform"
                                        :class="proxyRunning ? 'left-7' : 'left-1'" />
                                </button>
                            </div>
                        </div>

                        <div class="rounded-xl border p-4" :class="app.isDark ? 'border-slate-700/50 bg-slate-800/30' : 'border-slate-200 bg-slate-50'">
                            <h3 class="font-medium mb-4">默认绑定配置</h3>
                            <div class="flex items-center gap-3 flex-wrap">
                                <div class="flex items-center gap-2">
                                    <label class="text-xs whitespace-nowrap" :class="app.isDark ? 'text-slate-300' : 'text-slate-600'">IP</label>
                                    <input type="text" v-model="proxyIP"
                                        class="w-28 px-2 py-1.5 rounded-lg border text-sm font-mono outline-none shrink-0"
                                        :class="app.isDark ? 'bg-slate-900/60 border-slate-600 text-slate-200 focus:border-indigo-500' : 'bg-white border-slate-300 text-slate-700 focus:border-indigo-500'" />
                                </div>
                                <div class="flex items-center gap-2">
                                    <label class="text-xs whitespace-nowrap" :class="app.isDark ? 'text-slate-300' : 'text-slate-600'">端口</label>
                                    <input type="number" v-model.number="proxyPort" min="1" max="65535"
                                        class="w-20 px-2 py-1.5 rounded-lg border text-sm font-mono outline-none shrink-0"
                                        :class="app.isDark ? 'bg-slate-900/60 border-slate-600 text-slate-200 focus:border-indigo-500' : 'bg-white border-slate-300 text-slate-700 focus:border-indigo-500'" />
                                </div>
                                <button @click="saveProxyIP(); saveProxyPort()"
                                    class="ml-auto px-4 py-1.5 rounded-lg text-sm bg-indigo-600 text-white hover:bg-indigo-500 transition-colors shrink-0">
                                    保存
                                </button>
                            </div>
                            <p class="text-xs mt-2" :class="app.isDark ? 'text-slate-500' : 'text-slate-400'">
                                0.0.0.0 表示绑定所有网卡；127.0.0.1 仅允许本地连接
                            </p>
                        </div>

                        <div v-if="proxyRunning" class="rounded-xl border p-4" :class="app.isDark ? 'border-slate-700/50 bg-slate-800/30' : 'border-slate-200 bg-slate-50'">
                            <h3 class="font-medium mb-3">连接信息</h3>
                            <div class="flex items-center justify-between p-3 rounded-lg"
                                :class="app.isDark ? 'bg-slate-900/50' : 'bg-slate-100'">
                                <code class="font-mono text-sm text-indigo-400">socks5://{{ proxyIP === '0.0.0.0' ? '127.0.0.1' : proxyIP }}:{{ proxyPort }}</code>
                                <button @click="copyProxyAddr"
                                    class="px-3 py-1.5 rounded-lg text-xs transition-colors"
                                    :class="app.isDark ? 'hover:bg-slate-700' : 'hover:bg-slate-200'">
                                    复制
                                </button>
                            </div>
                        </div>
                    </div>

                    <div v-show="activeTab === 'about'" class="space-y-6">
                        <div class="rounded-xl border p-6 text-center" :class="app.isDark ? 'border-slate-700/50 bg-slate-800/30' : 'border-slate-200 bg-slate-50'">
                            <div class="w-16 h-16 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white text-2xl font-bold shadow-lg">
                                C
                            </div>
                            <h2 class="text-xl font-semibold mb-1">CCWT</h2>
                            <p class="text-sm" :class="app.isDark ? 'text-slate-400' : 'text-slate-500'">Claude Code Web Terminal</p>
                            <p class="text-xs mt-2" :class="app.isDark ? 'text-slate-500' : 'text-slate-400'">版本 1.0.0</p>
                        </div>

                        <div class="rounded-xl border p-4" :class="app.isDark ? 'border-slate-700/50 bg-slate-800/30' : 'border-slate-200 bg-slate-50'">
                            <h3 class="font-medium mb-3">功能说明</h3>
                            <ul class="space-y-2 text-sm" :class="app.isDark ? 'text-slate-300' : 'text-slate-600'">
                                <li class="flex items-start gap-2">
                                    <svg class="w-4 h-4 mt-0.5 text-green-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                                    </svg>
                                    基于 Web 的 Claude Code 终端
                                </li>
                                <li class="flex items-start gap-2">
                                    <svg class="w-4 h-4 mt-0.5 text-green-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                                    </svg>
                                    支持语音输入命令
                                </li>
                                <li class="flex items-start gap-2">
                                    <svg class="w-4 h-4 mt-0.5 text-green-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                                    </svg>
                                    内置 SOCKS5 代理
                                </li>
                                <li class="flex items-start gap-2">
                                    <svg class="w-4 h-4 mt-0.5 text-green-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                                    </svg>
                                    文件管理和在线编辑
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>