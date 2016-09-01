

build:
	
	@curl -L http://underscorejs.org/underscore-min.js -o underscore.js
	@git clone https://github.com/substack/minimist.git
	
	@go-bindata -pkg util -o modules.go underscore.js minimist/index.js
	@rm -rf minimist
	@rm -f underscore.js