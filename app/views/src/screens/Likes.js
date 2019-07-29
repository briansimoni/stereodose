import React from 'react';
import Octicon from 'react-octicon';

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
      <Octicon className={alreadyLiked ? 'liked' : ''} name="heart" />
      {playlist.likes.length}
    </span>
  );
}
