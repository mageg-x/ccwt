import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useMediaQuery } from '@vueuse/core'

export const useAppStore = defineStore('app', () => {
    const theme = ref('dark')
    const sidebarOpen = ref(true)
    const cmdPaletteOpen = ref(false)
    const isMobile = useMediaQuery('(max-width: 767px)')

    const isDark = computed(() => theme.value === 'dark' || theme.value === 'shell')

    function initTheme() {
        const saved = localStorage.getItem('ccwt-theme')
        if (saved) {
            theme.value = saved
        } else if (window.matchMedia('(prefers-color-scheme: light)').matches) {
            theme.value = 'light'
        }
        if (isMobile.value) {
            sidebarOpen.value = false
        }
    }

    function toggleTheme() {
        const themes = ['dark', 'light', 'shell']
        const idx = themes.indexOf(theme.value)
        theme.value = themes[(idx + 1) % themes.length]
        localStorage.setItem('ccwt-theme', theme.value)
    }

    function setTheme(t) {
        theme.value = t
        localStorage.setItem('ccwt-theme', theme.value)
    }

    function toggleSidebar() {
        sidebarOpen.value = !sidebarOpen.value
    }

    function toggleCmdPalette() {
        cmdPaletteOpen.value = !cmdPaletteOpen.value
    }

    return {
        theme, sidebarOpen, cmdPaletteOpen, isMobile, isDark,
        initTheme, toggleTheme, setTheme, toggleSidebar, toggleCmdPalette,
    }
})
