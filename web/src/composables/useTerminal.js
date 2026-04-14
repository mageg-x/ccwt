import { ref, onUnmounted, unref, watch } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'

const THEME_PRESETS = {
    shell: {
        background: '#11111b',
        foreground: '#cdd6f4',
        cursor: '#89dceb',
        cursorAccent: '#11111b',
        selectionBackground: 'rgba(137, 220, 235, 0.28)',
        black: '#45475a',
        red: '#f38ba8',
        green: '#a6e3a1',
        yellow: '#f9e2af',
        blue: '#89b4fa',
        magenta: '#f5c2e7',
        cyan: '#94e2d5',
        white: '#bac2de',
        brightBlack: '#585b70',
        brightRed: '#f38ba8',
        brightGreen: '#a6e3a1',
        brightYellow: '#f9e2af',
        brightBlue: '#89b4fa',
        brightMagenta: '#f5c2e7',
        brightCyan: '#94e2d5',
        brightWhite: '#a6adc8',
    },
    dark: {
        background: '#0b1020',
        foreground: '#dbeafe',
        cursor: '#22d3ee',
        cursorAccent: '#0b1020',
        selectionBackground: 'rgba(56, 189, 248, 0.22)',
        black: '#1f2937',
        red: '#f87171',
        green: '#4ade80',
        yellow: '#facc15',
        blue: '#60a5fa',
        magenta: '#c084fc',
        cyan: '#22d3ee',
        white: '#e2e8f0',
        brightBlack: '#64748b',
        brightRed: '#fb7185',
        brightGreen: '#86efac',
        brightYellow: '#fde047',
        brightBlue: '#93c5fd',
        brightMagenta: '#d8b4fe',
        brightCyan: '#67e8f9',
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
