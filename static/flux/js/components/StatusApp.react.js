var React = require('react');
var StatusStore = require('../stores/StatusStore').StatusStore;
var CarActions = require('../actions/StatusActions');
var Sidebar = require('./Sidebar.react');

function getAllStatuses(){
    return StatusStore.getAll()
}

var StatusApp = React.createClass({
    getInitialState: function(){
        return {
            stats: {
                id: '',
                update: {},
                last_request: null
            }
        }
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
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
        var content = <Sidebar stats={this.state.stats} />
        return (<div className={"body_mon"}>
                    {content}
                </div>)
    },
    _onChange: function(){
        this.setState({stats: getAllStatuses()});
    },
});

module.exports = StatusApp;
