import 'react-app-polyfill/ie11';
import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { detect } from 'detect-browser';

const browser = detect();
// handle the case where we don't detect the browser
if (browser) {
  const notSupportedMessage =
    'Spotify does not work on this os/browser. See https://developer.spotify.com/documentation/web-playback-sdk/#supported-browsers';
  if (browser.os === 'Android OS' && browser.name === 'firefox') {
    alert(notSupportedMessage);
  }
  if (browser.name === 'safari' && browser.os !== 'iOS') {
    alert(notSupportedMessage);
  }
  if (browser.os === 'iOS') {
    alert('We are currently working on supporting iOS 14. Check back soon.');
    window.location = 'https://apps.apple.com/us/app/id1518862133';
  }
}

ReactDOM.render(<App />, document.getElementById('root'));

// Only set the service worker once the user has logged in.
// This prevents certain login bugs where the PWA is switched to during the login process but doesn't receive auth cookies.
// It also only asks the user to install the app if they're logged in
isLoggedIn().then(loggedIn => {
  if (loggedIn) {
    if ('serviceWorker' in navigator) {
      window.addEventListener('load', function () {
        navigator.serviceWorker.register('/sw.js').then(
          function (registration) {
            // Registration was successful
            console.log('ServiceWorker registration successful with scope: ', registration.scope);
          },
          function (err) {
            // registration failed :(
            console.log('ServiceWorker registration failed: ', err);
          }
        );
      });
    }
  }
})


function isLoggedIn() {
  return new Promise((resolve) => {
    fetch('/api/users/me').then((response) => {
      if (response.status === 200) {
        resolve(true);
      }
      resolve(false);
    });
  });
}
