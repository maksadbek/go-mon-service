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

if(typeof(go_mon_fleet) !== "undefined"){
    userFleet = go_mon_fleet;
}

var StatusApp = React.createClass({
    getInitialState: function(){
        return {
            stats: {
                id: '',
                update: {"":[]},
                last_request: null
            },
            groupIndex: 0,
            searchPanelStyle: ""

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
        var speedButton,
            ignitionButton,
            satButton,
            wifiButton,
            fuelButton;



        if(monitoring_speed !== 0){
            speedButton = <button style={{"width":"34px"}}>
                            <img title="Скорость" src={"http://online.maxtrack.uz/i/monitoring/speed-header.png"}/> 
                        </button>
        }
        if(monitoring_actual_time !== 0){
            satButton =  <button style={{"width":"34px"}}>
                            <img src={"http://online.maxtrack.uz/i/monitoring/gsm-header.png"}/> 
                        </button>
        }
        if(status_ignition_object !== 0){
            ignitionButton = <button style={{"width":"34px"}}>
                                <img    title={"Статус зажигании объекта"} 
                                        src={"http://online.maxtrack.uz/i/monitoring/key-solid.png"} /> 
                             </button>
        } 
        if(monitoring_gprs_condition !== 0){
            wifiButton = <button style={{"width":"34px"}}>
                            <img title={"Актуальность позиции во времени и пространстве"} 
                                 src={"http://online.maxtrack.uz/i/monitoring/sat-header.png"} /> 
                        </button>
        }
        if(status_fuel !== 0){
            fuelButton = <button style={{"width":"34px"}}>
                            <img title={"Уровень топлива"} src={"http://online.maxtrack.uz/i/monitoring/fuel-header-tr.png"} /> 
                         </button>
        }
        sPanelStyle = this.state.searchPanelStyle;
        var update = this.state.stats.update;
        StatusStore.groupNames.forEach(function(group, id){
            groups.push(<option value={id}>{group}</option>);
        });
        for(var i in update){
            content.push(<Sidebar key={i} groupName={i} stats={update[i]}/>)
        }
        return (<div>   
                    <div className={"search_blocks x-menu x-menu-floating x-layer " + sPanelStyle} id="hide_serach">
                       <form onSubmit={this._onSearch} className="x-menu-list" style={{"height":"24px"}}>
                           <div className="x-menu-list-item ">
                               <div className="x-form-field-wrap x-form-field-trigger-wrap" style={{"width":"150px"}}>
                                   <input  ref="searchText" type="textfield" name="context" 
                                            style={{"vertical-align":"top", "width":"133px"}} 
                                            placeholder="Поиск обектов" 
                                            className="x-form-text x-form-field">
                                       <img src="/e/resources/images/default/s.gif" 
                                           style={{"position":"absolute", "right":"23px"}}
                                             onClick={this._onEmptySearch}
                                             className="x-form-trigger x-form-clear-trigger x-form-trigger-click" />
                                       <button style={{"background": "none", 
                                                        "border": "medium none", 
                                                        "left": "-6px", 
                                                        "padding": " 0", 
                                                        "position": "relative", 
                                                        "top":"0", 
                                                        "vertical-align": "top"}}
                                                      className="x-form-twin-triggers"
                                                      type="submit">
                                                  <img className="x-form-trigger x-form-search-trigger " 
                                                       src="/e/resources/images/default/s.gif" />
                                        </button>
                                   </input>
                               </div>
                           </div>
                       </form>
                    </div>
                    <div id={"west_side"}>
                        <div className="x-panel-tbar">
                            <div className="top_side x-toolbar x-small-editor x-toolbar-layout-ct">
                                <table style={{"float":"left"}}>
                                    <tr>
                                        <td>
                                            <select onChange={this._onGroupSelect} >
                                                {groups}
                                            </select>
                                        </td>
                                    </tr>
                                </table>
                                <table cellspacing={"0"} style={{"margin-top": "-7px", "float":"right"}}>
                                    <tbody>
                                        <tr>
                                            <td className="x-toolbar-cell">
                                                <span className="xtb-sep"></span>
                                            </td>
                                            <td className="x-toolbar-cell">
                                                <table  cellspacing="0" 
                                                        className="x-btn x-btn-icon" 
                                                        style={{"float": "left", "marginRight": "2px", "width": "40px"}}>
                                                    <tbody className="x-btn-small x-btn-icon-small-left">
                                                        <tr>
                                                            <td className="x-btn-tl">
                                                                <i>&nbsp;</i>
                                                            </td>
                                                            <td className="x-btn-tc"></td>
                                                            <td className="x-btn-tr">
                                                                <i>&nbsp;</i>
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td className="x-btn-ml">
                                                                <i>&nbsp;</i>
                                                            </td>
                                                            <td className="x-btn-mc">
                                                              <em onClick={this._onSearchClick} unselectable="on" className="x-btn-arrow">
                                                                    <button  
                                                                            id="search_show" 
                                                                            type="button" 
                                                                            className={"x-btn-text "+ sPanelStyle}>
                                                                            <img style={{"position": "relative", "right": "6px"}} 
                                                                            src="/i/monitoring/magnifier-zoom.png" />
                                                                    </button>
                                                              </em>
                                                            </td>
                                                            <td className="x-btn-mr">
                                                                <i>&nbsp;</i>
                                                            </td>
                                                        </tr>
                                                        <tr>
                                                            <td className="x-btn-bl">
                                                                <i>&nbsp;</i>
                                                            </td>
                                                            <td className="x-btn-bc"></td>
                                                            <td className="x-btn-br">
                                                                <i>&nbsp;</i>
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                         </div>
                        <div className={"bottom_side"}>
                            <table>
                                <tr>
                                    <td>
                                        <button id={"sort_button"}>Автомобиль</button>
                                    </td>
                                    <td>
                                        <div className={"button_monitoring"}>
                                        {speedButton}
                                        {wifiButton}
                                        {satButton}
                                        {ignitionButton}
                                        {fuelButton}
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
        // if MonReqToggler is false,
        // then, stop sending requests
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
        CarActions.SelectGroup({
            id: parseInt(event.target.value)
        });
    },
    _onSearchClick: function(event){
        if(this.state.searchPanelStyle === ""){
            this.setState({searchPanelStyle: "gomon-searchPenel-show"})
        } else {
            this.setState({searchPanelStyle: ""})
        }
    },
});

module.exports = StatusApp;
