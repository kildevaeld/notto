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

try {
	var r = sh.exec("lsds")	
	console.log('waiting 1 sec', r.stdout)
} catch (e) {
	console.log('GOT error', e)
}



sh.cd('../')
sh.exec('ls /').pipe(sh.echo)

sh.cd('_example')
var out = sh.cat('call.txt')
console.log(out)
/*fetch('http://api.livejazz.dk/concert?access=mobile').then(function (r) {
    return r.text()
}).then(function (e) {
    console.log(e)
})*/