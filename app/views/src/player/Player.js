import React, { Component, Fragment } from 'react';
import './Player.css';
import WebPlaybackReact from './WebPlaybackReact';
import Spotify from 'spotify-web-api-js';
import DisabledPlayer from './DisabledPlayer';

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
  }

  componentWillMount() {
    this.props.app.getAccessToken()
      .then((accessToken) => {
        this.setState({
          userAccessToken: accessToken
        });
      })
      .catch((error) => {
        this.setState({ authError: error });
      })
  }

  onPlayPause = async () => {
    let SDK = await this.getSDK();
    let options = { device_id: this.state.userDeviceId };
    let paused = this.state.playerState.paused;
    if (paused) {
      SDK.play(options);
    } else {
      SDK.pause(options);
    }
  }

  nextSong = async () => {
    const options = { device_id: this.state.userDeviceId };
    const SDK = await this.getSDK();
    SDK.skipToNext(options);
  }

  previousSong = async () => {
    const options = { device_id: this.state.userDeviceId };
    const SDK = await this.getSDK();
    SDK.skipToPrevious(options);
  }

  changeVolume = async (volume) => {
    const options = { device_id: this.state.userDeviceId };
    const SDK = await this.getSDK();
    SDK.setVolume(volume, options);
  }

  changeRepeatMode = async () => {
    const repeatMode = this.state.playerState.repeat_mode;
    const options = { device_id: this.state.userDeviceId };
    const SDK = await this.getSDK();

    switch (repeatMode) {
      case 0: // 0 is off
        SDK.setRepeat('context', options);
        break;
      case 1: // 1 is context
        SDK.setRepeat('track', options);
        break;
      default: // 2 is track
        SDK.setRepeat('off', options);
    }
  }

  // position is the desired percentage to seek to
  // duration is the total length in ms of the song.
  seek = async (position, duration) => {
    const SDK = await this.getSDK();
    const ms = Math.round((position / 100) * duration);
    return SDK.seek(ms, { device_id: this.state.userDeviceId });
  }

  getSDK = async () => {
    let SDK = new Spotify();
    let token;
    try {
      token = await this.props.app.getAccessToken();
    } catch (err) {
      alert(err.message);
    }
    SDK.setAccessToken(token);
    return SDK;
  }

  render() {
    let {
      userAccessToken,
      playerLoaded,
      playerSelected,
      playerState,
      authError
    } = this.state;

    // if (userDeviceId) {
    //   this.props.setDeviceID(userDeviceId);
    // }

    let webPlaybackSdkProps = {
      playerName: "Stereodose",
      playerInitialVolume: 1.0,
      playerRefreshRateMs: 100,
      playerAutoConnect: true,
      onPlayerRequestAccessToken: (() => this.props.app.getAccessToken()),
      onPlayerLoading: (() => this.setState({ playerLoaded: true })),
      onPlayerWaitingForDevice: ((data) => {
        this.setState({ playerSelected: false, userDeviceId: data.device_id })
        this.props.app.setState({deviceID: data.device_id});
      }),
      onPlayerDeviceSelected: (() => this.setState({ playerSelected: true })),
      onPlayerStateChange: (playerState => this.setState({ playerState: playerState })),
      onPlayerError: (playerError => alert(playerError))
    };

    return (
      <div>
        {authError &&

          <footer className="footer fixed-bottom">
            <div className="container-fluid">
              <h2 onClick={() => { window.location = "/auth/login" }} id="player-message-not-signed-in">{authError.message}</h2>
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
                  <DisabledPlayer />
                </div>
              </footer>
            }

            {playerLoaded && playerSelected && !playerState &&
              <footer className="footer fixed-bottom">
                <div className="container-fluid">
                  <DisabledPlayer />
                </div>
              </footer>
            }

            {playerLoaded && playerSelected && playerState &&
              <footer className="footer fixed-bottom">
                <div className="container-fluid">
                  <Fragment>
                    <NowPlayingScreen
                      playerState={playerState}
                      onPlayPause={this.onPlayPause}
                      onNextSong={this.nextSong}
                      onPreviousSong={this.previousSong}
                      onChangeVolume={this.changeVolume}
                      onChangeRepeat={this.changeRepeatMode}
                      onSeek={this.seek}
                    />
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