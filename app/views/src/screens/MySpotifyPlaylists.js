import React from 'react';
import Spotify from 'spotify-web-api-js';
import Playlist from './UserPlaylist';

class MySpotifyPlaylists extends React.Component {
	access_token

	constructor(props) {
		super(props);
		this.state = {
			error: null,
			isLoaded: false,
			playlists: []
		  };
	}

	componentDidMount() {
		let sdk = new Spotify();
		sdk.setAccessToken(this.props.access_token);
		sdk.getUserPlaylists()
		.then( (playlists) => {
			this.setState({
				playlists: playlists.items,
				isLoaded: true
			})
		}).catch( (err) => {
			this.setState({error: err});
		})
	}

	addPlaylistToStereodose(playlistID) {
		let data = {
			SpotifyID: playlistID,
			Category: "Weed",
			SubCategory: "Chill"
		}
		console.log(data);
		let ID = playlistID;
		fetch("/api/playlists/", {
			method: "POST",
			body: JSON.stringify({
				SpotifyID: ID,
				Category: "Weed",
				SubCategory: "Chill"
			})
		}).then((res) => {
			console.log("Made fetch call to POST playlist without error");
			console.log(res);
		})
	}

	render() {
		const { error, isLoaded, playlists } = this.state;
		if (error) {
			return <div>Error: {error.message}</div>;
		} else if (!isLoaded) {
			return <div>Loading...</div>;
		} else {
			return (
				<ul>
				{playlists.map(playlist => (
					<Playlist 
						add_playlist={this.addPlaylistToStereodose}
						spotifyid={playlist.id}
						key={playlist.id}
						name={playlist.name}>
					{playlist.name}
					</Playlist>
				))}
				</ul>
			);
		}
	}
}

export default MySpotifyPlaylists