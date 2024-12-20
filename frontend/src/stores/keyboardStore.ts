import { writable, derived, get } from 'svelte/store';
import type { KeyBinding, KeyBindingConfig, KeyboardContext } from '@/types/keyboard';
import { fileStore } from './fileStore';
import { OpenConfigFile } from '@/lib/wailsjs/go/main/App';

// Default keybindings configuration
const defaultKeybindings: KeyBindingConfig = {
    // Global shortcuts
    'command.showCommandPalette': {
        defaultBinding: {
            key: 'p',
            modifiers: ['ctrl', 'shift'],
            description: 'Show Command Palette',
            category: 'General',
            context: ['global']
        },
        action: () => {}  // This will be set by the command registration
    },
    'file.showFileFinder': {
        defaultBinding: {
            key: 'p',
            modifiers: ['ctrl'],
            description: 'Show File Finder',
            category: 'File',
            context: ['global']
        },
        action: () => {}
    },
    'modal.close': {
        defaultBinding: {
            key: 'Escape',
            description: 'Close Modal',
            category: 'General',
            context: ['commandPalette', 'fileFinder', 'aiAssistant']
        },
        action: () => {}
    },

    // Navigation mode shortcuts
    'vim.enableMode': {
        defaultBinding: {
            key: 'j',
            modifiers: ['alt'],
            description: 'Enable Vim Mode',
            category: 'Navigation',
            context: ['commandPalette', 'fileFinder']
        },
        action: () => {}
    },

    // Navigation shortcuts
    'navigation.goToEditor': {
        defaultBinding: {
            key: 'e',
            modifiers: ['ctrl', 'shift'],
            description: 'Go to Editor',
            category: 'Navigation',
            context: ['global']
        },
        action: () => {}
    },
    'navigation.goToSettings': {
        defaultBinding: {
            key: ',',
            modifiers: ['ctrl'],
            description: 'Open Settings',
            category: 'Navigation',
            context: ['global']
        },
        action: () => {}
    },
    'view.toggleLeftSidebar': {
        defaultBinding: {
            key: 'b',
            modifiers: ['ctrl'],
            description: 'Toggle Left Sidebar',
            category: 'View',
            context: ['global']
        },
        action: () => {}
    },
    'view.toggleRightSidebar': {
        defaultBinding: {
            key: 'b',
            modifiers: ['ctrl', 'alt'],
            description: 'Toggle Right Sidebar',
            category: 'View',
            context: ['global']
        },
        action: () => {}
    },

    // AI Assistant shortcuts
    'ai.sendMessage': {
        defaultBinding: {
            key: 'Enter',
            modifiers: ['ctrl'],
            description: 'Send Message',
            category: 'AI Assistant',
            context: ['aiAssistant']
        },
        action: () => {}
    },
    'ai.newConversation': {
        defaultBinding: {
            key: 'n',
            modifiers: ['ctrl'],
            description: 'New Conversation',
            category: 'AI Assistant',
            context: ['aiAssistant']
        },
        action: () => {}
    },

    // Git shortcuts
    'git.commit': {
        defaultBinding: {
            key: 'Enter',
            modifiers: ['ctrl'],
            description: 'Commit Changes',
            category: 'Git',
            context: ['git']
        },
        action: () => {}
    },
    'git.push': {
        defaultBinding: {
            key: 'p',
            modifiers: ['ctrl', 'alt'],
            description: 'Push Changes',
            category: 'Git',
            context: ['git']
        },
        action: () => {}
    },
    'git.pull': {
        defaultBinding: {
            key: 'l',
            modifiers: ['ctrl', 'alt'],
            description: 'Pull Changes',
            category: 'Git',
            context: ['git']
        },
        action: () => {}
    },
    'git.stash': {
        defaultBinding: {
            key: 's',
            modifiers: ['ctrl', 'alt'],
            description: 'Stash Changes',
            category: 'Git',
            context: ['git']
        },
        action: () => {}
    },

    // File Manager shortcuts
    'file.createFile': {
        defaultBinding: {
            key: 'n',
            modifiers: ['ctrl'],
            description: 'New File',
            category: 'File Manager',
            context: ['fileManager']
        },
        action: () => {}
    },
    'file.createFolder': {
        defaultBinding: {
            key: 'n',
            modifiers: ['ctrl', 'shift'],
            description: 'New Folder',
            category: 'File Manager',
            context: ['fileManager']
        },
        action: () => {}
    },
    'file.delete': {
        defaultBinding: {
            key: 'Delete',
            description: 'Delete Selected',
            category: 'File Manager',
            context: ['fileManager']
        },
        action: () => {}
    },
    'file.rename': {
        defaultBinding: {
            key: 'F2',
            description: 'Rename',
            category: 'File Manager',
            context: ['fileManager']
        },
        action: () => {}
    },
    'file.copy': {
        defaultBinding: {
            key: 'c',
            modifiers: ['ctrl'],
            description: 'Copy',
            category: 'File Manager',
            context: ['fileManager']
        },
        action: () => {}
    },
    'file.paste': {
        defaultBinding: {
            key: 'v',
            modifiers: ['ctrl'],
            description: 'Paste',
            category: 'File Manager',
            context: ['fileManager']
        },
        action: () => {}
    },
    'file.save': {
        defaultBinding: {
            key: 's',
            modifiers: ['ctrl'],
            description: 'Save File',
            category: 'File',
            context: ['global']
        },
        action: async () => {
            try {
                await fileStore.saveFile(fileStore.getActiveFilepath() || '');
            } catch (error) {
                console.error('Error saving file:', error);
            }
        }
    },
    'config.open': {
        defaultBinding: {
            key: '',
            modifiers: [],
            description: 'Open Editor Configuration',
            category: 'Settings',
            context: ['global']
        },
        action: async () => {
            try {
                const configPath = await OpenConfigFile();
                if (configPath) {
                    await fileStore.openFile(configPath);
                }
            } catch (error) {
                console.error('Error opening config file:', error);
            }
        }
    },
    // Tab cycling
    'cycleNextTab': {
        defaultBinding: {
            key: 'tab',
            modifiers: ['ctrl'],
            description: 'Cycle to next tab',
            category: 'Tabs',
            context: ['editor']
        },
        action: () => {
            const openFiles = get(fileStore).openFiles;
            const currentFile = fileStore.getActiveFilepath();
            
            if (!currentFile || openFiles.size <= 1) return;
            
            const files = Array.from(openFiles.keys());
            const currentIndex = files.indexOf(currentFile);
            const nextIndex = (currentIndex + 1) % files.length;
            
            fileStore.openFile(files[nextIndex]);
        }
    },
    'cyclePrevTab': {
        defaultBinding: {
            key: 'tab',
            modifiers: ['ctrl', 'shift'],
            description: 'Cycle to previous tab',
            category: 'Tabs',
            context: ['editor']
        },
        action: () => {
            const openFiles = get(fileStore).openFiles;
            const currentFile = fileStore.getActiveFilepath();
            
            if (!currentFile || openFiles.size <= 1) return;
            
            const files = Array.from(openFiles.keys());
            const currentIndex = files.indexOf(currentFile);
            const prevIndex = (currentIndex - 1 + files.length) % files.length;
            
            fileStore.openFile(files[prevIndex]);
        }
    },
    // Direct tab selection (max 9)
    ...Array.from({ length: Math.min(9, get(fileStore).openFiles.size) }, (_, i) => ({
        [`goToTab${i + 1}`]: {
            defaultBinding: {
                key: String(i + 1),
                modifiers: ['alt'],
                description: `Go to tab ${i + 1}`,
                category: 'Tabs',
                context: ['editor']
            },
            action: () => {
                const openFiles = get(fileStore).openFiles;
                const files = Array.from(openFiles.keys());
                
                if (i < files.length) {
                    fileStore.openFile(files[i]);
                }
            }
        }
    })).reduce((acc, binding) => ({ ...acc, ...binding }), {}),

    // Fuzzy finder quick selection (1-9)
    ...Array.from({ length: 9 }, (_, i) => ({
        [`fuzzyFinderSelect${i + 1}`]: {
            defaultBinding: {
                key: String(i + 1),
                modifiers: ['alt'],
                description: `Select result ${i + 1}`,
                category: 'File Finder',
                context: ['fileFinder']
            },
            action: () => {
                // Action will be handled by the FileFinder component
                // since it needs access to the current results
            }
        }
    })).reduce((acc, binding) => ({ ...acc, ...binding }), {}),
    // Terminal shortcuts
    'terminal.open': {
        defaultBinding: {
            key: 'j',
            modifiers: ['ctrl'],
            description: 'Open Terminal',
            category: 'Terminal',
            context: ['global']
        },
        action: () => {}
    },
    'terminal.returnToPrevious': {
        defaultBinding: {
            key: 'j',
            modifiers: ['alt'],
            description: 'Return to Previous Pane',
            category: 'Terminal',
            context: ['bottomPane']
        },
        action: () => {}
    },
};

    // Store for custom keybindings
    const customKeybindings = writable<Record<string, Partial<KeyBinding>>>({});

    // Store for current keyboard contexts (now supports multiple)
    export const activeContexts = writable<Set<KeyboardContext>>(new Set(['global']));

    // Function to add a keyboard context
    export function addKeyboardContext(context: KeyboardContext) {
        activeContexts.update(contexts => {
            // Create a new Set to ensure reactivity
            return new Set([...Array.from(contexts), context]);
        });
    }

    // Function to remove a keyboard context
    export function removeKeyboardContext(context: KeyboardContext) {
        if (context === 'global') return; // Never remove global context
        activeContexts.update(contexts => {
            // Create a new Set to ensure reactivity
            const newContexts = new Set(Array.from(contexts));
            newContexts.delete(context);
            return newContexts;
        });
    }

    // Function to check if a context is active
    export function hasKeyboardContext(context: KeyboardContext): boolean {
        return get(activeContexts).has(context);
    }

    // Derived store that combines default and custom keybindings
    export const keyBindings = derived(
        [customKeybindings, activeContexts],
        ([$customKeybindings, $activeContexts]) => {
            const bindings: Record<string, KeyBinding> = {};
            // Use Array.from to ensure proper Set iteration
            const activeContextArray = Array.from($activeContexts || new Set(['global']));

            for (const [command, config] of Object.entries(defaultKeybindings)) {
                // Check if any of the binding's contexts are currently active
                if (config.defaultBinding.context.some(ctx => activeContextArray.includes(ctx))) {
                    bindings[command] = {
                        ...config.defaultBinding,
                        action: config.action,
                        ...($customKeybindings[command] || {})
                    };
                }
            }

            return bindings;
        }
    );

    // Keep setKeyboardContext for backwards compatibility, but mark as deprecated
    /** @deprecated Use addKeyboardContext and removeKeyboardContext instead */
    export function setKeyboardContext(context: KeyboardContext) {
        // Create a new Set with just the global context and the new context
        activeContexts.set(new Set(['global', context]));
    }

    // Function to register a command's action
    export function registerCommand(command: string, action: () => void) {
        if (defaultKeybindings[command]) {
            defaultKeybindings[command].action = action;
            // Trigger store update
            customKeybindings.update(k => ({ ...k }));
        }
    }

    // Function to update custom keybinding
    export function updateKeybinding(command: string, binding: Partial<KeyBinding>) {
        customKeybindings.update(bindings => ({
            ...bindings,
            [command]: binding
        }));
    }

    // Helper function to format keybinding for display
    export function formatKeybinding(binding: KeyBinding): string {
        const modifiers = binding.modifiers || [];
        const parts = [
            ...modifiers.map(mod => mod.charAt(0).toUpperCase() + mod.slice(1)),
            binding.key.toUpperCase()
        ];
        return parts.join('+');
    }

    // Function to check if a keyboard event matches a keybinding
    export function matchesKeybinding(event: KeyboardEvent, binding: KeyBinding): boolean {
        const modifiers = binding.modifiers || [];

        const modifierMatch =
            modifiers.includes('ctrl') === event.ctrlKey &&
            modifiers.includes('alt') === event.altKey &&
            modifiers.includes('shift') === event.shiftKey &&
            modifiers.includes('meta') === event.metaKey;

        return modifierMatch && event.key.toLowerCase() === binding.key.toLowerCase();
    }

    // Function to handle global keyboard events
    export function handleKeyboardEvent(event: KeyboardEvent) {
      // Ignore enter when editor is focused
      if (event.target instanceof HTMLElement &&
          event.target.closest('.monaco-editor') &&
          event.key === 'Enter') {
          return;
      }

        const currentBindings = get(keyBindings);

        for (const [command, binding] of Object.entries(currentBindings)) {
            if (matchesKeybinding(event, binding)) {
                event.preventDefault();
                binding.action();
                return;
            }
        }
    }
