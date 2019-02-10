import React from "react";

class PickPlaylist extends React.Component {
  render() {
    const { playlists, onSelectPlaylist } = this.props;
    if (!playlists) {
      return <div></div>
    }
    return (
      <ol>
        {this.props.playlists.map((playlist) =>
          <li onClick={() => { onSelectPlaylist(playlist) }}
            key={playlist.id}>
            {playlist.name}
          </li>
        )}
      </ol>
    );
  }
}

export default PickPlaylist;