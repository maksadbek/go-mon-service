var React = require('react');
var CarActions = require('../actions/StatusActions');
var Status = require('./CarStatus.react');
var StatusStore = require('../stores/StatusStore');

var Sidebar = React.createClass({
    propTypes:{
        stats: React.PropTypes.array.isRequired,
        groupName: React.PropTypes.string.isRequired
    },
    getInitialState: function(){
        return { 
                style: "", 
                styleCheckAll: "",
                isChildChecked: false
        }
    },
    render: function(){
        count = 0;
        var statuses = [];
        var stat = this.props.stats;
        var checked = this.state.isChildChecked;
        stat.forEach(function(k){
            count++;
            statuses.push(<Status key={k.id} stat={k} isChecked={checked} />);
        })
        var group = this.props.groupName + " (" + count + ")";
        return (
            <div className={"body_monitoring"}>
                <div className={"show_panel "+ this.state.style + " " + this.state.styleCheckAll}> 
                    <form>
                        <label className="check_bock">
                            <input onChange={this._onCheckHandler} type="checkbox" name="checkAll" />
                        </label> 
                    </form>
                    <div id={"panel_1" } onClick={this._onClickHandler} >{group}</div>
                </div>
                <div id="hide_panel" className={"hide_panel " + this.state.style} > 
                    {statuses} 
                </div>
            </div>
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
        var style = "";
        if(this.state.styleCheckAll === ""){
            style = "checkAll";
        }
        this.setState({
                style:"active",
                styleCheckAll: style,
                isChildChecked: event.target.checked
        });
    },
    componentDidMount: function(){
    }
});

module.exports = Sidebar;
