var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var StatusConstants = require('../constants/StatusConstants');
var UserConstants = require('../constants/UserConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _carStatus = {};
var _clientInfo = {};
var _token = "";

function setClientInfo(info){
    _clientInfo.fleet = info.fleet;
    _clientInfo.login = info.login;
    _clientInfo.groups = info.groups;
    _clientInfo.hash = info.hash;
}

var UserStore = assign({}, EventEmitter.prototype, {
    auth: function(){
                var xhr = new XMLHttpRequest();
                xhr.open('POST', encodeURI("http://localhost:8080/signup"));
                xhr.setRequestHeader('Content-Type', 'application/json');
                xhr.onload = function() {
                        if (xhr.status === 200 ) {
                             _token = xhr.responseText;
                            console.log(_token);
                            UserStore.emitChange();
                        }
                        else if (xhr.status !== 200) {
                            UserStore.emitChange();
                            return _token;
                        }
                };
                xhr.send(JSON.stringify({
                                user: _clientInfo.login,
                                hash: _clientInfo.hash
                        })
                );
    },
    emitChange: function(){
            this.emit(CHANGE_EVENT);
    },
    addChangeListener: function(callback){
            this.on(CHANGE_EVENT, callback);
    },
    removeChangeListener: function(callback){
            this.removeListener(CHANGE_EVENT, callback);
    },
    dispatcherIndex: AppDispatcher.register(function(action){
            console.log(action)
            switch(action.actionType){
                case UserConstants.AUTH:
                    setClientInfo(action.info);
                    UserStore.auth();
                    break;
            }
            return true;
    })
});

var StatusStore = assign({}, EventEmitter.prototype, {
        sendAjax: function(){
                var xhr = new XMLHttpRequest();
                xhr.open('POST', encodeURI("http://localhost:8080/positions"));
                xhr.setRequestHeader('Content-Type', 'application/json');
                xhr.onload = function() {
                        if (xhr.status === 200 ) {
                            // parse by groups
                            _carStatus = JSON.parse(xhr.responseText);
                            console.log(_carStatus);
                            StatusStore.emitChange();
                        }
                        else if (xhr.status !== 200) {
                            StatusStore.emitChange();
                            return _carStatus;
                        }
                };
                xhr.send(JSON.stringify({
                        selectedFleetJs: "202",
                        user: "newmax",
                        groups: "1,2,3",
                        token: "l5/2LC1m6Hequ2NqV21z35I2XoCR8ju+3o1tZh+rU6Q="
                        })
                );
        },
        getAll: function(){
            return _carStatus;
        },
        emitChange: function(){
                this.emit(CHANGE_EVENT);
        },
        addChangeListener: function(callback){
                this.on(CHANGE_EVENT, callback);
        },
        removeChangeListener: function(callback){
                this.removeListener(CHANGE_EVENT, callback);
        },
        dispatcherIndex: AppDispatcher.register(function(action){
                switch(action.actionType){
                    case StatusConstants.SetClientInfo:
                        SetClientInfo(action.info);
                        StatusStore.emitChange();
                        break;
                }
                return true;
        })
});
module.exports = {
            StatusStore: StatusStore, 
            UserStore: UserStore
};
