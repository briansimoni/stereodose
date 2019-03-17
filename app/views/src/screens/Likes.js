import React from "react";
import Octicon from "react-octicon";

export default function Likes(props) {
  const { onLike } = props;
  return(
    <span onClick={onLike}>
      <Octicon name="heart"/>
      {props.number}
    </span>
  )
}