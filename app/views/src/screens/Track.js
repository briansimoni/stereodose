import React from "react";

// Many Tracks make up a Playlist component
export default function Track(props) {
  const { track, onPlay } = props;

  // some math to get the track duration in minutes:seconds
  const duration = track.duration / 1000
  const minutes = Math.floor(duration / 60);
  const seconds = Math.round(duration % 60);
  let displayTime = `${minutes}:${seconds}`;
  if (displayTime.split(":")[1].length === 1) {
    displayTime = `${minutes}:0${seconds}`;
  }

  return (
    <div className="row">

      <div className="col-2">
        <button className="track-play-button" onClick={onPlay}>play</button>
      </div>

      <div className="col-8">
        <h4 className="track-name">{track.name}</h4>
        <h5 className="artists">{track.artists}</h5>
      </div>

      <div className="col-2">
        <h5 className="track-duration">{displayTime}</h5>
      </div>
    </div>
  );
}