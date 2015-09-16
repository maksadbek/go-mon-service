var React = require('react');
var StatusApp = require('./components/Main.react');
var fleetListDOM = document.getElementById('fleet_list');
if(fleetListDOM !== null){ 
    React.render( <StatusApp/>, document.getElementById('fleet_list'));
} else {
    console.error("gomon: нет дом элемента");
}
