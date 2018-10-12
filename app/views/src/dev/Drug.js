import React from 'react';
import { Link } from 'react-router-dom';
import "./Drug.css";

// Drug renders the mood choices for the chosen Drug
// Weed -> Chill, Groovin, Thug Life
class Drug extends React.Component {
	constructor(props) {
		super(props);
		this.state = {categories: null, loading: true, error: null };
	}

	render() {
		let match = this.props.match;
			if (this.state.categories !== null && !this.state.loading) {
				return (
					<div className="row">
						<div className="col">
							<h2 className="mood-choice-header">Choose Your Mood</h2>
							<ul className="moods">
								{this.state.categories.map( (category, index) => 
									<li key={index}>
										<h4><Link to={`${match.url}/${category}`}>{category}</Link></h4>
									</li>
								)}
							</ul>
						</div>
					</div>
				)
			}

		return (
			<div>Loading...</div>
		)
	}

	componentDidMount() {
		fetch("/api/categories/", { credentials: "same-origin"} )
			.then((response) => {
				return response.json();
			})
			.then((json) => {
				let drug = this.props.match.params.drug;
				let categories = json[drug];
				this.setState({ loading: false, categories: categories });
			})
			.catch((err) => {
				console.log(err);
				this.setState({ loading: false, error: err });
			})
	}
}

export default Drug;