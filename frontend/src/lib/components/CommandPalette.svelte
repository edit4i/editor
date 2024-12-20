<script lang="ts">
    import { onMount, createEventDispatcher, afterUpdate, onDestroy } from 'svelte';
    import { Search } from 'lucide-svelte';
    import Input from './Input.svelte';
    import { fuzzySearch } from '@/lib/utils/fuzzySearch';
    import { commandStore, type Command } from '@/stores/commandStore';
    import { keyBindings, formatKeybinding, addKeyboardContext, removeKeyboardContext, type KeyBinding } from '@/stores/keyboardStore';
    import { focusStore } from '@/stores/focusStore';

    const dispatch = createEventDispatcher();

    export let show = false;
    let previousShow = show;

    let searchQuery = '';
    let selectedIndex = 0;
    let filteredCommands: Command[] = [];
    let inputElement: HTMLInputElement;
    let vimModeEnabled = false;

    let paletteId = focusStore.generateId('command-palette');

    // Convert keyboard bindings to commands
    $: allCommands = Object.entries($keyBindings).map(([id, binding]) => ({
        id,
        label: binding.description || id,
        category: binding.category,
        context: binding.context?.join(', ') || 'global',
        shortcut: formatKeybinding(binding),
        action: binding.action
    }));

    // Update filtered commands whenever commands or searchQuery changes
    $: {
        if (searchQuery?.length > 0) {
            filteredCommands = fuzzySearch(allCommands, searchQuery, (cmd) => `${cmd.label} ${cmd.category} ${cmd.context}`);
        } else {
            filteredCommands = [...allCommands];
        }
        selectedIndex = Math.min(selectedIndex, filteredCommands.length - 1);
    }

    // Initialize when opening
    $: if (show) {
        addKeyboardContext('commandPalette');
        filteredCommands = [...allCommands];
        selectedIndex = 0;
        focusStore.focus('command-palette', paletteId);
        // Focus input after a short delay to ensure DOM is ready
        setTimeout(() => {
            inputElement?.focus();
        }, 0);
    }

    // Reset state when closing
    $: if (!show && previousShow) {
        removeKeyboardContext('commandPalette');
        searchQuery = '';
        selectedIndex = 0;
        focusStore.restorePrevious();
    }

    let shortcuts: Record<string, string> = {};
    keyBindings.subscribe(bindings => {
        shortcuts = Object.entries(bindings).reduce((acc, [command, binding]) => ({
            ...acc,
            [binding.description]: formatKeybinding(binding)
        }), {});
    });

    afterUpdate(() => {
        // If command palette was showing and is now hidden, disable vim mode
        if (previousShow && !show) {
            vimModeEnabled = false;
        }
        previousShow = show;
    });


    function handleKeyDown(event: KeyboardEvent) {
        if (!show) return;

        // Enable vim mode when Alt+J are pressed together
        if (event.altKey && event.key.toLowerCase() === 'j') {
            event.preventDefault();
            vimModeEnabled = true;
            return;
        }

        switch(event.key) {
            case 'ArrowDown':
                event.preventDefault();
                selectedIndex = (selectedIndex + 1) % filteredCommands.length;
                break;
            case 'ArrowUp':
                event.preventDefault();
                selectedIndex = selectedIndex - 1 < 0 
                    ? filteredCommands.length - 1 
                    : selectedIndex - 1;
                break;
            case 'j':
                if (vimModeEnabled) {
                    event.preventDefault();
                    selectedIndex = (selectedIndex + 1) % filteredCommands.length;
                }
                break;
            case 'k':
                if (vimModeEnabled) {
                    event.preventDefault();
                    selectedIndex = selectedIndex - 1 < 0 
                        ? filteredCommands.length - 1 
                        : selectedIndex - 1;
                }
                break;
            case 'Enter':
                event.preventDefault();
                if (filteredCommands[selectedIndex]) {
                    executeCommand(filteredCommands[selectedIndex]);
                }
                break;
            case 'Escape':
                event.preventDefault();
                closeCommandPalette();
                break;
        }
    }

    function executeCommand(command: Command) {
        if (command.action) {
            command.action();
        }
        closeCommandPalette();
    }

    function closeCommandPalette() {
        show = false;
        dispatch('close');
    }

    function handleClickOutside() {
        closeCommandPalette();
    }

    function handleClickInside(event: MouseEvent) {
        event.stopPropagation();
    }

    onDestroy(() => {
        removeKeyboardContext('commandPalette');
    });

    onMount(() => {
        window.addEventListener('keydown', handleKeyDown);
        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        };
    });
</script>

{#if show}
    <button 
        class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-start justify-center pt-[20vh]"
        on:click={handleClickOutside}
    >
        <button 
            class="command-palette-content w-[600px] bg-gray-900 rounded-lg shadow-xl border border-gray-700"
            on:click={handleClickInside}
        >
            <div class="relative">
                <div class="pl-10">
                    <Input
                        bind:value={searchQuery}
                        placeholder="Type a command or search..."
                        bind:this={inputElement}
                        autofocus
                    />
                </div>
                <div class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
                    <Search size={16} />
                </div>
            </div>

            {#if filteredCommands.length > 0}
                <div class="max-h-[60vh] overflow-y-auto">
                    <div class="divide-y divide-gray-800">
                        {#each filteredCommands as command, index}
                            <button
                                class="w-full px-4 py-2 flex items-center justify-between text-left hover:bg-gray-800 
                                    {index === selectedIndex ? 'bg-gray-800' : ''}"
                                on:click={() => executeCommand(command)}
                            >
                                <div class="flex items-center space-x-2">
                                    <span class="text-gray-300">{command.label}</span>
                                    <div class="flex items-center gap-2">
                                        {#if command.category}
                                            <span class="text-xs text-gray-500 px-1.5 py-0.5 bg-gray-800 rounded">{command.category}</span>
                                        {/if}
                                        <span class="text-xs text-gray-600">{command.context}</span>
                                    </div>
                                </div>
                                {#if command.shortcut}
                                    <div class="flex items-center space-x-1">
                                        <span class="px-1.5 py-0.5 bg-gray-800 rounded text-xs text-gray-400 border border-gray-700">
                                            {command.shortcut}
                                        </span>
                                    </div>
                                {/if}
                            </button>
                        {/each}
                    </div>
                </div>
            {:else}
                <div class="px-4 py-8 text-center text-gray-500">
                    No commands found
                </div>
            {/if}
        </button>
    </button>
{/if}
