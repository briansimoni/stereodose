import React from "react";
import Playlist from './Playlist';
import Track from './Track'
// import Spotify from 'spotify-web-api-js';

class Playlists extends React.Component{
	device_id
	category
	subcategory
	access_token

	constructor(props) {
		super(props);
		this.state = {playlists: [{name: "poop"}], playlistData: null};

		this.getSongs = this.getSongs.bind(this);
	}

	render() {
		// for now...
		// lets pass down a function to playlist to get songs
		// when you click a playlist, return the songs up to this component's state
		// if the songs are not null, hide the playlists, render the songs and a back button
		let playlists = this.state.playlists;
		let playlistData = this.state.playlistData;
		if(playlists !== null) {
			if (playlistData !== null) {
				console.log(this.state.playlistData);
				console.log(this.state.playlistData !== null);
				return (
					<ul>
						<li>:)</li>
						{playlistData.Tracks.map(track => (
							<Track name={track.Name}></Track>
						))}
					</ul>
				);
			}
			return (
				<ul>
				{playlists.map(playlist => (
					<Playlist
						device_id={this.props.device_id}
						access_token={this.props.access_token}
						playlist_id={playlist.SpotifyID} 
						key={playlist.SpotifyID}
						name={playlist.name}
						getSongs={this.getSongs}>
					</Playlist>
				))}
				</ul>
			);
		}
		return <span>Loading</span>
	}

	componentDidMount() {
		fetch(`/api/playlists/?category=${this.props.category}&subcategory=${this.props.subcategory}`)
		.then(res => res.json())
		.then(
			(playlists) => {
				this.setState({
					playlists: playlists
				})
			},
			(error) => {
				alert(error);
			}
		);
	}

	getSongs(playlist_id) {
		fetch(`/api/playlists/${playlist_id}`)
		.then(res => res.json())
		.then(
			(playlistData) => {
				this.setState({playlistData: playlistData});
			},
			(error) => {
				alert(error);
			}
		);
	}


}

export default Playlists