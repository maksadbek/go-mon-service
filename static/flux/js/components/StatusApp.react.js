var React = require('react');
var StatusStore = require('../stores/StatusStore');
var CarActions = require('../actions/StatusActions');
var UserActions = require('../actions/UserActions');
var Sidebar = require('./Sidebar.react');
var UserStore = require('../stores/UserStore');

function getAllStatuses(){
    return StatusStore.getAll()
}
var userLogin   = "";
var userUID     = "";
var userHash    = "";
var userFleet   = "";

if(typeof(go_mon_login) !== "undefined"){
    userLogin = go_mon_login;
}

if(typeof(go_mon_uid) !== "undefined"){
    userUID = go_mon_uid;
}

if(typeof(go_mon_passw) !== "undefined"){
    userHash= go_mon_passw;
}

if(typeof(go_mon_login) !== "undefined"){
    go_mon_login = go_mon_fleet;
}

var StatusApp = React.createClass({
    getInitialState: function(){
        return {
            stats: {
                id: '',
                update: {"":[]},
                last_request: null
            },
            groupIndex: 0
        }
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
        UserStore.addChangeListener(this._onAuth);
    },

    componentWillMount: function(){
        UserActions.Auth({
            login:  userLogin,
            uid:    userUID,
            hash:   userHash,
            fleet:  userFleet,
            groups: "1,2,3" // TODO ochirib tashlash
        });
    },
    componentWillUnmount: function(){
        StatusStore.removeChangeListener(this._onChange);
        UserStore.removeChangeListener(this._onAuth);
    },

    render: function(){
        var content = [];
        var groups = [];
        var update = this.state.stats.update;
        StatusStore.groupNames.forEach(function(group, id){
            groups.push(<option value={id}>{group}</option>);
        });
        for(var i in update){
            content.push(<Sidebar key={i} groupName={i} stats={update[i]}/>)
        }
        return (<div>   
                    <select onChange={this._onGroupSelect} >
                        {groups}
                    </select>
                    <form onSubmit={this._onSearch}>
                       <input ref="searchText" type="textfield" name="context" /> 
                       <input type="submit" />
                    </form>
                    <button onClick={this._onEmptySearch}>X</button>
                    <div id={"west_side"}>
                        <div className={"bottom_side"}>
                            <table>
                                <tr>
                                    <td>
                                        <button id={"sort_button"}>Автомобиль</button>
                                    </td>
                                    <td>
                                        <div className={"button_monitoring"}>
                                            <button style={{"width":"28px", "margin-right":"0px"}}>
                                                <img title="Скорость" src={"http://online.maxtrack.uz/i/monitoring/speed-header.png"}/> 
                                            </button>
                                            <button style={{"width":"33px", "margin-right":"0px"}}>
                                                <img src={"http://online.maxtrack.uz/i/monitoring/gsm-header.png"}/> 
                                            </button>
                                            <button style={{"width":"34px", "margin-right":"-6"}}>
                                                <img    title={"Актуальность позиции во времени и пространстве"} 
                                                        src={"http://online.maxtrack.uz/i/monitoring/sat-header.png"}
                                                /> 
                                            </button>
                                            <button style={{"width":"26px", "margin-right":"-4"}}>
                                                <img title={"Статус зажигании объекта"} 
                                                     src={"http://online.maxtrack.uz/i/monitoring/key-solid.png"}
                                                /> 
                                            </button>
                                            <button style={{"width":"25px", "margin-right":"25px"}}>
                                                <img title={"Уровень топлива"} 
                                                     src={"http://online.maxtrack.uz/i/monitoring/fuel-header-tr.png"} 
                                                /> 
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                     		</table>
                        </div>
                    </div>
                    <div className={"body_mon"}>
                        {content}
                    </div>
                </div>)
    },
    _onChange: function(){
        this.setState({stats: getAllStatuses()});
        var loader = document.getElementById("gomon-loader");
        if(loader !== null){
            loader.remove();
        }
    },
    _onAuth: function(){
        StatusStore.sendAjax();
        setInterval(function(){
            StatusStore.sendAjax();
        }, 5000);
    },
    _onSearch: function(event){
        event.preventDefault();
        var target = event.target
        CarActions.SearchCar({
                name: target.context.value
        });
    },
    _onEmptySearch: function(){
        this.refs.searchText.getDOMNode().value = "";
        CarActions.DelSearchCon();
    },
    _onGroupSelect: function(event){
        console.log(event.target.value);
        CarActions.SelectGroup({
            id: parseInt(event.target.value)
        });
    }
});

module.exports = StatusApp;
