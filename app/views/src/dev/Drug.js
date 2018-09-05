import React from 'react';
import { Link } from 'react-router-dom';

class Drug extends React.Component {
	constructor(props) {
		super(props);
		this.state = {categories: null, loading: true, error: null };
	}

	render() {
		let match = this.props.match;
			if (this.state.categories !== null && !this.state.loading) {
				return (
					<ul>
						<li>{this.props.match.params.drug}</li>
						{this.state.categories.map( (category, index) => {
							return (
								<li key={index}>
									<Link to={`${match.url}/${category}`}>{category}</Link>
								</li>
							)
						})}
					</ul>
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