// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function CheckForUpdates():Promise<boolean>;

export function DeleteInstalledEnvironment(arg1:string,arg2:string,arg3:string,arg4:string):Promise<void>;

export function DoUpdate():Promise<void>;

export function GetAvailablePort():Promise<string>;

export function GetInstalledEnvironments():Promise<Array<main.Environment>>;

export function GetIp():Promise<string>;

export function GetKubernetesContexts():Promise<Array<string>>;

export function GetReleaseUrl():Promise<string>;

export function GetVersion():Promise<string>;

export function InstallEnvironment(arg1:string,arg2:main.EnvironmentSetup,arg3:Array<main.Section>,arg4:boolean,arg5:boolean):Promise<void>;

export function IsDockerInstalled():Promise<boolean>;

export function IsDockerRunning():Promise<boolean>;

export function IsEnvironmentInstalled(arg1:string,arg2:string,arg3:string,arg4:string):Promise<boolean>;

export function IsInternetConnected():Promise<boolean>;

export function IsKubernetesInstalled():Promise<boolean>;

export function IsPortAvailable(arg1:string):Promise<boolean>;

export function OpenFolderDialog(arg1:string):Promise<string>;

export function PopulateEnvironment(arg1:string,arg2:string,arg3:string,arg4:string):Promise<void>;

export function ReadEnvVariables(arg1:string):Promise<Array<main.Section>>;

export function SpecifyPlatformPath(arg1:string):Promise<string>;
