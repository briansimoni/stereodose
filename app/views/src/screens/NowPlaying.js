import React, { Component } from 'react';
import { Fragment } from 'react';

export default class NowPlaying extends Component {
  render() {
    let {
      playerState,
      playerState: { position: position_ms }
    } = this.props;
    let {
      id,
      uri: track_uri,
      name: track_name,
      duration_ms,
      artists: [{
        name: artist_name,
        uri: artist_uri
      }],
      album: {
        name: album_name,
        uri: album_uri,
        images: [{ url: album_image }]
      }
    } = playerState.track_window.current_track;

    return (
      <Fragment>
      <div className="row">
        <div className="col-3">
          <img id="album-image" src={album_image} alt={track_name} />
        </div>
        <div className="col">
          <p><a href={track_uri}>{track_name}</a> by <a href={artist_uri}>{artist_name}</a></p>
          <p><a href={album_uri}>{album_name}</a></p>
        </div>
      </div>

      <div className="row">
        <div className="col">
           <p>ID: {id} | Position: {position_ms} | Duration: {duration_ms}</p>
        </div>
      </div>
      </Fragment>
    );
  }
}