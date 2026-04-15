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
        // 优化打包文件数量
        rollupOptions: {
            output: {
                // 配置代码分割策略
                manualChunks(id) {
                    // 将第三方库打包成一个 chunk
                    if (id.includes('node_modules')) {
                        if (id.includes('monaco-editor')) {
                            return 'editor';
                        }
                        return 'vendor';
                    }
                },
                // 最小化 chunk 大小，合并小文件
                chunkFileNames: 'assets/[name]-[hash].js',
                entryFileNames: 'assets/[name]-[hash].js',
                assetFileNames: 'assets/[name]-[hash].[ext]'
            }
        },
        // 启用压缩
        minify: 'esbuild'
    },
})
