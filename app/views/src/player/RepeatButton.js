import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faRedo } from '@fortawesome/free-solid-svg-icons';

export default function RepeatButton(props) {
  const repeatMode = props.repeatMode;

  // off
  if (repeatMode === 0) {
    return (
      <span className="repeat-track-span">
        <FontAwesomeIcon onClick={props.onClick} icon={faRedo} className="repeat off" />
      </span>
    );
  }

  // context
  if (repeatMode === 1) {
    return (
      <span className="repeat-track-span">
        <FontAwesomeIcon onClick={props.onClick} icon={faRedo} className="repeat context" />
      </span>
    );
  }

  // track
  if (repeatMode === 2) {
    return (
      <span className="repeat-track-span">
        <FontAwesomeIcon onClick={props.onClick} icon={faRedo} className="repeat track" />
        <span className="repeat-track-number">1</span>
      </span>
    );
  }
  return (
    <span className="repeat-track-span">
      <FontAwesomeIcon onClick={props.onClick} icon={faRedo} className="repeat disabled" />
    </span>
  );
}
