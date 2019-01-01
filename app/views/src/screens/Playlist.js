import React from "react";
import Track from "./Track";

class Playlist extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      playlist: null,
      error: null
    };
  }

  render() {
    let {loading, playlist, error} = this.state;
    if (loading) {
      return <div></div>
    }
    if (error) {
      return <h3>{error.message}</h3>
    }
    if (playlist) {
      return (
        <div className="row">
          <div className="col">
            <div id="playlist-heading">
              <h2>{playlist.name}</h2>
              <img src={playlist.bucketImageURL} alt="playlist-artwork"/>
            </div>
            <ul className="list-group">
              {playlist.tracks.map( (track) => {
                return (
                  <li
                  className="list-group-item"
                  key={track.spotifyID}>
                    <Track track={track} playlist={playlist} onPlay={() => {this.playSong(playlist.URI, track.URI)}}/>
                  </li>
                )
              })}
            </ul>
          </div>
        </div>
      )
    }
  }

  // playSong makes an API call directly to Spotify
  async playSong(context, uri) {
    let data = {
      "context_uri": context,
      "offset": {
        "uri": uri
      }
    }

    try {
      const deviceID = await this.props.getDeviceID();
      const accessToken = await this.props.getAccessToken();

      const response = await fetch(`https://api.spotify.com/v1/me/player/play?device_id=${deviceID}`, {
        method: "PUT",
        headers : {
          "Authorization": `Bearer ${accessToken}`,
          "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
      });

      if (response.status < 200 || response.status >= 300) {
        const errorMessage = await response.text();
        throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
      }
    } catch(err) {
      alert(err.message);
    }
  }

  // grab the data for the particular playlist that the user requested
  // it can be used later to render a page populated with tracks
  async componentDidMount() {
    try {
      let playlistID = this.props.match.params.playlist

      const response = await fetch(`/api/playlists/${playlistID}`, { credentials: "same-origin" });
      if (response.status !== 200) {
        const errorMessage = await response.text();
        throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`)
      }

      const json = await response.json();

      this.setState({
        loading: false,
        playlist: json
      });
    } catch (err) {
      this.setState({
        loading: false,
        error: err
      })
    }
  }

}

// lets have this component have some function for getting the access token
// and have the player nested in this component

export default Playlist;