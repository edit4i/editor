<script lang="ts">
    import XtermComponent from '@/lib/terminal/XtermComponent.svelte';
    import { terminalStore, availableShells } from '@/stores/terminalStore';
    import { bottomPaneStore } from '@/stores/bottomPaneStore';
    import { Plus, X, ChevronLeft, ChevronRight } from 'lucide-svelte';
    import Button from '@/lib/components/Button.svelte';
    import Select from '@/lib/components/Select.svelte';
    import { get } from 'svelte/store';

    // Get the height from BottomPane
    export let height: number;

    let selectedShell = get(availableShells)[0];
    let tabsContainer: HTMLElement;
    let terminals: Record<string, XtermComponent> = {};

    // Watch for bottom pane state changes
    $: if ($bottomPaneStore.activeSection === 'terminal' && !$bottomPaneStore.collapsed) {
        // Find active terminal and focus it
        const activeTab = $terminalStore.find(tab => tab.active);
        if (activeTab) {
            terminals[activeTab.id]?.focus();
        }
    }

    // Handle tab actions
    function handleNewTab() {
        const newId = terminalStore.addTab(selectedShell);
        // Wait for the DOM to update before scrolling
        setTimeout(() => scrollToTab(newId), 0);
    }

    function handleTabClick(event: MouseEvent, id: string) {
        if (event.button === 1) { // Middle click
            event.preventDefault();
            terminalStore.removeTab(id);
        } else {
            terminalStore.setActiveTab(id);
            scrollToTab(id);
            // Focus the terminal when clicking its tab
            terminals[id]?.focus();
        }
    }

    function handleRemoveTab(event: MouseEvent, id: string) {
        event.stopPropagation();
        terminalStore.removeTab(id);
    }

    function scrollToTab(id: string) {
        if (tabsContainer) {
            // Find the tab element
            const tabElement = tabsContainer.querySelector(`[data-tab-id="${id}"]`);
            if (tabElement) {
                const containerWidth = tabsContainer.offsetWidth;
                const tabWidth = (tabElement as HTMLElement).offsetWidth;
                const tabLeft = (tabElement as HTMLElement).offsetLeft;
                const scrollLeft = tabLeft - (containerWidth - tabWidth) / 2;
                
                tabsContainer.scrollTo({
                    left: scrollLeft,
                    behavior: 'smooth'
                });
            }
        }
    }

    function scrollTabs(direction: 'left' | 'right') {
        if (tabsContainer) {
            const scrollAmount = 200;
            const newScrollLeft = tabsContainer.scrollLeft + (direction === 'left' ? -scrollAmount : scrollAmount);
            tabsContainer.scrollTo({
                left: newScrollLeft,
                behavior: 'smooth'
            });
        }
    }
</script>

<div class="h-full w-full bg-gray-800 overflow-hidden flex flex-col">
    <div class="flex items-center border-b border-gray-700">
        <div 
            bind:this={tabsContainer}
            class="flex-1 flex items-center overflow-x-scroll scrollbar-hide relative"
        >
            {#each $terminalStore as tab (tab.id)}
                <button
                    data-tab-id={tab.id}
                    class="flex items-center h-[35px] px-4 border-r border-gray-700 hover:bg-gray-700 transition-colors duration-200 gap-2 group whitespace-nowrap relative"
                    class:bg-gray-900={tab.active}
                    class:after:absolute={tab.active}
                    class:after:top-0={tab.active}
                    class:after:left-0={tab.active}
                    class:after:right-0={tab.active}
                    class:after:h-[2px]={tab.active}
                    class:after:bg-sky-500={tab.active}
                    on:click={(e) => handleTabClick(e, tab.id)}
                    on:mouseup={(e) => handleTabClick(e, tab.id)}
                >
                    {tab.name} ({tab.shell})
                    {#if $terminalStore.length > 1}
                        <button
                            class="opacity-0 group-hover:opacity-100 hover:text-sky-500 transition-opacity duration-200"
                            on:click={(e) => handleRemoveTab(e, tab.id)}
                        >
                            <X size={14} />
                        </button>
                    {/if}
                </button>
            {/each}
        </div>
        <div class="flex items-center gap-1 border-l border-gray-700 bg-gray-800 pl-1 sticky right-0">
            <button
                on:click={() => scrollTabs('left')}
                class="p-1.5 hover:bg-gray-700 transition-colors duration-200 text-gray-400 hover:text-gray-300"
            >
                <ChevronLeft size={16} />
            </button>
            <button
                on:click={() => scrollTabs('right')}
                class="p-1.5 hover:bg-gray-700 transition-colors duration-200 text-gray-400 hover:text-gray-300"
            >
                <ChevronRight size={16} />
            </button>
            <div class="h-[24px] border-l border-gray-700 mx-1" />
            <Select
                bind:value={selectedShell}
                options={$availableShells}
                variant="compact"
                class="w-36"
            />
            <Button
                variant="ghost"
                on:click={handleNewTab}
            >
                <Plus size={14} />
            </Button>
        </div>
    </div>

    <div class="flex-1 relative">
        {#each $terminalStore as tab (tab.id)}
            <div 
                class="absolute inset-0"
                class:hidden={!tab.active}
            >
                <XtermComponent 
                    bind:this={terminals[tab.id]} 
                    {height} 
                    id={tab.id} 
                    shell={tab.shell} 
                    active={tab.active} 
                />
            </div>
        {/each}
    </div>
</div>

<style>
    .scrollbar-hide {
        -ms-overflow-style: none;
        scrollbar-width: none;
    }
    .scrollbar-hide::-webkit-scrollbar {
        display: none;
    }
</style>
