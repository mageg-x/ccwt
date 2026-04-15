import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useDialogStore = defineStore('dialog', () => {
    const visible = ref(false)
    const type = ref('alert') // alert | confirm | prompt
    const title = ref('')
    const message = ref('')
    const okText = ref('确定')
    const cancelText = ref('取消')
    const placeholder = ref('')
    const inputValue = ref('')

    const queue = []
    let currentResolve = null

    function reset() {
        visible.value = false
        type.value = 'alert'
        title.value = ''
        message.value = ''
        okText.value = '确定'
        cancelText.value = '取消'
        placeholder.value = ''
        inputValue.value = ''
        currentResolve = null
    }

    function openNext() {
        if (visible.value) return
        const next = queue.shift()
        if (!next) return

        const payload = next.payload
        type.value = payload.type
        title.value = payload.title || ''
        message.value = payload.message || ''
        okText.value = payload.okText || '确定'
        cancelText.value = payload.cancelText || '取消'
        placeholder.value = payload.placeholder || ''
        inputValue.value = payload.defaultValue || ''
        currentResolve = next.resolve
        visible.value = true
    }

    function enqueue(payload) {
        return new Promise((resolve) => {
            queue.push({ payload, resolve })
            openNext()
        })
    }

    function closeAsCancel() {
        if (!visible.value || !currentResolve) return
        const resolve = currentResolve
        reset()
        resolve(false)
        openNext()
    }

    function closeAsConfirm() {
        if (!visible.value || !currentResolve) return
        const resolve = currentResolve
        const currentType = type.value
        const value = inputValue.value
        reset()
        if (currentType === 'prompt') {
            resolve(value)
        } else {
            resolve(true)
        }
        openNext()
    }

    function alert(message, options = {}) {
        return enqueue({
            type: 'alert',
            title: options.title || '提示',
            message,
            okText: options.okText || '知道了',
        })
    }

    function confirm(message, options = {}) {
        return enqueue({
            type: 'confirm',
            title: options.title || '请确认',
            message,
            okText: options.okText || '确认',
            cancelText: options.cancelText || '取消',
        })
    }

    function prompt(message, options = {}) {
        return enqueue({
            type: 'prompt',
            title: options.title || '请输入',
            message,
            okText: options.okText || '确认',
            cancelText: options.cancelText || '取消',
            placeholder: options.placeholder || '',
            defaultValue: options.defaultValue || '',
        }).then((result) => {
            if (result === false) return null
            return result
        })
    }

    return {
        visible,
        type,
        title,
        message,
        okText,
        cancelText,
        placeholder,
        inputValue,
        closeAsCancel,
        closeAsConfirm,
        alert,
        confirm,
        prompt,
    }
})

