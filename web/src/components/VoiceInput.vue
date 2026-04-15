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
const backendAvailable = ref(false)
const backendReason = ref('')
const engineLabel = ref('')

let mediaStream = null
let audioContext = null
let sourceNode = null
let processorNode = null
let silentGainNode = null
let inputSampleRate = 48000
let pcmChunks = []

async function fetchVoiceStatus() {
    try {
        const { data } = await voiceApi.getStatus()
        backendAvailable.value = !!data?.available
        backendReason.value = data?.reason || ''
    } catch {
        backendAvailable.value = false
        backendReason.value = 'status.request.failed'
    }
}

async function startRecording() {
    error.value = ''
    resultText.value = ''

    engineLabel.value = '百度在线识别'
    if (!backendAvailable.value) {
        const reasonMap = {
            'voice.disabled': '后端语音识别未启用',
            'baidu.token.failed': '百度鉴权失败或网络不可达',
            'status.request.failed': '无法获取后端语音状态',
        }
        error.value = `后端语音不可用：${reasonMap[backendReason.value] || backendReason.value || '未知原因'}`
        return
    }

    try {
        mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true })
        const AudioCtx = window.AudioContext || window.webkitAudioContext
        audioContext = new AudioCtx()
        inputSampleRate = audioContext.sampleRate || 48000
        pcmChunks = []

        sourceNode = audioContext.createMediaStreamSource(mediaStream)
        processorNode = audioContext.createScriptProcessor(4096, 1, 1)
        silentGainNode = audioContext.createGain()
        silentGainNode.gain.value = 0

        processorNode.onaudioprocess = (e) => {
            if (!recording.value) return
            const channel = e.inputBuffer.getChannelData(0)
            pcmChunks.push(new Float32Array(channel))
        }

        sourceNode.connect(processorNode)
        processorNode.connect(silentGainNode)
        silentGainNode.connect(audioContext.destination)

        recording.value = true
    } catch (e) {
        error.value = '无法访问麦克风: ' + (e?.message || e)
        cleanupRecorder()
    }
}

function cleanupRecorder() {
    try { sourceNode?.disconnect() } catch {}
    try { processorNode?.disconnect() } catch {}
    try { silentGainNode?.disconnect() } catch {}
    if (mediaStream) {
        mediaStream.getTracks().forEach(t => t.stop())
    }
    if (audioContext && audioContext.state !== 'closed') {
        audioContext.close().catch(() => {})
    }
    sourceNode = null
    processorNode = null
    silentGainNode = null
    mediaStream = null
    audioContext = null
}

function mergePCM(chunks) {
    let length = 0
    for (const c of chunks) length += c.length
    const merged = new Float32Array(length)
    let offset = 0
    for (const c of chunks) {
        merged.set(c, offset)
        offset += c.length
    }
    return merged
}

function downsamplePCM(buffer, fromRate, toRate) {
    if (fromRate === toRate) return buffer
    const ratio = fromRate / toRate
    const newLength = Math.round(buffer.length / ratio)
    const out = new Float32Array(newLength)
    let offset = 0
    for (let i = 0; i < newLength; i++) {
        const start = Math.round(i * ratio)
        const end = Math.round((i + 1) * ratio)
        let sum = 0
        let count = 0
        for (let j = start; j < end && j < buffer.length; j++) {
            sum += buffer[j]
            count++
        }
        out[offset++] = count > 0 ? sum / count : 0
    }
    return out
}

function encodeWav16kMono(float32Data) {
    const sampleRate = 16000
    const bytesPerSample = 2
    const blockAlign = bytesPerSample
    const byteRate = sampleRate * blockAlign
    const dataSize = float32Data.length * bytesPerSample
    const buffer = new ArrayBuffer(44 + dataSize)
    const view = new DataView(buffer)

    let offset = 0
    const writeString = (s) => {
        for (let i = 0; i < s.length; i++) view.setUint8(offset++, s.charCodeAt(i))
    }

    writeString('RIFF')
    view.setUint32(offset, 36 + dataSize, true); offset += 4
    writeString('WAVE')
    writeString('fmt ')
    view.setUint32(offset, 16, true); offset += 4
    view.setUint16(offset, 1, true); offset += 2
    view.setUint16(offset, 1, true); offset += 2
    view.setUint32(offset, sampleRate, true); offset += 4
    view.setUint32(offset, byteRate, true); offset += 4
    view.setUint16(offset, blockAlign, true); offset += 2
    view.setUint16(offset, 16, true); offset += 2
    writeString('data')
    view.setUint32(offset, dataSize, true); offset += 4

    for (let i = 0; i < float32Data.length; i++) {
        const s = Math.max(-1, Math.min(1, float32Data[i]))
        view.setInt16(offset, s < 0 ? s * 0x8000 : s * 0x7fff, true)
        offset += 2
    }
    return new Blob([buffer], { type: 'audio/wav' })
}

function stopRecording() {
    if (!recording.value) return
    recording.value = false

    const merged = mergePCM(pcmChunks)
    const downsampled = downsamplePCM(merged, inputSampleRate, 16000)
    const wavBlob = encodeWav16kMono(downsampled)
    cleanupRecorder()
    sendToBackend(wavBlob)
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

onMounted(fetchVoiceStatus)

onUnmounted(() => {
    cleanupRecorder()
})
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
            <p v-if="engineLabel" class="text-center text-xs text-slate-500 mb-2">
                当前引擎：{{ engineLabel }}
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
