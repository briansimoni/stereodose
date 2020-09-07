import React from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMusic } from '@fortawesome/free-solid-svg-icons';
import './GlobalShareButton.css';

// GlobalShareButton is just a link to the /profile/available page
// It is intended to almost always be visible as long as the user is logged in
function GlobalShareButton(props) {
  if (props.location.pathname.includes('profile/available')) {
    return <div />;
  }
  return (
    <div id="global-share-button">
      <Link to="/profile/available">
        <button type="button" className="btn btn-success">
          Share <FontAwesomeIcon icon={faMusic} />
        </button>
      </Link>
    </div>
  );
}

export default GlobalShareButton;
