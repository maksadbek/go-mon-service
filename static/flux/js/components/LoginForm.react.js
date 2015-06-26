var React = require('react');
var UserActions = require('../actions/UserActions');

var LoginForm = React.createClass({
    render: function(){
        return (
            <div>
                <label>login</label>
                <input ref="login" className="input-mini" type="text" />

                <label>fleet</label>
                <input  ref="fleet" className="input-mini" type="text"  />

                <label>hash</label>
                <input ref="hash" className="input-mini" type="text"  />

                <label>groups</label>
                <input ref="groups" className="input-mini" type="text" />

                <button onClick={this._onSubmit} className="btn btn-primary">Save changes</button>
                <button onClick={this._onCancel} className="btn">Cancel</button>

            </div>    
        );
    },
    _onSubmit: function(){
        console.log(React.findDOMNode(this.refs.login));
        UserActions.Authenticate({
                            login:React.findDOMNode(this.refs.login).value,
                            hash:React.findDOMNode(this.refs.hash).value,
                            fleet:React.findDOMNode(this.refs.fleet).value,
                            groups:React.findDOMNode(this.refs.groups).value
                            });
    },
    
    _onCancel: function(){
        this.refs.login = "";
        this.refs.hash= "";
        this.refs.fleet= "";
        this.refs.groups= "";
    },
});

module.exports = LoginForm;
