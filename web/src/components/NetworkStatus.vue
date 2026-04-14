<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const offline = ref(!navigator.onLine)

function onOnline() { offline.value = false }
function onOffline() { offline.value = true }

onMounted(() => {
    window.addEventListener('online', onOnline)
    window.addEventListener('offline', onOffline)
})
onUnmounted(() => {
    window.removeEventListener('online', onOnline)
    window.removeEventListener('offline', onOffline)
})
</script>

<template>
    <Transition name="fade">
        <div v-if="offline"
            class="fixed top-0 left-0 right-0 z-[100] flex items-center justify-center gap-2 py-2 bg-amber-500 text-amber-900 text-sm font-medium">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636a9 9 0 010 12.728M5.636 5.636a9 9 0 000 12.728M12 12h.01" />
            </svg>
            网络已断开，正在等待重新连接...
        </div>
    </Transition>
</template>
