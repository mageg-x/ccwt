<script setup>
import { useAppStore } from './stores/app'
import { useAuthStore } from './stores/auth'
import NetworkStatus from './components/NetworkStatus.vue'

const app = useAppStore()
const auth = useAuthStore()
app.initTheme()
auth.fetchMe()
</script>

<template>
    <div :class="{ 'theme-light': app.theme === 'light' }" class="h-full">
        <router-view />
        <NetworkStatus />
    </div>
</template>
