<script setup>
import { computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useDialogStore } from '../stores/dialog'

const dialog = useDialogStore()
const { visible, type, title, message, okText, cancelText, placeholder, inputValue } = storeToRefs(dialog)

const isPrompt = computed(() => type.value === 'prompt')
const hasCancel = computed(() => type.value === 'confirm' || type.value === 'prompt')

function onBackdrop() {
    dialog.closeAsCancel()
}

function onConfirm() {
    dialog.closeAsConfirm()
}

function onCancel() {
    dialog.closeAsCancel()
}

function onKeydown(e) {
    if (!visible.value) return
    if (e.key === 'Escape') {
        e.preventDefault()
        dialog.closeAsCancel()
    } else if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault()
        dialog.closeAsConfirm()
    }
}

onMounted(() => {
    window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
    window.removeEventListener('keydown', onKeydown)
})

watch(
    () => [visible.value, type.value],
    async ([v, t]) => {
        if (!v || t !== 'prompt') return
        await nextTick()
        document.querySelector('[data-dialog-input]')?.focus()
    }
)
</script>

<template>
    <Transition name="fade">
        <div
            v-if="visible"
            class="fixed inset-0 z-[120] flex items-center justify-center p-4"
            @click.self="onBackdrop"
        >
            <div class="absolute inset-0 bg-slate-950/55 backdrop-blur-sm"></div>
            <div class="relative w-full max-w-md rounded-2xl border shadow-2xl overflow-hidden bg-slate-900 border-slate-700 text-slate-100">
                <div class="px-5 py-4 bg-gradient-to-r from-sky-600/20 via-indigo-500/20 to-cyan-500/20 border-b border-slate-700/80">
                    <div class="text-sm font-semibold tracking-wide">{{ title }}</div>
                </div>
                <div class="px-5 py-4 space-y-3">
                    <p class="text-sm leading-6 text-slate-200 whitespace-pre-wrap break-words">{{ message }}</p>
                    <input
                        v-if="isPrompt"
                        data-dialog-input
                        v-model="inputValue"
                        :placeholder="placeholder"
                        class="w-full px-3 py-2.5 rounded-xl bg-slate-800 border border-slate-600 text-slate-100 placeholder-slate-400 outline-none focus:border-sky-400"
                    />
                </div>
                <div class="px-5 pb-5 flex justify-end gap-2">
                    <button
                        v-if="hasCancel"
                        @click="onCancel"
                        class="px-4 py-2 rounded-xl text-sm border border-slate-600 text-slate-200 hover:bg-slate-800 transition-colors"
                    >
                        {{ cancelText }}
                    </button>
                    <button
                        @click="onConfirm"
                        class="px-4 py-2 rounded-xl text-sm text-white bg-gradient-to-r from-sky-500 to-indigo-500 hover:from-sky-400 hover:to-indigo-400 transition-colors"
                    >
                        {{ okText }}
                    </button>
                </div>
            </div>
        </div>
    </Transition>
</template>
