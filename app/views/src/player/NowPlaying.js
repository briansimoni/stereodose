import React, { Component } from 'react';
import { Fragment } from 'react';
import ProgressBar from './ProgressBar';
import VolumeSlider from './VolumeSlider';
import RepeatButton from "./RepeatButton";
import ShuffleButton from "./ShuffleButton";

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

    return (
      <Fragment>
        <div className="row">
          <div className="d-none d-md-block col-md-2 col-lg-1">
            <img id="album-image" src={album_image} alt={track_name} />
          </div>
          <div className="col-sm-3 col-md-2 col-lg-2">
            <div className="row">
              <div className="col text-center">
                <span><a href={track_uri}>{track_name}</a> by <a href={artist_uri}>{artist_name}</a></span>
              </div>
              
            </div>
          </div>
          <div className="col-sm-7 col-md-7 col-lg-7 text-center">
            <br />
            <div className="controls">
              <ShuffleButton shuffle={this.props.playerState.shuffle} onClick={this.props.onShuffle}/>
              <RepeatButton onClick={this.props.onChangeRepeat} repeatMode={playerState.repeat_mode}/>
              <div onClick={this.props.onPreviousSong} className="arrow-left"></div>
              <div onClick={this.props.onPlayPause} id="play-pause" className={playerState.paused ? "button play": "button pause"} alt="play-pause-button"></div>
              <div onClick={this.props.onNextSong} className="arrow-right"></div>
              <VolumeSlider className="volume-slider" onChangeVolume={this.props.onChangeVolume} disabled={false}/>
            </div>
            <ProgressBar onSeek={this.props.onSeek} position={position_ms} duration={duration_ms} disabled={false}/>
          </div>
        </div>
      </Fragment>
    );
  }
}