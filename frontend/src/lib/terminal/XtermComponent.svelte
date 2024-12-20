<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { Terminal } from '@xterm/xterm';
    import '@xterm/xterm/css/xterm.css';

    export let height: number;
    export let id: string;

    let terminalElement: HTMLElement;
    let terminal: Terminal;

    const terminalTheme = {
        background: '#181818', // Darker background
        foreground: '#c5c8c6', // Slightly darker text color
        cursor: '#528bff', // Bright blue cursor
        selectionBackground: '#3e4451', // Selection color
        selectionForeground: '#d1d5db', // Selection color
        black: '#1e1e1e',
        red: '#e06c75',
        green: '#98c379',
        yellow: '#e5c07b',
        blue: '#61afef',
    };

    // Function to update terminal size
    function updateTerminalSize() {
        if (!terminal || !terminalElement) return;
        
        // Get the current dimensions of the container
        const computedStyle = window.getComputedStyle(terminalElement);
        const width = parseInt(computedStyle.width);
        const paddingLeft = parseInt(computedStyle.paddingLeft);
        const paddingRight = parseInt(computedStyle.paddingRight);
        const paddingTop = parseInt(computedStyle.paddingTop);
        const paddingBottom = parseInt(computedStyle.paddingBottom);

        // Calculate available space
        const availableWidth = width - paddingLeft - paddingRight;
        const availableHeight = height - paddingTop - paddingBottom;

        // Calculate rows and columns based on character size
        const charWidth = 9; // Approximate width of a character in pixels
        const charHeight = 17; // Approximate height of a character in pixels
        const cols = Math.floor(availableWidth / charWidth);
        const rows = Math.floor(availableHeight / charHeight);

        // Update terminal dimensions
        terminal.resize(cols, rows);
    }

    // Watch for height changes
    $: if (height) {
        updateTerminalSize();
    }

    onMount(() => {
        terminal = new Terminal({
            fontSize: 14,
            fontFamily: 'monospace',
            theme: terminalTheme
        });

        terminal.open(terminalElement);
        terminal.write(`Welcome to Terminal ${id}\r\n$ `);
        
        // Initial size update
        updateTerminalSize();

        // Listen for window resize
        window.addEventListener('resize', updateTerminalSize);
    });

    onDestroy(() => {
        if (terminal) {
            terminal.dispose();
        }
        window.removeEventListener('resize', updateTerminalSize);
    });
</script>

<div 
    class={`h-full w-full p-2`}
    style={`background-color: ${terminalTheme.background};`}
    bind:this={terminalElement} 
/>
