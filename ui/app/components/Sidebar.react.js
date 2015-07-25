var React = require('react');
var CarActions = require('../actions/StatusActions');
var Status = require('./CarStatus.react');
var StatusStore = require('../stores/StatusStore').StatusStore;

var Sidebar = React.createClass({
    propTypes:{
        stats: React.PropTypes.array.isRequired,
        groupName: React.PropTypes.string.isRequired
    },
    getInitialState: function(){
        return { style: "", isChildChecked: false}
    },
    render: function(){
        var statuses = [];
        var stat = this.props.stats;
        var group = this.props.groupName;
        var checked = this.state.isChildChecked;
        stat.forEach(function(k){
            statuses.push(<Status key={k.id} stat={k} isChecked={checked} />);
        })
        return (
            <tbody>
                <tr>
                    <td className={"mdl-data-table__cell--non-numeric"}>Number</td>
                    <td>Speed</td>
                    <td>Satellites</td>
                    <td>Was online</td>
                    <td>Ignition</td>
                    <td>Fuel</td>
                </tr>
                {statuses} 
            </tbody>
        );
    },

    _onClickHandler: function(){
        if(this.state.style == "") {
            this.setState({style:"active"});
        }else {
            this.setState({style: ""});
        }
    },
    _onCheckHandler: function(event){
        this.setState({isChildChecked: event.target.checked});
    },
    componentDidMount: function(){
    }
});

module.exports = Sidebar;
