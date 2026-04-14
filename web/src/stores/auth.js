import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(null)
    const initialized = ref(false)
    let fetchPromise = null

    const isLoggedIn = computed(() => !!user.value)
    const isAdmin = computed(() => user.value?.role === 'admin')

    async function fetchMe() {
        if (fetchPromise) {
            return fetchPromise
        }

        fetchPromise = (async () => {
            try {
                const { data } = await authApi.getMe()
                user.value = data.user
                return data.user
            } catch {
                user.value = null
                return null
            } finally {
                initialized.value = true
                fetchPromise = null
            }
        })()

        return fetchPromise
    }

    async function ensureAuthLoaded() {
        if (initialized.value) {
            return user.value
        }
        return fetchMe()
    }

    async function login(username, password) {
        const { data } = await authApi.login({ username, password })
        user.value = data.user
        initialized.value = true
        return data
    }

    async function register(username, password) {
        const { data } = await authApi.register({ username, password })
        user.value = data.user
        initialized.value = true
        return data
    }

    async function logout() {
        await authApi.logout()
        user.value = null
        initialized.value = true
    }

    return {
        user,
        initialized,
        isLoggedIn,
        isAdmin,
        fetchMe,
        ensureAuthLoaded,
        login,
        register,
        logout,
    }
})
