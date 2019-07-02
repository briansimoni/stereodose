import React from "react";
import Octicon from "react-octicon";

export default function ShuffleButton(props) {
  return <Octicon onClick={props.onClick} className="repeat off" name="git-pull-request"/>
}