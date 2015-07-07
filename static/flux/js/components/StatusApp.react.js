var React = require('react');
var StatusStore = require('../stores/StatusStore').StatusStore;
var CarActions = require('../actions/StatusActions');
var UserActions = require('../actions/UserActions');
var Sidebar = require('./Sidebar.react');
var UserStore = require('../stores/StatusStore').UserStore;

function getAllStatuses(){
    return StatusStore.getAll()
}

var StatusApp = React.createClass({
    getInitialState: function(){
        return {
            stats: {
                id: '',
                update: {"":[]},
                last_request: null
            }
        }
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
        UserStore.addChangeListener(this._onAuth);
    },

    componentWillMount: function(){
        UserActions.Auth({
            login: "",
            uid: "",
            hash: "",
            fleet: "",
            groups: ""
        });
    },
    componentWillUnmount: function(){
        StatusStore.removeChangeListener(this._onChange);
        UserStore.removeChangeListener(this._onAuth);
    },

    render: function(){
        var content = [];
        var update = this.state.stats.update;
        for(var i in update){
            content.push(<Sidebar key={i} groupName={i} stats={update[i]}/>)
        }
        return (<div className={"body_mon"}>
                    {content}
                </div>)
    },
    _onChange: function(){
        this.setState({stats: getAllStatuses()});
    },
    _onAuth: function(){
        StatusStore.sendAjax();
        setInterval(function(){
            StatusStore.sendAjax();
        }, 5000);
    }
});

module.exports = StatusApp;
