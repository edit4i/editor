<script lang="ts">
    import { ChevronLeft, ChevronRight, ChevronUp, ChevronDown } from 'lucide-svelte';
    import { fade } from 'svelte/transition';
    
    export let orientation: 'horizontal' | 'vertical' = 'vertical';
    export let side: 'left' | 'right' | 'top' | 'bottom' = 'right';
    export let size: number;
    export let minSize = 200;
    export let maxSize = 600;
    
    let isDragging = false;
    let isHovering = false;
    let showHandle = false;
    let startPos: number;
    let startSize: number;

    const handleConfig = {
        left: { icon: ChevronLeft, position: 'top-1/2 -translate-y-1/2 -translate-x-1/2 left-0' },
        right: { icon: ChevronRight, position: 'top-1/2 -translate-y-1/2 translate-x-1/2 right-0' },
        top: { icon: ChevronUp, position: 'left-1/2 -translate-x-1/2 -translate-y-1/2 top-0' },
        bottom: { icon: ChevronDown, position: 'left-1/2 -translate-x-1/2 translate-y-1/2 bottom-0' }
    };
    
    function handleMouseDown(e: MouseEvent) {
        e.preventDefault();
        isDragging = true;
        showHandle = true;
        startPos = orientation === 'vertical' ? e.clientX : e.clientY;
        startSize = size;
        document.body.style.userSelect = 'none';
        document.body.style.cursor = orientation === 'vertical' ? 'ew-resize' : 'ns-resize';
            
        const mouseMoveHandler = (e: MouseEvent) => {
            if (!isDragging) return;
            e.preventDefault();
            const delta = orientation === 'vertical' 
                ? e.clientX - startPos 
                : e.clientY - startPos;
            
            let newSize;
            if (side === 'left' || side === 'top') {
                newSize = Math.max(minSize, Math.min(maxSize, startSize - delta));
            } else {
                newSize = Math.max(minSize, Math.min(maxSize, startSize + delta));
            }
            size = newSize;
        };
        
        const mouseUpHandler = () => {
            isDragging = false;
            showHandle = false;
            document.removeEventListener('mousemove', mouseMoveHandler);
            document.removeEventListener('mouseup', mouseUpHandler);
            document.body.style.userSelect = '';
            document.body.style.cursor = '';
        };
        
        document.addEventListener('mousemove', mouseMoveHandler);
        document.addEventListener('mouseup', mouseUpHandler);
    }
    
    function handleMouseEnter() {
        isHovering = true;
        showHandle = true;
    }
    
    function handleMouseLeave() {
        if (!isDragging) {
            isHovering = false;
            showHandle = false;
        }
    }

    $: isVertical = orientation === 'vertical';
    $: resizeClasses = isVertical 
        ? 'w-1 hover:w-[5px] transition-all duration-200 ease-out cursor-ew-resize'
        : 'h-1 hover:h-[5px] transition-all duration-200 ease-out cursor-ns-resize';
    $: tooltipPositionClasses = isVertical
        ? `top-1/2 -translate-y-1/2 ${side === 'left' ? 'left-4' : 'right-4'}`
        : `left-1/2 -translate-x-1/2 ${side === 'top' ? 'top-4' : 'bottom-4'}`;
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    class={`relative group select-none z-[9999] ${resizeClasses}`}
    class:bg-sky-500={isDragging}
    class:border-sky-500={isDragging}
    class:border-2={isDragging}
    on:mousedown={handleMouseDown}
    on:mouseenter={handleMouseEnter}
    on:mouseleave={handleMouseLeave}
>
    <div 
        class="absolute inset-0 bg-gray-700 opacity-0 group-hover:opacity-20 transition-opacity duration-200"
        class:opacity-40={isDragging}
        transition:fade={{ duration: 150 }}
    />
    
    {#if showHandle || isDragging}
        <div 
            class={`absolute ${handleConfig[side].position} bg-gray-700 rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200 z-[9999]`}
            class:bg-sky-500={isDragging}
            class:opacity-100={isDragging}
        >
            <svelte:component this={handleConfig[side].icon} size={16} class="text-white" />
        </div>
    {/if}
    
    <div 
        class={`absolute whitespace-nowrap bg-gray-800 text-gray-300 text-xs py-1 px-2 rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none select-none z-[9999] ${tooltipPositionClasses}`}
        class:opacity-100={isDragging}
    >
        Drag to resize
    </div>
</div>
