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
    let { loading, playlist, error } = this.state;
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
              <img src={playlist.bucketImageURL} alt="playlist-artwork" />
            </div>
            <ul className="list-group">
              {playlist.tracks.map((track) => {
                return (
                  <li
                    className="list-group-item"
                    key={track.spotifyID}>
                    <Track track={track} playlist={playlist} onPlay={() => { this.playSong(playlist, track.URI) }} />
                  </li>
                )
              })}
            </ul>
          </div>
        </div>
      )
    }
  }

  // getContextURIs is designed so that we get an array of track URIs
  // For very large playlists, we need to get just a slice relative to the selected track
  // so that we can avoid HTTP 413 (request too large) errors
  getContextURIs(playlist, trackURI) {
    const trackURIs = playlist.tracks.map((track) => track.URI);
    // Taking a guess at the payload maximum size
    // With trial and error, length of 500 seems to be pretty safe
    // Only use slices in the case where the playlist is very large
    if (playlist.tracks.length < 500) {
      return trackURIs;
    }
    const trackIndex = trackURIs.indexOf(trackURI);
    return this.getSlice(trackURIs, trackIndex, 500);
  }

  // a is the array
  // i is the index of the selected element
  // l is the length of the desired slice
  getSlice = (a, i, l) => {
    const lowerDistance = Math.floor(l / 2);
    const upperDistance = Math.ceil(l / 2);

    // beginning
    if (i - lowerDistance < 0) {
      const firstHalf = a.slice(i - lowerDistance);
      const secondHalf = a.slice(0, l - firstHalf.length);
      return firstHalf.concat(secondHalf);
    }

    // end
    if (i + upperDistance > a.length) {
      const firstHalf = a.slice(i - lowerDistance, a.length);
      const secondHalf = a.slice(0, l - firstHalf.length);
      return firstHalf.concat(secondHalf);
    }

    // middle
    return a.slice(i - lowerDistance, i + upperDistance);
  }

  // playSong makes an API call directly to Spotify
  // playlist can simply be the playlist object from component state
  async playSong(playlist, selectedTrack) {
    const uris = this.getContextURIs(playlist, selectedTrack);
    let data = {
      "uris": uris,
      "offset": {
        "uri": selectedTrack
      }
    }

    try {
      const deviceID = await this.props.getDeviceID();
      const accessToken = await this.props.getAccessToken();

      const response = await fetch(`https://api.spotify.com/v1/me/player/play?device_id=${deviceID}`, {
        method: "PUT",
        headers: {
          "Authorization": `Bearer ${accessToken}`,
          "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
      });

      if (response.status < 200 || response.status >= 300) {
        const errorMessage = await response.text();
        throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
      }
    } catch (err) {
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