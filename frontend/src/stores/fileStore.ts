import { writable, get } from 'svelte/store';
import type { service } from '@/lib/wailsjs/go/models';
import { GetProjectFiles, GetFileContent, SaveFile, CreateFile, CreateDirectory, RenameFile, DeleteFile, LoadDirectoryContents } from '@/lib/wailsjs/go/main/App';
import { getLanguageFromPath } from '@/lib/utils/languageMap';

type FileNode = service.FileNode;

interface OpenFile {
    path: string;
    content: string;
    isDirty: boolean;
    language: string;
    cursor: { line: number; column: number };
}

interface FileState {
    fileTree: FileNode[] | null;
    activeFilePath: string | null;
    currentProjectPath: string | null;
    openFiles: Map<string, OpenFile>;
    loading: boolean;
    error: string | null;
}

// Load initial state from localStorage
const savedState = localStorage.getItem('fileState');
const initialState: FileState = savedState ? {
    ...JSON.parse(savedState),
    // Convert the openFiles object back to a Map
    openFiles: new Map(Object.entries(JSON.parse(savedState).openFiles || {}))
} : {
    fileTree: null,
    activeFilePath: null,
    currentProjectPath: null,
    openFiles: new Map(),
    loading: false,
    error: null,
};

function createFileStore() {
    const { subscribe, set, update } = writable<FileState>(initialState);

    // Save state changes to localStorage
    subscribe(state => {
        // Convert Map to object for JSON serialization
        const serializedState = {
            ...state,
            openFiles: Object.fromEntries(state.openFiles)
        };
        localStorage.setItem('fileState', JSON.stringify(serializedState));
    });

    return {
        subscribe,
        
        // Clear all state and localStorage
        clearState() {
            localStorage.removeItem('fileState');
            set({
                fileTree: null,
                activeFilePath: null,
                currentProjectPath: null,
                openFiles: new Map(),
                loading: false,
                error: null,
            });
        },

        // Set current project
        setCurrentProject(projectPath: string) {
            const state = get({ subscribe });
            
            // If changing projects, only keep open files from the new project
            if (projectPath !== state.currentProjectPath) {
                update(state => {
                    const newOpenFiles = new Map();
                    
                    // Only keep files that belong to the new project
                    state.openFiles.forEach((file, path) => {
                        if (path.startsWith(projectPath)) {
                            newOpenFiles.set(path, file);
                        }
                    });

                    return {
                        fileTree: null,
                        activeFilePath: Array.from(newOpenFiles.keys())[0] || null,
                        currentProjectPath: projectPath,
                        openFiles: newOpenFiles,
                        loading: false,
                        error: null,
                    };
                });
            }

            return this.loadProjectFiles();
        },

        // Load project files
        async loadProjectFiles(path?: string) {
            const state = get({ subscribe });
            if (!path && !state.currentProjectPath) return;

            path = path || (state.currentProjectPath ?? undefined);
            update(state => ({ ...state, loading: true, error: null }));

            try {
                const rootNode = await GetProjectFiles(path!);
                update(state => ({ 
                    ...state, 
                    fileTree: rootNode.children || [],
                    loading: false
                }));
            } catch (err) {
                update(state => ({
                    ...state,
                    error: err instanceof Error ? err.message : 'Failed to load project files',
                    loading: false
                }));
                throw err;
            }
        },

        // Open a file
        async openFile(path: string) {
            const state = get({ subscribe });
            if (state.openFiles.has(path)) {
                update(state => ({ ...state, activeFilePath: path }));
                return;
            }

            try {
                const content = await GetFileContent(path);
                update(state => {
                    const newOpenFiles = new Map(state.openFiles);
                    const openFile: OpenFile = {
                        path,
                        content,
                        isDirty: false,
                        language: getLanguageFromPath(path),
                        cursor: { line: 0, column: 0 }
                    };
                    newOpenFiles.set(path, openFile);
                    return {
                        ...state,
                        openFiles: newOpenFiles,
                        activeFilePath: path
                    };
                });
            } catch (err) {
                update(state => ({
                    ...state,
                    error: err instanceof Error ? err.message : 'Failed to open file'
                }));
            }
        },

        // Open a virtual file (for diffs, etc.)
        openVirtualFile: (path: string, content: string, language: string) => {
            update(state => {
                // Create virtual file
                const virtualFile: OpenFile = {
                    path,
                    content,
                    isDirty: false,
                    language,
                    cursor: { line: 0, column: 0 }
                };
                
                // Add to open files
                state.openFiles.set(path, virtualFile);
                state.activeFilePath = path;
                
                return state;
            });
        },

        // Close a file
        closeFile(path: string) {
            update(state => {
                const newOpenFiles = new Map(state.openFiles);
                newOpenFiles.delete(path);
                
                const activeFilePath = state.activeFilePath === path
                    ? Array.from(newOpenFiles.keys())[0] || null
                    : state.activeFilePath;
                
                return { ...state, openFiles: newOpenFiles, activeFilePath };
            });
        },

        // Refresh files
        async refreshFiles() {
            const state = get({ subscribe });
            if (!state.currentProjectPath) return;
            
            // Force a complete refresh of the project
            await this.loadProjectFiles(state.currentProjectPath);
        },

        // Set active file
        setActiveFile(path: string | null) {
            update(state => ({ ...state, activeFilePath: path }));
        },

        // Get currently open file
        getActiveFilepath() {
            return get({ subscribe }).activeFilePath;
        },

        // Mark file as dirty
        markAsDirty(path: string) {
            update(state => {
                const file = state.openFiles.get(path);
                if (!file || file.isDirty) return state; // Skip if already dirty

                const newOpenFiles = new Map(state.openFiles);
                newOpenFiles.set(path, { ...file, isDirty: true });
                return { ...state, openFiles: newOpenFiles };
            });
        },

        // Update file content
        updateFileContent(path: string, content: string, isDirty = true) {
            update(state => {
                const file = state.openFiles.get(path);
                if (!file) return state;

                const newOpenFiles = new Map(state.openFiles);
                newOpenFiles.set(path, { ...file, content, isDirty });
                return { ...state, openFiles: newOpenFiles };
            });
        },

        // Save file content
        async saveFile(path: string) {
            const state = get({ subscribe });
            const file = state.openFiles.get(path);
            if (!file) return;
            
            try {
                const content = file.content;
                await SaveFile(path, content);
                
                // Update the store to mark file as not dirty
                update(state => {
                    const file = state.openFiles.get(path);
                    if (file) {
                        const newOpenFiles = new Map(state.openFiles);
                        newOpenFiles.set(path, { 
                            ...file, 
                            isDirty: false 
                        });
                        return { ...state, openFiles: newOpenFiles };
                    }
                    return state;
                });
            } catch (err) {
                update(state => ({
                    ...state,
                    error: err instanceof Error ? err.message : 'Failed to save file'
                }));
            }
        },

        // Create new file
        async createFile(path: string): Promise<void> {
            try {
                // Get parent directory path
                const parentPath = path.substring(0, path.lastIndexOf("/"));
                
                await CreateFile(path);
                
                // Update only the parent directory
                const updatedNode = await LoadDirectoryContents(parentPath);
                if (updatedNode) {
                    update(state => {
                        if (!state.fileTree) return state;
                        
                        const newTree = [...state.fileTree];
                        const updateDir = (nodes: FileNode[]): boolean => {
                            for (const node of nodes) {
                                if (node.path === parentPath) {
                                    node.children = updatedNode.children;
                                    return true;
                                }
                                if (node.children && updateDir(node.children)) {
                                    return true;
                                }
                            }
                            return false;
                        };
                        updateDir(newTree);
                        return { ...state, fileTree: newTree };
                    });
                }
            } catch (error) {
                update(state => ({ ...state, error: `Failed to create file: ${error}` }));
                throw error;
            }
        },

        // Create new directory
        async createDirectory(path: string): Promise<void> {
            try {
                const parentPath = path.substring(0, path.lastIndexOf("/"));
                
                await CreateDirectory(path);
                
                // Update only the parent directory
                const updatedNode = await LoadDirectoryContents(parentPath);
                if (updatedNode) {
                    update(state => {
                        if (!state.fileTree) return state;
                        
                        const newTree = [...state.fileTree];
                        const updateDir = (nodes: FileNode[]): boolean => {
                            for (const node of nodes) {
                                if (node.path === parentPath) {
                                    node.children = updatedNode.children;
                                    return true;
                                }
                                if (node.children && updateDir(node.children)) {
                                    return true;
                                }
                            }
                            return false;
                        };
                        updateDir(newTree);
                        return { ...state, fileTree: newTree };
                    });
                }
            } catch (error) {
                update(state => ({ ...state, error: `Failed to create directory: ${error}` }));
                throw error;
            }
        },

        // Rename/Move file or directory
        async renameFile(oldPath: string, newPath: string): Promise<void> {
            try {
                const oldParentPath = oldPath.substring(0, oldPath.lastIndexOf("/"));
                const newParentPath = newPath.substring(0, newPath.lastIndexOf("/"));
                
                await RenameFile(oldPath, newPath);
                
                // Update both old and new parent directories
                const [oldUpdatedNode, newUpdatedNode] = await Promise.all([
                    LoadDirectoryContents(oldParentPath),
                    oldParentPath !== newParentPath ? LoadDirectoryContents(newParentPath) : null
                ]);

                update(state => {
                    if (!state.fileTree) return state;
                    
                    const newTree = [...state.fileTree];
                    const updateDir = (nodes: FileNode[], path: string, updatedChildren: FileNode[]): boolean => {
                        for (const node of nodes) {
                            if (node.path === path) {
                                node.children = updatedChildren;
                                return true;
                            }
                            if (node.children && updateDir(node.children, path, updatedChildren)) {
                                return true;
                            }
                        }
                        return false;
                    };

                    if (oldUpdatedNode) {
                        updateDir(newTree, oldParentPath, oldUpdatedNode.children);
                    }
                    if (newUpdatedNode) {
                        updateDir(newTree, newParentPath, newUpdatedNode.children);
                    }
                    return { ...state, fileTree: newTree };
                });
            } catch (error) {
                update(state => ({ ...state, error: `Failed to rename/move: ${error}` }));
                throw error;
            }
        },

        // Delete file or directory
        async deleteFile(path: string): Promise<void> {
            try {
                const state = get({ subscribe });
                const currentProjectPath = state.currentProjectPath;
                if (!currentProjectPath) return;

                // Get parent path, use project root if file is in root directory
                const lastSlashIndex = path.lastIndexOf("/");
                const parentPath = lastSlashIndex > 0 
                    ? path.substring(0, lastSlashIndex)
                    : currentProjectPath;
                
                // Optimistically remove from the tree first
                update(state => {
                    if (!state.fileTree) return state;
                    
                    const newTree = [...state.fileTree];
                    const removeNode = (nodes: FileNode[]): boolean => {
                        for (let i = 0; i < nodes.length; i++) {
                            if (nodes[i].path === path) {
                                nodes.splice(i, 1);
                                return true;
                            }
                            if (nodes[i].children && removeNode(nodes[i].children)) {
                                return true;
                            }
                        }
                        return false;
                    };
                    removeNode(newTree);
                    return { ...state, fileTree: newTree };
                });

                await DeleteFile(path);
                
                // Then update the parent directory to ensure proper sorting
                const updatedNode = await LoadDirectoryContents(parentPath);
                if (updatedNode) {
                    update(state => {
                        if (!state.fileTree) return state;
                        
                        // If it's the root directory, update the entire tree
                        if (parentPath === currentProjectPath) {
                            return { ...state, fileTree: updatedNode.children || [] };
                        }
                        
                        // Otherwise update the specific directory
                        const newTree = [...state.fileTree];
                        const updateDir = (nodes: FileNode[]): boolean => {
                            for (const node of nodes) {
                                if (node.path === parentPath) {
                                    node.children = updatedNode.children;
                                    return true;
                                }
                                if (node.children && updateDir(node.children)) {
                                    return true;
                                }
                            }
                            return false;
                        };
                        updateDir(newTree);
                        return { ...state, fileTree: newTree };
                    });
                }

                // Also remove from open files if it was open
                update(state => {
                    const newOpenFiles = new Map(state.openFiles);
                    newOpenFiles.delete(path);
                    
                    // If the deleted file was active, set a new active file
                    let newActiveFilePath = state.activeFilePath;
                    if (state.activeFilePath === path) {
                        const openFilePaths = Array.from(newOpenFiles.keys());
                        newActiveFilePath = openFilePaths.length > 0 ? openFilePaths[0] : null;
                    }
                    
                    return {
                        ...state,
                        openFiles: newOpenFiles,
                        activeFilePath: newActiveFilePath
                    };
                });
            } catch (error) {
                update(state => ({ ...state, error: `Failed to delete: ${error}` }));
                throw error;
            }
        },

        // Load directory contents
        async loadDirectoryContents(dirPath: string) {
            update(state => ({ ...state, loading: true, error: null }));
            try {
                const updatedNode = await LoadDirectoryContents(dirPath);
                if (!updatedNode) return;

                // Update the node in the file tree
                update(state => {
                    const updateNodeInTree = (nodes: FileNode[] | null): FileNode[] | null => {
                        if (!nodes) return null;
                        return nodes.map(node => {
                            if (node.path === dirPath) {
                                return { ...updatedNode, isLoaded: true };
                            }
                            if (node.type === "directory" && node.children) {
                                const updatedChildren = updateNodeInTree(node.children);
                                if (updatedChildren !== node.children) {
                                    return { ...node, children: updatedChildren };
                                }
                            }
                            return node;
                        });
                    };

                    const updatedTree = updateNodeInTree(state.fileTree);
                    return {
                        ...state,
                        fileTree: updatedTree,
                        loading: false
                    };
                });
            } catch (err) {
                update(state => ({
                    ...state,
                    error: err instanceof Error ? err.message : 'Failed to load directory contents',
                    loading: false
                }));
                throw err; // Re-throw to allow handling in the UI
            }
        },

        // Reset store
        reset() {
            localStorage.removeItem('fileState');
            set(initialState);
        }
    };
}

export const fileStore = createFileStore();
