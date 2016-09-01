
build: promise-impl.go

promise-impl.go: dist-promise/bundle.js
	go-bindata -pkg promise -o promise-impl.go -prefix dist-promise dist-promise

clean:
	rm -rf promise-impl.go
