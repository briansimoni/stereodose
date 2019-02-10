import React from "react";
import {Fragment} from "react";

class PickPlaylist extends React.Component {
  render() {
    const { playlists, onSelectPlaylist } = this.props;
    if (!playlists) {
      return <div></div>
    }
    return (
      <Fragment>
        <h2 id="tab-content-title">Playlists Available From Spotify</h2>
        <div className="list-group">
          {this.props.playlists.map((playlist) =>
            <button onClick={() => { onSelectPlaylist(playlist) }}
              key={playlist.id}
              type="button" className="list-group-item list-group-item-action">
              {playlist.name}
            </button>
          )}
        </div>
      </Fragment>
    );
  }
}

export default PickPlaylist;