"use strict";
var _1 = require('../');
var docker = require('docker');
var BuildTask = (function () {
    function BuildTask(mod) {
        this.mod = mod;
    }
    BuildTask.prototype.run = function (builder, env) {
        var options = this._getBuildOptions(this.mod, env);
        console.log(JSON.stringify(options, null, 2));
        return _1.Promise.resolve(_1.TaskState.Success);
    };
    BuildTask.prototype._getBuildOptions = function (mod, env) {
        var out = {};
        var exclude = ['name', 'postrun', 'prerun', 'prebuild', 'postbuild',
            'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize'];
        if (mod.phase) {
            if (!Array.isArray(mod.phase))
                mod.phase = [mod.phase];
            if (!!!~mod.phase.indexOf(env)) {
                //builder.emit('notification', n.Skipping, mod);
                return null;
            }
        }
        for (var key in mod) {
            var value = mod[key];
            if (key[0] == "$") {
                if (key.substr(1) === env) {
                    out = _.extend(out, mod[key]);
                }
                else if (key.substr(1) === process.platform) {
                    out = _.extend(out, mod[key]);
                }
            }
            else if (!!~exclude.indexOf(key)) {
                continue;
            }
            else {
                if (_.isObject(value) && !Array.isArray(value)) {
                    value = this._getBuildOptions(value, env);
                }
                if (_.isFunction(value))
                    continue;
                out[key] = value;
            }
        }
        return out;
    };
    return BuildTask;
}());
exports.BuildTask = BuildTask;
