import React from 'react';
// import Spotify from 'spotify-web-api-js';

class UserPlaylist extends React.Component{
	spotifyid
	name
	add_playlist

	// constructor(props) {
	// 	super(props);
	// }

	addToStereodose() {
		this.props.add_playlist(this.props.spotifyid);
	}

	removeFromStereodose() {
		this.props.removeFromStereodose(this.props.spotifyid);
	}

	render() {
		return <li onClick={() => {this.addToStereodose()}}>{this.props.name}</li>
	}
}

export default UserPlaylist