var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var StatusConstants = require('../constants/StatusConstants');
var UserConstants = require('../constants/UserConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _carStatus = {};
var _clientInfo = {};
var _token = "";
var _markersOnMap = {};

function setClientInfo(info){
    _clientInfo.fleet = info.fleet;
    _clientInfo.login = info.login;
    _clientInfo.groups = info.groups;
    _clientInfo.hash = info.hash;
    _clientInfo.uid = info.uid;
}

var UserStore = assign({}, EventEmitter.prototype, {
    auth: function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', encodeURI("http://"+go_mon_host+":8080/signup"));
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.onload = function() {
                if (xhr.status === 200 ) {
                     _token = xhr.responseText;
                    UserStore.emitChange();
                }
                else if (xhr.status !== 200) {
                    UserStore.emitChange();
                    return _token;
                }
        };
        xhr.send(JSON.stringify({
                        user: _clientInfo.login,
                        hash: _clientInfo.hash,
                        uid: _clientInfo.uid
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
        updateMarker: function(info){
            if(_markersOnMap[info.id] !== undefined){
                _markersOnMap[info.id].latitude= info.latitude;
                _markersOnMap[info.id].longitude= info.longitude;
                _markersOnMap[info.id].direction= info.direction;
                _markersOnMap[info.id].speed= info.speed;
                _markersOnMap[info.id].sat= info.sat;
                _markersOnMap[info.id].owner= info.owner;
                _markersOnMap[info.id].formatted_time= info.time;
                _markersOnMap[info.id].addparams= info.additional;
                _markersOnMap[info.id].action= '1';
            }
        },
        redrawMap: function(){
            mon.obj_array(_markersOnMap, true);
        },
        sendAjax: function(){
                var xhr = new XMLHttpRequest();
                xhr.open('POST', encodeURI("http://"+go_mon_host+":8080/positions"));
                xhr.setRequestHeader('Content-Type', 'application/json');
                xhr.onload = function() {
                        if (xhr.status === 200 ) {
                            _carStatus = JSON.parse(xhr.responseText);
                            StatusStore.emitChange();
                        }
                        else if (xhr.status !== 200) {
                            StatusStore.emitChange();
                            return _carStatus;
                        }
                };
                xhr.send(JSON.stringify({
                        selectedFleetJs: _clientInfo.fleet,
                        user: _clientInfo.login,
                        groups: _clientInfo.groups,
                        token: _token,
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
                    case StatusConstants.AddMarker:
                        // the structure of info must be:
                        // { id: "1234", pos: { lat: "123", lng:...}}
                        _markersOnMap[action.info.id] = action.info.pos;
                        mon.obj_array(_markersOnMap, true);
                        my_sm.push(action.info.id);
                        break;
                    case StatusConstants.DelMarker:
                        _markersOnMap[action.info.id].action = '-1';
                        mon.obj_array(_markersOnMap, true);
                        for(var i in my_sm){
                            if(my_sm[i] == action.info.id){
                                my_sm.pop(i);
                            }
                        }
                        break;
                }
                return true;
        })
});
module.exports = {
            StatusStore: StatusStore, 
            UserStore: UserStore
};
