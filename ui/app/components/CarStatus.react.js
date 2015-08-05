var React = require('react');
var StatusActions = require('../actions/StatusActions');
var StatusStore = require('../stores/StatusStore').StatusStore;
var Mui  = require('material-ui');
var ThemeManager = new Mui.Styles.ThemeManager();
var ListItem = Mui.ListItem;
var FontIcon = Mui.FontIcon;
var CheckBox = Mui.CheckBox; 

var CarStatus = React.createClass({
    propTypes:{
        stat: React.PropTypes.object.isRequired,
    },
    getInitialState: function(){
        return { isChecked: this.props.isChecked }
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
        mon.setCenterObj(this.props.stat.id);
    },
    render: function(){
        var stat = this.props.stat;
        StatusStore.updateMarker(stat);

        // set time
        var time = new Date(stat.time);
        var now = new Date(Date.now());
        var delta = Math.abs(now - time) / 1000000;
        var rangeInMinutes = Math.floor(delta / 60)
        var timeIndicator;
        var timeMsg = "";
        
        if(rangeInMinutes >= 24*60) {
            timeMsg = "1 day ago";
        }else if(rangeInMinutes > 60 && rangeInMinutes < 24*60){
            timeMsg = Math.ceil((rangeInMinutes / 60)) + " hours ago";
        }else if(rangeInMinutes > 20 && rangeInMinutes <= 60){
            timeMsg = "1 hour ago";
        }else if(rangeInMinutes >= 0 && rangeInMinutes <= 20){
            timeMsg = "within 20 minutes"; 
        }

        // set ignition indicator
        var ign;
        if (stat.fuel_val===0) {
            ign= "off";
        } else if (stat.fuel_val > 0) {
            ign= "on";
        } else {
            ign= "on";
        }
        return (
              <ListItem  primaryText={stat.number} />
        );
    }
});

module.exports = CarStatus;
