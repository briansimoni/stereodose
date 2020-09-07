/**
 * PlaylistController just contains the logic on which "Playlist" component to render.
 * If the user is logged in, it should render a component which grabs tracks from Spotify.
 * If the user is not logged in, it will render tracks from the Stereodose database.
 * We need the playlist to render even if the user is not logged in to help out search engines
 * and also increase the conversion rate for possible new users. The code will be very similar
 * between the "logged in" and "logged out" components
 */

import React from 'react';
import Playlist from './playlist/Playlist';
import LoggedOutPlaylist from './playlist/LoggedOutPlaylist';

export default function PlaylistController(props) {
  if (props.app.userLoggedIn()) {
    return <Playlist {...props}/>
  }
  return <LoggedOutPlaylist {...props}/>
}
