import React, { Component, Fragment } from 'react';
import './App.css';
import WebPlaybackReact from './Spotify/WebPlaybackReact';
import Playlists from './screens/Playlists';
import MySpotifyPlaylists from './screens/MySpotifyPlaylists';

import Login from './Spotify/Login.js';

// import Header from './layout/Header.js';
// import Footer from './layout/Footer.js';

import IntroScreen from './screens/Intro.js';
import NowPlayingScreen from './screens/NowPlaying.js';

export default class App extends Component {
	state = {
	  // User's session credentials
	  userDeviceId: null,
	  userAccessToken: null,
  
	  // Player state
	  playerLoaded: false,
	  playerSelected: false,
	  playerState: null
	}
  
	componentWillMount() {
		Login().then( (token) => {
			this.onSuccessfulAuthorization(token);
		})
	//   LoginCallback({
	// 	onSuccessfulAuthorization: this.onSuccessfulAuthorization.bind(this),
	// 	onAccessTokenExpiration: this.onAccessTokenExpiration.bind(this)
	//   });
	}
	
	onSuccessfulAuthorization(accessToken) {
	  this.setState({
		userAccessToken: accessToken
	  });
	}
	
	onAccessTokenExpiration() {
	  this.setState({
		userDeviceId: null,
		userAccessToken: null,
		playerLoaded: false,
		playerSelected: false,
		playerState: null
	  });
  
	  console.error("The user access token has expired.");
	}
	
	render() {
	  let {
		userDeviceId,
		userAccessToken,
		playerLoaded,
		playerSelected,
		playerState
	  } = this.state;
	  
	  let webPlaybackSdkProps = {
		playerName: "Spotify React Player",
		playerInitialVolume: 1.0,
		playerRefreshRateMs: 100,
		playerAutoConnect: true,
		onPlayerRequestAccessToken: (() => userAccessToken),
		onPlayerLoading: (() => this.setState({ playerLoaded: true })),
		onPlayerWaitingForDevice: (data => this.setState({ playerSelected: false, userDeviceId: data.device_id })),
		onPlayerDeviceSelected: (() => this.setState({ playerSelected: true })),
		onPlayerStateChange: (playerState => this.setState({ playerState: playerState })),
		onPlayerError: (playerError => console.error(playerError))
	  };
	  
	  return (
		<div className="App">
		  {/* <Header /> */}
		  <main>
			{!userAccessToken && <IntroScreen />}
			{userAccessToken &&
			  <WebPlaybackReact {...webPlaybackSdkProps}>
			  <h4>my playlists!</h4>
		  		<MySpotifyPlaylists access_token={userAccessToken}></MySpotifyPlaylists>

					<h4>Playlists on Stereodose</h4>
					<Playlists access_token={userAccessToken} category="Weed" subcategory="Chill" device_id={userDeviceId}></Playlists>
				{!playerLoaded &&
				  <h2 className="action-orange">Loading Player</h2>
				}
  
				{playerLoaded && !playerSelected && 
				  <Fragment>
					<h2 className="action-green">Loading Player</h2>
					<h2 className="action-orange">Waiting for device to be selected</h2>
				  </Fragment>
				}
  
				{playerLoaded && playerSelected && !playerState &&
				  <Fragment>
					{/* <h2 className="action-green">Loading Player</h2> */}
					{/* <h2 className="action-green">Waiting for device to be selected</h2> */}
					{/* <h2 className="action-orange">Start playing music ...</h2> */}
				  </Fragment>
				}
  
				{playerLoaded && playerSelected && playerState &&
				  <Fragment>
					<h2 className="action-green">Loading Player</h2>
					<h2 className="action-green">Waiting for device to be selected</h2>
					<h2 className="action-green">Start playing music!</h2>
					<NowPlayingScreen playerState={playerState} />
				  </Fragment>
				}
			  </WebPlaybackReact>
			}
		  </main>
  
		  {/* <Footer /> */}
		</div>
	  );
	}
  };