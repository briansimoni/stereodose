import React from 'react';
import ProgressBar from './ProgressBar';
import VolumeSlider from './VolumeSlider';
import RepeatButton from './RepeatButton';
import ShuffleButton from './ShuffleButton';

export default function DisabledPlayer(props) {
  return (
    <div className="row">
      <div className="col-sm-1">{/* <img id="album-image" src={album_image} alt={track_name} /> */}</div>
      <div className="col-md-1">
        <div className="row justify-content-center">
          {/* <span><a href={track_uri}>{track_name}</a> by <a href={artist_uri}>{artist_name}</a></span> */}
        </div>
      </div>
      <div className="col-md-8 text-center">
        <br />
        <div className="controls">
          <ShuffleButton />
          <RepeatButton />
          <div className="arrow-left disabled" />
          <div id="play-pause" className="button play disabled" alt="play-pause-button" />
          <div className="arrow-right disabled" />
          <VolumeSlider className="volume-slider disabled" disabled={true} />
        </div>
        <ProgressBar position={0} duration={100} disabled={true} />
      </div>
    </div>
  );
}
