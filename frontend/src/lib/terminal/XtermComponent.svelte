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

        const charWidth = 9;
        const charHeight = 17;
        const cols = Math.floor(availableWidth / charWidth);
        const rows = Math.floor(availableHeight / charHeight);

        terminal.resize(cols, rows);

        try {
            await ResizeTerminal(id, cols, rows);
        } catch (error) {
            console.error('Error resizing terminal:', error);
        }
    }

    // Handle terminal events from backend
    function handleTerminalEvent(event: any) {
        if (!terminal || isDestroyed) return;

        switch (event.Type) {
            case 0: // EventData
                if (event.Data) {
                    const data = new Uint8Array(event.Data);
                    terminal.write(data);
                }
                break;
            case 1: // EventResize
                terminal.resize(event.Cols, event.Rows);
                break;
            case 2: // EventCursor
                // Xterm.js handles cursor position automatically
                break;
            case 3: // EventExit
                isDestroyed = true;
                terminal.write('\r\nTerminal session ended.\r\n');
                break;
        }
    }

    // Watch for height changes
    $: if (height && terminal && !isDestroyed) {
        updateTerminalSize();
    }

    onMount(async () => {
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
                // Convert string to byte array
                const bytes = new TextEncoder().encode(data);
                HandleInput(id, Array.from(bytes));
            }
        });

        // Handle terminal resize
        terminal.onResize(({ cols, rows }) => {
            if (!isDestroyed) {
                ResizeTerminal(id, cols, rows);
            }
        });

        terminal.open(terminalElement);

        try {
            // Create terminal on backend
            await CreateTerminal(id, shell);

            // Initial size update
            await updateTerminalSize();

            // Subscribe to terminal events
            EventsOn(`terminal:${id}`, handleTerminalEvent);
        } catch (error) {
            console.error('Error initializing terminal:', error);
            terminal.write('Error initializing terminal\r\n');
        }

        // Listen for window resize
        window.addEventListener('resize', updateTerminalSize);
    });

    onDestroy(async () => {
        isDestroyed = true;
        window.removeEventListener('resize', updateTerminalSize);
        
        // Unsubscribe from terminal events
        EventsOff(`terminal:${id}`);

        if (terminal) {
            try {
                // Destroy terminal on backend
                await DestroyTerminal(id);
                terminal.dispose();
            } catch (error) {
                console.error('Error disposing terminal:', error);
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
