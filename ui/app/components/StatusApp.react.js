var React = require('react');
var StatusStore = require('../stores/StatusStore').StatusStore;
var CarActions = require('../actions/StatusActions');
var UserActions = require('../actions/UserActions');
var Sidebar = require('./Sidebar.react');
var UserStore = require('../stores/StatusStore').UserStore;
var Status = require('./CarStatus.react');

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
            },
            isChildChecked: false
        }
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
        UserStore.addChangeListener(this._onAuth);
    },

    componentWillMount: function(){
        UserActions.Auth({
            login: "taxi",
            uid: "taxi",
            hash: "b5ea8985533defbf1d08d5ed2ac8fe9b",
            fleet: "436",
            groups: "1,2,3" // TODO ochirib tashlash
        });
    },
    componentWillUnmount: function(){
        StatusStore.removeChangeListener(this._onChange);
        UserStore.removeChangeListener(this._onAuth);
    },

    render: function(){
        var content = [];
        var update = this.state.stats.update;
        var checked = this.state.isChildChecked;
        for(var i in update){
            update[i].forEach(function(k){
                content.push(<Status key={k.id} stat={k} isChecked={checked} />);
            });
            //content.push(<Sidebar key={i} groupName={i} stats={update[i]}/>)
        }
        return (   
                <table className={"mdl-data-table mdl-js-data-table  mdl-shadow--2dp"}>
                    <thead>
                        <tr>
                            <th>
                                <label className={"mdl-checkbox mdl-js-checkbox mdl-js-ripple-effect"} htmlFor={"checkbox-2"}>
                                    <input type={"checkbox"} id={"checkbox-2"} className={"mdl-checkbox__input"} />
                                </label>
                            </th>
                            <th className={"mdl-data-table__cell--non-numeric"}>Number</th>
                            <th>Speed</th>
                            <th>Satellites</th>
                            <th>Was online</th>
                            <th>Ignition</th>
                            <th>Fuel</th>
                        </tr>
                    </thead>
                    <tbody>
                    {content}
                    </tbody>
                </table>
                )
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
