var sh = require('sh')

var test = "World"
sh.echo("Call me").pipe('call.txt')
setTimeout(function () {
    console.log('Hello ', test);
}, 1000);

new Promise(function (resolve, reject) {
    setTimeout(function () {
        resolve('Hello from promsie')
    }, 600)
}).then(console.log)

var r = sh.exec("ls")

console.log('waiting 1 sec', r.stdout)


sh.cd('../')
sh.exec('ls').pipe(sh.echo)

/*fetch('http://api.livejazz.dk/concert?access=mobile').then(function (r) {
    return r.text()
}).then(function (e) {
    console.log(e)
})*/