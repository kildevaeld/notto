
import {Builder, TaskState, ITask} from '../builder';
import {IPromise, Module, runHook} from '../types'
import * as _ from 'underscore';
const docker = require('docker');

export function buildTask(mod:Module): ITask {
    return new BuildTask(mod);
}
 
export class BuildTask implements ITask {
    name = "Build";
    constructor(public mod: Module) { }
    run(builder: Builder, env: string): IPromise<TaskState>  {

        let options = this._getBuildOptions(this.mod, env);
        console.log('run this')
        return runHook('prebuild', this.mod)
        .then( e => {
            return docker.build(options);
        })
        .then( out => {
            return runHook('postbuild', this.mod);
        })
        .then((out) => {
            return TaskState.Success;
        });
    }


    _getBuildOptions(mod: Module, env: string) {
        var out = {}
        var exclude = ['name', 'postrun', 'prerun', 'prebuild', 'postbuild',
            'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize']
        /*if (mod.phase) {

            if (!Array.isArray(mod.phase)) mod.phase = [<string>mod.phase];
            if (!!!~mod.phase.indexOf(env)) {
                //builder.emit('notification', n.Skipping, mod);
                return null;
            }
        }*/
        
        for (var key in mod.build) {

            var value = mod.build[key];
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
                    value = this._getBuildOptions(value, env);
                }

                if (_.isFunction(value)) continue;
                out[key] = value
            }
        }

        out.name = this.mod.name + "-image";

        return out;
    }
}