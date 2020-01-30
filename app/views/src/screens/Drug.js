import React from 'react';
import { Link } from 'react-router-dom';
import NoMatch from '../404';
import { Route } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import './Screens.css';

// Drug renders the mood choices for the chosen Drug
// Weed -> Chill, Groovin, Thug Life
export default function Drug(props) {
  const drug = props.match.params.drug;
  const match = props.match;
  const categories = props.app.state.categories;

  if (categories !== null && !categories.find(category => category.name === drug)) {
    return <Route component={NoMatch} />;
  }

  if (categories !== null) {
    let subCategories = categories.find(category => category.name === drug).subcategories;
    return (
      <div>
        <div className="row">
          <div className="col">
            <h2 className="mood-choice-header">
              <Link to="/"><FontAwesomeIcon icon={faArrowLeft} /></Link>
              Choose Your Mood
            </h2>
          </div>
        </div>
        <div className="row">
          <div className="col">
            <ul id="moods" className="moods">
              {subCategories.map((subCategory, index) => (
                <li key={index}>
                  <h3>
                    <Link to={`${match.url}/${subCategory}/type`}>{subCategory}</Link>
                  </h3>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="row justify-content-md-center">
      <div className="spinner-grow text-success text-center" role="status">
        <span className="sr-only">Loading...</span>
      </div>
    </div>
  );
}
