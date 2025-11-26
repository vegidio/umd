export namespace fetch {
	
	export class Request {
	    Url: string;
	    FilePath: string;
	
	    static createFrom(source: any = {}) {
	        return new Request(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Url = source["Url"];
	        this.FilePath = source["FilePath"];
	    }
	}
	export class Response {
	    Request?: Request;
	    StatusCode: number;
	    Size: number;
	    Downloaded: number;
	    Progress: number;
	    Hash: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Request = this.convertValues(source["Request"], Request);
	        this.StatusCode = source["StatusCode"];
	        this.Size = source["Size"];
	        this.Downloaded = source["Downloaded"];
	        this.Progress = source["Progress"];
	        this.Hash = source["Hash"];
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

export namespace shared {
	
	export class Download {
	    Url: string;
	    FilePath: string;
	    Error: any;
	    IsSuccess: boolean;
	    Hash: string;
	
	    static createFrom(source: any = {}) {
	        return new Download(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Url = source["Url"];
	        this.FilePath = source["FilePath"];
	        this.Error = source["Error"];
	        this.IsSuccess = source["IsSuccess"];
	        this.Hash = source["Hash"];
	    }
	}

}

export namespace types {
	
	export class Media {
	    Url: string;
	    Extension: string;
	    Type: number;
	    Extractor: number;
	    Metadata: Record<string, any>;
	    Headers: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new Media(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Url = source["Url"];
	        this.Extension = source["Extension"];
	        this.Type = source["Type"];
	        this.Extractor = source["Extractor"];
	        this.Metadata = source["Metadata"];
	        this.Headers = source["Headers"];
	    }
	}

}

