import { writable } from 'svelte/store';

export interface TerminalTab {
    id: string;
    name: string;
    active: boolean;
    shell: string;
}

export const AVAILABLE_SHELLS = [
    'bash',
    'zsh',
    'fish',
    'sh'
] as const;

function createTerminalStore() {
    const { subscribe, update } = writable<TerminalTab[]>([
        { id: '1', name: 'Terminal 1', active: true, shell: 'bash' }
    ]);

    return {
        subscribe,
        addTab: (shell: string = 'bash') => {
            update(tabs => {
                // Deactivate all tabs
                const updatedTabs = tabs.map(tab => ({ ...tab, active: false }));
                // Add new tab
                const newId = (tabs.length + 1).toString();
                return [...updatedTabs, { 
                    id: newId, 
                    name: `Terminal ${newId}`, 
                    active: true,
                    shell
                }];
            });
        },
        removeTab: (id: string) => {
            update(tabs => {
                const index = tabs.findIndex(tab => tab.id === id);
                if (tabs.length === 1) return tabs; // Don't remove last tab
                
                const newTabs = tabs.filter(tab => tab.id !== id);
                // If we removed the active tab, activate the previous one (or the last one)
                if (tabs[index].active) {
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
