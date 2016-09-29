
var Builder = require('docker.builder');
var ui = require('ui');
var n = Builder.Notification;
var argv = require('minimist')(process.argv.slice(1));
console.log('Starting ...')
Builder.createBuilder(require('./doc'), 'staging')
.then( function(builder) {
    
    var p;
    builder.on(Builder.NotificationEvent, function (e, m) {
        var str = Builder.Notification[e] + " " + (Array.isArray(m) ? m.map(function(z) {return z.name}).join(', ') : m.name)
        switch (e) {
            case n.Building:
            case n.Starting:
            case n.Creating:
            case n.Stopping:
                p = ui.Process(" " + str + " ...");
                p.Start()
                break;
            case n.Build:
            case n.Started: 
            case n.Created:
            case n.Stopped:
                p.Success("done");
                p = null;
                break;
        }
        
    })
    
    
    //return builder.remove(true,true);
    return builder.remove(true,true)
    .then(function(e) {
        console.log('done')
    });
    

}).catch(function (e) {
    console.log(e)
})