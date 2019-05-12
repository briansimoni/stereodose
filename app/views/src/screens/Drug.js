import React from 'react';
import { Link } from 'react-router-dom';
import "./Screens.css";
import NoMatch from "../404";
import { Route } from 'react-router-dom';

// Drug renders the mood choices for the chosen Drug
// Weed -> Chill, Groovin, Thug Life
export default function Drug(props) {

  const drug = props.match.params.drug;
  const match = props.match;
  const categories = props.app.state.categories;

  if (categories !== null && !(drug in categories)) {
    return <Route component={NoMatch} />
  }

  if (categories !== null) {
    return (
      <div className="row">
        <div className="col">
          <h2 className="mood-choice-header">Choose Your Mood</h2>
          <ul id="moods" className="moods">
            {categories[drug].map((category, index) =>
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
    <div className="row justify-content-md-center">
      <div className="spinner-grow text-success text-center" role="status">
        <span className="sr-only">Loading...</span>
      </div>
    </div>
  )
}