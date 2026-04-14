import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
    plugins: [vue(), tailwindcss()],
    server: {
        port: 5173,
        proxy: {
            '/api': 'http://localhost:3000',
            '/ws': {
                target: 'ws://localhost:3000',
                ws: true,
            },
        },
    },
    build: {
        outDir: 'dist',
        emptyOutDir: true,
    },
})
