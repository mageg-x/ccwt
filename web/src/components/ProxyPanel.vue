<script setup>
import { ref, onMounted } from 'vue'
import { useAppStore } from '../stores/app'
import { useDialogStore } from '../stores/dialog'
import * as proxyApi from '../api/proxy'

const emit = defineEmits(['close'])
const app = useAppStore()
const dialog = useDialogStore()
const running = ref(false)
const address = ref('')
const ip = ref('127.0.0.1')
const port = ref(1080)
const loading = ref(false)

async function fetchStatus() {
    const { data } = await proxyApi.getStatus()
    running.value = data.running
    address.value = data.address
    ip.value = data.bind_host || data.client_ip || '127.0.0.1'
    if (typeof data.port === 'number' && data.port > 0) {
        port.value = data.port
    }
}

async function toggle() {
    loading.value = true
    try {
        if (running.value) {
            await proxyApi.stop()
        } else {
            await proxyApi.start(ip.value, port.value)
        }
        await fetchStatus()
    } catch (e) {
        await dialog.alert(e.response?.data?.error || '操作失败', { title: '代理操作失败' })
    } finally {
        loading.value = false
    }
}

function copyAddr() {
    navigator.clipboard.writeText(`socks5://${ip.value}:${port.value}`)
}

onMounted(fetchStatus)
</script>

<template>
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="emit('close')">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>
        <div class="relative w-full max-w-md rounded-2xl shadow-2xl border p-6"
            :class="app.isDark ? 'bg-slate-800 border-slate-700' : 'bg-white border-slate-200'">

            <h3 class="text-lg font-semibold mb-4 flex items-center gap-2">
                <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                </svg>
                SOCKS5 认证代理
            </h3>

            <!-- 状态 -->
            <div class="flex items-center gap-3 mb-4 p-3 rounded-xl"
                :class="app.isDark ? 'bg-slate-700/50' : 'bg-slate-100'">
                <div class="w-3 h-3 rounded-full" :class="running ? 'bg-green-400 animate-pulse' : 'bg-slate-500'"></div>
                <span class="text-sm">{{ running ? '运行中' : '已停止' }}</span>
                <span v-if="running && address" class="text-xs text-slate-400 ml-auto font-mono">{{ address }}</span>
            </div>

            <div class="mb-4 p-3 rounded-xl space-y-3"
                :class="app.isDark ? 'bg-slate-700/30' : 'bg-slate-100'">
                <div class="text-xs font-medium" :class="app.isDark ? 'text-slate-300' : 'text-slate-600'">浏览器 SOCKS5 配置</div>
                <div class="grid grid-cols-[72px_1fr] gap-2 items-center text-sm">
                    <span class="text-slate-400">IP</span>
                    <input v-model.trim="ip" :disabled="running || loading"
                        class="px-2.5 py-1.5 rounded-lg border font-mono"
                        :class="app.isDark
                            ? 'bg-slate-900/60 border-slate-600 text-slate-200 focus:border-indigo-500'
                            : 'bg-white border-slate-300 text-slate-700 focus:border-indigo-500'" />
                    <span class="text-slate-400">端口</span>
                    <input v-model.number="port" type="number" min="1" max="65535" :disabled="running || loading"
                        class="px-2.5 py-1.5 rounded-lg border font-mono outline-none"
                        :class="app.isDark
                            ? 'bg-slate-900/60 border-slate-600 text-slate-200 focus:border-indigo-500'
                            : 'bg-white border-slate-300 text-slate-700 focus:border-indigo-500'" />
                </div>
            </div>

            <!-- 操作 -->
            <div class="flex gap-3 mb-4">
                <button @click="toggle" :disabled="loading"
                    class="flex-1 py-2.5 rounded-xl text-sm font-medium transition-all"
                    :class="running
                        ? 'bg-red-500/20 text-red-400 hover:bg-red-500/30'
                        : 'bg-indigo-600 text-white hover:bg-indigo-500'">
                    {{ loading ? '处理中...' : running ? '停止代理' : '启动代理' }}
                </button>
                <button v-if="running" @click="copyAddr"
                    class="px-4 py-2.5 rounded-xl text-sm border transition-colors"
                    :class="app.isDark ? 'border-slate-600 hover:bg-slate-700' : 'border-slate-300 hover:bg-slate-100'">
                    复制地址
                </button>
            </div>

            <!-- 使用说明 -->
            <div class="text-xs space-y-2 p-3 rounded-xl"
                :class="app.isDark ? 'bg-slate-900/50 text-slate-400' : 'bg-slate-50 text-slate-500'">
                <p><span class="font-medium">IP：</span><code class="font-mono">{{ ip }}</code></p>
                <p><span class="font-medium">端口：</span><code class="font-mono">{{ port }}</code></p>
                <p>在浏览器代理设置中填入以上两项（SOCKS5）即可。</p>
            </div>

            <div class="flex justify-end mt-4">
                <button @click="emit('close')"
                    class="px-4 py-2 text-sm rounded-lg transition-colors"
                    :class="app.isDark ? 'hover:bg-slate-700' : 'hover:bg-slate-100'">
                    关闭
                </button>
            </div>
        </div>
    </div>
</template>
