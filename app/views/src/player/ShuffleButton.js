import React from 'react';
import Octicon from 'react-octicon';

export default function ShuffleButton(props) {
  if (props.shuffle === true) {
    return <Octicon onClick={props.onClick} className="shuffle on" name="git-pull-request" />;
  }

  if (props.shuffle === false) {
    return <Octicon onClick={props.onClick} className="shuffle off" name="git-pull-request" />;
  }

  return <Octicon onClick={props.onClick} className="shuffle disabled" name="git-pull-request" />;
}
