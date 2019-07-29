import React from 'react';
import Octicon from 'react-octicon';

export default function RepeatButton(props) {
  const repeatMode = props.repeatMode;

  // off
  if (repeatMode === 0) {
    return <Octicon onClick={props.onClick} className="repeat off" name="sync" />;
  }

  // context
  if (repeatMode === 1) {
    return <Octicon onClick={props.onClick} className="repeat context" name="sync" />;
  }

  // track
  if (repeatMode === 2) {
    return <Octicon onClick={props.onClick} className="repeat track" name="sync" />;
  }
  return <Octicon onClick={props.onClick} className="repeat disabled" name="sync" />;
}
