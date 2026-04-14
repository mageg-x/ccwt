<script setup>
import { useAppStore } from './stores/app'
import NetworkStatus from './components/NetworkStatus.vue'

const app = useAppStore()
app.initTheme()
</script>

<template>
    <div :class="{ 'theme-light': app.theme === 'light' }" class="h-full">
        <router-view />
        <NetworkStatus />
    </div>
</template>
