<script setup>
import { useAppStore } from './stores/app'
import NetworkStatus from './components/NetworkStatus.vue'
import DialogHost from './components/DialogHost.vue'

const app = useAppStore()
app.initTheme()
</script>

<template>
    <div :class="{ 'theme-light': app.theme === 'light' }" class="h-full">
        <router-view />
        <NetworkStatus />
        <DialogHost />
    </div>
</template>
