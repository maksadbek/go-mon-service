// TODO make inputbox for search and button for search cancel
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
            },
            searchCon: {}
        }
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
        UserStore.addChangeListener(this._onAuth);
    },

    componentWillMount: function(){
        UserActions.Auth({
            login: go_mon_login,
            uid: go_mon_uid,
            hash: go_mon_passw,
            fleet:go_mon_fleet,
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
        for(var i in update){
            content.push(<Sidebar key={i} groupName={i} stats={update[i]}/>)
        }
        return (<div>   
                    <form onSubmit={this._onSearch}>
                       <input type="textfield" name="context" /> 
                       <input type="submit" />
                    </form>
                    <button onChange={this._onEmptySearch}>X</button>
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
        if(leader !== null){
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
                name: target.value
        });
    },
    _onEmptySearch: function(){
        CarActions.DelSearchCon();
    }
});

module.exports = StatusApp;
