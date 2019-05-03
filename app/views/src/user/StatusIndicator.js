import React from "react";
import { Fragment } from "react";
import { Link } from "react-router-dom";
import "./Profile.css";

// UserStatusIndicator encapsulates the logic of the user's status:
// logged in or not; Spotify premium or not
class UserStatusIndicator extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      loggedIn: null,
      username: "",
      user: {}
    };

    let loggedIn = this.checkSessionCookie();
    this.state.loggedIn = loggedIn;
    this.props.isUserLoggedIn(loggedIn);
  }

  render() {
    if (this.state.loggedIn === null) {
      return <div><p>loading</p></div>
    }

    if (this.state.loggedIn === false) {
      return (
        <li className="nav-item">
          <span onClick={() => { this.logIn() }} className="nav-link">Sign In</span>
        </li>
      )
    }

    if (this.state.loggedIn === true) {
      return (
        <Fragment>
          <li className="nav-item">
            <Link className="nav-link" to="/profile">Profile</Link>
          </li>
          {/* Logout is a special case. Need to use a plain <a> tag instead of <Link>*/}
          <li className="nav-item">
            <a href="/auth/logout" className="nav-link">logout</a>
          </li>

        </Fragment>
      )
    }
  }

  logIn() {
    window.location = "/auth/login";
  }

  async componentDidMount() {
    // hooking into jQuery/Bootstrap 4 API to handle menu collapse
    if (window.$) {
      const jQuery = window.$;
      jQuery("nav a").click(() => {
        jQuery('#navbarSupportedContent').collapse('hide');
      });
    }

    if (this.state.loggedIn === true) {
      try {
        await this.fetchProfileData();
      } catch(err) {
        alert(err.message);
      }
    }
  }

  async fetchProfileData() {
    const response = await fetch("/api/users/me", { credentials: "same-origin" });
    if (response.status !== 200) {
      throw new Error(`Unable to fetch profile data ${response.statusText}`);
    }
    const user = await response.json();
    // check the user's product level: premium or not
    if (user.product !== "premium") {
      throw new Error("You do not have Spotify Premium. The web player will not work");
    }
    const displayName = user.displayName ? user.displayName : user.spotifyID;
    this.setState({ user });
  }

  // checkSessionCookie returns true if the user is logged in
  // false otherwise
  checkSessionCookie() {
    // stolen from Stack Overflow
    function getCookie(name) {
      var dc = document.cookie;
      var prefix = name + "=";
      var begin = dc.indexOf("; " + prefix);
      if (begin === -1) {
        begin = dc.indexOf(prefix);
        if (begin !== 0) return null;
      }
      else {
        begin += 2;
        var end = document.cookie.indexOf(";", begin);
        if (end === -1) {
          end = dc.length;
        }
      }
      // because unescape has been deprecated, replaced with decodeURI
      //return unescape(dc.substring(begin + prefix.length, end));
      return decodeURI(dc.substring(begin + prefix.length, end));
    }

    let cookie = getCookie("_stereodose-session");
    if (!cookie) {
      return false
    }

    return true;
  }
}

export default UserStatusIndicator;