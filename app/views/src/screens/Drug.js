import React from 'react';
import { Link } from 'react-router-dom';
import NoMatch from '../404';
import { Route } from 'react-router-dom';
import Helmet from 'react-helmet';
import { captializeFirstLetter } from '../util';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import './Screens.css';

// Drug renders the mood choices for the chosen Drug
// Weed -> Chill, Groovin, Thug Life
export default function Drug(props) {
  const drug = props.match.params.drug;
  const match = props.match;
  const categories = props.app.state.categories;

  if (categories !== null && !categories.find((category) => category.name === drug)) {
    return <Route component={NoMatch} />;
  }

  if (categories !== null) {
    let subCategories = categories.find((category) => category.name === drug).subcategories;
    return (
      <div>
        <Helmet>
          <title>Stereodose | {captializeFirstLetter(drug)} | Choose Your Mood</title>
          <meta name="Description" content={generateDescription(drug)}></meta>
        </Helmet>

        <div className="row">
          <div className="col">
            <h2 className="mood-choice-header">
              <Link to="/">
                <FontAwesomeIcon icon={faArrowLeft} />
              </Link>
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
      <div className="spinner-grow text-center" role="status">
        <span className="sr-only">Loading...</span>
      </div>
    </div>
  );
}

/**
 * Provide the drug name and generate a specific description for the meta tag
 * @param {String} drug
 * @returns {String}
 */
function generateDescription(drug) {
  switch(drug.toLowerCase()) {
    case 'weed':
      return 'Chill, Groovin, and Thug Life playlists'
    case 'ecstacy':
      return 'Dance, Floored, and Rolling Balls playlists'
    case 'shrooms':
      return 'Matrix, Shaman, and Space playlists'
    case 'lsd':
      return 'Calm, Trippy, and Rockstar playlists'
    default:
      return 'Choose your mood'
  }
}
