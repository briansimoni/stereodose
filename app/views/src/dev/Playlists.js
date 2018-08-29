import React from "react";
import { Link } from "react-router-dom";

class Playlists extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			loading: true,
			error: null,
			playlists: null
		}
	}

	render() {
		console.log('whole thign rendered');
		let {loading, err, playlists} = this.state;
		if (loading) {
			return <h3>Loading</h3>
		}
		
		if (err) {
			return <h3>Error: {err}</h3>
		}

		if (playlists) {
			let match = this.props.match;
			return (
				<div>
					<ul>
						{playlists.map( (playlist) => {
							return (
								<Link to={`${match.url}/${playlist.spotifyID}`}>
									<li key={playlist.spotifyID}>{playlist.name}</li>
								</Link>
							)
						})}
					</ul>
				</div>
			);
		}
	}

	componentDidMount() {
		let drug = this.props.match.params.drug;
		let subcategory = this.props.match.params.subcategory;

		fetch(`/api/playlists/?category=${drug}&subcategory=${subcategory}`)
		.then( (response) => {
			return response.json();
		})
		.then( (json) => {
			this.setState({
				loading: false,
				playlists: json
			})
		})
		.catch( (err) => {
			this.setState({
				loading: false,
				error: err
			})
		});
	}
}

export default Playlists