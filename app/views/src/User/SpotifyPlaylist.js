import React from "react";

// SpotifyPlaylist is used to make up a list of playlists on the user profile page
// Specifically, they are in a state where they are on Spotify but not Stereodose.
class SpotifyPlaylist extends React.Component {

	playlist // playlist name on Spotify, result of API call to Spotify
	categories // result of categories API call
	onUpdate // function passed in by parent

	constructor(props) {
		super(props);
		if (!props.categories || !props.playlist || !props.onUpdate) {
			throw new Error("SpotifyPlaylist requires playlist, categories, and onUpdate props");
		}

		this.state = {
			drug: "",
			mood: ""
		}

		this.onDrugSelection = this.onDrugSelection.bind(this);
		this.onMoodSelection = this.onMoodSelection.bind(this);
		this.onShareToStereodose = this.onShareToStereodose.bind(this);
	}

	render() {
		let playlist = this.props.playlist;
		let categories = this.props.categories;
		let selectedDrug = this.state.drug;
		let moodOptions = [{value: "", text: "Select Your Mood"}];
		if (selectedDrug !== "") {
			let moods = categories[selectedDrug];
			moods.forEach( (mood) => {
				moodOptions.push( {value: mood, text: mood});
			});
		}
		if (playlist && categories) {
			let drugs = Object.keys(categories);
			return (
				<tr>
					<td>{playlist.name}</td>
					<td>
						<select value={this.state.drug} onChange={this.onDrugSelection}>
							<option value="">Choose Your Drug</option>
							{drugs.map( (drug, index) => {
								return <option key={index} value={drug}>{drug}</option>
							})}
						</select>
					</td>
					<td><select value={this.state.mood} onChange={this.onMoodSelection}>
							{moodOptions.map( (mood, index) => {
								return <option key={index} value={mood.value}>{mood.text}</option>
							})}
						</select>
					</td>
					<td><button onClick={this.onShareToStereodose}>Share to Stereodose</button></td>
				</tr>
			)
		}
		return <tr></tr>
	}

	onMoodSelection(event) {
		this.setState({mood: event.target.value});
	}

	onDrugSelection(event) {
		this.setState({drug: event.target.value});
	}

	async onShareToStereodose() {
		let playlist = this.props.playlist
		let drug = this.state.drug
		let mood = this.state.mood;
		// if the user has not selected both drug and mood, do nothing, exit function
		if (!(drug && mood) ) {
			return;
		}
		let resp = await fetch(`/api/playlists/`, {
			method: "POST",
			body: JSON.stringify({
					SpotifyID: playlist.id,
					Category: drug,
					Subcategory: mood
				})
		});
		if (resp.status !== 201) {
			alert("error! " + resp.status + " " + resp.statusText);
		}
		this.props.onUpdate();
	}
}

export default SpotifyPlaylist;