import React from "react";
import Spotify from "spotify-web-api-js";
import SpotifyPlaylist from "./SpotifyPlaylist";
import StereodosePlaylist from "./StereodosePlaylist";

class UserProfile extends React.Component {

	constructor(props) {
		super(props);

		this.state = {
			spotifyPlaylists: null,
			stereodosePlaylists: null,
			categories: null,
			loading: true
		}

		this.checkPlaylists = this.checkPlaylists.bind(this);
	}

	render() {
		let {spotifyPlaylists, stereodosePlaylists, categories, loading} = this.state;
		if (spotifyPlaylists !== null && stereodosePlaylists !== null && !loading) {
			return (
				<div>
				<h2>Playlists Available From Spotify</h2>
				<table>
				<tbody>
				<tr>
					<th>Playlist Name</th>
					<th>Drug</th>
					<th>Mood</th>
				</tr>
				{spotifyPlaylists.map( (playlist) => {
					return <SpotifyPlaylist
							 	key={playlist.id} 
							 	categories={categories}
								playlist={playlist}
								onUpdate={ () => {this.checkPlaylists()}}
							   />
				})}
				</tbody>
				</table>

				<h2>Playlists Shared to Stereodose</h2>
				<table>
					<tbody>
						<tr>
							<th>Playlist Name</th>
							<th>Drug</th>
							<th>Mood</th>
						</tr>
					{stereodosePlaylists.map( (playlist) => {
						return <StereodosePlaylist 
									key={playlist.spotifyID} 
									playlist={playlist} 
									onUpdate={ () => {this.checkPlaylists()}}
								/>
					})}
					</tbody>
				</table>
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
		let resp = await fetch("/api/categories/", { credentials: "same-origin" });
		let categories = await resp.json();
		let state = this.state;
		state.categories = categories;
		this.setState(state);
		this.checkPlaylists();
	}

	async checkPlaylists() {
		let SDK = new Spotify();
		// TODO: catch errors here
		let token = await this.props.getAccessToken();
		SDK.setAccessToken(token);
		let userPlaylists = await SDK.getUserPlaylists();

		let response = await fetch("/api/playlists/me", {credentials: "same-origin"});
		let stereodosePlaylists = await response.json();

		let diffedSpotifyPlaylists = [];
		let diffedStereodosePlaylists = [];

		// so old school
		let spotifyPlaylists = userPlaylists.items;
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
		state.loading = false;
		this.setState(state);
	}
}

export default UserProfile