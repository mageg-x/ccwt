import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as fileApi from '../api/file'

export const useFileStore = defineStore('file', () => {
    const tree = ref(null)
    const loading = ref(false)
    const editingFile = ref(null) // { path, content, language }

    async function loadTree(path = '.') {
        loading.value = true
        try {
            const { data } = await fileApi.getTree(path)
            tree.value = data.tree
        } catch (e) {
            console.error('加载文件树失败', e)
        } finally {
            loading.value = false
        }
    }

    async function openFile(path) {
        const { data } = await fileApi.readFile(path)
        const ext = path.split('.').pop()
        const langMap = {
            js: 'javascript', ts: 'typescript', py: 'python',
            go: 'go', json: 'json', yaml: 'yaml', yml: 'yaml',
            md: 'markdown', html: 'html', css: 'css', vue: 'html',
            sh: 'shell', bash: 'shell', sql: 'sql', xml: 'xml',
            java: 'java', rs: 'rust', cpp: 'cpp', c: 'c', h: 'c',
            rb: 'ruby', php: 'php', swift: 'swift', kt: 'kotlin',
            toml: 'ini', env: 'ini', conf: 'ini', cfg: 'ini',
        }
        editingFile.value = {
            path,
            content: data.content,
            language: langMap[ext] || 'plaintext',
        }
    }

    function closeEditor() {
        editingFile.value = null
    }

    async function saveFile(path, content) {
        await fileApi.writeFile(path, content)
        if (editingFile.value?.path === path) {
            editingFile.value.content = content
        }
    }

    return { tree, loading, editingFile, loadTree, openFile, closeEditor, saveFile }
})
