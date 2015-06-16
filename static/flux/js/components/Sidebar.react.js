var React = require('react');
var CarActions = require('../actions/StatusActions');
var Status = require('./CarStatus.react');


var markers = [];
var Sidebar = React.createClass({
    propTypes:{
        stats: React.PropTypes.object.isRequired,
        map: React.PropTypes.object.isRequired,
        bounds: React.PropTypes.object.isRequired
    },
    render: function(){
        var statuses = [];
        var stat = this.props.stats.update;
        var map = this.props.map;
        var bounds = this.props.bounds;
        var shape = {
            coords: [1, 1, 1, 20, 18, 20, 18 , 1],
            type: 'circle'
        };
        if(markers.length === 0){
            for(var i in stat){
                var latLng = new google.maps.LatLng(
                    stat[i].latitude, stat[i].longitude
                );
                    var marker = new google.maps.Marker({
                        position: latLng,
                        map: map,
                        title: stat[i].number,
                        id: stat[i].id,
                        shape: shape,
                    });
                    markers.push(marker);
                bounds.extend(marker.position);
                statuses.push(
                    <Status key={stat[i].id} stat={stat[i]} />
                );
            }
        } else{
            markers.forEach(function(m){
                for(var i in stat){
                    if(stat[i].id === m.id){
                        m.setPosition(new google.maps.LatLng(stat[i].latitude, stat[i].longitude));
                        m.title = stat[i].number;
                        m.id = stat[i].id;
                        statuses.push(
                            <Status key={stat[i].id} stat={stat[i]} />
                        );
                    }
                }
            });
        }
        return (
            <nav className="menu">
                {statuses} 
            </nav>
        );
    },
    componentDidMount: function(){
        console.log(this.props.stats);
    }
});

module.exports = Sidebar;
