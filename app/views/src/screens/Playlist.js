import React from "react";
import Track from './Track';
// import Spotify from 'spotify-web-api-js';

class Playlist extends React.Component{
	device_id
	category
	subcategory
	access_token

	constructor(props) {
		super(props);
		this.state = {playlists: [{name: "poop"}]};
	}

	render() {
		let playlists = this.state.playlists;	
		if(playlists !== null) {
			return (
				<ul>
				{playlists.map(playlist => (
					<Track
						device_id={this.props.device_id}
						access_token={this.props.access_token}
						track_id={playlist.SpotifyID} 
						key={playlist.SpotifyID}
						name={playlist.name}>
					</Track>
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
		)
	}

	// getPlaylists makes an API call to the Stereodose backend
	// to get a list of playlists based on the category passed in
	// getPlaylists(category, subcategory) {
	// 	console.log('wat!');
	// 	return new Promise( (resolve, reject) => {
	// 		let req = new XMLHttpRequest();
	// 		req.open("GET", `/api/playlists/?category=${category}&subcategory=${subcategory}`);
	// 		req.addEventListener("readystatechange", () => {
	// 			if (this.readyState === 4) {
	// 				if (this.status === 200) {
	// 					try {
	// 						console.log(this.responseText);
	// 						return JSON.parse(this.responseText);
	// 					} catch (err) {
	// 						reject(err)
	// 					}
	// 				} else {
	// 					console.log(this.responseText);
	// 					let status = String(this.status);
	// 					let err = new Error(`${status} Unable to retrieve songs for ${category} : ${subcategory} - ${this.statusText}`);
	// 					reject(err)
	// 				}
	// 			}
	// 		});
	// 	})
	// }


}

export default Playlist