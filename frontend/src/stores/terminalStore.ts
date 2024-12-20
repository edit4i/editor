import { writable, get } from 'svelte/store';
import { projectStore } from './project';
import { editorConfigStore } from '@/stores/editorConfigStore';
import { GetAvailableShells } from '@/lib/wailsjs/go/main/App';

export interface TerminalTab {
    id: string;
    name: string;
    active: boolean;
    shell: string;
}

// Initialize with a default shell that will be updated
export const availableShells = writable<string[]>(['/bin/bash']);

// Load available shells on startup
GetAvailableShells().then(shells => {
    availableShells.set(shells);
}).catch(err => {
    console.error('Failed to get available shells:', err);
});

function createTerminalStore() {
    const { subscribe, update, set } = writable<TerminalTab[]>([]);

    return {
        subscribe,
        addTab: (shell: string = '') => {
            let newId: string = '';
            update(tabs => {
                // Get default shell from config or use first available shell
                const config = get(editorConfigStore);
                const currentShells = get(availableShells);
                const defaultShell = config.terminal.DefaultShell || currentShells[0] || '/bin/sh';
                const shellToUse = shell || defaultShell;
                
                // Deactivate all tabs
                const updatedTabs = tabs.map(tab => ({ ...tab, active: false }));
                // Add new tab
                newId = String(Date.now()); // Use timestamp for unique IDs
                return [...updatedTabs, { 
                    id: newId, 
                    name: `Terminal ${tabs.length + 1}`, 
                    active: true,
                    shell: shellToUse
                }];
            });
            return newId;
        },
        removeTab: (id: string) => {
            update(tabs => {
                if (tabs.length === 1) return tabs; // Don't remove last tab
                
                const index = tabs.findIndex(tab => tab.id === id);
                const wasActive = tabs[index]?.active;
                
                // Remove the tab
                const newTabs = tabs.filter(tab => tab.id !== id);
                
                // If we removed the active tab, activate the previous one (or the last one)
                if (wasActive && newTabs.length > 0) {
                    const newActiveIndex = Math.min(index, newTabs.length - 1);
                    newTabs[newActiveIndex].active = true;
                }
                
                return newTabs;
            });
        },
        setActiveTab: (id: string) => {
            update(tabs => 
                tabs.map(tab => ({
                    ...tab,
                    active: tab.id === id
                }))
            );
        }
    };
}

export const terminalStore = createTerminalStore();
