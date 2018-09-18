import React from 'react';
import Drugs from './dev/Drugs';
import Drug from './dev/Drug';
import Playlists from './dev/Playlists';
import Playlist from './dev/Playlist';
import Player from './Player';
import {HashRouter, Route, Switch} from 'react-router-dom';
import UserStatusIndicator from './User/StatusIndicator';
import UserProfile from './User/Profile';
import Header from './layout/Header';


class App extends React.Component {

	accessToken
	deviceIDPromise
	deviceIDResolver

	// loggedInPromise resolves at a later time with the user's logged in status (true/false)
	loggedInPromise
	// callback function that can be called to resolve when we know that the user is logged in or not
	loggedInPromiseResolver

	constructor(props) {
		super(props);

		this.state = {
			accessToken: null,
			loggedIn: false
		}

		this.deviceIDPromise = new Promise( (resolve, reject) => {
			resolve = resolve.bind(this);
			this.deviceIDResolver = resolve;
		});

		this.loggedInPromise = new Promise( (resolve) => {
				this.loggedInPromiseResolver = resolve;
		})
		

		// TODO: figure out how an arrow function could eliminate this line
		this.getAccessToken = this.getAccessToken.bind(this);
		this.setDeviceID = this.setDeviceID.bind(this);
		this.isUserLoggedIn = this.isUserLoggedIn.bind(this);
	}

	isUserLoggedIn(loggedIn) {
		let state = this.state;
		state.loggedIn = loggedIn;
		this.setState(state);
	}

	render() {
		return (
				<HashRouter>
					<div>
						<Header>
						<Route 
							path="/" 
							render={ (props) => 
								<UserStatusIndicator
									{...props}
									isUserLoggedIn={ (loggedIn) => this.loggedInPromiseResolver(loggedIn)}
								/>
						}/>
						</Header>

						<main role="main" className="container">
						{/* Routes wrapped in a Switch match only the first route for ambiguous matches*/}
						<Switch>
							<Route exact path="/profile"
								render={ (props)=>
									<UserProfile
										{...props}
										getAccessToken={ ()=> this.getAccessToken()}
									/>
								}
							/>

							<Route exact path="/" component={Drugs} />
							<Route exact path="/:drug" component={Drug} />
							<Route exact path="/:drug/:subcategory" component={Playlists} />
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
						</Switch>
						</main>
						
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
					</div>
				</HashRouter>
		)
	}

	// pass setDeviceID to the player component so we can lift "state" up
	// and then move it over to peers
	setDeviceID(deviceID) {
		this.deviceIDResolver(deviceID);
	}

	// getAccessToken will return a Promise to either get the access token
	// Should be able to pass this function around as a prop to components that need a token
	// i.e. <Player> and <Playlist>
	async getAccessToken() {
		let loggedIn = await this.loggedInPromise;
		if (loggedIn === false) {
			throw new Error("The user is not logged in");
		}

		let response =  await fetch("/auth/token", {credentials: "same-origin"});
		let token = await response.json();
		let expiresOn = token.expiry;
		let now = new Date();
		let expiresDate = new Date(expiresOn);
		if(now < expiresDate) {
			this.accessToken = token.access_token;
			return token.access_token;
		}
		response = await fetch("/auth/refresh", { credentials: "same-origin" });
		token = await response.json();
		this.accessToken = token.access_token;
		return token.access_token;

	}
}

export default App;