import React from "react";

class UserStatusIndicator extends React.Component{
	constructor(props) {
		super(props);

		this.state = {
			loggedIn: null,
			username: null,
			image: null
		};

		let loggedIn = this.checkSessionCookie();
		this.state.loggedIn = loggedIn;
		this.props.isUserLoggedIn(loggedIn);
	}

	render() {
		if (this.state.loggedIn === null) {
			return <div><p>loading</p></div>
		}

		if (this.state.loggedIn === false) {
			return (
				<div>
					<h3 onClick={ () => {this.logIn()}}>Log In With Spotify</h3>
				</div>
			)
		}

		if (this.state.loggedIn === true) {
			return <div><p>I'm logged in!</p></div>
		}
	}

	logIn() {
		window.location = "/auth/login";
	}

	checkSessionCookie() {
		// stolen from Stack Overflow
		function getCookie(name) {
			var dc = document.cookie;
			var prefix = name + "=";
			var begin = dc.indexOf("; " + prefix);
			if (begin === -1) {
				begin = dc.indexOf(prefix);
				if (begin !== 0) return null;
			}
			else {
				begin += 2;
				var end = document.cookie.indexOf(";", begin);
				if (end === -1) {
					end = dc.length;
				}
			}
			// because unescape has been deprecated, replaced with decodeURI
			//return unescape(dc.substring(begin + prefix.length, end));
			return decodeURI(dc.substring(begin + prefix.length, end));
		}
	
		let cookie = getCookie("_stereodose-session");
		if (!cookie) {
			return false
		}

		return true;
	}
}

export default UserStatusIndicator;