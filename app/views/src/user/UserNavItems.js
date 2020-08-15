import React from 'react';
import { Fragment } from 'react';
import { Link } from 'react-router-dom';
import './Profile.css';

class UserNavItems extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      loggedIn: null,
      username: '',
      user: {}
    };

    let loggedIn = this.props.app.userLoggedIn();
    this.state.loggedIn = loggedIn;
  }

  render() {
    if (this.state.loggedIn === null) {
      return (
        <div className="row justify-content-md-center">
          <div className="spinner-grow text-center" role="status">
            <span className="sr-only">Loading...</span>
          </div>
        </div>
      );
    }

    if (this.state.loggedIn === false) {
      return (
        <li className="nav-item">
          <span
            onClick={() => {
              this.logIn();
            }}
            className="nav-link"
          >
            Sign In
          </span>
        </li>
      );
    }

    if (this.state.loggedIn === true) {
      return (
        <Fragment>
          <li className="nav-item">
            <Link className="nav-link" to="/profile">
              Profile
            </Link>
          </li>
          {/* Logout is a special case. Need to use a plain <a> tag instead of <Link>*/}
          <li className="nav-item">
            <a href="/auth/logout" className="nav-link">
              Logout
            </a>
          </li>
        </Fragment>
      );
    }
  }

  logIn() {
    window.location = `/auth/login?path=${window.location.pathname}`;
  }

  async componentDidMount() {
    // hooking into jQuery/Bootstrap 4 API to handle menu collapse
    if (window.$) {
      const jQuery = window.$;
      jQuery('nav a').click(() => {
        jQuery('#navbarSupportedContent').collapse('hide');
      });
    }

    if (this.state.loggedIn === true) {
      try {
        await this.fetchProfileData();
      } catch (err) {
        alert(err.message);
      }
    }
  }

  async fetchProfileData() {
    const response = await fetch('/api/users/me', { credentials: 'same-origin' });
    if (response.status !== 200) {
      throw new Error(`Unable to fetch profile data ${response.statusText}`);
    }
    const user = await response.json();
    // check the user's product level: premium or not
    if (user.product !== 'premium') {
      throw new Error('You do not have Spotify Premium. The web player will not work');
    }
    const displayName = user.displayName ? user.displayName : user.spotifyID;
    this.setState({ user, username: displayName });
  }
}

export default UserNavItems;
