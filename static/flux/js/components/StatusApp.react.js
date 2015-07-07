var React = require('react');
var StatusStore = require('../stores/StatusStore').StatusStore;
var CarActions = require('../actions/StatusActions');
var Sidebar = require('./Sidebar.react');

function setUserInfo(){
    CarActions.SetUserInfo({
            login: "Lizing",
            fleet: "585",
            groups: "1,2,3"
    });
};

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
});

module.exports = StatusApp;
