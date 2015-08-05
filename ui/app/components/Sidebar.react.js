var React = require('react');
var CarActions = require('../actions/StatusActions');
var Status = require('./CarStatus.react');
var StatusStore = require('../stores/StatusStore').StatusStore;
var Mui  = require('material-ui');
var ThemeManager = new Mui.Styles.ThemeManager();
var ListItem = Mui.ListItem;
var Checkbox = Mui.Checkbox;

var Sidebar = React.createClass({
    childContextTypes: {
          muiTheme: React.PropTypes.object
    },
    getChildContext: function() {
        return {
            muiTheme: ThemeManager.getCurrentTheme()
        };
    },
    propTypes:{
        stats: React.PropTypes.array.isRequired,
        groupName: React.PropTypes.string.isRequired
    },
    getInitialState: function(){
        return { style: "", isChildChecked: false}
    },
    componentWillReceiveProps: function(nextProps) {
        if(nextProps.isChecked === this.props.isChecked){
            return;
        }
        this.setState({ isChecked: nextProps.isChecked});
        if(nextProps.isChecked){
            this._addMarker();
        }else{
            this._delMarker();
        }
    },
    _addMarker: function(){
        StatusActions.AddMarkerToMap({
            id: this.props.stat.id,
            pos: {
                id: this.props.stat.id,
                latitude: this.props.stat.latitude,
                longitude: this.props.stat.longitude,
                direction: this.props.stat.direction,
                speed: this.props.stat.speed,
                sat: this.props.stat.sat,
                owner: this.props.stat.owner,
                formatted_time: this.props.stat.time,
                addparams: this.props.stat.additional,
                car_name: this.props.stat.number,
                action: '2'
            }
        });
    },
    _delMarker: function(){
        StatusActions.DelMarkerFromMap({
            id: this.props.stat.id
        });
    },
    _onTick: function(event){
        if(this.state.isChecked){
            this.setState({isChecked : false});
            this._delMarker();
        } else {
            this.setState({isChecked : true});
            this._addMarker();
        }
    },
    _onTitleClick: function(){
        // on click to the title, center the marker on the map
    },
    render: function(){
        var statuses = [];
        var stat = this.props.stats;
        var group = this.props.groupName;
        var checked = this.state.isChildChecked;
        stat.forEach(function(k){
            //statuses.push(<Status key={k.id} stat={k} isChecked={checked} />);
            var item = <ListItem  key={k.id} primaryText={k.number} leftCheckbox={<Checkbox name="checkbox"></Checkbox>}>
                        </ListItem>
            statuses.push(item);
        })
        return (
            <ListItem open={true} primaryText={group} leftCheckbox={<Checkbox name="checkbox"></Checkbox>} disabled={false}>
                    {statuses} 
            </ListItem>
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
        this.setState({isChildChecked: event.target.checked});
    },
});

module.exports = Sidebar;
