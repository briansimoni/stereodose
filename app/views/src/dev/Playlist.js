import React from "react";

class Playlist extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			loading: true,
			playlist: null,
			error: null
		};
	}

	render() {
		let {loading, playlist, error} = this.state;
		if (loading) {
			return <h3>Loading</h3>
		}
		if (error) {
			console.log(error);
			return <h3>eeeeeee</h3>
		}
		if (playlist) {
			return (
				<div>
					{playlist.tracks.map( (track) => {
						return (
							<li 
							key={track.spotifyID} onClick={() => this.playSong(playlist.URI, track.URI)}>
								{track.name}
							</li>
						)
					})}
				</div>
			)
		}
	}

	playSong(context, uri) {
		let data = {
			"context_uri": context,
			"offset": {
				"uri": uri
			}
		}

		this.props.getDeviceID().then( (deviceID) => {
			this.props.getAccessToken().then( (accessToken) => {
				fetch(`https://api.spotify.com/v1/me/player/play?device_id=${deviceID}`,{
			
					method: "PUT",
					headers : {
						"Authorization": `Bearer ${accessToken.token.access_token}`,
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
			})
		})
	}

	componentDidMount() {

		let playlistID = this.props.match.params.playlist
		fetch(`/api/playlists/${playlistID}`)
		.then( (response) => {
			console.log(response);
			return response.json();
		})
		.then( (json) => {
			this.setState({
				loading: false,
				playlist: json
			});
		})
		.catch( (err) => {
			console.log(err);
			this.setState({
				loading: false,
				error: err
			});
		})
	}
}

// lets have this component have some function for getting the access token
// and have the player nested in this component

export default Playlist;