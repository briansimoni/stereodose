import React from 'react';
import Feedback from './Feedback';
// import appStoreImage from '../images/Download_on_the_App_Store_Badge_US-UK_RGB_blk_092917.svg';
import braveImage from '../images/brave_logo_2color_reversed.svg';
import qrCode from '../images/btc-qr-code.png'
import Helmet from 'react-helmet';

export default function About() {
  return (
    <div className="row">
      <Helmet>
        <title>Stereodose - Music Inspired By Drugs</title>
        <meta
          name="Description"
          content="Stereodose is the psychedelic music discovery app that you never knew you needed. Elevate your musical taste as you achieve a new level in auditory experience."
        ></meta>
      </Helmet>

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
        <h1>Support Stereodose</h1>
        <p>There is just one person maintaining the whole website and servers aren't free.</p>
        <p>I accept Bitcoin donations.</p>
        <img id="qr-code" alt="bitcoin-QR-code" src={qrCode}></img>
        <p>35YJZMv1QBkLdhFVEqTYtoMu8jk7CnXfC7</p>
        <p>
          You can also send me cash that I will inevitably use on beer or coffee or whatever.
          <a href="https://www.buymeacoffee.com/stereodose" target="_blank">
            <img id="buy-me-a-coffee-button" src="https://cdn.buymeacoffee.com/buttons/v2/default-violet.png" alt="Buy Me A Coffee"></img>
          </a>
        </p>
        <p>
          Support Stereodose by supporting yourself. Stop sending Apple and Google all of your personal information.
          <a href="https://brave.com/ste942" rel="noopener noreferrer" target="_blank">
            {' '}
            Download the Brave web browser.
          </a>
        </p>
        <p>
          Brave is a fork of Google Chrome with privacy and security central to it's design. It blocks ads and third
          party trackers by default. If you choose to view ads, <strong>Brave will pay you in crypto currency!</strong>
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
