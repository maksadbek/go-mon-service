var AppDispatcher = require('../dispatcher/AppDispatcher');
var StatusConstants = require('../constants/StatusConstants');

var StatusActions = {
    SetUserInfo: function(info){
        AppDispatcher.dispatch({
                actionType: StatusConstants.SetClientInfo,
                info: info
        });
    },
    AddMarkerToMap: function(info){
        AppDispatcher.dispatch({
                actionType: StatusConstants.AddMarker,
                info: info
        });
    },
    DelMarkerFromMap: function(info){
        AppDispatcher.dispatch({
                actionType: StatusConstants.DelMarker,
                info: info
        });
    }
};

module.exports = StatusActions;
