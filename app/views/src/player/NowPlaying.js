import React, { Component } from 'react';
import { Fragment } from 'react';
import ProgressBar from './ProgressBar';
import VolumeSlider from './VolumeSlider';
import RepeatButton from './RepeatButton';
import ShuffleButton from './ShuffleButton';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlay } from '@fortawesome/free-solid-svg-icons';
import { faPause } from '@fortawesome/free-solid-svg-icons';
import { faStepForward } from '@fortawesome/free-solid-svg-icons';
import { faStepBackward } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';

export default class NowPlaying extends Component {
  render() {
    let {
      playerState,
      playerState: { position: position_ms }
    } = this.props;
    let {
      name: track_name,
      duration_ms,
      artists: [{ name: artist_name }],
      album: {
        images: [{ url: album_image }]
      }
    } = playerState.track_window.current_track;

    return (
      <Fragment>
        <div className="row">
          <div className="d-none d-md-block col-md-2 col-lg-1">
            {this.props.app.state.currentPlaylist && (
              <Link to={this.props.app.state.currentPlaylist}>
                <img id="album-image" src={album_image} alt={track_name} />
              </Link>
            )}
          </div>
          <div className="col-sm-3 col-md-2 col-lg-2">
            <div className="row">
              <div className="col text-center">
                <span>
                  {this.props.app.state.currentPlaylist && (
                    <Link id="current-track-link" to={this.props.app.state.currentPlaylist}>
                      <span>{track_name}</span>
                      <span id="current-track-link-separator"> by </span>
                      <span>{artist_name}</span>
                    </Link>
                  )}
                </span>
              </div>
            </div>
          </div>
          <div className="col-sm-7 col-md-7 col-lg-7 text-center">
            {/* <br /> */}
            <div className="controls">
              <ShuffleButton shuffle={this.props.playerState.shuffle} onClick={this.props.onShuffle} />
              <RepeatButton onClick={this.props.onChangeRepeat} repeatMode={playerState.repeat_mode} />
              <FontAwesomeIcon onClick={this.props.onPreviousSong} icon={faStepBackward} />
              {playerState.paused && <FontAwesomeIcon onClick={this.props.onPlayPause} icon={faPlay} />}
              {!playerState.paused && <FontAwesomeIcon onClick={this.props.onPlayPause} icon={faPause} />}
              <FontAwesomeIcon onClick={this.props.onNextSong} icon={faStepForward} />
              <VolumeSlider className="volume-slider" onChangeVolume={this.props.onChangeVolume} disabled={false} />
            </div>
            <ProgressBar onSeek={this.props.onSeek} position={position_ms} duration={duration_ms} disabled={false} />
          </div>
        </div>
      </Fragment>
    );
  }
}
