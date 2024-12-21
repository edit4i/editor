<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { Terminal, type ITerminalOptions, type ITheme } from '@xterm/xterm';
    import '@xterm/xterm/css/xterm.css';
    import { CreateTerminal, DestroyTerminal, HandleInput, ResizeTerminal } from '@/lib/wailsjs/go/main/App';
    import { EventsOn, EventsOff } from '@/lib/wailsjs/runtime/runtime';
    import { projectStore } from '@/stores/project';
    import { editorConfigStore } from '@/stores/editorConfigStore';
    import { get } from 'svelte/store';

    export let height: number;
    export let id: string;
    export let shell: string;
    export let active: boolean = false;

    let terminalElement: HTMLElement;
    let terminal: Terminal | null = null;
    let isDestroyed = false;
    let isInitialized = false;
    let resizeTimeout: number | null = null;

    console.log('[Terminal] Initializing with id:', id, 'shell:', shell);

    // Get terminal config
    const config = get(editorConfigStore);
    const terminalConfig = config.terminal;

    const terminalTheme: ITheme = {
        background: terminalConfig.theme.background,
        foreground: terminalConfig.theme.foreground,
        cursor: terminalConfig.theme.cursor,
        selectionBackground: terminalConfig.theme.selectionBackground,
        selectionForeground: terminalConfig.theme.selectionForeground,
    };

    // Function to update terminal size with debouncing
    async function updateTerminalSize() {
        if (!terminal || !terminalElement || isDestroyed) return;

        if (resizeTimeout) {
            clearTimeout(resizeTimeout);
        }

        resizeTimeout = window.setTimeout(async () => {
            const computedStyle = window.getComputedStyle(terminalElement);
            const width = parseInt(computedStyle.width);
            const paddingLeft = parseInt(computedStyle.paddingLeft);
            const paddingRight = parseInt(computedStyle.paddingRight);
            const paddingTop = parseInt(computedStyle.paddingTop);
            const paddingBottom = parseInt(computedStyle.paddingBottom);

            const availableWidth = width - paddingLeft - paddingRight;
            const availableHeight = height - paddingTop - paddingBottom;

            // TODO: Account for bottom status bar without hardcoding
            const bottomBarRows = 3;

            const charWidth = 9;
            const charHeight = 17;
            const cols = Math.floor(availableWidth / charWidth);
            const rows = Math.floor(availableHeight / charHeight) - bottomBarRows;

            console.log('[Terminal] Resizing to:', { cols, rows });
            terminal!.resize(cols, rows);

            try {
                await ResizeTerminal(id, cols, rows);
                console.log('[Terminal] Backend resize successful');
            } catch (error) {
                console.error('[Terminal] Error resizing:', error);
            }
            resizeTimeout = null;
        }, 100); // Debounce resize events by 100ms
    }

    // Handle terminal events from backend
    function handleTerminalEvent(event: any) {
        console.log('[Terminal] Received event:', event);
        if (!terminal || isDestroyed) return;

        switch (event.Type) {
            case 0: // EventData
                if (event.Data) {
                    // Decode base64 data
                    const base64Data = event.Data;
                    const binaryStr = atob(base64Data);
                    const bytes = Uint8Array.from(binaryStr, c => c.charCodeAt(0));
                    
                    console.log('[Terminal] Writing decoded data:', new TextDecoder().decode(bytes));
                    terminal.write(bytes);
                }
                break;
            case 1: // EventResize
                console.log('[Terminal] Resize event:', event);
                terminal.resize(event.Cols, event.Rows);
                break;
            case 2: // EventCursor
                console.log('[Terminal] Cursor event:', event);
                break;
            case 3: // EventExit
                console.log('[Terminal] Exit event received');
                isDestroyed = true;
                terminal.write('\r\nTerminal session ended.\r\n');
                break;
        }
    }

    // Watch for height changes
    $: if (height && terminal && !isDestroyed) {
        console.log('[Terminal] Height changed:', height);
        updateTerminalSize();
    }

    // Watch for active state changes
    $: if (active && terminal) {
        terminal.focus();
        // Trigger resize to ensure proper layout after tab switch
        updateTerminalSize();
    }

    async function initializeTerminal() {
        if (isInitialized || isDestroyed) return;
        
        console.log('[Terminal] Mounting component');

        console.log('[Terminal] config:', config);
        
        // Create xterm.js instance
        terminal = new Terminal({
            theme: terminalTheme,
            cursorBlink: true,
            allowProposedApi: true,
            scrollback: 10000,
            fontSize: terminalConfig.fontSize,
            fontFamily: terminalConfig.fontFamily,
            // Prevent terminal from handling our keyboard shortcuts
            allowTransparency: true
        });

        terminal.attachCustomKeyEventHandler((event: KeyboardEvent) => {
            // Don't handle Ctrl+Alt+J
            if (event.ctrlKey && event.altKey && event.key.toLowerCase() === 'j') {
                return false;
            }
            return true;
        });

        // Attach terminal to DOM
        terminal.open(terminalElement);

        // Handle terminal input
        terminal.onData((data) => {
            console.log('[Terminal] Input received:', data);
            if (isDestroyed) return;

            // Special handling for Enter key
            if (data === '\r') {
                console.log('[Terminal] Enter key pressed');
            }

            try {
                // @ts-ignore
                HandleInput(id, btoa(data));
            } catch (error) {
                console.error('[Terminal] Error handling input:', error);
            }
        });

        // Get current project path
        const projectPath = get(projectStore).currentProject?.Path || '';

        // Create backend terminal
        console.log('[Terminal] Creating backend terminal');
        try {
            await CreateTerminal(id, shell, projectPath);
            console.log('[Terminal] Backend terminal created');

            // Subscribe to terminal events
            console.log('[Terminal] Subscribing to events');
            EventsOn(`terminal:${id}`, handleTerminalEvent);

            // Initial resize
            updateTerminalSize();
            isInitialized = true;
        } catch (error) {
            console.error('[Terminal] Error creating terminal:', error);
            isDestroyed = true;
        }
    }

    // Expose focus method
    export function focus() {
        if (terminal) {
            terminal.focus();
        }
    }

    onMount(() => {
        initializeTerminal();
    });

    onDestroy(() => {
        if (!isDestroyed) {
            console.log('[Terminal] Destroying terminal');
            if (resizeTimeout) {
                clearTimeout(resizeTimeout);
            }
            EventsOff(`terminal:${id}`);
            DestroyTerminal(id).catch(error => {
                console.error('[Terminal] Error destroying terminal:', error);
            });
            if (terminal) {
                terminal.dispose();
                terminal = null;
            }
            isDestroyed = true;
        }
    });
</script>

<div 
    class={`h-full w-full p-2`}
    style={`background-color: ${terminalTheme.background};`}
    bind:this={terminalElement}
/>
