
import {ITask, TaskState, Builder} from '../builder';
import {Module, runHook, IPromise} from '../types';
const docker = require('docker');


export function startTask(mod: Module): ITask {
    return new StartTask(mod);
}

export class StartTask implements ITask {
    name = "Create";
    constructor(public mod: Module) {

    }

    run(builder: Builder, env: string): IPromise<TaskState> {
        return docker.start(this.mod.name).then( e => TaskState.Success)
        .catch(e => {
            throw new Error(`start ${this.mod.name}: ${e}`)
        })
    }

    
}