import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHeart } from '@fortawesome/free-solid-svg-icons';

export default function Likes(props) {
  const { onLike, user, playlist } = props;

  // Check if the user has already clicked like.
  // If they have, render the button with a different color
  let alreadyLiked = false;
  if (user) {
    if (user.likes.find(like => like.playlistID === playlist.spotifyID)) {
      alreadyLiked = true;
    }
  }

  return (
    <span onClick={onLike}>
      <FontAwesomeIcon icon={faHeart} className={alreadyLiked ? 'liked' : ''}/>
      {playlist.likes.length}
    </span>
  );
}
