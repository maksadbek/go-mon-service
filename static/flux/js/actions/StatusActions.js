var AppDispatcher = require('../dispatcher/AppDispatcher');
var TodoConstants = require('../constants/StatusConstants');

var StatusActions = {
    SetUserInfo: function(info){
        AppDispatcher.dispatch({
                actionType: TodoConstants.SetClientInfo,
                info: info
        });
    }
};
module.exports = StatusActions;
