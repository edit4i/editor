export namespace db {
	
	export class Project {
	    ID: number;
	    Name: string;
	    Path: string;
	    LastOpened: sql.NullTime;
	    CreatedAt: sql.NullTime;
	    UpdatedAt: sql.NullTime;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Path = source["Path"];
	        this.LastOpened = this.convertValues(source["LastOpened"], sql.NullTime);
	        this.CreatedAt = this.convertValues(source["CreatedAt"], sql.NullTime);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], sql.NullTime);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace service {
	
	export class EditorConfig {
	    // Go type: struct { Theme string "json:\"theme\" mapstructure:\"theme\""; FontSize int "json:\"fontSize\" mapstructure:\"fontSize\""; TabSize int "json:\"tabSize\" mapstructure:\"tabSize\""; WordWrap bool "json:\"wordWrap\" mapstructure:\"wordWrap\""; LineNumbers bool "json:\"lineNumbers\" mapstructure:\"lineNumbers\""; RelativeLines bool "json:\"relativeLines\" mapstructure:\"relativeLines\""; Minimap bool "json:\"minimap\" mapstructure:\"minimap\""; StickyScroll bool "json:\"stickyScroll\" mapstructure:\"stickyScroll\""; Vim struct { Enabled bool "json:\"enabled\" mapstructure:\"enabled\""; DefaultMode string "json:\"defaultMode\" mapstructure:\"defaultMode\"" } "json:\"vim\" mapstructure:\"vim\"" }
	    editor: any;
	    // Go type: struct { DefaultShell string "json:\"defaultShell\" mapstructure:\"defaultShell\""; FontSize int "json:\"fontSize\" mapstructure:\"fontSize\""; FontFamily string "json:\"fontFamily\" mapstructure:\"fontFamily\""; Theme struct { Background string "json:\"background\" mapstructure:\"background\""; Foreground string "json:\"foreground\" mapstructure:\"foreground\""; Cursor string "json:\"cursor\" mapstructure:\"cursor\""; SelectionBackground string "json:\"selectionBackground\" mapstructure:\"selectionBackground\""; SelectionForeground string "json:\"selectionForeground\" mapstructure:\"selectionForeground\"" } "json:\"theme\" mapstructure:\"theme\"" }
	    terminal: any;
	    keyboard: struct { CustomBindings map[string]service.;
	
	    static createFrom(source: any = {}) {
	        return new EditorConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.editor = this.convertValues(source["editor"], Object);
	        this.terminal = this.convertValues(source["terminal"], Object);
	        this.keyboard = this.convertValues(source["keyboard"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileNode {
	    name: string;
	    path: string;
	    type: string;
	    size?: number;
	    // Go type: time
	    lastModified: any;
	    children?: FileNode[];
	    isLoaded: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FileNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.type = source["type"];
	        this.size = source["size"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.children = this.convertValues(source["children"], FileNode);
	        this.isLoaded = source["isLoaded"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class KeyBinding {
	    key: string;
	    modifiers: string[];
	
	    static createFrom(source: any = {}) {
	        return new KeyBinding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.modifiers = source["modifiers"];
	    }
	}

}

export namespace sql {
	
	export class NullTime {
	    // Go type: time
	    Time: any;
	    Valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NullTime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Time = this.convertValues(source["Time"], null);
	        this.Valid = source["Valid"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace struct { CustomBindings map[string]service {
	
	export class  {
	    customBindings: {[key: string]: service.KeyBinding};
	
	    static createFrom(source: any = {}) {
	        return new (source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.customBindings = this.convertValues(source["customBindings"], service.KeyBinding, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

