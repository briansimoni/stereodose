import React from 'react';
import Track from './Track';
import Visualizer from './Visualizer';
import Spotify from 'spotify-web-api-js';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye } from '@fortawesome/free-solid-svg-icons';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';

// RandomPlaylist is almost a copy of Playlist
// The main difference is that likes and comments are removed
// and the playlist data comes from the /api/playlists/random endpoint
class RandomPlaylist extends React.Component {
  likePending = false;

  constructor(props) {
    super(props);

    this.state = {
      visualizerLoading: false,
      visualizerShown: false,
      trackAnalysis: null,
      loading: true,
      playlist: null,
      error: null
    };
  }

  render() {
    let { loading, playlist, error } = this.state;
    if (loading) {
      return (
        <div className="row justify-content-md-center">
          <div className="spinner-grow text-success text-center" role="status">
            <span className="sr-only">Loading...</span>
          </div>
        </div>
      );
    }
    if (error) {
      return <h3>{error.message}</h3>;
    }

    const { drug, subcategory } = this.props.match.params;
    let albumImageUrl = null;
    if (this.props.app.state.currentTrack) {
      albumImageUrl = this.props.app.state.currentTrack.album.images[0].url || null;
    }

    return (
      <div className="row">
        <div className="col">
          {this.state.visualizerShown && this.state.trackAnalysis && (
            <Visualizer
              app={this.props.app}
              analysis={this.state.trackAnalysis}
              toggleVisualizer={this.toggleVisualizer}
            />
          )}
          <div id="playlist-heading">
            <h2>
              {/* The header contains the playlist name and a back button*/}
              <Link to={`/${this.props.match.params.drug}/${this.props.match.params.subcategory}/type`}><FontAwesomeIcon icon={faArrowLeft} /></Link>
              {`${drug}: ${subcategory}`}
            </h2>
            {albumImageUrl && <img src={albumImageUrl} alt="playlist-artwork" />}
            {!albumImageUrl && <div id="random-playlist-image-placeholder" alt="playlist-artwork" />}
          </div>
          <span>
            {!this.state.visualizerLoading && (
              <span>
                <FontAwesomeIcon onClick={this.toggleVisualizer} icon={faEye} />
                Visualizer - Alpha
              </span>
            )}

            {this.state.visualizerLoading && (
              <div id="visualizer-loading-spinner" className="spinner-border spinner-border-md text-info" role="status">
                <span className="sr-only">Loading...</span>
              </div>
            )}
          </span>

          <ul className="list-group playlist">
            {playlist.tracks &&
              playlist.tracks.map((track, index) => {
                return (
                  <li className="list-group-item" key={index}>
                    <Track
                      currentlyPlayingTrack={this.props.app.state.currentTrack}
                      track={track}
                      playlist={playlist}
                      paused={this.props.app.state.paused}
                      onPlay={() => {
                        this.playSong(playlist, track.URI);
                      }}
                    />
                  </li>
                );
              })}
          </ul>
        </div>
      </div>
    );
  }

  toggleVisualizer = async () => {
    if (!this.state.visualizerShown) {
      try {
        this.setState({ visualizerLoading: true });
        const accessToken = await this.props.app.getAccessToken();
        const SDK = new Spotify();
        SDK.setAccessToken(accessToken);
        const playerState = await this.props.app.player.getCurrentState();
        const trackId = playerState.track_window.current_track.id;
        const analysis = await SDK.getAudioAnalysisForTrack(trackId);
        this.setState({ trackAnalysis: analysis, visualizerLoading: false });
      } catch (error) {
        console.error(error);
        alert(error.message);
      }
    }
    this.setState({ visualizerShown: !this.state.visualizerShown });
    console.log(this.state.visualizerShown);
  };

  // getContextURIs is designed so that we get an array of track URIs
  // For very large playlists, we need to get just a slice relative to the selected track
  // so that we can avoid HTTP 413 (request too large) errors
  getContextURIs(playlist, trackURI) {
    const trackURIs = playlist.tracks.map(track => track.URI);
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
  };

  // playSong makes an API call directly to Spotify
  // playlist can simply be the playlist object from component state
  async playSong(playlist, selectedTrack) {
    // first, if the selectedTrack is currently playing, we actually need to pause instead
    if (this.props.app.state.currentTrack) {
      const currentTrackId = this.props.app.state.currentTrack.linked_from.id || this.props.app.state.currentTrack.id;
      if (
        selectedTrack.includes(currentTrackId)
      ) {
        await this.props.app.player.togglePlay();
        return;
      }
    }

    const uris = this.getContextURIs(playlist, selectedTrack);
    let data = {
      uris: uris,
      offset: {
        uri: selectedTrack
      }
    };

    try {
      const deviceID = this.props.app.state.deviceID;
      if (!deviceID) {
        return;
      }
      const accessToken = await this.props.app.getAccessToken();

      const response = await fetch(`https://api.spotify.com/v1/me/player/play?device_id=${deviceID}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      });

      // When the user clicks play, it has to be from some kind of playlist
      // the current path is added to the app component so other components can
      // create links. In this way, users can come back to the playlist that contains
      // the currently playing track.
      this.props.app.setState({ currentPlaylist: this.props.location.pathname });

      if (response.status < 200 || response.status >= 300) {
        const errorMessage = await response.text();
        throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
      }
    } catch (err) {
      alert(err.message);
    }
  }

  updatePlaylistState = async () => {
    const { drug, subcategory } = this.props.match.params;
    // if the user hasn't selected a different drug/mood combination, reuse the random playlist that was already downloaded
    // this makes the links back to the random playlist (e.g. the track name in the player) bring you back to the same set
    // of tracks that you were looking at before 
    if (this.props.app.state.randomPlaylistData) {
      if (this.props.app.state.randomPlaylistData.drug === drug && this.props.app.state.randomPlaylistData.subcategory === subcategory) {
        this.setState({
          loading: false,
          playlist: this.props.app.state.randomPlaylistData.playlist
        });
        return;
      }
    }
    const response = await fetch(`/api/playlists/random?category=${drug}&subcategory=${subcategory}`, { credentials: 'same-origin' });
    if (response.status !== 200) {
      const errorMessage = await response.text();
      throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
    }

    const playlist = await response.json();

    this.props.app.setState({
      randomPlaylistData: {
        drug,
        subcategory,
        playlist: playlist
      }
    });

    this.setState({
      loading: false,
      playlist: playlist
    });
  };

  async componentDidMount() {
    try {
      await this.updatePlaylistState();
    } catch (err) {
      this.setState({
        loading: false,
        error: err
      });
    }
  }
}

export default RandomPlaylist;
