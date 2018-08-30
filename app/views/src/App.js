import React from 'react';
import Drugs from './dev/Drugs';
import Drug from './dev/Drug';
import Playlists from './dev/Playlists';
import Playlist from './dev/Playlist';
import Player from './Player';
import {BrowserRouter as Router, Route} from 'react-router-dom';


class App extends React.Component {

	accessToken
	deviceIDPromise
	deviceIDResolver

	constructor(props) {
		super(props);

		this.state = {accessToken: null}

		this.deviceIDPromise =  new Promise( (resolve, reject) => {
			resolve = resolve.bind(this);
			this.deviceIDResolver = resolve;
		})
		

		// TODO: figure out how an arrow function could eliminate this line
		this.getAccessToken = this.getAccessToken.bind(this);
		this.setDeviceID = this.setDeviceID.bind(this);
	}

	render() {
		return (
			<div>
				<h1 onClick={ () => {this.getAccessToken()} }>Header</h1>
				<Router>
					<div>
						<Route 
							path="/" 
							render={ (props) => 
								<Player 
								{...props} 
								getAccessToken={ () => this.getAccessToken()}
								setDeviceID={(deviceID) => this.setDeviceID(deviceID)}>
								</Player>
							}
						/>
						<Route exact path="/" component={Drugs} />
						<Route exact path="/:drug" component={Drug} />
						<Route exact path="/:drug/:subcategory" component={Playlists} />
						{/* <Route path="/:drug/:subcategory/:playlist" component={Playlist} /> */}
						<Route 
							path="/:drug/:subcategory/:playlist"
							render={(props) => 
							<Playlist
							{...props} 
							getAccessToken={ () => this.getAccessToken()} 
							getDeviceID={ () => this.deviceIDPromise }
							/>
						}
						/>
						
					</div>
				</Router>
			</div>
		)
	}

	// pass setDeviceID to the player component so we can lift "state" up
	// and then move it over to peers
	setDeviceID(deviceID) {
		this.deviceIDResolver(deviceID);
	}

	// getAccessToken will return a Promise to either get the access token or will redirect
	// the user to Login
	// Should be able to pass this function around as a prop to components that need a token
	// i.e. <Player> and <Playlist>
	async getAccessToken() {
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

		try {
			let response =  await fetch("/auth/token");
			let token = await response.json();
			let expiresOn = token.expiry;
			let now = new Date();
			let expiresDate = new Date(expiresOn);
			if(now < expiresDate) {
				this.accessToken = token.access_token;
				return token.access_token;
			}
			response = await fetch("/auth/refresh");
			token = await response.json();
			this.accessToken = token.access_token;
			return token.access_token;
		} catch(err) {
			return err;
		}
	}
}

export default App;