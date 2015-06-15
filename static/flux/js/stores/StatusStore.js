var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var StatusConstants = require('../constants/StatusConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _carStatus = {};
var _clientInfo = {};

function SetClientInfo(info){
    _clientInfo.fleet = info.fleet;
    _clientInfo.login = info.login;
    _clientInfo.groups = info.groups;
}

var StatusStore = assign({}, EventEmitter.prototype, {
        sendAjax: function(){
                var xhr = new XMLHttpRequest();
                xhr.open('POST', encodeURI("http://localhost:8080/positions"));
                xhr.setRequestHeader('Content-Type', 'application/json');
                xhr.onload = function() {
                        if (xhr.status === 200 ) {
                             _carStatus = JSON.parse(xhr.responseText);
                            console.log(_carStatus);
                            StatusStore.emitChange();

                        }
                        else if (xhr.status !== 200) {
                            return _carStatus;
                            StatusStore.emitChange();
                        }
                };
                xhr.send(JSON.stringify({
                                selectedFleetJs: _clientInfo.fleet,
                                user: _clientInfo.login,
                                groups: _clientInfo.groups
                                }));
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
                        console.log(action.info);
                        SetClientInfo(action.info);
                        StatusStore.emitChange();
                        break;
                }
                return true;
        })
});
module.exports = StatusStore;
