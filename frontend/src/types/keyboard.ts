export type ModifierKey = 'ctrl' | 'alt' | 'shift' | 'meta';

export type KeyboardContext = 
    | 'global' 
    | 'commandPalette' 
    | 'fileFinder' 
    | 'editor'
    | 'aiAssistant'
    | 'git'
    | 'fileManager'
    | 'bottomPane';

export type KeyBinding = {
    key: string;
    modifiers?: ModifierKey[];
    description: string;
    category?: string;
    context: KeyboardContext[];
    action: () => void;
};

export type KeyBindingConfig = {
    [command: string]: {
        defaultBinding: Omit<KeyBinding, 'action'>;
        action: () => void;
    };
};
