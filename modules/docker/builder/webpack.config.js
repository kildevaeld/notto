const path = require('path');
module.exports = {
    entry: path.resolve(process.cwd(), "src/index.ts"),
    output: {
        path:path.resolve(process.cwd(), '..'),
        filename: 'builder.js',
        library: 'builder',
        libraryTarget: 'commonjs2'
    },
    resolve: {
    // Add `.ts` and `.tsx` as a resolvable extension.
    extensions: ['', '.webpack.js', '.web.js', '.ts', '.tsx', '.js']
    },
    module: {
        loaders: [
            { test: /\.ts$/, loader: 'ts-loader'}
        ]
    },
    externals: {
        underscore: 'underscore',
        docker: 'docker',
        events: 'events',
        util: 'util'
    }
}