export namespace main {
	
	export class EposAccessPoints {
	    apiGateway: string;
	    dataPortal: string;
	
	    static createFrom(source: any = {}) {
	        return new EposAccessPoints(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.apiGateway = source["apiGateway"];
	        this.dataPortal = source["dataPortal"];
	    }
	}
	export class Section {
	    name: string;
	    variables: {[key: string]: string};
	
	    static createFrom(source: any = {}) {
	        return new Section(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.variables = source["variables"];
	    }
	}
	export class EnvironmentSetup {
	    name: string;
	    version: string;
	    context: string;
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentSetup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.context = source["context"];
	    }
	}
	export class Environment {
	    platform: string;
	    environmentSetup: EnvironmentSetup;
	    variables: Section[];
	    accessPoints: EposAccessPoints;
	
	    static createFrom(source: any = {}) {
	        return new Environment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.platform = source["platform"];
	        this.environmentSetup = this.convertValues(source["environmentSetup"], EnvironmentSetup);
	        this.variables = this.convertValues(source["variables"], Section);
	        this.accessPoints = this.convertValues(source["accessPoints"], EposAccessPoints);
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

