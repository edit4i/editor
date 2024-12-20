<script lang="ts">
    import XtermComponent from '@/lib/terminal/XtermComponent.svelte';

    // Terminal component to be used inside BottomPane
    // Get the height from BottomPane
    export let height: number;
</script>

<div class="h-full w-full bg-gray-800 overflow-hidden">
    <XtermComponent {height} />
</div>
