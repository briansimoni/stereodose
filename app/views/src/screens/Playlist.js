import React from 'react';

class Playlist extends React.Component {

	constructor(props) {
		super(props);
		this.state = {test: "omg"};
	}
	name
	playlist_id
	device_id
	access_token
	// getSongs is passed in by the parent, and triggers the Tracks to render
	getSongs


	render() {
		return <li onClick={() => {this.props.getSongs(this.props.playlist_id)}}>{this.props.name}</li>
	}
}

export default Playlist