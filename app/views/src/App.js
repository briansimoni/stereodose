import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import Main from './temp';

class App extends Component {
  render() {


	// Start: garbage temporary JS for POC purposes

	// this needs to log us in
	// if we have a cookie, call the /refresh endpoint
	// if not, window location to /auth/login
	// returns an Access Token
	let checkLoginStatus = function() {
		function getCookie(name) {
			var dc = document.cookie;
			var prefix = name + "=";
			var begin = dc.indexOf("; " + prefix);
			if (begin == -1) {
				begin = dc.indexOf(prefix);
				if (begin != 0) return null;
			}
			else
			{
				begin += 2;
				var end = document.cookie.indexOf(";", begin);
				if (end == -1) {
				end = dc.length;
				}
			}
			// because unescape has been deprecated, replaced with decodeURI
			//return unescape(dc.substring(begin + prefix.length, end));
			return decodeURI(dc.substring(begin + prefix.length, end));
		}
		let cookie = getCookie("_stereodose-session");
		if (!cookie) {
			// throw new Error("No cookie boi");
			window.location = "/auth/login";
			return;
		}

		return new Promise( (resolve, reject) => {
			let req = new XMLHttpRequest();
			req.open("GET", "/auth/refresh");
			req.addEventListener("readystatechange", function () {
				if (this.readyState === 4) {
					if (this.status === 200) {
						let data = JSON.parse(this.responseText);
						console.log(data);
						resolve(data.access_token);
					} else {
						reject(new Error("failed to get the token boi"));
					}
				}

			})
			req.send();
		});
	}
	checkLoginStatus().then(function(token) {
		console.log(token);

		let refreshButton = document.getElementById("refresh");
		refreshButton.addEventListener("click", function() {
			window.location.reload();
		});
        // let Token = {{.AccessToken}};
        window.onSpotifyWebPlaybackSDKReady = () => {
        //   const token = {{.AccessToken}};
          const player = new window.Spotify.Player({
            name: 'Web Playback SDK Quick Start Player',
			getOAuthToken: cb => { cb(token); }
          });
    
          // Error handling
          player.addListener('initialization_error', ({ message }) => { console.error(message); });
          player.addListener('authentication_error', ({ message }) => { console.error(message); });
          player.addListener('account_error', ({ message }) => { console.error(message); });
          player.addListener('playback_error', ({ message }) => { console.error(message); });
    
          // Playback status updates
          player.addListener('player_state_changed', state => { console.log(state); });
    
          // Ready
          player.addListener('ready', ({ device_id }) => {
            Main(token, device_id);
            console.log('Ready with Device ID', device_id);
          });
    
          // Not Ready
          player.addListener('not_ready', ({ device_id }) => {
            console.log('Device ID has gone offline', device_id);
          });
    
          // Connect to the player!
          player.connect();
        };
	})
	// END: garbage code


	
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>
      </div>
    );
  }
}

export default App;
