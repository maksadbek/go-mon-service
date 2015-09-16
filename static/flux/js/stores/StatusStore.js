var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var StatusConstants = require('../constants/StatusConstants');
var UserConstants = require('../constants/UserConstants');
var UserStore = require('./UserStore');
var assign = require('object-assign');
var lunr = require('lunr');

var CHANGE_EVENT = 'change';
var UNCHECK_EVENT = 'uncheck_event';

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

var searchIdx = lunr(function(){
    this.field('number');
    this.ref('id');
});

var indexed = false;
var isFirstMarkerDrawed = false;
var StatusStore = assign({}, EventEmitter.prototype, {
    groupNames: ["все"],
    groupIndex: 0,
    centerMarker: function(id){
        if(_markersOnMap[id].onMap){
            mon.setCenterObj(id);
        }
    },
    uncheckAllMarkers: function(){
        StatusStore.emitUncheckChange();
    },
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
    redrawMap: function(zoom){
        // mon is global object
        // can be used to control the Map
        // mapLoaded is global variable, true if the google maps was loaded
        if(typeof(mon) !== "undefined" && mapLoaded){
            mon.obj_array(_markersOnMap, zoom);
        }
        
    },
    sendAjax: function(){
        if(MonReqToggler !== 1){
            return _carStatus;
        }
        var xhr = new XMLHttpRequest();
        xhr.open('POST', encodeURI(positionURL));
        xhr.setRequestHeader('Content-Type','application/json');
        xhr.onload = function() {
            if (xhr.status === 200 ) {
                _carStatus = JSON.parse(xhr.responseText);
                // if search index container is empty, 
                // then fill it and groups container by the way
                if(!indexed){
                    window.uncheckAllMarkers = StatusStore.uncheckAllMarkers;
                    _carStatus.update.forEach(function(group){
                        StatusStore.groupNames.push(group.groupName);
                        group.data.forEach(function(v){
                            searchIdx.add({
                                id: v.id, 
                                number: v.number
                            });
                        });
                    });
                    indexed = true;
                    // TODO, this is for test
                    // _markersOnMap[Object.keys(_markersOnMap)[0]].action = '2';
                    // StatusStore.redrawMap(false);
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
        xhr.setRequestHeader("X-Access-Token", UserStore.token);
        xhr.send(JSON.stringify({
            fleetID: go_mon_fleet, // TODO use UserStore.clientInfo.fleet,
            userName: UserStore.clientInfo.login,
            })
        );
    },
    getAll: function(){
        // if user filtered
        if(_search){
            var filteredData = [];     // list of groups and its values
            var foundCar;   // car with required criteria
            // iterate over all found items
            _carStatus.update.forEach(function(group){
                // iterate over all items in the group
                var res = {groupName: group.groupName, data: []}; // result
                group.data.forEach(function(car){
                    _searchRes.forEach(function(foundRef){
                        if(car.id === parseInt(foundRef.ref)){
                            res.data.push(car);
                        }
                    });
                })
                if(res.data.length !== 0){
                    filteredData.push(res);
                }
            });
            return {
                 // replace values of car list with a list that found items
                 id: _carStatus.id,
                 update: filteredData 
            }
        }

        // if user filtered by group, the return only that group
        if(StatusStore.groupIndex !== 0){
            var filteredStatuses = [];
            filteredStatuses.push(_carStatus.update[StatusStore.groupIndex]);
            return {
                 id: _carStatus.id,
                 update: filteredStatuses
            }
        }
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
    // uncheck listeners
    emitUncheckChange: function(){
        this.emit(UNCHECK_EVENT);
    },
    addUncheckListener: function(callback){
        this.on(UNCHECK_EVENT, callback);
    },
    removeUncheckListener: function(callback){
        this.removeListener(UNCHECK_EVENT, callback);
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
                _markersOnMap[action.info.id] = {
                    id: action.info.id,
                    latitude: action.info.stat.latitude,
                    longitude: action.info.stat.longitude,
                    direction: action.info.stat.direction,
                    speed: action.info.stat.speed,
                    sat: action.info.stat.sat,
                    owner: action.info.stat.owner,
                    formatted_time: action.info.stat.time,
                    addparams: action.info.stat.additional,
                    car_name: action.info.stat.number,
                    action: '2',
                    onMap: true
                }
                // check for existence of id in my_sm
                // if it does not exist the push it
                id = parseInt(action.info.id);
                if(my_sm.indexOf(id) === -1){
                    my_sm.push(id);
                }
                // pass true to zoom the map
                StatusStore.redrawMap(true);
                break;
            case StatusConstants.DelMarker:
                _markersOnMap[action.info.id].action = '-1';
                _markersOnMap[action.info.id].onMap = false;
                // remove car id from my_sm array
                my_sm.remove(parseInt(action.info.id))
                StatusStore.redrawMap(true);
                break;
            case StatusConstants.SearchCar:
                var number = action.info.name;
                _searchRes = searchIdx.search(number);
                _search = true;
                StatusStore.emitChange();
                break;
            case StatusConstants.DelSearchCon:
                _search = false;
                StatusStore.emitChange();
                break;
            case StatusConstants.SelectGroup:
                StatusStore.groupIndex = action.info.id;
                StatusStore.emitChange();
                break;
        }
        return true;
    })
});
module.exports =  StatusStore;
