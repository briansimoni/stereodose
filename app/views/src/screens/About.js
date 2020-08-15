import React from 'react';
import Feedback from './Feedback';
import appStoreImage from '../images/Download_on_the_App_Store_Badge_US-UK_RGB_blk_092917.svg';
import braveImage from '../images/brave_logo_2color_reversed.svg';

export default function About() {
  return (
    <div className="row">
      <div className="col">
        <h1 className="about-header">About</h1>
        <p>
          Stereodose is the psychedelic music discovery app that you never knew you needed. Elevate your musical taste
          as you achieve a new level auditory experience. Share, like, listen, and comment on playlists shared by an
          enlightened community. Genres can only go so far in describing how a song feels, and generic moods like
          “happy” are too broad in scope to properly match a song’s essence. These limitations have brought us to
          categorize our music by substance names, which capture many intangible aspects about music that other sorting
          methods usually miss. Stereodose is not here to promote illegal drug use, but to bring you music historically
          defined by counter-culture.
        </p>
        <p className="about-body">
          Stereodose is a reincarnation of the web/mobile application that closed down back in 2016. It provided a way
          for people to discover music that mainstream services just don't offer.
          <strong>
            {' '}
            Stereodose is open source so you can see progress and even contribute on
            <a href="https://github.com/briansimoni/stereodose"> GitHub</a>
          </strong>
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
        <p>Stereodose is now available for iOS!</p>
        <a href="https://apps.apple.com/us/app/id1518862133">
          <img id="apple-app-store-button" alt="apple-app-store-button" src={appStoreImage}></img>
        </a>

        <h1>Support</h1>
        <p>
          Supoprt Stereodose by supporting yourself. Stop sending Apple and Google all of your personal information.
          <a href="https://brave.com/ste942" rel="noopener noreferrer" target="_blank">
            {' '}
            Download the Brave web browser.
          </a>
        </p>
        <p>
          Brave is a fork of Google Chrome with privacy and security central to it's design. It blocks ads and third
          party trackers by default. If you choose to view adds, <strong>Brave will pay you in crypto currency!</strong>
        </p>
        <a href="https://brave.com/ste942" rel="noopener noreferrer" target="_blank">
          <img id="brave-browser-button" alt="brave-browser-button" src={braveImage}></img>
        </a>

        <h1>Legal</h1>
        <p>
          <a href="/terms-and-conditions">Terms And Conditions</a>
        </p>
        <p>
          <a href="/privacy-policy">Privacy Policy</a>
        </p>
        <Feedback />
      </div>
    </div>
  );
}
