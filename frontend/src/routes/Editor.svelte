<script lang="ts">
    import { X, Circle, ChevronLeft, ChevronRight } from "lucide-svelte";
    import { onMount, onDestroy } from 'svelte';
    import type { SidebarState } from "@/types/ui";
    import LeftSidebar from "@/lib/editor/LeftSidebar.svelte";
    import RightSidebar from "@/lib/editor/RightSidebar.svelte";
    import ResizeHandle from "@/lib/editor/ResizeHandle.svelte";
    import Topbar from "@/lib/editor/Topbar.svelte";
    import BottomBar from "@/lib/editor/BottomBar.svelte";
    import { fileStore } from '@/stores/fileStore';
    import { projectStore } from '@/stores/project';
    import { registerCommand, setKeyboardContext } from '@/stores/keyboardStore';
    import { get } from 'svelte/store';
    import Editor from "@/lib/editor/Editor.svelte";
    import FileFinder from "@/lib/components/FileFinder.svelte";
    import Modal from "@/lib/components/Modal.svelte";
    import BottomPane from "@/lib/editor/panes/BottomPane.svelte";
    import { bottomPaneStore } from '@/stores/bottomPaneStore';
    import { focusStore } from '@/stores/focusStore';

    // Convert open files to tabs
    $: tabs = Array.from($fileStore.openFiles.entries()).map(([path, file]) => ({
        id: path,
        name: path.split('/').pop() || '',
        active: path === $fileStore.activeFilePath,
        isDirty: file.isDirty
    }));

    // Sidebar states
    let leftSidebarState: SidebarState = {
        collapsed: false,
        activeSection: 'files',
        isAllCollapsed: false
    };

    let rightSidebarCollapsed = false;
    let isSourceControlActive = false;
    let isExplorerActive = true;
    let showCommandPalette = false;
    let showFileFinder = false;

    // Bottom pane state
    let bottomPaneState = $bottomPaneStore;

    // Sidebar widths and heights
    let leftSidebarWidth = 300;
    let rightSidebarWidth = 600;
    let bottomPaneHeight = 300;

    // Source control state
    let modifiedFilesCount = 2;

    let tabsContainer: HTMLDivElement;

    let showCloseConfirmModal = false;
    let fileToClose: string | null = null;

    function setActiveTab(id: string) {
        fileStore.setActiveFile(id);
        scrollToTab(id);
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

    // Watch for active file changes
    $: if ($fileStore.activeFilePath) {
        // Wait for the DOM to update before scrolling
        setTimeout(() => scrollToTab($fileStore.activeFilePath), 0);
    }

    function handleCloseTab(id: string) {
        const file = $fileStore.openFiles.get(id);
        if (file?.isDirty) {
            fileToClose = id;
            showCloseConfirmModal = true;
        } else {
            closeTab(id);
        }
    }

    function confirmCloseTab() {
        if (fileToClose) {
            closeTab(fileToClose);
            fileToClose = null;
            showCloseConfirmModal = false;
        }
    }

    function closeTab(id: string) {
        fileStore.closeFile(id);
    }

    function handleResize() {
        editor?.layout();
    }

    function scrollTabs(direction: 'left' | 'right') {
        if (tabsContainer) {
            const scrollAmount = 200;
            const targetScroll = tabsContainer.scrollLeft + (direction === 'left' ? -scrollAmount : scrollAmount);
            tabsContainer.scrollTo({
                left: targetScroll,
                behavior: 'smooth'
            });
        }
    }

    function handleTabClick(event: MouseEvent, id: string) {
        if (event.button === 1) { // Middle click
            event.preventDefault();
            handleCloseTab(id);
        }
    }

    onMount(() => {
        const state = get(projectStore);
        if (state.currentProject?.Path) {
            fileStore.loadProjectFiles(state.currentProject.Path);
        }
        
        setKeyboardContext('global');

        // Register terminal shortcut
        registerCommand('terminal.open', () => {
            // Save current focus
            focusStore.focus('editor', $fileStore.activeFilePath || 'editor');
            
            // Show terminal
            bottomPaneStore.update(state => ({
                ...state,
                collapsed: false,
                activeSection: 'terminal'
            }));
        });

        // Register return to previous shortcut
        registerCommand('terminal.returnToPrevious', () => {
            const previous = get(focusStore).previousContext;
            if (previous && previous.component === 'editor') {
                // Collapse terminal
                bottomPaneStore.update(state => ({
                    ...state,
                    collapsed: true
                }));

                // Focus editor and file
                focusStore.restorePrevious();
                if (previous.id !== 'editor') {
                    fileStore.setActiveFile(previous.id);
                }
            }
        });

        registerCommand('file.showFileFinder', () => showFileFinder = true);
        
        // Register sidebar toggle commands
        registerCommand('view.toggleLeftSidebar', () => {
            leftSidebarState.collapsed = !leftSidebarState.collapsed;
        });
        registerCommand('view.toggleRightSidebar', () => {
            rightSidebarCollapsed = !rightSidebarCollapsed;
        });
    });

    onDestroy(() => {
        setKeyboardContext('global');
    });

</script>

<div class="flex flex-col h-screen bg-gray-900 text-gray-300">
    <Topbar 
        bind:isLeftSidebarCollapsed={leftSidebarState.collapsed}
        bind:isRightSidebarCollapsed={rightSidebarCollapsed}
        bind:isSourceControlActive
        bind:isExplorerActive
        {modifiedFilesCount}
        onToggleLeftSidebar={() => leftSidebarState.collapsed = !leftSidebarState.collapsed}
        onToggleRightSidebar={() => rightSidebarCollapsed = !rightSidebarCollapsed}
        onToggleSourceControl={() => isSourceControlActive = !isSourceControlActive}
        onToggleExplorer={() => isExplorerActive = !isExplorerActive}
        showCommandPalette={() => showCommandPalette = true}
        showFileFinder={() => showFileFinder = true}
    />
    
    <div class="flex flex-1 overflow-hidden">
        {#if !leftSidebarState.collapsed}
            <div class="h-full" style="width: {leftSidebarWidth}px">
                <LeftSidebar state={leftSidebarState} />
            </div>
            <ResizeHandle
                orientation="vertical"
                side="right"
                bind:size={leftSidebarWidth}
                minSize={200}
                maxSize={600}
            />
        {/if}

        
        <main class="flex-1 flex flex-col min-w-0 overflow-hidden relative">
            <div class="flex items-center border-b border-gray-800 bg-gray-900 relative">
                <div 
                    bind:this={tabsContainer}
                    class="flex overflow-x-scroll scrollbar-hide relative flex-1"
                >
                    {#each tabs as tab (tab.id)}
                        <button
                            data-tab-id={tab.id}
                            class="flex items-center h-[34px] px-4 border-r border-gray-800 cursor-pointer relative
                                {tab.active
                                    ? 'bg-gray-900 before:absolute before:top-0 before:left-0 before:right-0 before:h-[2px] before:bg-sky-500'
                                    : 'bg-gray-800 hover:bg-gray-700'} 
                                transition-colors duration-200"
                            on:click={() => setActiveTab(tab.id)}
                            on:mouseup={(e) => handleTabClick(e, tab.id)}
                        >
                            <span class="flex items-center gap-2">
                                {#if tab.isDirty}
                                    <Circle size={8} class="fill-current text-gray-300" />
                                {/if}
                                {tab.name}
                            </span>
                            <button
                                on:click|stopPropagation={() => handleCloseTab(tab.id)}
                                class="ml-2 text-gray-400 hover:text-gray-100 transition-colors duration-200"
                                aria-label="Close {tab.name}"
                            >
                                <X size={14} />
                            </button>
                        </button>
                    {/each}
                </div>
                <div class="flex items-center gap-1 border-gray-800 bg-gray-900 pl-1 sticky right-0">
                    <button
                        on:click={() => scrollTabs('left')}
                        class="p-1.5 hover:bg-gray-800 transition-colors duration-200 rounded-md border border-gray-700"
                        aria-label="Scroll tabs left"
                    >
                        <ChevronLeft size={16} />
                    </button>
                    <button
                        on:click={() => scrollTabs('right')}
                        class="p-1.5 hover:bg-gray-800 transition-colors duration-200 rounded-md border border-gray-700"
                        aria-label="Scroll tabs right"
                    >
                        <ChevronRight size={16} />
                    </button>
                </div>
            </div>

            <div class="flex-1 relative overflow-hidden">
                <Editor />
            </div>
            {#if !$bottomPaneStore.collapsed}
                <ResizeHandle 
                    orientation="horizontal" 
                    side="top"
                    bind:size={bottomPaneHeight}
                    minSize={100} 
                    maxSize={800}
                />
                <BottomPane state={$bottomPaneStore} height={bottomPaneHeight} />
            {/if}
        </main>
        
        {#if !rightSidebarCollapsed}
            <ResizeHandle
                orientation="vertical"
                side="left"
                bind:size={rightSidebarWidth}
                minSize={200}
                maxSize={800}
            />
            <div class="h-full" style="width: {rightSidebarWidth}px">
                <RightSidebar collapsed={rightSidebarCollapsed} />
            </div>
        {/if}
    </div>
    
    <BottomBar />

    <FileFinder 
        bind:show={showFileFinder} 
        on:close={() => showFileFinder = false} 
        on:select={({ detail }) => scrollToTab(detail.path)}
    />
    
    <Modal
        bind:show={showCloseConfirmModal}
        title="Unsaved Changes"
        confirmText="Close without saving"
        cancelText="Cancel"
        confirmButtonClass="bg-red-600 hover:bg-red-700"
        on:confirm={confirmCloseTab}
        on:close={() => showCloseConfirmModal = false}
    >
        <p>You have unsaved changes in this file. Are you sure you want to close it?</p>
    </Modal>
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