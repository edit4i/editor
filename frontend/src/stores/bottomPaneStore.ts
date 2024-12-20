import { writable } from 'svelte/store';

// Store for bottom pane state
export const bottomPaneStore = writable({
    collapsed: false,
    activeSection: 'terminal',
    isAllCollapsed: false
});
