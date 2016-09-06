
var __private_ui = module.__private_ui;

var ui = {
    input: function (msg) {
        return __private_ui.form({
            type: 'input',
            message: msg
        });
    },

    password: function (msg) {
        return __private_ui.form({
            type: 'password',
            message: msg
        });
    },

    chain: function (opts) {
        var out = {};
        for (var key in opts) {
            var val = opts[key];
            if (typeof val === 'string') {
                val = {
                    type: val,
                    message: key
                }
            } else if (!val.message) {
                val.message = key
            }
            out[key] = __private_ui.form(val);
        }
        return out;
    },

    process: function (msg, promise) {
        var p = __private_ui.process(msg);
        p.Start()
        return promise.then(function (r) {
            p.Success("done");
            return r;
        }).catch(function (e) {
            p.Failure("error")
            throw e;
        })
    }
}

module.exports = ui;