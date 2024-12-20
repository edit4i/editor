<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { Terminal } from '@xterm/xterm';
    import '@xterm/xterm/css/xterm.css';
    import { CreateTerminal, DestroyTerminal, HandleInput, ResizeTerminal } from '@/lib/wailsjs/go/main/App';
    import { EventsOn, EventsOff } from '@/lib/wailsjs/runtime/runtime';

    export let height: number;
    export let id: string;
    export let shell: string;

    let terminalElement: HTMLElement;
    let terminal: Terminal | null = null;
    let isDestroyed = false;

    console.log('[Terminal] Initializing with id:', id, 'shell:', shell);

    const terminalTheme = {
        background: '#181818',
        foreground: '#c5c8c6',
        cursor: '#528bff',
        selectionBackground: '#3e4451',
        selectionForeground: '#d1d5db',
        black: '#1e1e1e',
        red: '#e06c75',
        green: '#98c379',
        yellow: '#e5c07b',
        blue: '#61afef',
    };

    // Function to update terminal size
    async function updateTerminalSize() {
        if (!terminal || !terminalElement || isDestroyed) return;
        
        const computedStyle = window.getComputedStyle(terminalElement);
        const width = parseInt(computedStyle.width);
        const paddingLeft = parseInt(computedStyle.paddingLeft);
        const paddingRight = parseInt(computedStyle.paddingRight);
        const paddingTop = parseInt(computedStyle.paddingTop);
        const paddingBottom = parseInt(computedStyle.paddingBottom);

        const availableWidth = width - paddingLeft - paddingRight;
        const availableHeight = height - paddingTop - paddingBottom;

        // TODO: Account for bottom status bar without hardcoding
        const bottomBarRows = 2;

        const charWidth = 9;
        const charHeight = 17;
        const cols = Math.floor(availableWidth / charWidth);
        const rows = Math.floor(availableHeight / charHeight) - bottomBarRows;

        console.log('[Terminal] Resizing to:', { cols, rows });
        terminal.resize(cols, rows);

        try {
            await ResizeTerminal(id, cols, rows);
            console.log('[Terminal] Backend resize successful');
        } catch (error) {
            console.error('[Terminal] Error resizing:', error);
        }
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

    onMount(async () => {
        console.log('[Terminal] Mounting component');
        if (isDestroyed) return;

        terminal = new Terminal({
            fontSize: 14,
            fontFamily: 'monospace',
            theme: terminalTheme,
            cursorBlink: true,
        });

        // Handle terminal input
        terminal.onData((data) => {
            if (!isDestroyed) {
                console.log('[Terminal] Input received:', data);
                // Handle Enter key
                if (data === '\r') {
                    console.log('[Terminal] Enter key pressed');
                    // Send carriage return
                    const bytes = new TextEncoder().encode('\r');
                    HandleInput(id, Array.from(bytes));
                } else {
                    // Convert string to byte array and send
                    const bytes = new TextEncoder().encode(data);
                    HandleInput(id, Array.from(bytes));
                }
            }
        });

        // Handle terminal resize
        terminal.onResize(({ cols, rows }) => {
            if (!isDestroyed) {
                console.log('[Terminal] Terminal resize:', { cols, rows });
                ResizeTerminal(id, cols, rows);
            }
        });

        terminal.open(terminalElement);
        console.log('[Terminal] Terminal opened');

        try {
            console.log('[Terminal] Creating backend terminal');
            // Create terminal on backend
            await CreateTerminal(id, shell);
            console.log('[Terminal] Backend terminal created');

            // Initial size update
            await updateTerminalSize();

            // Subscribe to terminal events
            console.log('[Terminal] Subscribing to events');
            EventsOn(`terminal:${id}`, handleTerminalEvent);
        } catch (error) {
            console.error('[Terminal] Error initializing:', error);
            terminal.write('Error initializing terminal\r\n');
        }

        // Listen for window resize
        window.addEventListener('resize', updateTerminalSize);
    });

    onDestroy(async () => {
        console.log('[Terminal] Destroying terminal');
        isDestroyed = true;
        window.removeEventListener('resize', updateTerminalSize);
        
        // Unsubscribe from terminal events
        EventsOff(`terminal:${id}`);

        if (terminal) {
            try {
                // Destroy terminal on backend
                await DestroyTerminal(id);
                terminal.dispose();
                console.log('[Terminal] Terminal destroyed');
            } catch (error) {
                console.error('[Terminal] Error disposing:', error);
            }
            terminal = null;
        }
    });
</script>

<div 
    class={`h-full w-full p-2`}
    style={`background-color: ${terminalTheme.background};`}
    bind:this={terminalElement} 
/>
