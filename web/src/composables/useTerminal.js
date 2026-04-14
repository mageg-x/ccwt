import { ref, onUnmounted, unref, watch } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'

const THEME_PRESETS = {
    shell: {
        background: '#000000',
        foreground: '#d0d0d0',
        cursor: '#f5f5f5',
        selectionBackground: 'rgba(255, 255, 255, 0.25)',
        black: '#000000',
        red: '#cd3131',
        green: '#0dbc79',
        yellow: '#e5e510',
        blue: '#2472c8',
        magenta: '#bc3fbc',
        cyan: '#11a8cd',
        white: '#e5e5e5',
        brightBlack: '#666666',
        brightRed: '#f14c4c',
        brightGreen: '#23d18b',
        brightYellow: '#f5f543',
        brightBlue: '#3b8eea',
        brightMagenta: '#d670d6',
        brightCyan: '#29b8db',
        brightWhite: '#ffffff',
    },
    dark: {
        background: '#0f172a',
        foreground: '#e2e8f0',
        cursor: '#6366f1',
        selectionBackground: '#334155',
        black: '#1e293b',
        red: '#ef4444',
        green: '#22c55e',
        yellow: '#f59e0b',
        blue: '#3b82f6',
        magenta: '#a855f7',
        cyan: '#06b6d4',
        white: '#e2e8f0',
        brightBlack: '#475569',
        brightRed: '#f87171',
        brightGreen: '#4ade80',
        brightYellow: '#fbbf24',
        brightBlue: '#60a5fa',
        brightMagenta: '#c084fc',
        brightCyan: '#22d3ee',
        brightWhite: '#f8fafc',
    },
    light: {
        background: '#ffffff',
        foreground: '#1e293b',
        cursor: '#6366f1',
        selectionBackground: '#dbeafe',
        black: '#1e293b',
        red: '#dc2626',
        green: '#16a34a',
        yellow: '#ca8a04',
        blue: '#2563eb',
        magenta: '#9333ea',
        cyan: '#0891b2',
        white: '#f8fafc',
        brightBlack: '#64748b',
        brightRed: '#ef4444',
        brightGreen: '#22c55e',
        brightYellow: '#eab308',
        brightBlue: '#3b82f6',
        brightMagenta: '#a855f7',
        brightCyan: '#06b6d4',
        brightWhite: '#0f172a',
    },
}

function resolveTheme(themeValue) {
    const name = themeValue === 'shell' ? 'shell' : themeValue === 'light' ? 'light' : 'dark'
    return THEME_PRESETS[name]
}

function copyText(text) {
    if (!text) return

    if (navigator.clipboard?.writeText) {
        navigator.clipboard.writeText(text).catch(() => {})
        return
    }

    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.setAttribute('readonly', '')
    textarea.style.position = 'absolute'
    textarea.style.left = '-9999px'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
}

export function useTerminal(containerRef, { onData, onResize, theme = 'dark' } = {}) {
    const term = ref(null)
    const fitAddon = ref(null)
    let resizeObs = null

    function createTerminal() {
        const currentTheme = unref(theme)
        const t = new Terminal({
            cursorBlink: true,
            fontSize: 14,
            fontFamily: "'JetBrains Mono', 'Fira Code', 'Cascadia Code', Menlo, Monaco, monospace",
            lineHeight: 1.2,
            scrollback: 10000,
            allowTransparency: true,
            theme: resolveTheme(currentTheme),
            allowProposedApi: true,
        })

        t.attachCustomKeyEventHandler((ev) => {
            const isCopy = (ev.ctrlKey || ev.metaKey) && !ev.altKey && (ev.key === 'c' || ev.key === 'C')
            if (isCopy && term.value?.hasSelection()) {
                if (ev.type === 'keydown') {
                    copyText(term.value.getSelection())
                    term.value.clearSelection()
                }
                return false
            }
            return true
        })

        const fit = new FitAddon()
        t.loadAddon(fit)
        t.loadAddon(new WebLinksAddon())

        term.value = t
        fitAddon.value = fit

        t.onData((data) => {
            onData?.(data)
        })

        t.onResize(({ rows, cols }) => {
            onResize?.(rows, cols)
        })

        return t
    }

    function mount() {
        if (!containerRef.value || term.value) return
        const t = createTerminal()
        t.open(containerRef.value)
        fitAddon.value.fit()

        resizeObs = new ResizeObserver(() => {
            fitAddon.value?.fit()
        })
        resizeObs.observe(containerRef.value)
    }

    function write(data) {
        term.value?.write(data instanceof ArrayBuffer ? new Uint8Array(data) : data)
    }

    function fit() {
        fitAddon.value?.fit()
    }

    function focus() {
        term.value?.focus()
    }

    function dispose() {
        resizeObs?.disconnect()
        term.value?.dispose()
        term.value = null
    }

    watch(
        () => unref(theme),
        (nextTheme) => {
            if (term.value) {
                term.value.options.theme = resolveTheme(nextTheme)
            }
        }
    )

    onUnmounted(dispose)

    return { term, mount, write, fit, focus, dispose }
}
