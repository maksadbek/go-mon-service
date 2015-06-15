var React = require('react');
var StatusStore = require('../stores/StatusStore');
var CarActions = require('../actions/StatusActions');
var Sidebar = require('./Sidebar.react');

function setUserInfo(){
    CarActions.SetUserInfo({
            login: "newmax",
            fleet: "202",
            groups: "1,2,3"
    });
};

function getAllStatuses(){
    return StatusStore.getAll()
}
var StatusApp = React.createClass({
    getInitialState: function(){
        return {stats: {
                id: '',
                update: {},
                last_request: null
            }}
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
    },

    componentWillMount: function(){
        setUserInfo();
        StatusStore.sendAjax();
        setInterval(function(){
            StatusStore.sendAjax();
        }, 5000);
    },
    componentWillUnmount: function(){
        StatusStore.removeChangeListener(this._onChange);
    },

    render: function(){
        return (
                <Sidebar stats={this.state.stats} />
        )
    },
    _onChange: function(){
        this.setState({stats: getAllStatuses()});
    }
});

module.exports = StatusApp;
