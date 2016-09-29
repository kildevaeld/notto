
var fs = require('__private_fs')


exports.readdir = function(path) {
    return new Promise(function (resolve, reject) {
        fs.readdir(path, function (err, result) {
            if (err) return reject(err);
            resolve(result);
        })
    })
}

exports.readFile = function(path) {
    return new Promise(function (resolve, reject) {
        fs.readFile(path, function (err, result) {
            if (err) return reject(err);
            resolve(result);
        })
    })
}

exports.writeFile = function(path, data) {
    return new Promise(function (resolve, reject) {
        fs.writeFile(path, data, function (err, result) {
            if (err) return reject(err);
            resolve(result);
        })
    })
}

exports.mkdir = function(path) {
    return new Promise(function (resolve, reject) {
        fs.mkdir(path, function (err, result) {
            if (err) return reject(err);
            resolve(result);
        })
    })
}

exports.rename = function(source, target) {
    return new Promise(function (resolve, reject) {
        fs.rename(source, target, function (err, result) {
            if (err) return reject(err);
            resolve(result);
        })
    })
}

exports.stat = function (source) {
    return new Promise(function (resolve, reject) {
        fs.stat(source, function (err, result) {
            if (err) return reject(err)
            resolve(result);
        })
    })
}

exports.unlink = function (source) {
    return new Promise(function (resolve, reject) {
        fs.unlink(source, function (err, result) {
            if (err) return reject(err)
            resolve(result);
        })
    })
}