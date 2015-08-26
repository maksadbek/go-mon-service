var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var UserConstants = require('../constants/UserConstants');
var assign = require('object-assign');

var _clientInfo = {};
var _token = "";
var host = "217.29.118.23";
if(typeof(go_mon_host) !== "undefined"){
    host = go_mon_host;
}
var authURL = "http://"+host+":8080/signup";

var CHANGE_EVENT = 'change';

function setClientInfo(info){
    _clientInfo.fleet = info.fleet;
    _clientInfo.login = info.login;
    _clientInfo.groups = info.groups;
    _clientInfo.hash = info.hash;
    _clientInfo.uid = info.uid;
}

var UserStore = assign({}, EventEmitter.prototype, {
    clientInfo: _clientInfo,
    token: _token,
    auth: function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', encodeURI(authURL));
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.onload = function() {
            if (xhr.status === 200 ) {
                resp = JSON.parse(xhr.responseText)
                UserStore.token = resp.token;
                UserStore.emitChange();
            }
            else if (xhr.status !== 200) {
                resp = JSON.parse(xhr.responseText)
                console.error(resp.message);
                UserStore.emitChange();
                return UserStore.token;
            }
        };
        xhr.send(
            JSON.stringify({
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

module.exports = UserStore;
