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
    },
    SearchCar: function(info){
        AppDispatcher.dispatch({
                actionType: StatusConstants.SearchCar,
                info: info
        });
    },
    DelSearchCon: function(info){
        AppDispatcher.dispatch({
                actionType: StatusConstants.DelSearchCon
        });
    },
    SelectGroup: function(info){
        AppDispatcher.dispatch({
                actionType: StatusConstants.SelectGroup,
                info: info
        });
    }
};

module.exports = StatusActions;
