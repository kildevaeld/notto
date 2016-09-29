"use strict";
var eachAsync = function (list, cb) {
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
};
