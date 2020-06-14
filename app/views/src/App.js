import React from 'react';
import Drugs from './screens/Drugs';
import Drug from './screens/Drug';
import Playlists from './screens/Playlists';
import Playlist from './screens/Playlist';
import Player from './player/Player';
import About from './screens/About';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import UserProfile from './user/UserProfile';
import Profile from './user/Profile';
import Header from './user/Header';
import NoMatch from './404';
import ChoosePlaylistType from './screens/ChoosePlaylistType';
import RandomPlaylist from './screens/RandomPlaylist';

// App is the top level component for Stereodose.
// A reference to itself is passed to child components for an inversion of control.
// Simply stated, this allows for child components to use App's high level methods
// such as getAccessToken(), getCategories(), and Player methods without making wasteful HTTP calls
// since both of those pieces of data are held in App's memory
class App extends React.Component {
  // accessToken is a Spotify OAuth token
  accessToken = null;
  // player is an instance of https://developer.spotify.com/documentation/web-playback-sdk/reference/#api-spotify-player
  // The WebPlaybackReact component sets this property.
  // It's methods seem to perform better than making API calls, however the SDK is still experimental
  player = null;

  state = {
    // added currentTrack because the Visualizer ultimately needs it
    // The Player Component should setState({ currentTrack }) if the track ever changes
    currentTrack: null,
    paused: false,
    categories: null
  };

  render() {
    return (
      <BrowserRouter>
        <div>
          <Route path="/" render={props => <Header {...props} app={this} />} />

          <main role="main" className="container">
            {/* Routes wrapped in a Switch match only the first route for ambiguous matches*/}
            <Switch>
              <Route exact path="/profile" render={props => <Profile {...props} app={this} />} />

              <Route exact path="/profile/shared" render={props => <UserProfile {...props} app={this} />} />

              <Route exact path="/profile/available" render={props => <UserProfile {...props} app={this} />} />

              <Route exact path="/" render={props => <Drugs {...props} app={this} />} />

              <Route exact path="/about" render={props => <About {...props} app={this} />} />

              <Route exact path="/:drug" render={props => <Drug {...props} app={this} />} />

              <Route exact path="/:drug/:subcategory" component={Playlists} />

              <Route exact path="/:drug/:subcategory/type" render={props => <ChoosePlaylistType {...props} app={this} />} />

              <Route exact path="/:drug/:subcategory/random" render={props => <RandomPlaylist {...props} app={this} />} />

              <Route exact path="/:drug/:subcategory/:playlist" render={props => <Playlist {...props} app={this} />} />

              <Route component={NoMatch} />
            </Switch>

          </main>

          <Route path="/" render={props => <Player {...props} app={this} />} />
        </div>
      </BrowserRouter>
    );
  }

  async componentDidMount() {
    try {
      await this.getCategories();
    } catch (err) {
      alert(err.message);
    }
  }

  getCategories = async () => {
    const response = await fetch('/api/categories/', { credentials: 'same-origin' });
    if (response.status !== 200) {
      throw new Error('Unable to fetch categories');
    }
    const categories = await response.json();
    this.setState({ categories });
  };

  // getAccessToken will return a Promise to resolve to a Spotify API access_token
  // The token is cached in the member variable of this object and updated upon expiry
  getAccessToken = async () => {
    if (!this.userLoggedIn()) {
      throw new Error('Sign in with Spotify Premium to Play Music');
    }
    let token;
    if (this.accessToken === null) {
      const response = await fetch('/auth/token', { credentials: 'same-origin' });
      if (response.status !== 200) {
        throw new Error(`Unable to fetch Spotify access token: ${response.status} ${response.statusText}`);
      }
      token = await response.json();
    } else {
      token = this.accessToken;
    }

    if (!this.tokenIsExpired(token)) {
      this.accessToken = token;
      return token.access_token;
    }
    const response = await fetch('/auth/refresh', { credentials: 'same-origin' });
    if (response.status !== 200) {
      throw new Error(`Unable to refresh Spotify access token: ${response.status} ${response.statusText}`);
    }
    token = await response.json();
    // the token from the refresh endpoint does not have an expiry
    // it has "expires_in" in seconds (probably 3600)
    token.expiry = new Date().setSeconds(token.expires_in);
    this.accessToken = token;
    return token.access_token;
  };

  tokenIsExpired(token) {
    const expiresOn = token.expiry;
    const now = new Date();
    const expiresDate = new Date(expiresOn);
    if (now < expiresDate) {
      return false;
    }
    return true;
  }

  // userLoggedIn returns true if the user is logged in, false otherwise
  userLoggedIn() {
    // stolen from Stack Overflow
    function getCookie(name) {
      var dc = document.cookie;
      var prefix = name + '=';
      var begin = dc.indexOf('; ' + prefix);
      if (begin === -1) {
        begin = dc.indexOf(prefix);
        if (begin !== 0) return null;
      } else {
        begin += 2;
        var end = document.cookie.indexOf(';', begin);
        if (end === -1) {
          end = dc.length;
        }
      }
      // because unescape has been deprecated, replaced with decodeURI
      //return unescape(dc.substring(begin + prefix.length, end));
      return decodeURI(dc.substring(begin + prefix.length, end));
    }

    let cookie = getCookie('stereodose_session');
    if (!cookie) {
      return false;
    }

    return true;
  }
}

export default App;
