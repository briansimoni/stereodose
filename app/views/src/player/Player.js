import React, { Component, Fragment } from 'react';
import './Player.css';
import WebPlaybackReact from './WebPlaybackReact';
import Spotify from 'spotify-web-api-js';
import DisabledPlayer from './DisabledPlayer';
import GlobalShareButton from '../user/sharing/GlobalShareButton';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSignInAlt } from '@fortawesome/free-solid-svg-icons';
import spotifyIcon from '../images/Spotify_Icon_RGB_Black.png';

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

      authError: null
    };
  }

  componentDidMount() {
    this.props.app
      .getAccessToken()
      .then(accessToken => {
        this.setState({
          userAccessToken: accessToken
        });
      })
      .catch(error => {
        this.setState({ authError: error });
      });
  }

  onPlayPause = () => {
    const paused = this.state.playerState.paused;
    const { player } = this.props.app;
    if (paused) {
      return player.resume();
    } else {
      return player.pause();
    }
  };

  nextSong = () => {
    this.props.app.player.nextTrack();
  };

  previousSong = async () => {
    this.props.app.player.previousTrack();
  };

  changeVolume = async volume => {
    this.props.app.player.setVolume(volume);
  };

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
      default:
        // 2 is track
        SDK.setRepeat('off', options);
    }
  };

  shuffle = async () => {
    const SDK = await this.getSDK();
    const options = { device_id: this.state.userDeviceId };
    return SDK.setShuffle(!this.state.playerState.shuffle, options);
  };

  // position is the desired percentage to seek to
  // duration is the total length in ms of the song.
  seek = async (position, duration) => {
    const ms = Math.round((position / 100) * duration);
    await this.props.app.player.seek(ms);
  };

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
  };

  render() {
    let { userAccessToken, playerLoaded, playerSelected, playerState, authError } = this.state;

    // if (userDeviceId) {
    //   this.props.setDeviceID(userDeviceId);
    // }

    let webPlaybackSdkProps = {
      app: this.props.app,
      playerName: 'Stereodose',
      playerInitialVolume: 0.5,
      playerRefreshRateMs: 100,
      playerAutoConnect: true,
      onPlayerRequestAccessToken: () => this.props.app.getAccessToken(),
      onPlayerLoading: () => this.setState({ playerLoaded: true }),
      onPlayerWaitingForDevice: data => {
        this.setState({ playerSelected: false, userDeviceId: data.device_id });
        this.props.app.setState({ deviceID: data.device_id });
      },
      onPlayerDeviceSelected: () => this.setState({ playerSelected: true }),
      onPlayerStateChange: playerState => {
        this.setState({ playerState: playerState });
        // we check the id string because if we checked the object, it will reference different places in memory
        // i.e. they will never be equal
        // so it will always setState() which will always re-render, which means it re-renders at least once a second
        // which means the performance will be absolutely terrible
        let appTrackID = '';
        if (this.props.app.state.currentTrack) {
          appTrackID = this.props.app.state.currentTrack.id;
        }
        if (appTrackID !== playerState.track_window.current_track.id) {
          this.props.app.setState({
            currentTrack: playerState.track_window.current_track,
          });
        }
        if (playerState.paused !== this.props.app.state.paused) {
          this.props.app.setState({ paused: playerState.paused });
        }
      },
      onPlayerError: playerError => alert(playerError)
    };

    return (
      <div>
        {authError && (
          <footer className="footer fixed-bottom">
            <div className="container-fluid sign-in-required-container">
              <p>Spotify Premium is required to play music</p>
              <button
                onClick={() => {
                  window.location = `/auth/login?path=${window.location.pathname}`;
                }}
                id="player-message-not-signed-in"
              >
                <img alt="spotify-logo" src={spotifyIcon}></img>
                Sign In
                <FontAwesomeIcon icon={faSignInAlt}/>
              </button>
            </div>
          </footer>
        )}
        {userAccessToken && (
          <WebPlaybackReact {...webPlaybackSdkProps}>
            {!playerLoaded && <h2 className="action-orange">Loading Player</h2>}

            {!playerSelected && (
              <footer className="footer fixed-bottom">
                <div className="container-fluid">
                  <GlobalShareButton location={this.props.location} />
                  <DisabledPlayer />
                </div>
              </footer>
            )}

            {playerLoaded && playerSelected && !playerState && (
              <footer className="footer fixed-bottom">
                <div className="container-fluid">
                  <GlobalShareButton location={this.props.location} />
                  <DisabledPlayer />
                </div>
              </footer>
            )}

            {playerLoaded && playerSelected && playerState && (
              <footer className="footer fixed-bottom">
                <div className="container-fluid">
                  <GlobalShareButton location={this.props.location} />
                  <Fragment>
                    <NowPlayingScreen
                      app={this.props.app}
                      playerState={playerState}
                      onPlayPause={this.onPlayPause}
                      onNextSong={this.nextSong}
                      onPreviousSong={this.previousSong}
                      onChangeVolume={this.changeVolume}
                      onChangeRepeat={this.changeRepeatMode}
                      onShuffle={this.shuffle}
                      onSeek={this.seek}
                    />
                  </Fragment>
                </div>
              </footer>
            )}
          </WebPlaybackReact>
        )}
      </div>
    );
  }
}
