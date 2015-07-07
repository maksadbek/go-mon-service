var React = require('react');
var TodoActions = require('../actions/StatusActions');

var CarStatus = React.createClass({
    propTypes:{
        stat: React.PropTypes.object.isRequired
    },
    getInitialState: function(){
        return this.props.stat ;
    },
    render: function(){
        var stat = this.props.stat;
	/*
        var info =  "id: "+stat.id+ "\n"+
                    "latitude: "+stat.latitude+"\n"+
                    "longitude:"+ stat.longitude+ "\n"+
                    "time:"+stat.time+"\n"+
                    "owner:"+ stat.owner+"\n"+ 
                    "number:"+ stat.number+"\n"+
                    "direction:"+stat.direction+"\n"+
                    "speed:"+stat.speed+"\n"+
                    "sat:"+stat.sat+"\n"+
                    "ignition:"+stat.ignition+
                    "gsmsignal:"+stat.gsmsignal+"\n"+
                    "battery66:"+stat.battery66+"\n"+
                    "seat:"+stat.seat+"\n"+
                    "batterylvl:"+stat.batterylvl+"\n"+
                    "fuel:"+stat.fuel+"\n"+
                    "fuel_val:"+stat.fuel_val+"\n"+
                    "mu_additional:"+stat.mu_additional+"\n"+
                    "customization:"+stat.customization+"\n"+
                    "additional:"+stat.additional+"\n"+
                    "action:"+stat.action+ ""+"\n";
	*/
        // set speed
        var speed;
        if(stat.ignition === 0 && stat.speed === 0){
            speed = <img src={"./images/parking.jpg"} />
        } else if(stat.speed >= 0 && stat.speed <= 5){
            speed = <span style={{ color:"black" }}><b>{stat.speed}</b></span>
        } else if (stat.speed > 5 && stat.speed < 80) {
            speed = <span style={{color:"blue"}}><b>{stat.speed}</b></span>;
        } else if (speed >= 80) {
            speed = <span style={{color:"red"}}><b>{stat.speed}</b></span>;
        }

        // set time
        var time = new Date(stat.time);
        var now = new Date(Date.now());
        var delta = Math.abs(now - time) / 1000;
        var rangeInMinutes = Math.floor(delta / 60)
        var timeIndicator;
	var timeMsg = ""
        
        if(rangeInMinutes >= 24*60) {
	    timeMsg = "Позиция определена 1 дней назад" 
            timeIndicator = "./images/gsm-4.png";
        }else if(rangeInMinutes > 60 && rangeInMinutes < 24*60){
	    timeMsg = "Позиция определена" + (rangeInMinutes / 60) + " час  назад" 
            timeIndicator = "./images/gsm-1.png";
        }else if(rangeInMinutes > 20 && rangeInMinutes <= 60){
	    timeMsg = "Позиция определена 1 час  назад" 
            timeIndicator = "./images/gsm-2.png";
        }else if(rangeInMinutes >= 0 && rangeInMinutes <= 20){
	    timeMsg = "Позиция определена 20 минут  назад" 
            timeIndicator = "./images/gsm-3.png";
        }

        // set satellite indicator
        var satIndicator;
	var satTitle = "количество спутников" + stat.sat
        if (stat.sat==6767) {
                satIndicator = "/minus-shield.png";
        } else {
                if (stat.sat >= 0 && stat.sat <=2) {
                    satIndicator = "./images/sat-1.png";
                } else if (stat.sat >=3 && stat.sat <=4) {
                    satIndicator = "./images/sat-2.png";
                } else if (stat.sat >=5 && stat.sat <=6) {
                    satIndicator = "./images/sat-4.png";
                } else {
                    satIndicator = "./images/sat-3.png";
                }
        }

        // set ignition indicator
        var ignIndicator;
	var ignTitle = "";
        if (stat.fuel_val===0) {
            ignIndicator = "./images/key-off.png";
	    ignTitle = "зажигания обьекта отключена";
        } else if (stat.fuel_val > 0) {
            ignIndicator= "./images/key-on.png";
	    ignTitle = "зажигания обьекта включена";
        } else {
            ignIndicator = "./images/key-no.png";
        }
	/*
        if (stat.gsmsignal !== '' && stat.gsmsignal !== '-1') {           
            ignIndicator = "./images/android.png";
        }
	*/

        // set fuel indicator
        var fuelIndicator;
	var fuelTitle = "Объем топлива" + stat.fuel_val + "  литр"
        if (stat.fuel_val>=0 && stat.fuel_val<25) {
            fuelIndicator = "./images/fuel-0.png";
        } else if (stat.fuel_val >= 25 && stat.fuel_val < 50) {
            fuelIndicator = "./images/fuel-25.png";
        } else if (stat.fuel_val>=50 && stat.fuel_val<75) {
            fuelIndicator = "./images/fuel-50.png";
        } else if (stat.fuel_val>=75 && stat.fuel_val<95) {
            fuelIndicator = "./images/fuel-75.png";
        }else{
            fuelIndicator = "./images/fuel-100.png";
        }

        return (
            <div className="bottom_side">
                <table>
                  <tr>
                    <td>
                        <label className="check_bock">
                            <input type="checkbox" name="checkAll" />
                        </label> 
                        <span id="title_moni">{stat.number}</span>
                    </td>
                    <td>
                      <div className="button_monitoring">
                        <table>
                          <tr>
                            <td>{speed}</td>
                            <td style={{paddingRight:"11px"}}><img title={timeMsg} style={{marginTop:"6px"}} src={timeIndicator} /></td>
                            <td style={{paddingRight:"9px"}}><img title={satTitle} style={{marginTop:"3px"}} src={satIndicator} /></td>
                            <td style={{paddingRight:"11px"}}><img title={ignTitle} style={{marginTop:"5px"}} src={ignIndicator} /></td>
                            <td style={{paddingRight:"12px"}}><img title={fuelTitle} style={{marginTop:"9px"}} src={fuelIndicator} /></td>
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
