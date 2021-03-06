import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlay } from '@fortawesome/free-solid-svg-icons';
import { faPause } from '@fortawesome/free-solid-svg-icons';

// Many Tracks make up a Playlist component
export default function Track(props) {
  const { track, onPlay } = props;

  // some math to get the track duration in minutes:seconds
  const duration = track.duration_ms / 1000;
  const minutes = Math.floor(duration / 60);
  const seconds = Math.round(duration % 60);
  let displayTime = `${minutes}:${seconds}`;
  if (displayTime.split(':')[1].length === 1) {
    displayTime = `${minutes}:0${seconds}`;
  }

  let artistsString = track.artists.reduce((str, artist) => {
    return (str += artist.name + ', ');
  }, '');
  artistsString = artistsString.slice(0, artistsString.length - 2);

  // Apparently the Spotify ID from the web API doesn't necessarily match the ID in the web player
  // Sometimes there is linked_from object that contains the ID that you would find in the web API
  let currentTrackId;
  if (props.currentlyPlayingTrack) {
    currentTrackId = props.currentlyPlayingTrack.linked_from.id || props.currentlyPlayingTrack.id;
  }

  let thisTrackId;
  if (track.linked_from) {
    thisTrackId = track.linked_from.id;
  } else {
    thisTrackId = track.id;
  }

  return (
    <div className="row">
      <div className="col-2">
        <button
          // added data elements to correlate to events in Google Tag Manager
          data-track-id={track.id}
          data-track-name={track.name}
          className="track-play-button btn"
          onClick={onPlay}
        >
          {currentTrackId === thisTrackId && !props.paused && <FontAwesomeIcon icon={faPause} />}
          {(currentTrackId !== thisTrackId || props.paused) && <FontAwesomeIcon icon={faPlay} />}
        </button>
      </div>

      <div className="col-8">
        <h5 onClick={onPlay} className="track-name">
          {track.name}
        </h5>
        <h6 className="artists">{artistsString}</h6>
      </div>

      <div className="col-2">
        <h6 className="track-duration">{displayTime}</h6>
      </div>
    </div>
  );
}
