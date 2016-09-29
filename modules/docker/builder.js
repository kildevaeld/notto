module.exports =
/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};

/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {

/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId])
/******/ 			return installedModules[moduleId].exports;

/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			exports: {},
/******/ 			id: moduleId,
/******/ 			loaded: false
/******/ 		};

/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);

/******/ 		// Flag the module as loaded
/******/ 		module.loaded = true;

/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}


/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;

/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;

/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";

/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(0);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ function(module, exports, __webpack_require__) {

	"use strict";
	function __export(m) {
	    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
	}
	__export(__webpack_require__(1));


/***/ },
/* 1 */
/***/ function(module, exports, __webpack_require__) {

	"use strict";
	var __extends = (this && this.__extends) || function (d, b) {
	    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
	    function __() { this.constructor = d; }
	    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
	};
	var events_1 = __webpack_require__(2);
	var _ = __webpack_require__(3);
	var util_1 = __webpack_require__(4);
	var index_1 = __webpack_require__(5);
	var types_1 = __webpack_require__(7);
	var config_1 = __webpack_require__(11);
	var docker = __webpack_require__(8);
	(function (Notification) {
	    Notification[Notification["Starting"] = 0] = "Starting";
	    Notification[Notification["Started"] = 1] = "Started";
	    Notification[Notification["StartError"] = 2] = "StartError";
	    Notification[Notification["Stopping"] = 3] = "Stopping";
	    Notification[Notification["Stopped"] = 4] = "Stopped";
	    Notification[Notification["AlreadyStarted"] = 5] = "AlreadyStarted";
	    Notification[Notification["Building"] = 6] = "Building";
	    Notification[Notification["Build"] = 7] = "Build";
	    Notification[Notification["BuildError"] = 8] = "BuildError";
	    Notification[Notification["Skipping"] = 9] = "Skipping";
	    Notification[Notification["Removing"] = 10] = "Removing";
	    Notification[Notification["Removed"] = 11] = "Removed";
	    Notification[Notification["RemovingImage"] = 12] = "RemovingImage";
	    Notification[Notification["RemovedImage"] = 13] = "RemovedImage";
	    Notification[Notification["Creating"] = 14] = "Creating";
	    Notification[Notification["Created"] = 15] = "Created";
	})(exports.Notification || (exports.Notification = {}));
	var Notification = exports.Notification;
	;
	exports.NotificationEvent = "notification";
	var n = Notification;
	function series(arr, iterator, callback) {
	    callback = callback || function () { };
	    if (!Array.isArray(arr) || !arr.length) {
	        return Promise.resolve();
	    }
	    var completed = 0;
	    return new Promise(function (resolve, reject) {
	        var iterate = function () {
	            iterator(arr[completed], function (err) {
	                if (err) {
	                    callback(err);
	                    callback = function () { };
	                    reject(err);
	                }
	                else {
	                    ++completed;
	                    if (completed >= arr.length) {
	                        resolve();
	                    }
	                    else {
	                        iterate();
	                    }
	                }
	            });
	        };
	        iterate();
	    });
	}
	;
	(function (TaskState) {
	    TaskState[TaskState["Success"] = 0] = "Success";
	    TaskState[TaskState["Error"] = 1] = "Error";
	    TaskState[TaskState["Skipping"] = 2] = "Skipping";
	})(exports.TaskState || (exports.TaskState = {}));
	var TaskState = exports.TaskState;
	function test(config, platform, env) {
	    return config_1.sanitize(config, platform, env);
	}
	exports.test = test;
	var Builder = (function (_super) {
	    __extends(Builder, _super);
	    function Builder(modules, env) {
	        _super.call(this);
	        this.modules = modules;
	        this.env = env;
	    }
	    Builder.prototype.build = function (force) {
	        if (force === void 0) { force = false; }
	        return this._build(force);
	    };
	    Builder.prototype._build = function (force) {
	        var _this = this;
	        var builds = this.modules.map(function (step) {
	            if (!step.build)
	                return null;
	            return step;
	        }).filter(function (step) { return step != null; });
	        if (builds.length == 0) {
	            return Promise.resolve(TaskState.Skipping);
	        }
	        var tasks;
	        if (force) {
	            tasks = builds.map(function (m) { return index_1.buildTask(m); });
	        }
	        else {
	            tasks = builds.map(function (m) {
	                if (docker.hasImage(m.image, true)) {
	                    return index_1.buildTask(m);
	                }
	                return null;
	            }).filter(function (m) { return m != null; });
	        }
	        if (tasks.length == 0)
	            return Promise.resolve(TaskState.Skipping);
	        this.emit(exports.NotificationEvent, n.Building, tasks.map(function (m) { return m.mod; }));
	        return Promise.all(tasks.map(function (m) { return m.run(_this, _this.env); }))
	            .then(function (states) {
	            _this.emit(exports.NotificationEvent, n.Build, builds.map);
	            return TaskState.Success;
	        });
	    };
	    Builder.prototype.start = function (autoBuild) {
	        var _this = this;
	        if (autoBuild === void 0) { autoBuild = false; }
	        var ret = this.modules.map(function (m) { return docker.hasImage(m.image, true); });
	        var tbb = ret.filter(function (m) { return m === false; });
	        if (tbb.length > 0 && !autoBuild) {
	            throw new Error('You have to build first');
	        }
	        else if (tbb.length > 0) {
	            return this.build()
	                .then(function () { _this._start(); });
	        }
	        return Promise.resolve(this._start());
	    };
	    Builder.prototype.stop = function () {
	        var _this = this;
	        return new Promise(function (resolve, reject) {
	            _this.modules.forEach(function (m) {
	                if (!docker.isRunning(m.name, true))
	                    return;
	                _this.emit(exports.NotificationEvent, n.Stopping, m);
	                types_1.runHook('prestop', m, true);
	                docker.stop(m.name, true);
	                _this.emit(exports.NotificationEvent, n.Stopped, m);
	                types_1.runHook('poststop', m, true);
	            });
	        });
	    };
	    Builder.prototype.remove = function (force, images) {
	        var _this = this;
	        if (force === void 0) { force = false; }
	        if (images === void 0) { images = false; }
	        var c = this.modules.map(function (m) { return docker.hasContainer(m.name, true); });
	        c = c.map(function (m, i) {
	            if (m == false)
	                return Promise.resolve();
	            m = _this.modules[i];
	            _this.emit(exports.NotificationEvent, n.Removing, m);
	            types_1.runHook('preremove', m, true);
	            return docker.remove(m.name, force)
	                .then(function () {
	                _this.emit(exports.NotificationEvent, n.Removed, m);
	                types_1.runHook('postremove', m, true);
	                if (images && m.build) {
	                    return docker.removeImage(m.image);
	                }
	            });
	        });
	        return Promise.all(c);
	    };
	    Builder.prototype._create = function (mod) {
	        var exclude = ['postrun', 'prerun', 'prebuild', 'postbuild',
	            'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize'];
	        var o = _.omit(mod, exclude);
	        if (o.volume) {
	            o.binds = o.volume;
	        }
	        /*if (mod.build) {
	            o.image = mod.image
	        } else if (mod.build) {
	            o.image = mod.name + "-image";
	        } else {
	            return Promise.reject(new Error("no image: " + o.name))
	        }*/
	        o.image = mod.image;
	        o.pull = true;
	        return docker.create(o, true);
	    };
	    Builder.prototype._start = function () {
	        for (var i = 0, ii = this.modules.length; i < ii; i++) {
	            var mod = this.modules[i];
	            if (mod.phase) {
	                if (!Array.isArray(mod.phase))
	                    mod.phase = [mod.phase];
	                if (!!!~mod.phase.indexOf(this.env)) {
	                    this.emit(exports.NotificationEvent, n.Skipping, mod);
	                    continue;
	                }
	            }
	            var name = mod.name;
	            var hasContainer = docker.hasContainer(name, true), isRunning = docker.isRunning(name, true);
	            if (isRunning) {
	                this.emit(exports.NotificationEvent, n.AlreadyStarted, mod);
	                continue;
	            }
	            if (!hasContainer) {
	                this.emit(exports.NotificationEvent, n.Creating, mod);
	                types_1.runHook('precreate', mod, true);
	                this._create(mod);
	                this.emit(exports.NotificationEvent, n.Created, mod);
	                types_1.runHook('postcreate', mod, true);
	            }
	            types_1.runHook('prestart', mod, true);
	            this.emit(exports.NotificationEvent, n.Starting, mod);
	            docker.start(mod, true);
	            if (mod.check) {
	                var i_1 = docker.inspect(mod.name, true);
	                if (!i_1)
	                    throw new Error('container not started');
	                var addr = i_1.NetworkSettings.IPAddress;
	                if (process.platform == 'darwin') {
	                    addr = dockermachine().trim();
	                }
	                if (!docker.check(addr + ":" + mod.check, 60)) {
	                    throw new Error("Module " + mod.name + " timed out");
	                }
	            }
	            types_1.runHook('poststart', mod, true);
	            this.emit(exports.NotificationEvent, n.Started, mod);
	        }
	    };
	    return Builder;
	}(events_1.EventEmitter));
	exports.Builder = Builder;
	function createBuilder(mod, env) {
	    if (env === void 0) { env = "development"; }
	    if (_.isFunction(mod)) {
	        mod = mod();
	    }
	    var known_modules = {};
	    return Promise.resolve(mod)
	        .then(function (options) {
	        if (typeof options.initialize === 'function') {
	            return Promise.resolve(options.initialize(options))
	                .then(function () { return options; });
	        }
	        return options;
	    })
	        .then(function (options) {
	        options = config_1.sanitize(options, process.platform, env);
	        parseModule(options, known_modules);
	        for (var key in known_modules) {
	            known_modules[key] = config_1.sanitize(known_modules[key], process.platform, env);
	        }
	        var out = [];
	        if (options.dependencies != null) {
	            resolveDependencies(options.dependencies, known_modules, out);
	        }
	        out.push(options);
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
	    out.attachStdout = true;
	    out.attachStderr = true;
	    //stdinOnce: true,
	    out.attachStdin = true;
	    out.stdinOnce = true;
	    out.openStdin = true;
	    return out;
	}
	function parseModule(options, known_modules) {
	    var name = options.name;
	    if (!known_modules[name]) {
	        known_modules[name] = options;
	    }
	    else {
	        known_modules[name] = config_1.deepObjectExtend(known_modules[name], options);
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


/***/ },
/* 2 */
/***/ function(module, exports) {

	module.exports = require("events");

/***/ },
/* 3 */
/***/ function(module, exports) {

	module.exports = require("underscore");

/***/ },
/* 4 */
/***/ function(module, exports) {

	module.exports = require("util");

/***/ },
/* 5 */
/***/ function(module, exports, __webpack_require__) {

	"use strict";
	function __export(m) {
	    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
	}
	__export(__webpack_require__(6));
	__export(__webpack_require__(9));
	__export(__webpack_require__(10));


/***/ },
/* 6 */
/***/ function(module, exports, __webpack_require__) {

	"use strict";
	var builder_1 = __webpack_require__(1);
	var types_1 = __webpack_require__(7);
	var _ = __webpack_require__(3);
	var docker = __webpack_require__(8);
	function buildTask(mod) {
	    return new BuildTask(mod);
	}
	exports.buildTask = buildTask;
	var BuildTask = (function () {
	    function BuildTask(mod) {
	        this.mod = mod;
	        this.name = "Build";
	    }
	    BuildTask.prototype.run = function (builder, env) {
	        var _this = this;
	        var options = this._getBuildOptions(this.mod, env);
	        options.name = this.mod.image;
	        return types_1.runHook('prebuild', this.mod)
	            .then(function (e) {
	            return docker.build(options);
	        })
	            .then(function (out) {
	            return types_1.runHook('postbuild', _this.mod);
	        })
	            .then(function (out) {
	            return builder_1.TaskState.Success;
	        });
	    };
	    BuildTask.prototype._getBuildOptions = function (mod, env) {
	        var out = {};
	        var exclude = ['name', 'postrun', 'prerun', 'prebuild', 'postbuild',
	            'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize'];
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
	        //out.name = this.mod.name + "-image";
	        return out;
	    };
	    return BuildTask;
	}());
	exports.BuildTask = BuildTask;


/***/ },
/* 7 */
/***/ function(module, exports) {

	"use strict";
	function runHook(hook, mod, sync) {
	    if (sync === void 0) { sync = false; }
	    if (sync) {
	        return typeof mod[hook] === 'function' ? mod[hook].call(mod, mod) : void 0;
	    }
	    if (typeof mod[hook] === 'function')
	        return Promise.resolve(mod[hook].call(mod, mod));
	    return Promise.resolve();
	}
	exports.runHook = runHook;
	function eachAsync(list, cb) {
	    var i = 0, l = list.length, result = [];
	    return next().then(function () {
	        return result;
	    });
	    function next() {
	        if (i == l)
	            return Promise.resolve(result);
	        return Promise.resolve(cb(list[i++])).then(function (ret) {
	            result.push(ret);
	            return next();
	        });
	    }
	    ;
	}
	exports.eachAsync = eachAsync;
	function eachAsyncMap(list, cb) {
	    var i = 0, l = list.length, result = [];
	    return next().then(function () {
	        return result;
	    });
	    function next() {
	        if (i == l)
	            return Promise.resolve(result);
	        var ret;
	        try {
	            ret = cb(list[i++]);
	        }
	        catch (e) {
	            return Promise.reject(e);
	        }
	        return Promise.resolve(ret).then(function (ret) {
	            if (ret == null)
	                return next();
	            result.push(ret);
	            return next();
	        });
	    }
	    ;
	}
	exports.eachAsyncMap = eachAsyncMap;


/***/ },
/* 8 */
/***/ function(module, exports) {

	module.exports = require("docker");

/***/ },
/* 9 */
/***/ function(module, exports, __webpack_require__) {

	"use strict";
	var builder_1 = __webpack_require__(1);
	var _ = __webpack_require__(3);
	var docker = __webpack_require__(8);
	function createTask(mod) {
	    return new CreateTask(mod);
	}
	exports.createTask = createTask;
	var exclude = ['postrun', 'prerun', 'prebuild', 'postbuild',
	    'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize'];
	var CreateTask = (function () {
	    function CreateTask(mod) {
	        this.mod = mod;
	        this.name = "Create";
	    }
	    CreateTask.prototype.run = function (builder, env) {
	        var _this = this;
	        var o = _.omit(this.mod, exclude);
	        if (o.volume) {
	            o.binds = o.volume;
	        }
	        if (this.mod.image) {
	            o.image = this.mod.image;
	        }
	        else if (this.mod.build) {
	            o.image = this.mod.name + "-image";
	        }
	        else {
	            return Promise.reject(new Error("no image: " + o.name));
	        }
	        o.pull = true;
	        return docker.create(o).then(function (e) { return builder_1.TaskState.Success; })
	            .catch(function (e) {
	            throw new Error("create " + _this.mod.name + ": " + e);
	        });
	    };
	    CreateTask.prototype._getCreateOptions = function (mod, env) {
	        var out = {};
	        var exclude = ['name', 'postrun', 'prerun', 'prebuild', 'postbuild',
	            'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize'];
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
	                    value = this._getCreateOptions(value, env);
	                }
	                if (_.isFunction(value))
	                    continue;
	                out[key] = value;
	            }
	        }
	        return out;
	    };
	    return CreateTask;
	}());
	exports.CreateTask = CreateTask;


/***/ },
/* 10 */
/***/ function(module, exports, __webpack_require__) {

	"use strict";
	var builder_1 = __webpack_require__(1);
	var docker = __webpack_require__(8);
	function startTask(mod) {
	    return new StartTask(mod);
	}
	exports.startTask = startTask;
	var StartTask = (function () {
	    function StartTask(mod) {
	        this.mod = mod;
	        this.name = "Create";
	    }
	    StartTask.prototype.run = function (builder, env) {
	        var _this = this;
	        return docker.start(this.mod.name).then(function (e) { return builder_1.TaskState.Success; })
	            .then(function (e) {
	            if (!_this.mod.check) {
	                return e;
	            }
	            var i = docker.inspect(_this.mod.name, true);
	            if (!i)
	                throw new Error('container not started');
	            if (!docker.check(i.NetworkSettings.IPAddress + ":" + _this.mod.check, 60)) {
	                throw new Error("Module " + _this.mod.name + " timed out");
	            }
	            return e;
	        })
	            .catch(function (e) {
	            throw new Error("start " + _this.mod.name + ": " + e);
	        });
	    };
	    return StartTask;
	}());
	exports.StartTask = StartTask;


/***/ },
/* 11 */
/***/ function(module, exports, __webpack_require__) {

	"use strict";
	var _ = __webpack_require__(3);
	function sanitize(mod, platform, env) {
	    if (mod.phase && mod.phase.indexOf(env) == -1) {
	        //console.log(mod.name, mod.build, env, mod.phase)
	        return null;
	    }
	    var out = interpolate(mod, platform, env);
	    if (mod.build && !mod.image) {
	        out.image = out.name + "-image";
	    }
	    if (!mod.dependencies) {
	        return out;
	    }
	    out.dependencies = [];
	    for (var i = 0, ii = mod.dependencies.length; i < ii; i++) {
	        var dep = sanitize(mod.dependencies[i], platform, env);
	        if (dep == null)
	            continue;
	        out.dependencies.push(dep);
	    }
	    //console.log(out.name, out.build, env)
	    return out;
	}
	exports.sanitize = sanitize;
	function deepObjectExtend(target, source) {
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
	exports.deepObjectExtend = deepObjectExtend;
	function interpolate(obj, platform, env) {
	    var out = {};
	    for (var key in obj) {
	        if (key[0] === '$') {
	            continue;
	        }
	        else {
	            if (_.isObject(obj[key]) && !Array.isArray(obj[key]) && !_.isFunction(obj[key])) {
	                out[key] = interpolate(obj[key], platform, env);
	            }
	            else {
	                out[key] = obj[key];
	            }
	        }
	    }
	    for (var key in obj) {
	        if (key[0] !== '$') {
	            continue;
	        }
	        var s = key.substr(1);
	        if (s !== platform && s !== env) {
	            continue;
	        }
	        var o = obj[key];
	        deepObjectExtend(out, o);
	    }
	    return out;
	}


/***/ }
/******/ ]);