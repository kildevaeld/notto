

function toPromise() {
    var args = Array.prototype.slice.call(arguments);
    return new Promise(function (resolve, reject) {
        var m = args.shift()
        args.push(function (err, result) {
            //if (err) return reject(err)
            return resolve(result);
        })
        module.__private_crypto[m].apply(m, args);
    })
}


exports.randomBytes = function (length) {
    return new Promise(function (resolve, reject) {
        module.__private_crypto.randomBytes(length, function (err, result) {
            if (err) return reject(err);
            resolve(result);
        })
    })
    module.__private_crypto.randomBytes(length, function () {
        console.log('rparaprap')
    })
    
} 