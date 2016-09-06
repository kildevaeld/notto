
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

exports.rename = function(source, target) {
    return new Promise(function (resolve, reject) {
        fs.rename(source, target, function (err, result) {
            if (err) return reject(err);
            resolve(result);
        })
    })
}
