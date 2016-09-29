"use strict";
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
var events_1 = require('events');
var _ = require('underscore');
var util_1 = require('util');
var eachAsync = function (list, cb) {
    var i = 0, l = list.length, result = [];
    return next().then(function () {
        return result;
    });
    function next() {
        if (i == l)
            return exports.Promise.resolve(result);
        return exports.Promise.resolve(cb(list[i++])).then(function (ret) {
            result.push(ret);
            return next();
        });
    }
    ;
};
function eachAsyncMap(list, cb) {
    var i = 0, l = list.length, result = [];
    return next().then(function () {
        return result;
    });
    function next() {
        if (i == l)
            return exports.Promise.resolve(result);
        return exports.Promise.resolve(cb(list[i++])).then(function (ret) {
            if (ret == null)
                return next();
            result.push(ret);
            return next();
        });
    }
    ;
}
var docker = require('docker');
(function (Notifications) {
    Notifications[Notifications["Starting"] = 0] = "Starting";
    Notifications[Notifications["Started"] = 1] = "Started";
    Notifications[Notifications["StartError"] = 2] = "StartError";
    Notifications[Notifications["Stopping"] = 3] = "Stopping";
    Notifications[Notifications["Stopped"] = 4] = "Stopped";
    Notifications[Notifications["AlreadyStarted"] = 5] = "AlreadyStarted";
    Notifications[Notifications["Building"] = 6] = "Building";
    Notifications[Notifications["Build"] = 7] = "Build";
    Notifications[Notifications["BuildError"] = 8] = "BuildError";
    Notifications[Notifications["Skipping"] = 9] = "Skipping";
    Notifications[Notifications["Removing"] = 10] = "Removing";
    Notifications[Notifications["Removed"] = 11] = "Removed";
    Notifications[Notifications["RemovingImage"] = 12] = "RemovingImage";
    Notifications[Notifications["RemovedImage"] = 13] = "RemovedImage";
})(exports.Notifications || (exports.Notifications = {}));
var Notifications = exports.Notifications;
;
var NotificationEvent = "notification";
var n = Notifications;
function runHook(hook, mod) {
    if (typeof mod[hook] === 'function')
        return exports.Promise.resolve(mod[hook].call(mod, mod));
    return exports.Promise.resolve();
}
(function (TaskState) {
    TaskState[TaskState["Success"] = 0] = "Success";
    TaskState[TaskState["Error"] = 1] = "Error";
    TaskState[TaskState["Skipping"] = 2] = "Skipping";
})(exports.TaskState || (exports.TaskState = {}));
var TaskState = exports.TaskState;
var Builder = (function (_super) {
    __extends(Builder, _super);
    function Builder(modules, env) {
        _super.call(this);
        this.modules = modules;
        this.env = env;
    }
    Builder.prototype.build = function (force) {
        var _this = this;
        if (force === void 0) { force = false; }
        var builds = this.modules.map(function (step) {
            if (!step.build && !step.image)
                return null;
            return step; //_.pick(step, ['build', 'prebuild', 'postbuild'])
        }).filter(function (step) { return step != null; });
        /*.map(function (mod) {
            

            return this._buildModule(self, mod);
        });*/
        return eachAsyncMap(builds, function (m) {
            return docker.hasImage(m.name).then(function (ok) {
                if (ok) {
                    _this.emit(NotificationEvent, n.Skipping, m);
                    return null;
                }
                return m;
            });
        }).then(function (builds) {
            console.log(JSON.stringify(builds, null, 2));
        });
        return exports.Promise.all(builds);
    };
    Builder.prototype.start = function (autoBuild) {
        var _this = this;
        if (autoBuild === void 0) { autoBuild = false; }
        var self = this;
        return eachAsync(this.modules, function (mod) {
            if (mod.phase) {
                if (!Array.isArray(mod.phase))
                    mod.phase = [mod.phase];
                if (!!!~mod.phase.indexOf(self.env)) {
                    _this.emit('notification', n.Skipping, mod);
                    return;
                }
            }
            var name = mod.name;
            var promises = [
                docker.hasContainer(name),
                docker.isRunning(name),
                docker.hasImage(name)
            ];
            return exports.Promise.all(promises)
                .then(function (_a) {
                var hasContainer = _a[0], isRunning = _a[1], hasImage = _a[2];
                _this.emit('notification', n.Starting, mod);
                if (isRunning) {
                    _this.emit("notification", n.AlreadyStarted, mod);
                    ;
                    return false;
                }
                else if (hasContainer) {
                    return runHook('prestart', mod)
                        .then(function () {
                        return docker.start(mod.name);
                    })
                        .then(function () {
                        return true;
                    });
                }
                else if (!hasImage && autoBuild == true) {
                    return self._buildModule.call(self, mod)
                        .then(function () {
                        return runHook('prestart', mod);
                    }).then(function () {
                        var options = getCreateOptions(this, mod, self.env);
                        return docker.create(options);
                    }); /*.then(function () {
                        return true;
                    })*/
                }
                else {
                    var options = getCreateOptions(mod, self.env);
                    return runHook('prestart', mod)
                        .then(function () {
                        return docker.create(name, mod.image || name, options)
                            .then(function () {
                            return true;
                        });
                    });
                }
            }).then(function (started) {
                if (started) {
                    self.trigger('notification', n.started, mod);
                    return runHook('poststart', mod);
                }
            });
        }).then(function (results) {
            return results;
        });
    };
    Builder.prototype._buildModule = function (mod) {
        this.emit('notification', n.Building, mod);
        var build = mod.build;
        if (build == undefined)
            return exports.Promise.resolve();
        if (build.path == null && mod.image == null) {
            this.emit('notification', n.Build, mod);
            return exports.Promise.resolve();
        }
        return runHook('prebuild', mod).then(function () {
            return docker.build({
                path: build.path,
                name: mod.name + "-" + this.env + "-image"
            });
        })
            .then(function (out) {
            return runHook('postbuild', mod)
                .then(function () {
                this.emit('notification', n.Build, mod);
                return out;
            });
        });
    };
    return Builder;
}(events_1.EventEmitter));
exports.Builder = Builder;
function createBuilder(mod, env) {
    if (_.isFunction(mod)) {
        mod = mod();
    }
    var known_modules = {};
    return exports.Promise.resolve(mod)
        .then(function (options) {
        if (typeof options.initialize === 'function') {
            return exports.Promise.resolve(options.initialize(options))
                .then(function () { return options; });
        }
        return options;
    })
        .then(function (options) {
        parseModule(options, known_modules);
        var out = [];
        if (options.dependencies != null) {
            resolveDependencies(options.dependencies, known_modules, out);
        }
        out.push(options);
        var builds = out.map(function (step) {
            if (!step.build)
                return null;
            return _.pick(step, ['build', 'prebuild', 'postbuild']);
        }).filter(function (step) { return step != null; });
        return new Builder(out, env);
    });
}
exports.createBuilder = createBuilder;
function getCreateOptions(builder, mod, env) {
    var out = {};
    var exclude = ['name', 'postrun', 'prerun', 'prebuild', 'postbuild',
        'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize'];
    if (mod.phase) {
        if (!Array.isArray(mod.phase))
            mod.phase = [mod.phase];
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
                value = getCreateOptions(builder, value, env);
            }
            if (_.isFunction(value))
                continue;
            out[key] = value;
        }
    }
    return out;
}
function parseModule(options, known_modules) {
    var name = options.name;
    if (!known_modules[name]) {
        known_modules[name] = options;
    }
    else if (Object.keys(known_modules[name]).length < Object.keys(options).length) {
        known_modules[name] = options;
    }
    if (options.dependencies) {
        var deps = options.dependencies;
        for (var i = 0, ii = deps.length; i < ii; i++) {
            var mod = deps[i];
            parseModule(mod, known_modules);
        }
    }
}
function resolveDependencies(dependencies, known_modules, out) {
    var first = true;
    for (var i = 0, ii = dependencies.length; i < ii; i++) {
        var mod = dependencies[i];
        mod = known_modules[mod.name];
        var deps = mod.dependencies;
        var found = _.find(out, function (v) {
            return v.name == mod.name;
        });
        if (found != null)
            continue;
        if (!deps) {
            var found = _.find(out, function (v) {
                return v.name == mod.name;
            });
            if (!found)
                out.push(known_modules[mod.name]);
            continue;
        }
        var sdep = _.find(deps, function (v) {
            var m = _.find(known_modules, function (vv) {
                return vv.name == v.name;
            });
            if (m && m.dependencies) {
                return !!_.find(m.dependencies, function (vvv) {
                    return vvv.name == mod.name;
                });
            }
            return false;
        });
        if (sdep) {
            throw new Error(util_1.format('circle dep: %s and %s depends on eachother', sdep.name, mod.name));
        }
        resolveDependencies(deps, known_modules, out);
        var found = _.find(out, function (v) {
            return v.name == mod.name;
        });
        if (!found)
            out.push(mod);
    }
}
