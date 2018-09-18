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
			return <h3>{error.message}</h3>
		}
		if (playlist) {
			return (
				<ul className="list-group">
					{playlist.tracks.map( (track) => {
						return (
							<li 
							className="list-group-item"
							key={track.spotifyID} onClick={() => this.playSong(playlist.URI, track.URI)}>
								{track.name}
							</li>
						)
					})}
				</ul>
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
						"Authorization": `Bearer ${accessToken}`,
						"Content-Type": "application/json"
					},
					body: JSON.stringify(data)
				})
				.catch( (error) => {
					console.error(error);
				})
			})
		})
	}

	componentDidMount() {

		let playlistID = this.props.match.params.playlist
		fetch(`/api/playlists/${playlistID}`, { credentials: "same-origin" })
		.then( (response) => {
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