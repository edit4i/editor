<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { Terminal } from '@xterm/xterm';
    import '@xterm/xterm/css/xterm.css';

    export let height: number;

    let terminalElement: HTMLElement;
    let terminal: Terminal;

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
            theme: {
                background: '#1f2937', // Tailwind gray-800
                foreground: '#d1d5db', // Tailwind gray-300
                cursor: '#60a5fa', // Tailwind blue-400
            }
        });

        terminal.open(terminalElement);
        terminal.write('Welcome to Edit4i Terminal\r\n$ ');
        
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
    class="h-full w-full" 
    bind:this={terminalElement} 
    style="height: {height}px"
/>
