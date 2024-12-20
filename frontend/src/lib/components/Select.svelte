<script lang="ts">
    import { twMerge } from "tailwind-merge";

    export let value: string;
    export let options: string[];
    export let disabled: boolean = false;
    export let variant: 'default' | 'compact' = 'default';

    const baseClasses = `
        w-full
        bg-gray-800
        text-gray-300
        text-sm
        rounded
        border
        border-gray-700
        focus:outline-none
        focus:ring-1
        focus:ring-blue-500
        focus:border-blue-500
        disabled:opacity-50
        disabled:cursor-not-allowed
        appearance-none
    `;

    const variantClasses = {
        default: 'p-2',
        compact: 'py-0.5 px-2'
    };

    $: classes = twMerge(baseClasses, variantClasses[variant], $$props.class || '');
</script>

<div class="relative">
    <select
        bind:value
        {disabled}
        class={classes}
        style="background-color: rgb(31, 41, 55)" 
    >
        {#each options as option}
            <option value={option}>{option}</option>
        {/each}
    </select>
    <div class="absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none text-gray-400">
        ▼
    </div>
</div>

<style>
    select option {
        background-color: rgb(31, 41, 55);
        color: rgb(209, 213, 219);
    }
</style>
