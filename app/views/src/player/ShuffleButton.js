import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faRandom } from '@fortawesome/free-solid-svg-icons';

export default function ShuffleButton(props) {
  if (props.shuffle === true) {
    return <FontAwesomeIcon onClick={props.onClick} icon={faRandom} className="shuffle on"/>;
  }

  if (props.shuffle === false) {
    return <FontAwesomeIcon onClick={props.onClick} icon={faRandom} className="shuffle off"/>;
  }

  return <FontAwesomeIcon onClick={props.onClick} icon={faRandom} className="shuffle disabled"/>;
}
