// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {db} from '../models';
import {service} from '../models';

export function AddProject(arg1:string,arg2:string):Promise<db.Project>;

export function Commit(arg1:string,arg2:string):Promise<void>;

export function CreateDirectory(arg1:string):Promise<void>;

export function CreateFile(arg1:string):Promise<void>;

export function CreateTerminal(arg1:string,arg2:string,arg3:string):Promise<void>;

export function DeleteFile(arg1:string):Promise<void>;

export function DestroyTerminal(arg1:string):Promise<void>;

export function DiscardChanges(arg1:string,arg2:string):Promise<void>;

export function GetAvailableShells():Promise<Array<string>>;

export function GetCurrentBranch(arg1:string):Promise<string>;

export function GetEditorConfig():Promise<service.EditorConfig>;

export function GetFileContent(arg1:string):Promise<string>;

export function GetGitStatus(arg1:string):Promise<Array<service.FileStatus>>;

export function GetProjectFiles(arg1:string):Promise<service.FileNode>;

export function GetRecentProjects():Promise<Array<db.Project>>;

export function Greet(arg1:string):Promise<string>;

export function HandleInput(arg1:string,arg2:Array<number>):Promise<void>;

export function InitGitRepository(arg1:string):Promise<void>;

export function IsGitRepository(arg1:string):Promise<boolean>;

export function ListBranches(arg1:string):Promise<Array<service.BranchInfo>>;

export function LoadDirectoryContents(arg1:string):Promise<service.FileNode>;

export function OpenConfigFile():Promise<string>;

export function OpenProjectFolder():Promise<string>;

export function RenameFile(arg1:string,arg2:string):Promise<void>;

export function ResizeTerminal(arg1:string,arg2:number,arg3:number):Promise<void>;

export function SaveFile(arg1:string,arg2:string):Promise<void>;

export function SearchFiles(arg1:string,arg2:string):Promise<Array<service.FileNode>>;

export function StageFile(arg1:string,arg2:string):Promise<void>;

export function UnstageFile(arg1:string,arg2:string):Promise<void>;
