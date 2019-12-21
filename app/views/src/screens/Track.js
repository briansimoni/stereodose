import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlay } from '@fortawesome/free-solid-svg-icons';
import { faPause } from '@fortawesome/free-solid-svg-icons';

// Many Tracks make up a Playlist component
export default function Track(props) {
  const { track, onPlay } = props;

  // some math to get the track duration in minutes:seconds
  const duration = track.duration / 1000;
  const minutes = Math.floor(duration / 60);
  const seconds = Math.round(duration % 60);
  let displayTime = `${minutes}:${seconds}`;
  if (displayTime.split(':')[1].length === 1) {
    displayTime = `${minutes}:0${seconds}`;
  }

  // Apparently the Spotify ID from the web API doesn't necessarilly match the ID in the web player
  // Somteimtes there is linked_from object that contains the ID that you would find in the web API
  let currentTrackId;
  if (props.currentlyPlayingTrack) {
    currentTrackId = props.currentlyPlayingTrack.linked_from.id || props.currentlyPlayingTrack.id;
  }

  return (
    <div className="row">
      <div className="col-2">
        <button
          // added data elements to correlate to events in Google Tag Manager
          data-track-id={track.spotifyID}
          data-track-name={track.name}
          className="track-play-button btn"
          onClick={onPlay}
        >
          {currentTrackId === track.spotifyID && !props.paused && <FontAwesomeIcon icon={faPause} />}
          {((currentTrackId !== track.spotifyID) || props.paused) && <FontAwesomeIcon icon={faPlay} />}
        </button>
      </div>

      <div className="col-8">
        <h5 className="track-name">{track.name}</h5>
        <h6 className="artists">{track.artists}</h6>
      </div>

      <div className="col-2">
        <h6 className="track-duration">{displayTime}</h6>
      </div>
    </div>
  );
}
