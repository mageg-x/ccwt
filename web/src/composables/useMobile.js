import { computed } from 'vue'
import { useAppStore } from '../stores/app'

export function useMobile() {
    const app = useAppStore()
    return { isMobile: computed(() => app.isMobile) }
}
