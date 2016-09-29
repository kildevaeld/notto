var _ = require('underscore');

var docker = module.__docker.create();

var wrap = function (method, options, sync) {
    if (sync) {
        return docker[method](options, true);
    }
    return new Promise(function (resolve, reject) {

        docker[method](options, function (err, result) {
            if (err) return reject(err);
            resolve(result);
        });
    });
}

exports.create = function (options, sync) {
    return wrap("Create", options, sync)
}

exports.build = function (options, sync) {
    return wrap("Build", options, sync)
}

exports.start = function (name, sync) {
    return wrap('Start', {
        name: name
    }, sync);
}

exports.exec = function (options, sync) {
    return wrap('Exec', options, sync);
}

exports.attach = function (options, sync) {
    return wrap('Attach', options, sync);
}

exports.stop = function (name, sync) {
    return wrap('Stop', {
        name: name
    }, sync);
}

exports.removeImage = function (nameOrId, sync) {
    return wrap("RemoveImage", {
        name: nameOrId
    }, sync)
}

exports.remove = function (nameOrID, force, sync) {
    return wrap("RemoveContainer", {
        name: nameOrID,
        force: force || false
    }, sync)
}

// Utilitiies
exports.hasImage = function (name, sync) {
    return wrap('HasImage', {
        name: name
    }, sync)
}

exports.hasContainer = function (nameOrId, sync) {
    return wrap('HasContainer', {
        name: nameOrId
    }, sync)
}

exports.isRunning = function (nameOrId, sync) {
    return wrap('IsRunning', {
        name: nameOrId
    }, sync)
}