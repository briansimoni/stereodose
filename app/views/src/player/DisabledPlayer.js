import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlay } from '@fortawesome/free-solid-svg-icons';
import { faStepForward } from '@fortawesome/free-solid-svg-icons';
import { faStepBackward } from '@fortawesome/free-solid-svg-icons';
import ProgressBar from './ProgressBar';
import VolumeSlider from './VolumeSlider';
import RepeatButton from './RepeatButton';
import ShuffleButton from './ShuffleButton';

export default function DisabledPlayer(props) {
  return (
    <div className="row" id="disabled-player-row">
      <div className="d-none d-md-block col-md-2 col-lg-1"></div>
      <div className="col-sm-3 col-md-2 col-lg-2">
        <div className="row text-center">
          {/* <span><a href={track_uri}>{track_name}</a> by <a href={artist_uri}>{artist_name}</a></span> */}
        </div>
      </div>
      <div className="col-sm-7 col-md-7 col-lg-7 text-center">
        <div className="controls">
          <ShuffleButton />
          <RepeatButton />
          <FontAwesomeIcon icon={faStepBackward} className="disabled" />
          <FontAwesomeIcon icon={faPlay} className="disabled" />
          <FontAwesomeIcon icon={faStepForward} className="disabled" />
          <VolumeSlider className="volume-slider" disabled={true} />
        </div>
        {/* onSeek is a callback function that does nothing to prevent errors when the player is disabled*/}
        <ProgressBar className="disabled" onSeek={() => {}} position={0} duration={100} disabled={true} />
      </div>
    </div>
  );
}
