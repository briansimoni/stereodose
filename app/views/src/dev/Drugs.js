import React from 'react';
import {
	Link
} from 'react-router-dom'

class Drugs extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			loading: true
		}
	}

	render() {
		if (this.state.loading) {
			return <p>Loading...</p>
		}

		if (this.state.categories !== null) {
			let drugNames = Object.keys(this.state.categories);

			return (
				<ul>
					{
						drugNames.map((drug, index) => {
							return <li key={index}>
								<Link to={`/${drug}`}>{drug}</Link>
								{/* <Route
									path={"/"+drug}
									component={Drug}
								/> */}
							</li>
						})
					}
				</ul>
				
			)
		}

		if (this.state.error) {
			return <p>{this.state.error}</p>
		}
	}

	componentDidMount() {
		fetch("/api/categories/")
			.then((response) => {
				return response.json();
			})
			.then((json) => {
				this.setState({ loading: false, categories: json });
			})
			.catch((err) => {
				this.setState({ loading: false, error: err });
			})
	}
}

export default Drugs;