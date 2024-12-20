<script lang="ts">
    import { RefreshCw, ChevronsUp, FilesIcon } from 'lucide-svelte';
    import Button from '../../components/Button.svelte';
    import FileTree from '../../components/FileTree.svelte';
    import { fileStore } from '@/stores/fileStore';
    import type { service } from '@/lib/wailsjs/go/models';

    type FileNode = service.FileNode;
    
    export let fileTree: FileNode[] = [];
    let isAllCollapsed = false;
    let collapseCounter = 0;  // Add a counter to force re-collapse

    function handleCollapseAll() {
        collapseCounter++;
        isAllCollapsed = true;
        // Force a re-collapse after a short delay
        setTimeout(() => {
            isAllCollapsed = false;
        }, 100);
    }
</script>

<div class="flex flex-col h-full">
    <div class="flex items-center justify-between h-[35px] px-4 border-b border-gray-800">
        <div class="flex items-center space-x-2">
            <FilesIcon size={16} />
            <span class="text-sm font-medium">Explorer</span>
        </div>
        <div class="flex items-center space-x-1">
            <Button
                variant="ghost"
                title="Refresh"
                on:click={() => fileStore.refreshFiles()}
            >
                <RefreshCw size={14} />
            </Button>
            <Button
                variant="ghost"
                title="Collapse All"
                on:click={handleCollapseAll}
            >
                <ChevronsUp size={14} />
            </Button>
        </div>
    </div>
    <div class="flex-1 overflow-auto">
        <FileTree {fileTree} bind:isAllCollapsed key={collapseCounter} />
    </div>
</div>
