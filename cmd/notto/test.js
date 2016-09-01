

var p = Promise.delay(1000).then(function () {
	console.log('delay')
})

var p1 = fetch('http://jsonip.com').then(function (res) { return res.json()});

Promise.all([p,p1]).then(function (results) {
	console.log(JSON.stringify(results[1]))
}).timeout(1200).catch(function (e) {
	console.log(e)
})