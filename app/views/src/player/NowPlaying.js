import React, { Component } from 'react';
import { Fragment } from 'react';
import ProgressBar from './ProgressBar';
import VolumeSlider from './VolumeSlider';
import RepeatButton from "./RepeatButton";

export default class NowPlaying extends Component {
  render() {
    let {
      playerState,
      playerState: { position: position_ms }
    } = this.props;
    let {
      uri: track_uri,
      name: track_name,
      duration_ms,
      artists: [{
        name: artist_name,
        uri: artist_uri
      }],
      album: {
        images: [{ url: album_image }]
      }
    } = playerState.track_window.current_track;

    // playerState.repeat_mode // integer
    // playerState.shuffle // bool

    return (
      <Fragment>
        <div className="row">
          <div className="col-sm-1">
            <img id="album-image" src={album_image} alt={track_name} />
          </div>
          <div className="col-md-1">
            <div className="row justify-content-center">
              <span><a href={track_uri}>{track_name}</a> by <a href={artist_uri}>{artist_name}</a></span>
            </div>
            {/* <div className="row justify-content-center">
              <span><a href={album_uri}>{album_name}</a></span>
            </div> */}
          </div>
          <div className="col-md-8 text-center">
            <br />
            <div className="controls">
              <RepeatButton onClick={this.props.onChangeRepeat} repeatMode={playerState.repeat_mode}/>
              <div onClick={this.props.onPreviousSong} className="arrow-left"></div>
              <div onClick={this.props.onPlayPause} id="play-pause" className={playerState.paused ? "button play": "button pause"} alt="play-pause-button"></div>
              <div onClick={this.props.onNextSong} className="arrow-right"></div>
              <VolumeSlider className="volume-slider" onChangeVolume={this.props.onChangeVolume}/>
            </div>
            <ProgressBar position={position_ms} duration={duration_ms}/>
          </div>
        </div>
      </Fragment>
    );
  }
}