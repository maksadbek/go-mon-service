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
                            <td><span>{stat.speed}</span></td>
                            <td style={{paddingRight:"11px"}}><img style={{marginTop:"6px"}} src={"./images/default/link.png"} /></td>
                            <td style={{paddingRight:"9px"}}><img style={{marginTop:"3px"}} src={"./images/default/network.png"} /></td>
                            <td style={{paddingRight:"11px"}}><img style={{marginTop:"5px"}} src={"./images/default/key.png"} /></td>
                            <td style={{paddingRight:"12px"}}><img style={{marginTop:"9px"}} src={"./images/default/battery.png"} /></td>
                            <td style={{paddingRight:"8px"}}><img style={{marginTop:"6px"}} src={"./images/default/exit.png"} /></td>
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
