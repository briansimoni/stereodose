import React from 'react';
import { Link } from 'react-router-dom';
import "./Screens.css";
import NoMatch from "../404";
import { Route } from 'react-router-dom';

// Drug renders the mood choices for the chosen Drug
// Weed -> Chill, Groovin, Thug Life
class Drug extends React.Component {
  constructor(props) {
    super(props);
    this.state = {categories: null, loading: true, error: null };
  }

  render() {
    let drug = this.props.match.params.drug;
    let match = this.props.match;
    if (this.state.categories !== null && !this.state.loading && !(drug in this.state.categories)) {
      return <Route component={NoMatch}/>
    }
    if (this.state.categories !== null && !this.state.loading) {
      return (
        <div className="row">
          <div className="col">
            <h2 className="mood-choice-header">Choose Your Mood</h2>
            <ul id="moods" className="moods">
              {this.state.categories[drug].map( (category, index) => 
                <li key={index}>
                  <h3><Link to={`${match.url}/${category}`}>{category}</Link></h3>
                </li>
              )}
            </ul>
          </div>
        </div>
      )
    }

    return (
      <div></div>
    )
  }

  async componentDidMount() {
    try {
      const response = await fetch("/api/categories/", { credentials: "same-origin"} );
      if (response.status !== 200) {
        throw new Error(response.status + ": Error fetching categories")
      }
      const categories = await response.json();
      this.setState({ loading: false, categories: categories });
    } catch (err) {
      console.log(err);
      this.setState({ loading: false, error: err });
    }
  }

}

export default Drug;