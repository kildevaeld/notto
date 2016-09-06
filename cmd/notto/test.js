var ui = require('ui');

var out = ui.chain({
	name: 'input',
	password: 'password'
})

console.dir(out)

var p = Promise.delay(1000).then(function () {
	console.log('delay')
})


var p1 = fetch('http://jsonip.com').then(function (res) { return res.json()});

ui.process("Computing ...",Promise.all([p,p1])).then(function (results) {
	console.log(JSON.stringify(results[1]))
}).timeout(1200).catch(function (e) {
	console.log(e)
})

