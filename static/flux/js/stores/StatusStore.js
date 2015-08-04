var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var StatusConstants = require('../constants/StatusConstants');
var UserConstants = require('../constants/UserConstants');
var UserStore = require('./UserStore');
var assign = require('object-assign');
var _ = require('lodash');

var CHANGE_EVENT = 'change';

var _carStatus = {};
var _markersOnMap = {};
var _search = false;
var _searchCase = [];
var _searchRes;

var host = "217.29.118.23";
if(typeof(go_mon_host) !== "undefined"){
    host = go_mon_host;
}
var positionURL = "http://"+host+":8080/positions";

var StatusStore = assign({}, EventEmitter.prototype, {
    groupNames: ["all"],
    groupIndex: 0,
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
        // mon is global object
        // can be used to control the Map
        if(typeof(mon) !== "undefined"){
            mon.obj_array(_markersOnMap, true);
        }
        
    },
    sendAjax: function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', encodeURI(positionURL));
        xhr.setRequestHeader('Content-Type','application/json');
        xhr.onload = function() {
            if (xhr.status === 200 ) {
                _carStatus = JSON.parse(xhr.responseText);
                // if search is on, then filter incoming data 
                // by criteria from _searchRes
                if(_search){
                    var res = [];
                    var m = {};
                    foundCar = _carStatus.update[_searchRes.group][_searchRes.id];
                    res.push(foundCar);
                    m[_searchRes.group] = res;
                    _carStatus.update = m;
                }
                // if search index container is empty, 
                // then fill it and groups container by the way
                if(_searchCase.length === 0){
                    for(var groupName in _carStatus.update){
                        StatusStore.groupNames.push(groupName);
                        _carStatus.update[groupName]
                        .forEach(function(v, index){
                            _searchCase.push({
                                group: groupName,
                                id: index, 
                                name: v.name,
                                number: v.number
                            });
                        });
                    }
                }
                if(StatusStore.groupIndex !== 0){
                    var groupName = StatusStore.groupNames[StatusStore.groupIndex];
                    var filteredStatuses = {};
                    filteredStatuses[groupName] = _carStatus.update[groupName];
                    _carStatus.update = filteredStatuses;
                }
                StatusStore.emitChange();
                return _carStatus;
            }
            else if (xhr.status !== 200) {
                StatusStore.emitChange();
                return _carStatus;
            }
            StatusStore.emitChange();
            return _carStatus;
        };
        xhr.send(JSON.stringify({
            selectedFleetJs: UserStore.clientInfo.fleet,
            user: UserStore.clientInfo.login,
            groups: UserStore.clientInfo.groups,
            token: UserStore.token,
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
                for(var i in my_sm){
                    if(my_sm[i] === action.info.id){
                        return;
                    }
                }
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
            case StatusConstants.SearchCar:
                var number = action.info.name;
                _searchRes = _.find(_searchCase, {'number': number});
                _search = true;
                break;
            case StatusConstants.DelSearchCon:
                _search = false;
                break;
            case StatusConstants.SelectGroup:
                console.log("dispatch", action.info);
                StatusStore.groupIndex = action.info.id;
                break;
        }
        return true;
    })
});
module.exports =  StatusStore;
