
var name = "livejazz-admin";
var mysql = "livejazz-mysql-" + (process.env.RUNN_ENV||'development');

module.exports = {
    name: name,
    link: {
        $staging: {
            mysql: mysql
        }
    },
    $development: {
        publish: "3000:3000"
    },
    envFile: "",
    env: {

    },
    prestart: function () {
        console.log('prestart')
    },
    build: {
        path: process.cwd() + "/admin"
    },
    dependencies: [
        {
            name: mysql,
            phase: ["development", "staging"],
            $darwin: {
                publish: '3306:3306'
            },
            prebuild: function () {

            },
            build: {
                path: process.cwd() + "/mysql"
            },
            dependencies: [
                {name: 'nginx'}
            ]
        }, {
            name: "nginx",
            phase: ["production", "staging"],
            volumne: [
                "./:/etc/confd"
            ],
            prestart: function () {
                console.log('prestart')
            },
            publish: ["80:80","443:443"]

        }
    ]
}