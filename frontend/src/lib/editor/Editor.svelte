<script lang="ts">
    import * as monaco from 'monaco-editor';
    import { onMount, onDestroy } from 'svelte';
    import { fileStore } from '@/stores/fileStore';
    import { editorConfigStore } from '@/stores/editorConfigStore';
    import { editorStateStore } from '@/stores/editorStateStore';
    import { initVimMode, VimMode } from 'monaco-vim';
    import { createEventDispatcher } from 'svelte';
    import { focusStore } from '@/stores/focusStore';
    import Breadcrumbs from "@/lib/editor/Breadcrumbs.svelte";

    const dispatch = createEventDispatcher();

    let editorContainer: HTMLElement;
    let editor: monaco.editor.IStandaloneCodeEditor;
    let currentModel: monaco.editor.ITextModel | null = null;
    let vimMode: { dispose: () => void } | null = null;
    let vimStatusBar: HTMLElement;
    let vimEnabled = false;
    let editorId = focusStore.generateId('editor');
    let previousPath: string | null = null;
    let editorReady = false;
    let editorInitialized = false;

    // Load editor config on mount
    onMount(async () => {
        await editorConfigStore.loadConfig();

        // Create editor with initial config
        const config = $editorConfigStore.editor;

        editor = monaco.editor.create(editorContainer, {
            theme: config.theme,
            fontSize: config.fontSize,
            tabSize: config.tabSize,
            wordWrap: config.wordWrap ? 'on' : 'off',
            lineNumbers: config.lineNumbers ? (config.relativeLines ? 'relative' : 'on') : 'off',
            minimap: {
                enabled: config.minimap
            },
            automaticLayout: true,
            scrollBeyondLastLine: false,
            glyphMargin: true,
            folding: true,
            lineDecorationsWidth: 10,
            renderLineHighlight: 'all',
            scrollbar: {
                verticalScrollbarSize: 10,
                horizontalScrollbarSize: 10
            },
            stickyScroll: {
                enabled: config.stickyScroll
            }
        });

        // Initialize vim mode if enabled in config
        vimEnabled = config.vim?.enabled || false;
        if (vimEnabled) {
            vimMode = initVimMode(editor, vimStatusBar);
            
            // Define custom :w command for saving
            VimMode.Vim.defineEx('write', 'w', () => {
                fileStore.saveFile(fileStore.getActiveFilepath() || '');
            });

            // Setup Ctrl+P for file finder using action system
            VimMode.Vim.defineAction('showFileFinder', (_cm: any) => {
                dispatch('showFileFinder');
            });
            VimMode.Vim.mapCommand('<C-p>', 'action', 'showFileFinder', {}, { context: 'normal' });
        }

        // Wait a bit for editor to fully initialize
        setTimeout(() => {
            editorInitialized = true;
            editorReady = true;

            // Set initial focus
            focusStore.focus('editor', editorId);
            editor.focus();

            // Initial file state restore if needed
            if ($fileStore.activeFilePath) {
                handleFileChange($fileStore.activeFilePath);
            }
        }, 100);
    });

    // Watch for file changes and save/restore editor state
    function handleFileChange(newPath: string) {
        if (!editor || !editorReady || !editorInitialized || newPath === previousPath) return;

        // Save state of previous file before switching
        if (previousPath) {
            editorStateStore.saveState(previousPath, editor);
        }

        // Restore state of current file after a short delay to ensure model is ready
        setTimeout(() => {
            if (editor && editorInitialized) {
                editorStateStore.restoreState(newPath, editor);
                previousPath = newPath;
            }
        }, 100);
    }

    // Watch for focus changes
    $: if (editor && $focusStore.activeContext?.component === 'editor') {
        // Focus editor when it becomes active
        editor.focus();
    }

    // Watch for file changes
    $: if ($fileStore.activeFilePath) {
        handleFileChange($fileStore.activeFilePath);
    }

    // Watch for config changes
    $: if (editor && $editorConfigStore) {
        const config = $editorConfigStore.editor;
        editor.updateOptions({
            theme: config.theme,
            fontSize: config.fontSize,
            tabSize: config.tabSize,
            wordWrap: config.wordWrap ? 'on' : 'off',
            lineNumbers: config.lineNumbers ? (config.relativeLines ? 'relative' : 'on') : 'off',
            minimap: {
                enabled: config.minimap
            }
        });

        // Watch for vim config changes
        if (vimStatusBar && $editorConfigStore.editor.vim?.enabled !== vimEnabled) {
            vimEnabled = $editorConfigStore.editor.vim?.enabled || false;
            if (vimEnabled && !vimMode) {
                // Enable Vim mode
                vimMode = initVimMode(editor, vimStatusBar);
            } else if (!vimEnabled && vimMode) {
                // Disable Vim mode
                vimMode.dispose();
                vimMode = null;
            }
        }
    }

    function getLanguageFromPath(path: string): string {
        const ext = path.split('.').pop()?.toLowerCase() || '';
        switch (ext) {
            case 'js':
                return 'javascript';
            case 'jsx':
                return 'javascript';
            case 'ts':
                return 'typescript';
            case 'tsx':
                return 'typescript';
            case 'html':
            case 'svelte':
                return 'html';
            case 'php':
                return 'php';
            case 'go':
                return 'go';
            case 'py':
                return 'python';
            case 'json':
                return 'json';
            case 'md':
                return 'markdown';
            case 'css':
                return 'css';
            default:
                return 'plaintext';
        }
    }

    // Subscribe to active file changes
    $: if ($fileStore.activeFilePath && $fileStore.openFiles.has($fileStore.activeFilePath)) {
        const file = $fileStore.openFiles.get($fileStore.activeFilePath)!;
        updateEditor(file);
    }

    async function updateEditor(file: { path: string; content: string; language: string }) {
        if (!editor) return;

        // Create or get model for this file
        let model = monaco.editor.getModel(monaco.Uri.file(file.path));
        if (!model) {
            model = monaco.editor.createModel(
                file.content,
                getLanguageFromPath(file.path),
                monaco.Uri.file(file.path)
            );
        } else {
            const currentValue = model.getValue();
            // Only update if content actually changed
            if (currentValue !== file.content) {
                model.setValue(file.content);
            }
        }

        // Set the model as current
        editor.setModel(model);
        
        // Focus the editor after setting the model
        editor.focus();

        // Set as current model
        currentModel = model;

        // Listen for content changes
        const disposable = model.onDidChangeContent(() => {
            const currentContent = model.getValue();
            if (currentContent !== file.content) {
                fileStore.updateFileContent(file.path, currentContent, true);
            }
        });

        // Add command for saving
        editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, async () => {
            if (file.path) {
                try {
                    await fileStore.saveFile(file.path);
                } catch (error) {
                    console.error('Error saving file:', error);
                }
            }
        });

        return () => disposable.dispose();
    }

    // Focus handling
    function handleEditorFocus() {
        if (editor) {
            focusStore.focus('editor', editorId);
            editor.focus();
        }
    }

    onDestroy(() => {
        if (vimMode) {
            vimMode.dispose();
        }
        if (editor) {
            editor.dispose();
        }
    });
</script>

<div class="h-full relative flex flex-col">
    {#if $fileStore.activeFilePath}
        <Breadcrumbs filepath={$fileStore.activeFilePath} />
    {/if}
    
    <div class="flex-1 relative">
        <div
            bind:this={editorContainer}
            class="w-full h-full flex-1"
            on:focus={handleEditorFocus}
        />
    </div>

    {#if $editorConfigStore.editor.vim?.enabled}
        <div 
            bind:this={vimStatusBar}
            class="h-6 bg-gray-800 border-t border-gray-700 px-2 flex items-center text-sm absolute bottom-0 right-0 left-0 z-10"
        />
    {/if}
</div>
