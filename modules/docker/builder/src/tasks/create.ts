
import {ITask, TaskState, Builder} from '../builder';
import {Module, runHook, IPromise} from '../types';
import  * as _ from 'underscore';
const docker = require('docker');

export function createTask(mod: Module): ITask {
    return new CreateTask(mod);
}

var exclude = [ 'postrun', 'prerun', 'prebuild', 'postbuild',
            'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize']

export class CreateTask implements ITask {
    name = "Create";
    constructor(public mod: Module) {

    }

    run(builder: Builder, env: string): IPromise<TaskState> {
        let o = _.omit(this.mod, exclude);
        if (o.volume) {
            o.binds = o.volume
        }

        if (this.mod.image) {
            o.image = this.mod.image
        } else if (this.mod.build) {
            o.image = this.mod.name + "-image";
        } else {
            return Promise.reject(new Error("no image: " + o.name))
        }
        o.pull = true
        
        return docker.create(o).then( e => TaskState.Success)
        .catch(e => {
            throw new Error(`create ${this.mod.name}: ${e}`)
        })
        
    }

    _getCreateOptions(mod: Module, env: string) {
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
                    value = this._getCreateOptions(value, env);
                }

                if (_.isFunction(value)) continue;
                out[key] = value
            }
        }

        

        

        return out;
    }
}