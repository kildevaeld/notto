import {EventEmitter} from 'events';
import * as _ from 'underscore';
import {format} from 'util';
import {BuildTask, buildTask, createTask, CreateTask, startTask} from './tasks/index';
import {IPromise, eachAsyncMap, PromiseConstructor, eachAsync, ModuleMap, Module, runHook} from './types'
import {sanitize, deepObjectExtend} from './config'

declare var Promise: PromiseConstructor
const docker = require('docker')

export enum Notification {
    Starting, Started, StartError,
    Stopping, Stopped, AlreadyStarted,
    Building, Build, BuildError,
    Skipping, Removing, Removed,
    RemovingImage, RemovedImage,
    Creating, Created
};

export const NotificationEvent = "notification";


const n = Notification;


function series<T>(arr: T[], iterator: (a: T, cb: (err: Error) => void) => void, callback?): IPromise<any> {
    callback = callback || function () { };
    if (!Array.isArray(arr) || !arr.length) {
        return Promise.resolve();
    }
    var completed = 0;
    return new Promise<void>(function (resolve, reject) {
        var iterate = function () {
            iterator(arr[completed], function (err) {
                if (err) {
                    callback(err);
                    callback = function () { };
                    reject(err)
                }
                else {
                    ++completed;
                    if (completed >= arr.length) { resolve(); }
                    else { iterate(); }
                }
            });
        };
        iterate();
    })
};




export enum TaskState {
    Success, Error, Skipping
}

export interface ITask {
    mod: Module;
    name: string
    run(builder: Builder, env: string): IPromise<TaskState>
}


export function test(config, platform, env) {
    return sanitize(config, platform, env)
}

export class Builder extends EventEmitter {

    constructor(public modules: Module[], public env: string) {
        super();
    }

    build(force: boolean = false) {
        return this._build(force)
    }

    private _build(force) {
        var builds = this.modules.map(step => {
            if (!step.build) return null;
            return step;
        }).filter(function (step) { return step != null })

        if (builds.length == 0) {
            return Promise.resolve(TaskState.Skipping);
        }

        let tasks: ITask[];
        if (force) {
            tasks = builds.map(m => buildTask(m));
        } else {
            tasks = builds.map(m => {
                if (docker.hasImage(m.image, true)) {
                    return buildTask(m);
                }
                return null
            }).filter(m => m != null);

        }

        if (tasks.length == 0) return Promise.resolve(TaskState.Skipping);

        this.emit(NotificationEvent, n.Building, tasks.map(m => m.mod));
        return Promise.all(tasks.map(m => m.run(this, this.env)))
            .then(states => {
                this.emit(NotificationEvent, n.Build, builds.map)
                return TaskState.Success;
            });
    }



    start(autoBuild: boolean = false) {
        let ret = this.modules.map(m => docker.hasImage(m.image, true));
        let tbb = ret.filter(m => m === false);
        if (tbb.length > 0 && !autoBuild) {
            throw new Error('You have to build first');
        } else if (tbb.length > 0) {

            return this.build()
                .then(() => { this._start() })
        }
        return Promise.resolve<any>(this._start())


    }

    stop() {
        return new Promise((resolve, reject) => {
            this.modules.forEach(m => {
                if (!docker.isRunning(m.name, true)) return;
                this.emit(NotificationEvent, n.Stopping, m)
                runHook('prestop', m, true);
                docker.stop(m.name, true)
                this.emit(NotificationEvent, n.Stopped, m)
                runHook('poststop', m, true);
            });
        });

    }

    remove(force: boolean = false, images: boolean = false) {
        let c = this.modules.map(m => docker.hasContainer(m.name, true));

        c = c.map((m, i) => {
            if (m == false) return Promise.resolve();
            m = this.modules[i];
            this.emit(NotificationEvent, n.Removing, m)
            runHook('preremove', m, true);
            return docker.remove(m.name, force)
                .then(() => {
                    this.emit(NotificationEvent, n.Removed, m);
                    runHook('postremove', m, true);
                    if (images && m.build) {
                        return docker.removeImage(m.image)
                    }
                });
        })
        return Promise.all(c);

    }

    _create(mod: Module) {
        var exclude = ['postrun', 'prerun', 'prebuild', 'postbuild',
            'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize']
        let o = _.omit(mod, exclude);
        if (o.volume) {
            o.binds = o.volume
        }

        /*if (mod.build) {
            o.image = mod.image
        } else if (mod.build) {
            o.image = mod.name + "-image";
        } else {
            return Promise.reject(new Error("no image: " + o.name))
        }*/
        o.image = mod.image;
        o.pull = true

        return docker.create(o, true);

    }

    _start() {

        for (let i = 0, ii = this.modules.length; i < ii; i++) {
            let mod = this.modules[i];
            if (mod.phase) {
                if (!Array.isArray(mod.phase)) mod.phase = [<any>mod.phase];

                if (!!!~mod.phase.indexOf(this.env)) {
                    this.emit(NotificationEvent, n.Skipping, mod);
                    continue;
                }
            }

            var name = mod.name

            var hasContainer = docker.hasContainer(name, true),
                isRunning = docker.isRunning(name, true);

            if (isRunning) {
                this.emit(NotificationEvent, n.AlreadyStarted, mod);
                continue;
            }
            if (!hasContainer) {
                this.emit(NotificationEvent, n.Creating, mod);
                runHook('precreate', mod, true)
                this._create(mod);
                this.emit(NotificationEvent, n.Created, mod);
                runHook('postcreate', mod, true)
            }

            runHook('prestart', mod, true)
            this.emit(NotificationEvent, n.Starting, mod);
            docker.start(mod, true)
            if (mod.check) {
                let i = docker.inspect(mod.name, true);
                if (!i) throw new Error('container not started');
                let addr = i.NetworkSettings.IPAddress;
                
                if (process.platform == 'darwin') {
                    addr = module.__machine().trim();
                }
                if (!docker.check(addr + ":" + mod.check, 60)) {
                    throw new Error("Module " + mod.name + " timed out")
                }
            }
            runHook('poststart', mod, true)
            this.emit(NotificationEvent, n.Started, mod);
        }

    }

}

export function createBuilder(mod: Module | (() => Module | IPromise<Module>), env: string = "development") {
    if (_.isFunction(mod)) {
        mod = (<any>mod)();
    }

    var known_modules: { [key: string]: Module } = {};

    return Promise.resolve<Module>(<any>mod)
        .then(function (options) {

            if (typeof options.initialize === 'function') {
                return Promise.resolve(options.initialize(options))
                    .then(function () { return options })
            }
            return options;
        })
        .then(function (options) {

            options = sanitize(options, process.platform, env)

            parseModule(options, known_modules)

            for (let key in known_modules) {
                known_modules[key] = sanitize(known_modules[key], process.platform, env)
            }

            var out = [];
            if (options.dependencies != null) {
                resolveDependencies(options.dependencies, known_modules, out);
            }

            out.push(options);


            return new Builder(out, env);

        });
}



function getCreateOptions(builder: Builder, mod: Module, env: string) {
    var out: any = {}
    var exclude = ['name', 'postrun', 'prerun', 'prebuild', 'postbuild',
        'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize']
    if (mod.phase) {

        if (!Array.isArray(mod.phase)) mod.phase = [<string>mod.phase];
        if (!!!~mod.phase.indexOf(env)) {
            builder.emit('notification', n.Skipping, mod);
            return out;
        }
    }
    for (var key in mod) {

        var value = mod[key];
        if (key[0] == "$") {
            if (key.substr(1) === env) {
                out = _.extend(out, mod[key]);
            } else if (key.substr(1) === process.platform) {
                out = _.extend(out, mod[key]);
            }
        } else if (!!~exclude.indexOf(key)) {

            continue;
        } else {
            if (_.isObject(value) && !Array.isArray(value)) {
                value = getCreateOptions(builder, value, env);
            }

            if (_.isFunction(value)) continue;
            out[key] = value
        }
    }

    out.attachStdout = true
    out.attachStderr = true
    //stdinOnce: true,
    out.attachStdin = true;
    out.stdinOnce = true
    out.openStdin = true
    return out;
}

function parseModule(options: Module, known_modules: { [key: string]: Module }) {
    var name = options.name;

    if (!known_modules[name]) {

        known_modules[name] = options;
    } else /*if (Object.keys(known_modules[name]).length < Object.keys(options).length)*/ {

        known_modules[name] = deepObjectExtend(known_modules[name], options);
    }
    if (options.dependencies) {
        var deps = options.dependencies
        for (var i = 0, ii = deps.length; i < ii; i++) {
            var mod = deps[i]
            parseModule(mod, known_modules)
        }
    }
}


function resolveDependencies(dependencies: Module[], known_modules: { [key: string]: Module }, out: Module[]) {
    var first = true;
    for (var i = 0, ii = dependencies.length; i < ii; i++) {
        var mod = dependencies[i];
        mod = known_modules[mod.name];
        var deps = mod.dependencies;

        var found = _.find(out, function (v) {
            return v.name == mod.name;
        })

        if (found != null) continue;

        if (!deps) {
            var found = _.find(out, function (v) {
                return v.name == mod.name;
            })

            if (!found) out.push(known_modules[mod.name]);

            continue;
        }

        var sdep = _.find(deps, function (v) {
            var m = _.find(known_modules, function (vv) {
                return vv.name == v.name
            })
            if (m && m.dependencies) {
                return !!_.find(m.dependencies, function (vvv) {
                    return vvv.name == mod.name
                })
            }
            return false;
        })

        if (sdep) {
            throw new Error(format('circle dep: %s and %s depends on eachother', sdep.name, mod.name))
        }

        resolveDependencies(deps, known_modules, out)

        var found = _.find(out, function (v) {
            return v.name == mod.name;
        })

        if (!found) out.push(mod);

    }

}