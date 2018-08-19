import React from 'react'

class Track extends React.Component {
	name
	context
	URI
	accessToken
	deviceID

	constructor(props) {
		super(props);

		this.playTrack = this.playTrack.bind(this);
	}
	
	render() {
		return <li onClick={ () => { this.playTrack() } } >{this.props.name}</li>
	}

	playTrack() {
		// let SDK = new Spotify();
		// SDK.setAccessToken(this.props.accessToken);
		// SDK.play({
		// 	deviceID: this.props.deviceID,
		// 	context: this.props.context,
		// 	offset : {
		// 		uri: this.props.URI
		// 	}
		// })

		let data = {
			"context_uri": `${this.props.context}`,
			"offset": {
				"uri": this.props.URI
			}
		}

		fetch(`https://api.spotify.com/v1/me/player/play?device_id=${this.props.deviceID}`,{
			
			method: "PUT",
			headers : {
				"Authorization": `Bearer ${this.props.accessToken}`,
				"Content-Type": "application/json"
			},
			body: JSON.stringify(data)
		})
		.then( (result) => {
			console.log("play track result");
			console.log(result);
		})
		.catch( (error) => {
			console.warn(error);
		})
	}
}

export default Track;