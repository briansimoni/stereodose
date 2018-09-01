import React from "react";
import Spotify from "spotify-web-api-js";

class UserProfile extends React.Component {

	constructor(props) {
		super(props);

		this.state = {
			spotifyPlaylists: null,
			stereodosePlaylists: null
		}

		this.checkPlaylists = this.checkPlaylists.bind(this);
	}

	render() {
		let {spotifyPlaylists, stereodosePlaylists} = this.state;
		if (spotifyPlaylists !== null && stereodosePlaylists !== null) {
			return (
				<div>
				<h2>Playlists Available From Spotify</h2>
				{spotifyPlaylists.map( (playlist) => {
					return <li key={playlist.id}>{playlist.name}</li>
				})}

				<h2>Playlists Shared to Stereodose</h2>
				{stereodosePlaylists.map( (playlist) => {
					return <li key={playlist.spotifyID}>{playlist.name}</li>
				})}
			</div>
			)
		}
		return (
			<div>
				<h2>Playlists From Spotify</h2>

				<h2>Playlists Shared to Stereodose</h2>
			</div>
		);
	}

	async componentDidMount() {
		this.checkPlaylists();
	}

	async checkPlaylists() {
		let SDK = new Spotify();
		let token = await this.props.getAccessToken();
		SDK.setAccessToken(token);
		let userPlaylists = await SDK.getUserPlaylists();

		let response = await fetch("/api/playlists/me", {credentials: "same-origin"});
		let stereodosePlaylists = await response.json();

		let diffedSpotifyPlaylists = [];
		let diffedStereodosePlaylists = [];

		// so old school
		let spotifyPlaylists = userPlaylists.items;
		console.log(spotifyPlaylists);
		for (let i = 0; i < spotifyPlaylists.length; i++) {
			let match = false;
			for (let j = 0; j < stereodosePlaylists.length; j++) {
				if (spotifyPlaylists[i].id === stereodosePlaylists[j].spotifyID) {
					diffedStereodosePlaylists.push(stereodosePlaylists[j]);
					match = true;
					break;
				}
			}
			if (match === false) {
				diffedSpotifyPlaylists.push(spotifyPlaylists[i]);
			}
		}
		
		let state = this.state;
		state.spotifyPlaylists = diffedSpotifyPlaylists;
		state.stereodosePlaylists = diffedStereodosePlaylists;
		this.setState(state);
	}
}

export default UserProfile