import React from "react";

export default function SharePlaylistButton(props) {
  const { disabled, onShareStereodose, inFlight } = props;
  return (
    <button
      onClick={onShareStereodose}
      className="btn btn-primary"
      disabled={disabled}>
      {inFlight ? "Uploading" : "Share"}
    </button>
  );
}