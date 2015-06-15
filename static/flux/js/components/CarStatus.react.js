var React = require('react');
var TodoActions = require('../actions/StatusActions');

var CarStatus = React.createClass({
    propTypes:{
        stat: React.PropTypes.object.isRequired
    },
    getInitialState: function(){
        return this.props.stat ;
    },
    render: function(){
        var stat = this.props.stat;
        var info =  "id: "+stat.id+ "\n"+
                    "latitude: "+stat.latitude+"\n"+
                    "longitude:"+ stat.longitude+ "\n"+
                    "time:"+stat.time+"\n"+
                    "owner:"+ stat.owner+"\n"+ 
                    "number:"+ stat.number+"\n"+
                    "direction:"+stat.direction+"\n"+
                    "speed:"+stat.speed+"\n"+
                    "sat:"+stat.sat+"\n"+
                    "ignition:"+stat.ignition+
                    "gsmsignal:"+stat.gsmsignal+"\n"+
                    "battery66:"+stat.battery66+"\n"+
                    "seat:"+stat.seat+"\n"+
                    "batterylvl:"+stat.batterylvl+"\n"+
                    "fuel:"+stat.fuel+"\n"+
                    "fuel_val:"+stat.fuel_val+"\n"+
                    "mu_additional:"+stat.mu_additional+"\n"+
                    "customization:"+stat.customization+"\n"+
                    //"additional:"+stat.additional+"\n"+
                    "action:"+stat.action+ ""+"\n";

        return (
            <a 
                className="menu-item tooltipped tooltipped-e" 
                aria-label={info}>{stat.number}
            </a>
        );
    }
});

module.exports = CarStatus;
