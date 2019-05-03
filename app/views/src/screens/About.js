import React from "react";

export default function About(props) {
  return (
    <div className="row">
      <div className="col">
        <h2 className="about-header">About</h2>
        <p className="about-body">
          Stereodose is a reincarnation of the web/mobile application that closed down back in 2016.
          It provided a way for people to discover music that mainstream services just don't offer.
          You can share your playlists from Spotify and let the community vote the best playlists to the top.
          <strong>Stereodose is currently in beta. It is open source so you can see progress and even contribute on 
            <a href="https://github.com/briansimoni/stereodose"> GitHub</a>
          </strong>
        </p>

        <p>
          Please report issues to <a href="https://github.com/briansimoni/stereodose/issues">https://github.com/briansimoni/stereodose/issues</a>
        </p>
      </div>
    </div>
  )
}