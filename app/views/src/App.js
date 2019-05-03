import React from 'react';
import Drugs from './screens/Drugs';
import Drug from './screens/Drug';
import Playlists from './screens/Playlists';
import Playlist from './screens/Playlist';
import Player from './player/Player';
import About from './screens/About';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import UserProfile from './user/Profile';
import Header from './user/Header';
import NoMatch from './404';

class App extends React.Component {

  accessToken = null

  // loggedInPromise resolves at a later time with the user's logged in status (true/false)
  loggedInPromise
  // callback function that can be called to resolve when we know that the user is logged in or not
  loggedInPromiseResolver

  constructor(props) {
    super(props);

    this.state = {
      accessToken: null,
      categories: null,
      loggedIn: false,
    }

    this.loggedInPromise = new Promise((resolve) => {
      this.loggedInPromiseResolver = resolve;
    })
  }

  isUserLoggedIn = loggedIn => {
    let state = this.state;
    state.loggedIn = loggedIn;
    this.setState(state);
  }

  render() {
    return (
      <BrowserRouter>
        <div>
          <Route
            path="/"
            render={(props) =>
              <Header
                {...props}
                app={this}
                isUserLoggedIn={(loggedIn) =>
                  this.loggedInPromiseResolver(loggedIn)
                }
              />
            }
          />

          <main role="main" className="container">
            {/* Routes wrapped in a Switch match only the first route for ambiguous matches*/}
            <Switch>
              <Route exact path="/profile"
                render={(props) =>
                  <UserProfile
                    {...props}
                    app={this}
                  />
                }
              />

              <Route exact path="/profile/shared"
                render={(props) =>
                  <UserProfile
                    {...props}
                    app={this}
                  />
                }
              />

              <Route exact path="/profile/available"
                render={(props) =>
                  <UserProfile
                    {...props}
                    app={this}
                  />
                }
              />

              <Route exact path="/"
                render={(props) =>
                  <Drugs
                    {...props}
                    app={this}
                  />
                }
              />

              <Route exact path="/about"
                render={(props) =>
                  <About
                    {...props}
                    app={this}
                  />
                }
              />

              {/* <Route exact path="/:drug" component={Drug} /> */}
              <Route exact path="/:drug"
                render={(props) =>
                  <Drug
                    {...props}
                    app={this}
                  />
                }
              />


              <Route exact path="/:drug/:subcategory" component={Playlists} />
              <Route
                exact
                path="/:drug/:subcategory/:playlist"
                render={(props) =>
                  <Playlist
                    {...props}
                    app={this}
                  />
                }
              />

              <Route component={NoMatch} />
            </Switch>
          </main>

          <Route
            path="/"
            render={(props) =>
              <Player
                {...props}
                app={this}
              >
              </Player>
            }
          />
        </div>
      </BrowserRouter >
    )
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
      throw new Error("Unable to fetch categories");
    }
    const categories = await response.json();
    this.setState({ categories: categories });
  }

  // pass setDeviceID to the player component so we can lift "state" up
  // and then move it over to peers
  setDeviceID = deviceID => {
    this.deviceIDResolver(deviceID);
  }

  // getAccessToken will return a Promise to resolve to a Spotify API access_token
  // The token is cached in the member variable of this object and updated upon expiry
  // Should be able to pass this function around as a prop to components that need a token
  // i.e. <Player> and <Playlist>
  getAccessToken = async () => {
    let loggedIn = await this.loggedInPromise;
    if (loggedIn === false) {
      throw new Error("Sign in with Spotify Premium to Play Music");
    }
    let token;
    if (this.accessToken === null) {
      let response = await fetch("/auth/token", { credentials: "same-origin" });
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
    let response = await fetch("/auth/refresh", { credentials: "same-origin" });
    if (response.status !== 200) {
      throw new Error(`Unable to refresh Spotify access token: ${response.status} ${response.statusText}`);
    }
    token = await response.json();
    // the token from the refresh endpoint does not have an expiry
    // it has "expires_in" in seconds (probably 3600)
    token.expiry = new Date().setSeconds(token.expires_in);
    this.accessToken = token;
    return token.access_token;
  }

  tokenIsExpired(token) {
    let expiresOn = token.expiry;
    let now = new Date();
    let expiresDate = new Date(expiresOn);
    if (now < expiresDate) {
      return false;
    }
    return true;
  }
}

export default App;