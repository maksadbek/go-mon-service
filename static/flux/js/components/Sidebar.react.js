var React = require('react');
var CarActions = require('../actions/StatusActions');
var Status = require('./CarStatus.react');

var Sidebar = React.createClass({
    propTypes:{
        stats: React.PropTypes.object.isRequired
    },
    render: function(){
        var statuses = [];
        var stat = this.props.stats.update;
        for(var i in stat){
           statuses.push(
            <Status key={stat[i].id} stat={stat[i]} />
           );
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
