import React, { Component } from 'react';

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
      <div className="panel panel-default">
        <div className="panel-heading">Now Playing View</div>
        <div className="panel-body">
          <img src={album_image} alt={track_name} />

          <h4><a href={track_uri}>{track_name}</a> by <a href={artist_uri}>{artist_name}</a></h4>
          <h4><a href={album_uri}>{album_name}</a></h4>
          <h4>ID: {id} | Position: {position_ms} | Duration: {duration_ms}</h4>
        </div>
      </div>
    );
  }
}