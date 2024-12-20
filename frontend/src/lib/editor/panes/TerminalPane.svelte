<script lang="ts">
    import XtermComponent from '@/lib/terminal/XtermComponent.svelte';
    import { terminalStore, AVAILABLE_SHELLS } from '@/stores/terminalStore';
    import { Plus, X } from 'lucide-svelte';
    import Button from '@/lib/components/Button.svelte';
    import Select from '@/lib/components/Select.svelte';

    // Get the height from BottomPane
    export let height: number;

    let selectedShell = AVAILABLE_SHELLS[0];

    // Handle tab actions
    function handleTabClick(event: MouseEvent, id: string) {
        if (event.button === 1) { // Middle click
            event.preventDefault();
            terminalStore.removeTab(id);
        } else {
            terminalStore.setActiveTab(id);
        }
    }
</script>

<div class="h-full w-full bg-gray-800 overflow-hidden flex flex-col">
    <div class="flex items-center border-b border-gray-700">
        <div class="flex-1 flex items-center overflow-x-auto">
            {#each $terminalStore as tab (tab.id)}
                <button
                    class="flex items-center h-[35px] px-4 border-r border-gray-700 hover:bg-gray-700 transition-colors duration-200 gap-2 group"
                    class:bg-gray-900={tab.active}
                    class:before:absolute={tab.active}
                    class:before:top-0={tab.active}
                    class:before:left-0={tab.active}
                    class:before:right-0={tab.active}
                    class:before:h-[2px]={tab.active}
                    class:before:bg-sky-500={tab.active}
                    on:click={(e) => handleTabClick(e, tab.id)}
                    on:mouseup={(e) => handleTabClick(e, tab.id)}
                >
                    {tab.name} ({tab.shell})
                    {#if $terminalStore.length > 1}
                        <button
                            class="opacity-0 group-hover:opacity-100 hover:text-sky-500 transition-opacity duration-200"
                            on:click|stopPropagation={() => terminalStore.removeTab(tab.id)}
                        >
                            <X size={14} />
                        </button>
                    {/if}
                </button>
            {/each}
        </div>
        <div class="flex items-center gap-2 px-2">
            <Select
                bind:value={selectedShell}
                options={AVAILABLE_SHELLS}
                variant="compact"
                class="w-36"
            />
            <Button
                variant="ghost"
                on:click={() => terminalStore.addTab(selectedShell)}
            >
                <Plus size={14} />
            </Button>
        </div>
    </div>

    {#each $terminalStore as tab (tab.id)}
        {#if tab.active}
            <div class="flex-1">
                <XtermComponent {height} id={tab.id} shell={tab.shell} />
            </div>
        {/if}
    {/each}
</div>
