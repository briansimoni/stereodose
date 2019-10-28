import React from 'react';

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

  // added data elements to correlate to events in Google Tag Manager
  return (
    <div className="row">
      <div className="col-2">
        <button data-track-id={track.spotifyID} data-track-name={track.name} className="track-play-button btn" onClick={onPlay}>
          play
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
