import React from 'react';
import Drugs from './dev/Drugs';
import Drug from './dev/Drug';
import Playlists from './dev/Playlists';
import Playlist from './dev/Playlist';
import {BrowserRouter as Router, Route} from 'react-router-dom';
import oauth2 from 'simple-oauth2';


class App extends React.Component {

	accessToken

	constructor(props) {
		super(props);

		this.state = {accessToken: null}

		// TODO: figure out how an arrow function could eliminate this line
		this.getAccessToken = this.getAccessToken.bind(this);
		this.doSomething = this.doSomething.bind(this);
	}

	doSomething() {
		console.log("SOMETHING");
	}

	render() {
		return (
			<div>
				<h1 onClick={ () => {this.getAccessToken()} }>Header</h1>
				<Router>
					<div>
						<Route exact path="/" component={Drugs} />
						<Route exact path="/:drug" component={Drug} />
						<Route exact path="/:drug/:subcategory" component={Playlists} />
						{/* <Route path="/:drug/:subcategory/:playlist" component={Playlist} /> */}
						<Route 
							path="/:drug/:subcategory/:playlist"
							render={(props) => <Playlist {...props} getAccessToken={ () => {this.getAccessToken()}} />}
						/>
						
					</div>
				</Router>
			</div>
		)
	}

	// getAccessToken will return a Promise to either get the access token or will redirect
	// the user to Login
	// Should be able to pass this function around as a prop to components that need a token
	// i.e. <Player> and <Playlist>
	async getAccessToken() {
		let token = this.state.accessToken;
		if (token && !token.expired()) {
			return token;
		}
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
			window.location = "/auth/login";
			return;
		}

		let call = new Promise((resolve, reject) => {
			let req = new XMLHttpRequest();
			req.open("GET", "/auth/refresh");
			req.addEventListener("readystatechange", function () {
				if (this.readyState === 4) {
					if (this.status === 200) {
						let data = JSON.parse(this.responseText);
						resolve(data);
					} else {
						reject(new Error("failed to get the token boi " + this.responseText));
					}
				}

			})
			req.send();
		});

		let data = await call;

		// Using the npm library only for the convenient helper AccessToken class
		// Should probably write the logic to track expiration myself and remove this dependency
		console.log(data);
		const fake = {
			client: {
			  id: '<client-id>',
			  secret: '<client-secret>'
			},
			auth: {
			  tokenHost: 'https://api.oauth.com'
			}
		  };
		let oauth = oauth2.create(fake);

		token = oauth.accessToken.create(data);
		
		this.accessToken = token;
		return token;
	}
}

export default App;