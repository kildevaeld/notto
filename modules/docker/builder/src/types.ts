declare var Promise: PromiseConstructor


export interface IPromise<T> extends Thenable<T> { }

export interface Thenable<R> {
    then<U>(onFulfilled?: (value: R) => U | Thenable<U>, onRejected?: (error: any) => U | Thenable<U>): Thenable<U>;
    then<U>(onFulfilled?: (value: R) => U | Thenable<U>, onRejected?: (error: any) => void): Thenable<U>;
    catch<U>(onRejected?: (error: any) => U | Thenable<U>): IPromise<U>;
}

export interface PromiseConstructor {
    new <R>(callback: (resolve: (value?: R | Thenable<R>) => void, reject: (error?: any) => void) => void): IPromise<R>
    resolve<R>(value?: R | Thenable<R>): IPromise<R>;
    reject(error: any): IPromise<any>;
    all<R>(promises: (R | Thenable<R>)[]): IPromise<R[]>;
    race<R>(promises: (R | Thenable<R>)[]): IPromise<R>;
}

export type Env = { [key: string]: string };
export type Link = { [key: string]: string };
export interface Prop<T> {
    $development?: T;
    $staging?: T;
    $production?: T;
    $testing?: T;
    $darwin?: T;
    $linux?: T;
    $windows?: T
    [key: string]: any;
}


export interface Build {
    path: string;
}

export interface Module extends Prop<Module> {
    name?: string;
    image?: string;
    link?: Prop<Link>;
    env?: Prop<Env>;
    phase: string[] | string;
    initialize?: (Module) => void;
    prestart?: () => void;
    poststart?: () => void;
    prebuild?: () => void;
    build?: Build
    postbuild?: () => void;
    publish: string[];
    volume: string[];
    dependencies?: Module[]
}

export type ModuleMap = { [key: string]: Module }


export function runHook(hook: string, mod: Module, sync:boolean = false) {
    if (sync) {
        return typeof mod[hook] === 'function' ? mod[hook].call(mod, mod) : void 0;
    }
    if (typeof mod[hook] === 'function') return Promise.resolve(mod[hook].call(mod, mod));
    return Promise.resolve();
}

export function eachAsync(list, cb): IPromise<void> {
    var i = 0, l = list.length, result = [];

    return next().then(function () {
        return result;
    });

    function next() {
        if (i == l) return Promise.resolve(result);
        return Promise.resolve(cb(list[i++])).then(function (ret) {
            result.push(ret);
            return next()
        });
    };
}

export function eachAsyncMap<T, U>(list: T[], cb: (a: T) => (IPromise<U> | U)): IPromise<U[]> {
    var i = 0, l = list.length, result = [];

    return next().then(function () {
        return result;
    });

    function next() {
        if (i == l) return Promise.resolve(result);
        let ret
        try {
            ret = cb(list[i++]);
        } catch (e) {
            return Promise.reject(e);
        }
        return Promise.resolve(ret).then(function (ret) {
            if (ret == null) return next();
            result.push(ret);
            return next()
        });

    };
}
