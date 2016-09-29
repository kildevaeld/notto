
import {Module} from './types';
import * as _ from 'underscore';

export function sanitize(mod: Module, platform: string, env: string): Module {
    if (mod.phase && mod.phase.indexOf(env) == -1) {
        return null;
    }
    let out: Module = interpolate(mod, platform, env);

    if (!mod.dependencies) {
        return out;
    }
    out.dependencies = [];
    for (let i = 0, ii = mod.dependencies.length; i < ii; i++) {
        
        let dep = sanitize(mod.dependencies[i],platform, env);
        if (dep == null) continue;
        out.dependencies.push(dep);
    }

    return out;
}

function deepObjectExtend (target, source) {
    for (var prop in source) {
        if (source.hasOwnProperty(prop)) {
            if (target[prop] && typeof source[prop] === 'object') {
                deepObjectExtend(target[prop], source[prop]);
            }
            else {
                target[prop] = source[prop];
            }
        }
    }

    return target;
}

function interpolate(obj, platform, env) {
    var out: any = {};
    for (let key in obj) {
        if (key[0] === '$') {
            continue
        } else {
            if (_.isObject(obj[key]) && !Array.isArray(obj[key])) {
                out[key] = interpolate(obj[key], platform, env)
            } else {
                out[key] = obj[key];
            }
        }
    }
    for (let key in obj) {
        if (key[0] !== '$') {
            continue;
        }
        let s = key.substr(1)
        
        if (s !== platform && s !== env) {
            continue;
        }
        let o = obj[key];
        
        deepObjectExtend(out, o);
    }

    return out;
}