var React = require('react');
var StatusActions = require('../actions/StatusActions');
var StatusStore = require('../stores/StatusStore');

var CarStatus = React.createClass({
    propTypes:{
        stat: React.PropTypes.object.isRequired,
        isChecked: React.PropTypes.bool
    },
    getInitialState: function(){
        return { 
                isChecked: this.props.isChecked,
                toolTipStyle: "none"
               }
    },
    _addMarker: function(){
        StatusActions.AddMarkerToMap({
            id: this.props.stat.id,
            stat:{
                latitude: this.props.stat.latitude,
                longitude: this.props.stat.longitude,
                direction: this.props.stat.direction,
                speed: this.props.stat.speed,
                sat: this.props.stat.sat,
                owner: this.props.stat.owner,
                formatted_time: this.props.stat.time,
                addparams: this.props.stat.additional,
                car_name: this.props.stat.number
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
        StatusStore.centerMarker(this.props.stat.id);
    },
    _onTitleMouseOver: function(){
        this.setState({toolTipStyle: "block"})
    },
    _onTitleMouseOut: function(){
        this.setState({toolTipStyle: "none"})
    },
    render: function(){
        var stat = this.props.stat;
        
       if(mapLoaded){
           StatusStore.updateMarker(stat);
           StatusStore.redrawMap(false);
       }

        var host = "beta.maxtrack.uz";
        if(typeof(go_mon_site) !== "undefined"){
            host = go_mon_site;
        }

        var speedStatus,
            timeStatus,
            ignitionStatus,
            satStatus,
            batteryStatus,
            fuelStatus;

        if(monitoring_actual_time !== 0){
            // set satellite indicator
            var satIndicator;
            var satTitle = "количество спутников" + stat.sat
            if (stat.sat==6767) {
                    satIndicator = "http://"+host+"/i/monitoring/shield.png";
            } else {
                if (stat.sat >= 0 && stat.sat <=2) {
                    satIndicator = "http://"+host+"/i/monitoring/sat-1.png";
                } else if (stat.sat >=3 && stat.sat <=4) {
                    satIndicator = "http://"+host+"/i/monitoring/sat-2.png";
                } else if (stat.sat >=5 && stat.sat <=6) {
                    satIndicator = "http://"+host+"/i/monitoring/sat-4.png";
                } else {
                    satIndicator = "http://"+host+"/i/monitoring/sat-3.png";
                }
            }
            satStatus = <td><span className="hide_tooltip">{satTitle}</span><img src={satIndicator}/></td>
        } 

        if (monitoring_speed !== 0) {
            var speed;
            if(stat.ignition === 0 && stat.speed === 0){
                speed = <img src={"http://"+host+"/i/monitoring/parking-monitor.jpg"} />
            } else if(stat.speed >= 0 && stat.speed <= 5){
                speed = <span style={{ color:"black" }}><b>{stat.speed}</b></span>
            } else if (stat.speed > 5 && stat.speed < 80) {
                speed = <span style={{color:"blue"}}><b>{stat.speed}</b></span>;
            } else if (speed >= 80) {
                speed = <span style={{color:"red"}}><b>{stat.speed}</b></span>;
            }
            speedStatus = <td>{speed}</td>
        }

        if (monitoring_gprs_condition !== 0) {
            // set time
            stat.time = stat.time.replace(/-/g, " ")
            var time = new Date(stat.time);
            var now = new Date();
            var range = Math.floor((now.getTime() - time.getTime())/ 60000);
            var timeindicator;
            var timeMsg = "";
            if(range >= 24*60 && range < 2*24*60) {
                timeMsg = "Позиция определена 1 дней назад" 
                timeIndicator = "http://"+host+"/i/monitoring/gsm-4.png";
            }else if(range > 60 && range < 24*60){
                timeMsg = "Позиция определена " + Math.ceil((range / 60)) + " час  назад" 
                timeIndicator = "http://"+host+"/i/monitoring/gsm-1.png";
            }else if(range > 20 && range <= 60){
                timeMsg = "Позиция определена 1 час  назад" 
                timeIndicator = "http://"+host+"/i/monitoring/gsm-2.png";
            }else if(range >= 0 && range <= 20){
                timeMsg = "Позиция определена 20 минут  назад" 
                timeIndicator = "http://"+host+"/i/monitoring/gsm-3.png";
            }else {
                timeMsg = "Позиция определена " + Math.floor(range/24/60) + " дней назад" 
                timeIndicator = "http://"+host+"/i/monitoring/gsm-5.png";
            }
            timeStatus = <td><span className="hide_tooltip">{timeMsg}</span><img src={timeIndicator} /></td>
        } 

        if(status_ignition_object !== 0) { 
            // set ignition indicator
            var ignIndicator;
            var ignTitle = "";
            if (stat.fuel_val===0) {
                ignIndicator = "http://"+host+"/i/monitoring/key-off.png";
                ignTitle = "зажигания обьекта отключена";
            } else if (stat.fuel_val > 0) {
                ignIndicator= "http://"+host+"/i/monitoring/key-on.png";
                ignTitle = "зажигания обьекта включена";
            } else {
                ignIndicator = "http://"+host+"/i/monitoring/key-no.png";
            }
            ignitionStatus = <td><span className="hide_tooltip">{ignTitle}</span><img src={ignIndicator} /></td>
        }         

        if(status_fuel !== 0){
            // set fuel indicator
            var fuelIndicator;
            var fuelTitle = "Объем топлива" + stat.fuel_val + "  литр"
            if (stat.fuel_val>=0 && stat.fuel_val<25) {
                fuelIndicator = "http://"+host+"/i/monitoring/fuel-0.png";
            } else if (stat.fuel_val >= 25 && stat.fuel_val < 50) {
                fuelIndicator = "http://"+host+"/i/monitoring/fuel-25.png";
            } else if (stat.fuel_val>=50 && stat.fuel_val<75) {
                fuelIndicator = "http://"+host+"/i/monitoring/fuel-50.png";
            } else if (stat.fuel_val>=75 && stat.fuel_val<95) {
                fuelIndicator = "http://"+host+"/i/monitoring/fuel-75.png";
            }else{
                fuelIndicator = "http://"+host+"/i/monitoring/fuel-100.png";
            }

            fuelStatus = <td><span className="hide_tooltip">{fuelTitle}</span><img src={fuelIndicator} /></td>
        }

        if(status_battery !== 0){
            // set battery indicator
            var batteryIndicator;
            var battery = Math.ceil(stat.battery66 / 1000);
            var batteryTitle = "Питание " + battery + " вольт"
            if (battery > 0 ) {
                batteryIndicator = "http://"+host+"/i/monitoring/battery-full.png";
            } else {
                batteryIndicator = "http://"+host+"/i/monitoring/battery-low.png";
            }
            batteryStatus = <td><span className="hide_tooltip">{batteryTitle}</span><img src={batteryIndicator} /></td>
        }
        var toolTipStyle = this.state.toolTipStyle;

        return (
            <div className="bottom_side">
                <table>
                  <tr>
                    <td onMouseOver={this._onTitleMouseOver} onMouseOut={this._onTitleMouseOut}>
                        <label className="check_bock">
                            <input  onChange={this._onTick} 
                                    value={stat.id} 
                                    checked={this.state.isChecked} 
                                    type="checkbox" name="checkAll" />
                        </label> 
                        <span onClick={this._onTitleClick} id="title_moni">{stat.number+" "+stat.name}</span>
                    </td>
                    <td>
                      <div className="button_monitoring">
                        <table>
                          <tr>
                            {speedStatus}
                            {timeStatus}
                            {satStatus}                    
                            {ignitionStatus}
                            {fuelStatus}
                            {batteryStatus}
                          </tr>
                        </table>
                           <div className="hoverBlock" style={{display: toolTipStyle}}>
                                <table   cellspacing="0">
                                    <tr>
                                        <td><strong>Объект</strong></td>
                                        <td><span>{stat.number}</span></td>
                                    </tr>
                                    <tr>
                                        <td><strong>Имя водителья</strong></td>
                                        <td><span>{stat.name}</span></td>
                                    </tr>
                                    <tr>
                                        <td><strong>Скорость</strong></td>
                                        <td><span>{stat.speed}</span></td>
                                    </tr>
                                    <tr>
                                        <td><strong>Время</strong></td>
                                        <td><span>{stat.time}</span></td>
                                    </tr>
                                    <tr>
                                        <td><span>Долгота</span></td>
                                        <td><span>{stat.longitude}</span></td>
                                    </tr>
                                    <tr>
                                        <td><span>Ширина</span></td>
                                        <td><span>{stat.latitude}</span></td>
                                    </tr>
                                    <tr>
                                        <td><span>Спутник</span></td>
                                        <td><span>{stat.sat}</span></td>
                                    </tr>
                                    <tr>
                                        <td><span>батарейка</span></td>
                                        <td><span>{stat.batterylvl}</span></td>
                                    </tr>
                                </table>
                            </div> 
                      </div>
                    </td>
                  </tr>
                </table>
            </div>
        );
    }
});

module.exports = CarStatus;
