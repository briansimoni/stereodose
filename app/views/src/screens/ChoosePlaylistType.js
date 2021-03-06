import React from 'react';
import { Link } from 'react-router-dom';
import NoMatch from '../404';
import { Route } from 'react-router-dom';
import Helmet from 'react-helmet';
import { captializeFirstLetter } from '../util';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import './Screens.css';

// ChoosePlaylist type renders a simple interface which presents the user with two options
// randomly selected playlist or a user-created playlist
export default function ChoosePlaylistType(props) {
  const drug = props.match.params.drug;
  const categories = props.app.state.categories;
  const subcategory = props.match.params.subcategory;

  if (categories !== null && !categories.find((category) => category.name === drug)) {
    return <Route component={NoMatch} />;
  }

  if (categories !== null) {
    return (
      <div>
        <Helmet>
          <title>Stereodose | {captializeFirstLetter(drug)} {captializeFirstLetter(subcategory)} | Choose Type</title>
          <meta name="Description" content="Choose between user-generated and random playlists"></meta>
        </Helmet>

        <div className="row">
          <div className="col">
            <h2 className="mood-choice-header">
              <Link to={`/${drug}`}>
                <FontAwesomeIcon icon={faArrowLeft} />
              </Link>
              {`${drug}: ${subcategory}`}
            </h2>
            <h3 className="text-center">
              <Link to={`/${drug}/${subcategory}`}>User Created</Link>
            </h3>
            <h3 className="text-center">
              <Link to={`/${drug}/${subcategory}/random`}>Random</Link>
            </h3>
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
