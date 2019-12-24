import React from 'react';
import Track from './Track';
import Comments from './Comments';
import Likes from './Likes';
import Visualizer from './Visualizer';
import Spotify from 'spotify-web-api-js';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye } from '@fortawesome/free-solid-svg-icons';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';

// Playlist is the parent component that controls the entire display for a particular playlist.
// It is a composite of likes, comments, tracks, and playlist image.
// For likes and comments to work, it also keeps track of user state without parent or peer
// component dependencies. In other words, it makes API calls to /api/users/me
class Playlist extends React.Component {
  // likePending is not part of state because it would cause race conditions
  // React likes to do things like batch state updates
  likePending = false;

  constructor(props) {
    super(props);

    this.state = {
      visualizerLoading: false,
      visualizerShown: false,
      trackAnalysis: null,
      loading: true,
      showComments: false,
      playlist: null,
      user: null,
      error: null
    };
  }

  render() {
    let { loading, showComments, playlist, error } = this.state;
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
              <Link to={`/${this.props.match.params.drug}/${this.props.match.params.subcategory}`}><FontAwesomeIcon icon={faArrowLeft} /></Link>
              {playlist.name}
            </h2>
            <img src={playlist.bucketImageURL} alt="playlist-artwork" />
          </div>
          <button className="btn btn-warning comment-toggle" onClick={this.toggleComments}>
            {showComments ? 'Show songs' : `Comments (${playlist.comments.length})`}
          </button>
          <Likes onLike={this.like} playlist={playlist} user={this.state.user} />
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

          {/* Conditionally render either the comments or playlist tracks */}
          {!showComments ? (
            <ul className="list-group playlist">
              {playlist.tracks &&
                playlist.tracks.map(track => {
                  return (
                    <li className="list-group-item" key={track.spotifyID}>
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
          ) : (
              <Comments
                comments={playlist.comments}
                onSubmitComment={this.submitComment}
                onDeleteComment={this.deleteComment}
                user={this.state.user}
              />
            )}
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

  toggleComments = () => {
    this.setState({ showComments: !this.state.showComments });
  };

  // if 401 need to alert user
  submitComment = async text => {
    const options = {
      method: 'POST',
      body: JSON.stringify({
        text: text
      }),
      credentials: 'same-origin'
    };
    const response = await fetch(`/api/playlists/${this.state.playlist.spotifyID}/comments`, options);
    if (response.status !== 201) {
      const errorMessage = await response.text();
      throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
    }

    try {
      const comment = await response.json();
      const playlist = this.state.playlist;
      playlist.comments.push(comment);
      this.setState({ playlist: playlist });
    } catch (err) {
      alert(err);
    }
  };

  deleteComment = async commentID => {
    const options = {
      method: 'DELETE',
      credentials: 'same-origin'
    };

    const playlist = this.state.playlist;

    const response = await fetch(`/api/playlists/${playlist.spotifyID}/comments/${commentID}`, options);
    if (response.status !== 200) {
      const errorMessage = await response.text();
      throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
    }

    // Instead of calling this.updatePlaylistState, simply remove the comment from state immediately
    // We know it was deleted from the database because the response was 200
    // This removes a network call and makes the app more responsive
    playlist.comments = playlist.comments.filter(comment => comment.ID !== commentID);
    this.setState({
      playlist: playlist
    });
  };

  // there is some condition that is possible to reach such that the like button stops working
  // the user has liked the playlist or not
  like = async () => {
    const { playlist, user } = this.state;
    const likePending = this.likePending;
    if (user === null || likePending) {
      return;
    }
    this.likePending = true;

    // The user already liked this playlist. Unlike.
    const like = user.likes.find(like => like.playlistID === playlist.spotifyID);
    if (like) {
      try {
        await this.unlike(like.ID);
        user.likes = user.likes.filter(l => l.ID !== like.ID);
        this.setState({
          user: user
        });
        this.likePending = false;
        return;
      } catch (err) {
        this.setState({
          loading: false,
          error: err
        });
        return;
      }
    }

    const options = {
      method: 'POST',
      credentials: 'same-origin'
    };

    const response = await fetch(`/api/playlists/${playlist.spotifyID}/likes`, options);
    if (response.status !== 201) {
      const errorMessage = await response.text();
      throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
    }

    const newLike = await response.json();
    playlist.likes.push(newLike);
    this.setState({
      playlist: playlist
    });
    await this.updateUserState();
    this.likePending = false;
  };

  unlike = async likeID => {
    const options = {
      method: 'DELETE',
      credentials: 'same-origin'
    };

    const playlist = this.state.playlist;

    const response = await fetch(`/api/playlists/${playlist.spotifyID}/likes/${likeID}`, options);
    if (response.status !== 200) {
      const errorMessage = await response.text();
      throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
    }

    playlist.likes = playlist.likes.filter(l => l.ID !== likeID);
    this.setState({ playlist: playlist });
  };

  updateUserState = async () => {
    // getting an access token implicitly tells us that the user is logged in
    try {
      await this.props.app.getAccessToken();
    } catch (err) {
      if (err.message === 'Sign in with Spotify Premium to Play Music') {
        this.setState({ user: null });
        return;
      }
    }

    const response = await fetch('/api/users/me');
    if (response.status !== 200) {
      const errorMessage = await response.text();
      throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
    }
    const user = await response.json();
    this.setState({ user: user });
  };

  updatePlaylistState = async () => {
    let playlistID = this.props.match.params.playlist;

    const response = await fetch(`/api/playlists/${playlistID}`, { credentials: 'same-origin' });
    if (response.status !== 200) {
      const errorMessage = await response.text();
      throw new Error(`${errorMessage}, ${response.status}, ${response.statusText}`);
    }

    const playlist = await response.json();

    // sort comments by time created
    playlist.comments.sort((a, b) => {
      const playlistADate = new Date(a.CreatedAt);
      const playlistBDate = new Date(b.CreatedAt);
      if (playlistADate < playlistBDate) {
        return -1;
      }
      if (playlistADate > playlistBDate) {
        return 1;
      }
      return 0;
    });

    this.setState({
      loading: false,
      playlist: playlist
    });
  };

  componentDidMount() {
    try {
      this.updatePlaylistState();
      this.updateUserState();
    } catch (err) {
      this.setState({
        loading: false,
        error: err
      });
    }
  }
}

export default Playlist;
