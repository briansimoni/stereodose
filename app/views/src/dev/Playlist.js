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
					<h4 onClick={ () => this.props.getAccessToken()}>Playlist Get access token</h4>
					{playlist.Tracks.map( (track) => {
						return (
							<li key={track.SpotifyID}>
								{track.Name}
							</li>
						)
					})}
				</div>
			)
		}
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