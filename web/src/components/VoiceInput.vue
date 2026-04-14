<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAppStore } from '../stores/app'
import * as voiceApi from '../api/voice'

const emit = defineEmits(['close', 'result'])
const app = useAppStore()

const recording = ref(false)
const processing = ref(false)
const error = ref('')
const resultText = ref('')
let mediaRecorder = null
let audioChunks = []

// 尝试使用 Web Speech API
const hasSpeechAPI = ref('webkitSpeechRecognition' in window || 'SpeechRecognition' in window)

async function startRecording() {
    error.value = ''
    resultText.value = ''

    // 优先使用 Web Speech API
    if (hasSpeechAPI.value) {
        startSpeechRecognition()
        return
    }

    // 降级到录音 + 后端识别
    try {
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
        audioChunks = []
        mediaRecorder = new MediaRecorder(stream)
        mediaRecorder.ondataavailable = (e) => audioChunks.push(e.data)
        mediaRecorder.onstop = async () => {
            stream.getTracks().forEach(t => t.stop())
            const blob = new Blob(audioChunks, { type: 'audio/wav' })
            await sendToBackend(blob)
        }
        mediaRecorder.start()
        recording.value = true
    } catch (e) {
        error.value = '无法访问麦克风: ' + e.message
    }
}

function stopRecording() {
    if (mediaRecorder?.state === 'recording') {
        mediaRecorder.stop()
    }
    recording.value = false
}

async function sendToBackend(blob) {
    processing.value = true
    try {
        const { data } = await voiceApi.recognize(blob)
        resultText.value = data.text
        emit('result', data.text)
    } catch (e) {
        error.value = e.response?.data?.error || '语音识别失败'
    } finally {
        processing.value = false
    }
}

function startSpeechRecognition() {
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition
    const recognition = new SpeechRecognition()
    recognition.lang = 'zh-CN'
    recognition.interimResults = false
    recognition.maxAlternatives = 1

    recognition.onresult = (e) => {
        const text = e.results[0][0].transcript
        resultText.value = text
        emit('result', text)
        recording.value = false
    }
    recognition.onerror = (e) => {
        error.value = '语音识别错误: ' + e.error
        recording.value = false
    }
    recognition.onend = () => {
        recording.value = false
    }

    recognition.start()
    recording.value = true
}
</script>

<template>
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="emit('close')">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>
        <div class="relative w-full max-w-sm rounded-2xl shadow-2xl border p-6"
            :class="app.isDark ? 'bg-slate-800 border-slate-700' : 'bg-white border-slate-200'">

            <h3 class="text-lg font-semibold mb-4 text-center">语音输入</h3>

            <!-- 录音按钮 -->
            <div class="flex justify-center mb-4">
                <button
                    @click="recording ? stopRecording() : startRecording()"
                    class="w-20 h-20 rounded-full flex items-center justify-center transition-all"
                    :class="recording
                        ? 'bg-red-500 animate-pulse shadow-lg shadow-red-500/50'
                        : 'bg-indigo-600 hover:bg-indigo-500 shadow-lg shadow-indigo-500/30'"
                >
                    <svg v-if="!recording" class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
                    </svg>
                    <svg v-else class="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24">
                        <rect x="6" y="6" width="12" height="12" rx="2" />
                    </svg>
                </button>
            </div>

            <p class="text-center text-sm text-slate-400 mb-3">
                {{ recording ? '正在录音，点击停止...' : processing ? '识别中...' : '点击开始录音' }}
            </p>

            <!-- 波形动画 -->
            <div v-if="recording" class="flex items-center justify-center gap-1 h-8 mb-3">
                <div v-for="i in 5" :key="i"
                    class="w-1 bg-indigo-400 rounded-full animate-bounce"
                    :style="{ height: Math.random() * 24 + 8 + 'px', animationDelay: i * 0.1 + 's' }">
                </div>
            </div>

            <!-- 处理中 -->
            <div v-if="processing" class="flex items-center justify-center gap-2 mb-3">
                <svg class="animate-spin w-4 h-4 text-indigo-400" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                </svg>
                <span class="text-sm">正在识别...</span>
            </div>

            <!-- 结果 -->
            <div v-if="resultText" class="p-3 rounded-xl mb-3"
                :class="app.isDark ? 'bg-slate-700/50' : 'bg-slate-100'">
                <p class="text-sm">{{ resultText }}</p>
            </div>

            <!-- 错误 -->
            <div v-if="error" class="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-400 text-sm mb-3">
                {{ error }}
            </div>

            <div class="flex justify-end">
                <button @click="emit('close')"
                    class="px-4 py-2 text-sm rounded-lg transition-colors"
                    :class="app.isDark ? 'hover:bg-slate-700' : 'hover:bg-slate-100'">
                    关闭
                </button>
            </div>
        </div>
    </div>
</template>
