var React = require('react');
var CarActions = require('../actions/StatusActions');
var Status = require('./CarStatus.react');

var markers = [];
var Sidebar = React.createClass({
    propTypes:{
        stats: React.PropTypes.array.isRequired,
        groupName: React.PropTypes.string.isRequired
    },
    getInitialState: function(){
        return {
                style: ""
                }
    },
    render: function(){
        var statuses = [];
        var stat = this.props.stats;
        var group = this.props.groupName;
        stat.forEach(function(k){
            statuses.push(<Status key={k.id} stat={k} />);
        })
        
        return (
            <div className={"body_monitoring"}>
                <div className={"show_panel " + this.state.style} onClick={this._onClickHandler} id={"panel_1"}> 
                    <form>
                        <label className="check_bock"><input type="checkbox" name="checkAll" />{group}</label> 
                    </form>
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
    componentDidMount: function(){
    }
});

module.exports = Sidebar;
