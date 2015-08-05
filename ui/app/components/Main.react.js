var React = require('react');
var StatusStore = require('../stores/StatusStore').StatusStore;
var CarActions = require('../actions/StatusActions');
var UserActions = require('../actions/UserActions');
var Sidebar = require('./Sidebar.react');
var UserStore = require('../stores/StatusStore').UserStore;
var Status = require('./CarStatus.react');
var Mui  = require('material-ui');
var ThemeManager = new Mui.Styles.ThemeManager();
mui = require('material-ui')

var AppBar = Mui.AppBar,
    MenuItem= Mui.MenuItem, 
    IconButton= Mui.IconButton, 
    List  = Mui.List,
    LeftNav= Mui.LeftNav;


menuItems = [
    { 
        type: MenuItem.Types.SUBHEADER, 
        text: 'Resources' 
    },
    { 
        type: MenuItem.Types.LINK, 
        payload: 'https://github.com', 
        text: 'GitHub' 
    },
    { 
        text: 'Disabled', 
        disabled: true 
    },
    { 
        type: MenuItem.Types.LINK, 
        payload: 'https://www.google.com', 
        text: 'Disabled Link', 
        disabled: true 
    },
];


function getAllStatuses(){
    return StatusStore.getAll()
}

var StatusApp = React.createClass({
    childContextTypes: {
          muiTheme: React.PropTypes.object
    },
    getChildContext: function() {
        return {
            muiTheme: ThemeManager.getCurrentTheme()
        };
    },
    getInitialState: function(){
        return {
            stats: {
                id: '',
                update: {"":[]},
                last_request: null
            },
            isChildChecked: false
        }
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
        UserStore.addChangeListener(this._onAuth);
    },

    componentWillMount: function(){
       UserActions.Auth({
           login: "taxi",
           uid: "taxi",
           hash: "b5ea8985533defbf1d08d5ed2ac8fe9b",
           fleet: "436",
           groups: "1,2,3" // TODO ochirib tashlash
       });
    },
    componentWillUnmount: function(){
        StatusStore.removeChangeListener(this._onChange);
        UserStore.removeChangeListener(this._onAuth);
    },
    toggleLeftNav: function(){
        console.log(this.refs.leftNav);
        React.findDOMNode(this.refs.leftNav).toggle()
    },
    render: function(){
        var content = [];
        var update = this.state.stats.update;
        var checked = this.state.isChildChecked;
        for(var i in update){
            content.push(<Sidebar key={i} groupName={i} stats={update[i]}/>)
        }
        return (   
            <div>
                <LeftNav ref="leftNav" docked={false} menuItems={menuItems} />
                <AppBar
                    onLeftIconButtonTouchTap={this.toggleLeftNav}
                    title="Wherepo"
                    iconElementLeft={<IconButton></IconButton>}
                />
                <List>
                    {content}
                </List>
            </div>
            )
    },
    _onChange: function(){
        this.setState({stats: getAllStatuses()});
    },
    _onAuth: function(){
        StatusStore.sendAjax();
        setInterval(function(){
            StatusStore.sendAjax();
        }, 5000);
    }
});

module.exports = StatusApp;
