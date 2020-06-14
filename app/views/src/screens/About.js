import React from 'react';
import Feedback from './Feedback';


export default function About() {
  return (
    <div className="row">
      <div className="col">
        <h1 className="about-header">About</h1>
        <p className="about-body">
          Stereodose is a reincarnation of the web/mobile application that closed down back in 2016. It provided a way
          for people to discover music that mainstream services just don't offer. You can share your playlists from
          Spotify and let the community vote the best playlists to the top.
          <strong>
            {' '}
            Stereodose is currently in beta. It is open source so you can see progress and even contribute on
            <a href="https://github.com/briansimoni/stereodose"> GitHub</a>
          </strong>
        </p>

        <p>
          Please report issues to{' '}
          <a href="https://github.com/briansimoni/stereodose/issues">
            https://github.com/briansimoni/stereodose/issues
          </a>
        </p>

        <h1>Compatibility</h1>
        <p>
          You need to have Spotify Premium for the player to function. Additionally, the Spotify Web SDK only supports
          certain browsers. While it does seem to work okay on many mobile browsers, it isn't officially supported. See
          <a href="https://developer.spotify.com/documentation/web-playback-sdk/#supported-browsers">
            {' '}
            https://developer.spotify.com/documentation/web-playback-sdk/#supported-browsers
          </a>
        </p>
        <p>Stereodose is coming to the iOS app store in the next few months. Be sure to check back!</p>
        <h1>Legal</h1>
        <ul>
          <li><a href="/terms-and-conditions">Terms And Conditions</a></li>
          <li><a href="/privacy-policy">Privacy Policy</a></li>
        </ul>
        <Feedback />
      </div>
    </div>
  );
}

