import React from "react"

// ProgressBar represents how far into the song you are
// It displays visually like a loading bar
export default function ProgressBar(props) {
  let progress = props.position / props.duration;
  let percentage = Math.round(progress * 1000) / 10;

  return (
    <div className="progress">
      <div className="progress-bar" role="progressbar" style={{ width: percentage + '%' }} aria-valuenow="25" aria-valuemin="0" aria-valuemax="100"></div>
    </div>
  )
}