import { computed } from 'vue'
import { useAppStore } from '../stores/app'

export function useTheme() {
    const app = useAppStore()
    const isDark = computed(() => app.theme === 'dark')

    return { isDark, theme: computed(() => app.theme), toggleTheme: app.toggleTheme }
}
