import React, { Component, Fragment } from 'react';
import './Player.css';
import WebPlaybackReact from './WebPlaybackReact';
import Spotify from 'spotify-web-api-js';

import NowPlayingScreen from './NowPlaying';

export default class Player extends Component {

  constructor(props) {
    super(props);
    this.state = {
      // User's session credentials
      userDeviceId: null,
      userAccessToken: null,

      // Player state
      playerLoaded: false,
      playerSelected: false,
      playerState: null,

      authError: null,
    };

    this.onPlayPause = this.onPlayPause.bind(this);
  }

  componentWillMount() {
    this.props.getAccessToken()
      .then((accessToken) => {
        this.setState({
          userAccessToken: accessToken
        });
      })
      .catch((error) => {
        this.setState({ authError: error });
      })
  }

  async onPlayPause() {
    let options = {device_id: this.state.userDeviceId};
    let SDK = new Spotify();
    let token;
    try {
      token = await this.props.getAccessToken();
    } catch (err) {
      alert(err.message);
    }
    SDK.setAccessToken(token);
    let paused = this.state.playerState.paused;
    if (paused) {
      SDK.play(options);
    } else {
      SDK.pause(options);
    }
  }

  render() {
    let {
      userDeviceId,
      userAccessToken,
      playerLoaded,
      playerSelected,
      playerState,
      authError
    } = this.state;

    if (userDeviceId) {
      this.props.setDeviceID(userDeviceId);
    }

    let webPlaybackSdkProps = {
      playerName: "Stereodose",
      playerInitialVolume: 1.0,
      playerRefreshRateMs: 100,
      playerAutoConnect: true,
      onPlayerRequestAccessToken: (() => this.props.getAccessToken()),
      onPlayerLoading: (() => this.setState({ playerLoaded: true })),
      onPlayerWaitingForDevice: (data => this.setState({ playerSelected: false, userDeviceId: data.device_id })),
      onPlayerDeviceSelected: (() => this.setState({ playerSelected: true })),
      onPlayerStateChange: (playerState => this.setState({ playerState: playerState })),
      onPlayerError: (playerError => alert(playerError))
    };

    return (
      <div>
        {authError &&

          <footer className="footer fixed-bottom">
            <div className="container-fluid">
              <h2>{authError.message}</h2>
            </div>
          </footer>
        }
        {userAccessToken &&
          <WebPlaybackReact {...webPlaybackSdkProps}>

            {!playerLoaded &&
              <h2 className="action-orange">Loading Player</h2>
            }

            {!playerSelected &&
              <footer className="footer fixed-bottom">
                <div className="container-fluid">
                </div>
              </footer>
            }

            {playerLoaded && playerSelected && !playerState &&
              <footer className="footer fixed-bottom">
                <div className="container-fluid">
                </div>
              </footer>
            }

            {playerLoaded && playerSelected && playerState &&
              <footer className="footer fixed-bottom">
                <div className="container-fluid">
                  <Fragment>
                    <NowPlayingScreen playerState={playerState} onPlayPause={this.onPlayPause} />
                  </Fragment>
                </div>
              </footer>
            }
          </WebPlaybackReact>
        }
      </div>
    );
  }
};