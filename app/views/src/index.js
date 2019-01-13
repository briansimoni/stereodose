import 'react-app-polyfill/ie11';
import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { detect } from 'detect-browser';

const browser = detect();
// handle the case where we don't detect the browser
if (browser) {
  console.log(browser.name);
  console.log(browser.version);
  console.log(browser.os);
  const notSupportedMessage = 'Spotify does not work on this os/browser. Try using Chrome. See https://developer.spotify.com/documentation/web-playback-sdk/#supported-browsers';
  if (browser.os === 'Android OS' && browser.name === 'firefox') {
    alert(notSupportedMessage);
  }

  if (browser.name === 'safari') {
    alert(notSupportedMessage);
  }
  if (browser.os === 'iOS' && browser.name === 'chrome') {
    alert(notSupportedMessage);
  }
}

ReactDOM.render(<App />, document.getElementById('root'));

if ('serviceWorker' in navigator) {
  window.addEventListener('load', function () {
    navigator.serviceWorker.register('/sw.js').then(function (registration) {
      // Registration was successful
      console.log('ServiceWorker registration successful with scope: ', registration.scope);
    }, function (err) {
      // registration failed :(
      console.log('ServiceWorker registration failed: ', err);
    });
  });
}