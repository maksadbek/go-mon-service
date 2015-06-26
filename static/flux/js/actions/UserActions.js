var AppDispatcher = require('../dispatcher/AppDispatcher');
var UserConstants= require('../constants/UserConstants');

var UserActions = {
    Authenticate: function(info){
        AppDispatcher.dispatch({
                actionType: UserConstants.AUTH,
                info: info
        });
    }
};

module.exports = UserActions;
