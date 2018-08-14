import React from 'react';
import Spotify from 'spotify-web-api-js';

class Track extends React.Component {
	name
	track_id
	device_id
	access_token

	play() {
		let device_id = this.props.device_id;
		let track_id = this.props.track_id;
		
		let SDK = new Spotify();
		SDK.setAccessToken(this.props.access_token);

		let options = {
			device_id: device_id,
			context_uri: playlistID,
			offset: {
				uri: "spotify:track:" + trackID
			}
		}
		SDK.play(options)
		.then( (data) => {
			console.log(data);
			console.log("should be playing");
		})
		.catch( (err) => {
			console.warn(err);
		})
	}

	render() {
		return <li onClick={() => {this.play()}}>{this.props.name}</li>
	}
}

export default Track