import React from 'react';
import logo from './../logo.svg';

export default () => {
  return (
    <header className="App-header">
      <img src={logo} className="App-logo" alt="logo" />
      <span className="visually-hidden">Spotify for Developers</span>
    </header>
  );
};