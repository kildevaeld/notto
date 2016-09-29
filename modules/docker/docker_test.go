package docker

import (
	"testing"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules"
)

var rootDir = "./"

func mustError(result error) func(err error) error {
	return func(err error) error {
		if err != nil {
			return multierror.Append(result, err)
		}
		return nil
	}
}

func TestBuilder(t *testing.T) {

	vm := notto.New()

	/*var result error
	result = mustError(result)(shell.Define(vm, false))
	result = mustError(result)(process.Define(vm))
	result = mustError(result)(util.Define(vm))
	result = mustError(result)(promise.Define(vm))
	result = mustError(result)(fetch.Define(vm))
	result = mustError(result)(global.Define(vm))
	result = mustError(result)(fs.Define(vm))
	result = mustError(result)(ui.Define(vm))
	result = mustError(result)(fsm.Define(vm))
	result = mustError(result)(s3.Define(vm))
	result = mustError(result)(archive.Define(vm))
	result = mustError(result)(Define(vm))*/
	result := modules.Define(vm)
	result = Define2(vm, nil)
	if result != nil {
		t.Fatal(result)
	}

	var s = `
        var docker = require('docker.builder');
        var config = require('./test');
        docker.createBuilder(config,'production')
        .then(function (builder) {
			//console.log(JSON.stringify(builder.modules, null, 2));
            /*builder.on('notification', function (e, m) {
                console.log(e,m.name)
            });*/
			builder.build(true)
            //return builder.start(true);
            /*builder.build()
            .then(function () {
                return builder.start();
            }).catch(console.log)*/
        }).catch(function (e) {
			console.log('error',e)
		})
        
    `

	_, e := vm.RunScript(s, "")
	if e != nil {
		t.Fatal(e)
	}
}

func TDocker(t *testing.T) {

	vm := notto.New()

	result := modules.Define(vm)
	result = Define2(vm, nil)
	if result != nil {
		t.Fatal(result)
	}

	s := `
		/*var docker = require('docker').create();

		docker.Start({name: "nginx"}, function (err) {
			console.log('started: ', err)
			docker.HasImage({name: 'sha256:1caaa506d2e83f60781476e5732a7b70e0443ab848cc484085fb2984b40fa805', timeout:1}, function (e, ok) {
				console.log(e, ok)
			});
		});
	
		docker.HasContainer({name:'nginx'}, function (err, ok) {
			console.log(err, ok)
		});

		docker.RemoveContainer({name: 'nginx', force:true}, function (err) {
			console.log(err)
		});

		docker.ListContainers({}, function (err, list) {
			console.log(JSON.stringify(list))
		});*/

		var docker = require('docker');
		var builder = require('docker.builder');

		/*docker.create({
			name:'nginx', 
			image:'alpine', 
			cmd:['sh'], 
			tty: true,
			env: {
				YOU: "ME"
			},
			autoRemove: true,
			attachStdout: true, 
			attachStderr: true, 
			openStdin: true, 
			attachStdin: false, 
			stdinOnce: true
		})
		.then(function (result) {
			
			docker.start({name: 'nginx'})
			.then(function () {
				console.log('done)
			})
		});*/

		docker.remove('something', true)
		.then(function () {
			return docker.build({
				name: "test-image",
				path: "./docket-test"
			})
		}).then(function () {
			return docker.create({
				name: 'something',
				image: 'test-image',
				attachStdout: true, 
				attachStderr: true,
				attachStdin: true,
				tty: true
			})
		}).then(function () {
			return docker.start('something')
		}).then(function () {
			return docker.hasContainer('something').then(function (ok) {
				console.log(ok)
			})
		})
		.catch(function (e) {
			console.error(e)
		})

	
		
		

	`

	if _, e := vm.RunScript(s, ""); e != nil {
		t.Fatal(e)
	}

}
