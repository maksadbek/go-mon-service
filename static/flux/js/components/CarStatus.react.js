var React = require('react');
var StatusActions = require('../actions/StatusActions');
var StatusStore = require('../stores/StatusStore').StatusStore;

var CarStatus = React.createClass({
    propTypes:{
        stat: React.PropTypes.object.isRequired,
        isChecked: React.PropTypes.bool

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
        StatusStore.redrawMap();

        // set speed
        var speed;
        if(stat.ignition === 0 && stat.speed === 0){
            speed = <img src={"http://"+go_mon_site+"/i/monitoring/parking-monitor.jpg"} />
        } else if(stat.speed >= 0 && stat.speed <= 5){
            speed = <span style={{ color:"black" }}><b>{stat.speed}</b></span>
        } else if (stat.speed > 5 && stat.speed < 80) {
            speed = <span style={{color:"blue"}}><b>{stat.speed}</b></span>;
        } else if (speed >= 80) {
            speed = <span style={{color:"red"}}><b>{stat.speed}</b></span>;
        }

        // set time
        stat.time = stat.time.replace(/-/g, " ")
        var time = new Date(stat.time);
        var now = new Date(Date.now());
        var delta = Math.abs(now - time) / 1000000;
        var rangeInMinutes = Math.floor(delta / 60)
        var timeIndicator;
        var timeMsg = "";
        
        if(rangeInMinutes >= 24*60) {
        timeMsg = "Позиция определена 1 дней назад" 
            timeIndicator = "http://"+go_mon_site+"/i/monitoring/gsm-4.png";
        }else if(rangeInMinutes > 60 && rangeInMinutes < 24*60){
        timeMsg = "Позиция определена" + Math.ceil((rangeInMinutes / 60)) + " час  назад" 
            timeIndicator = "http://"+go_mon_site+"/i/monitoring/gsm-1.png";
        }else if(rangeInMinutes > 20 && rangeInMinutes <= 60){
        timeMsg = "Позиция определена 1 час  назад" 
            timeIndicator = "http://"+go_mon_site+"/i/monitoring/gsm-2.png";
        }else if(rangeInMinutes >= 0 && rangeInMinutes <= 20){
        timeMsg = "Позиция определена 20 минут  назад" 
            timeIndicator = "http://"+go_mon_site+"/i/monitoring/gsm-3.png";
        }

        // set satellite indicator
        var satIndicator;
        var satTitle = "количество спутников" + stat.sat
        if (stat.sat==6767) {
                satIndicator = "http://"+go_mon_site+"/i/monitoring/shield.png";
        } else {
                if (stat.sat >= 0 && stat.sat <=2) {
                    satIndicator = "http://"+go_mon_site+"/i/monitoring/sat-1.png";
                } else if (stat.sat >=3 && stat.sat <=4) {
                    satIndicator = "http://"+go_mon_site+"/i/monitoring/sat-2.png";
                } else if (stat.sat >=5 && stat.sat <=6) {
                    satIndicator = "http://"+go_mon_site+"/i/monitoring/sat-4.png";
                } else {
                    satIndicator = "http://"+go_mon_site+"/i/monitoring/sat-3.png";
                }
        }

        // set ignition indicator
        var ignIndicator;
        var ignTitle = "";
        if (stat.fuel_val===0) {
            ignIndicator = "http://"+go_mon_site+"/i/monitoring/key-off.png";
        ignTitle = "зажигания обьекта отключена";
        } else if (stat.fuel_val > 0) {
            ignIndicator= "http://"+go_mon_site+"/i/monitoring/key-on.png";
        ignTitle = "зажигания обьекта включена";
        } else {
            ignIndicator = "http://"+go_mon_site+"/i/monitoring/key-no.png";
        }

        // set fuel indicator
        var fuelIndicator;
        var fuelTitle = "Объем топлива" + stat.fuel_val + "  литр"
        if (stat.fuel_val>=0 && stat.fuel_val<25) {
            fuelIndicator = "http://"+go_mon_site+"/i/monitoring/fuel-0.png";
        } else if (stat.fuel_val >= 25 && stat.fuel_val < 50) {
            fuelIndicator = "http://"+go_mon_site+"/i/monitoring/fuel-25.png";
        } else if (stat.fuel_val>=50 && stat.fuel_val<75) {
            fuelIndicator = "http://"+go_mon_site+"/i/monitoring/fuel-50.png";
        } else if (stat.fuel_val>=75 && stat.fuel_val<95) {
            fuelIndicator = "http://"+go_mon_site+"/i/monitoring/fuel-75.png";
        }else{
            fuelIndicator = "http://"+go_mon_site+"/i/monitoring/fuel-100.png";
        }

        
        return (
            <div className="bottom_side">
                <table>
                  <tr>
                    <td>
                        <label className="check_bock">
                            <input onChange={this._onTick} value={stat.id} checked={this.state.isChecked} type="checkbox" name="checkAll" />
                        </label> 
                        <span onClick={this._onTitleClick} id="title_moni">{stat.number}</span>
                    </td>
                    <td>
                      <div className="button_monitoring">
                        <table>
                          <tr>
                            <td>{speed}</td>
                            <td style={{paddingRight:"11px"}}><img title={timeMsg}  src={timeIndicator} /></td>
                            <td style={{paddingRight:"9px"}}><img title={satTitle}  src={satIndicator} /></td>
                            <td style={{paddingRight:"11px"}}><img title={ignTitle}  src={ignIndicator} /></td>
                            <td style={{paddingRight:"12px"}}><img title={fuelTitle} src={fuelIndicator} /></td>
                          </tr>
                        </table>
                      </div>
                    </td>
                  </tr>
                </table>
            </div>
        );
    }
});

module.exports = CarStatus;
