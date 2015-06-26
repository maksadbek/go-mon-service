var React = require('react');
var StatusStore = require('../stores/StatusStore').StatusStore;
var UserStore = require('../stores/StatusStore').UserStore;
var CarActions = require('../actions/StatusActions');
var Sidebar = require('./Sidebar.react');
var LoginForm = require('./LoginForm.react');

function getAllStatuses(){
    return StatusStore.getAll()
}

var StatusApp = React.createClass({
    getInitialState: function(){
        var bounds = new google.maps.LatLngBounds();
        var shape = {
            coords: [1, 1, 1, 20, 18, 20, 18 , 1],
            type: 'poly'
        };
        var myLatlng = new google.maps.LatLng(-25.363882,131.044922);
        var mapOptions = { zoom: 10 };
        var map = new google.maps.Map( document.getElementById('map-canvas'), mapOptions);
        map.fitBounds(bounds);
        return {
            authorized: false,
            map: map,
            bounds: bounds,
            stats: {
                id: '',
                update: {},
                last_request: null
            }
        }
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
        UserStore.addChangeListener(this._onAuth);
    },

    componentWillMount: function(){
        StatusStore.sendAjax();
        setInterval(function(){
            StatusStore.sendAjax();
        }, 5000);
    },
    componentWillUnmount: function(){
        StatusStore.removeChangeListener(this._onChange);
    },

    render: function(){
        var content;
        if(this.state.authorized){
            content = <Sidebar  bounds={this.state.bounds} stats={this.state.stats} map={this.state.map} />
        } else {
           content = <LoginForm/> 
        }
        return (
                <div>
                    {content}
                </div>
        )
    },
    _onChange: function(){
        this.setState({stats: getAllStatuses()});
    },
    _onAuth: function(){
        this.setState({authorized : true});
        this.forceUpdate();
    },
});

module.exports = StatusApp;
