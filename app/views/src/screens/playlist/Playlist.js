import React from 'react';
import Track from './Track';
import Comments from './Comments';
import Likes from './Likes';
// import Visualizer from './Visualizer';
// import Data2D from './Visualizer/Data2D';
import Data3D from './Visualizer/Data3D';
import Helmet from 'react-helmet';
import { captializeFirstLetter } from '../../util';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye } from '@fortawesome/free-solid-svg-icons';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';
import SpotifyWebApi from 'spotify-web-api-js';

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
    let { loading, showComments, playlist, spotifyPlaylist, error } = this.state;
    if (loading) {
      return (
        <div className="row justify-content-md-center">
          <div className="spinner-grow text-center" role="status">
            <span className="sr-only">Loading...</span>
          </div>
        </div>
      );
    }
    if (error) {
      return <h3>{error.message}</h3>;
    }

    const { drug, subcategory } = this.props.match.params;

    return (
      <div className="row">
        <Helmet>
          <title>
            {playlist.name} | {captializeFirstLetter(drug)} {captializeFirstLetter(subcategory)}
          </title>
        </Helmet>

        <div className="col">
          {this.state.visualizerShown && (
            // <Visualizer app={this.props.app} toggleVisualizer={this.toggleVisualizer} />
            // <Data2D app={this.props.app} toggleVisualizer={this.toggleVisualizer} />
            <Data3D app={this.props.app} toggleVisualizer={this.toggleVisualizer} />
            // <Visualizer2 app={this.props.app} toggleVisualizer={this.toggleVisualizer} />
          )}
          <div id="playlist-heading">
            <h2>
              {/* The header contains the playlist name and a back button*/}
              <Link to={`/${this.props.match.params.drug}/${this.props.match.params.subcategory}`}>
                <FontAwesomeIcon icon={faArrowLeft} />
              </Link>
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
          </span>

          {/* Conditionally render either the comments or playlist tracks */}
          {!showComments ? (
            <ul className="list-group playlist">
              {spotifyPlaylist &&
                spotifyPlaylist.tracks.map((track, index) => {
                  return (
                    <li className="list-group-item" key={index}>
                      <Track
                        currentlyPlayingTrack={this.props.app.state.currentTrack}
                        track={track.track}
                        playlist={spotifyPlaylist}
                        paused={this.props.app.state.paused}
                        onPlay={() => {
                          this.playSong(track, index);
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
      if (!this.props.app.state.currentTrack) {
        alert('Start playing music to use the visualizer');
        return;
      }
    }
    this.setState({ visualizerShown: !this.state.visualizerShown });
  };

  // playSong makes an API call directly to Spotify
  // playlist can simply be the playlist object from component state
  async playSong(selectedTrack) {
    const playlistId = this.props.match.params.playlist;
    // first, if the selectedTrack is currently playing, we actually need to pause instead
    if (this.props.app.state.currentTrack) {
      const currentTrackId = this.props.app.state.currentTrack.linked_from.id || this.props.app.state.currentTrack.id;
      if (selectedTrack.track.id.includes(currentTrackId)) {
        await this.props.app.player.togglePlay();
        return;
      }
    }

    const accessToken = await this.props.app.getAccessToken();
    const spotify = new SpotifyWebApi();
    spotify.setAccessToken(accessToken);

    const deviceID = this.props.app.state.deviceID;
    if (!deviceID) {
      return;
    }

    let selectedTrackId;
    if (selectedTrack.track.linked_from) {
      selectedTrackId = selectedTrack.track.linked_from.id;
    } else {
      selectedTrackId = selectedTrack.track.id;
    }

    try {
      await spotify.play({
        device_id: deviceID,
        context_uri: `spotify:playlist:${playlistId}`,
        offset: {
          uri: `spotify:track:${selectedTrackId}`
        }
      });

      // When the user clicks play, it has to be from some kind of playlist
      // the current path is added to the app component so other components can
      // create links. In this way, users can come back to the playlist that contains
      // the currently playing track.
      this.props.app.setState({ currentPlaylist: this.props.location.pathname });
    } catch (err) {
      alert('Something went wrong. Try refreshing the page.');
      console.error(err.message);
    }
  }

  toggleComments = () => {
    this.setState({ showComments: !this.state.showComments });
  };

  // if 401 need to alert user
  submitComment = async (text) => {
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

  deleteComment = async (commentID) => {
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
    playlist.comments = playlist.comments.filter((comment) => comment.ID !== commentID);
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
    const like = user.likes.find((like) => like.playlistID === playlist.spotifyID);
    if (like) {
      try {
        await this.unlike(like.ID);
        user.likes = user.likes.filter((l) => l.ID !== like.ID);
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

  unlike = async (likeID) => {
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

    playlist.likes = playlist.likes.filter((l) => l.ID !== likeID);
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
      playlist: playlist
    });
  };

  /**
   * gets the full spotify playlists including all of the tracks
   * and then adds it to component state
   */
  async getSpotifyPlaylist() {
    const playlistId = this.props.match.params.playlist;

    const spotify = new SpotifyWebApi();
    const accessToken = await this.props.app.getAccessToken();
    spotify.setAccessToken(accessToken);

    let country = 'US';
    if (this.props.app.userLoggedIn()) {
      if (!this.props.app.state.spotifyUser) {
        await this.props.app.getSpotifyUser();
        country = this.props.app.state.spotifyUser.country || 'US';
      }
    }

    const spotifyPlaylist = await spotify.getPlaylist(playlistId);

    let tracks = [];
    let trackPage = await spotify.getPlaylistTracks(playlistId, {
      market: country
    });

    tracks = tracks.concat(trackPage.items);
    while (tracks.length < trackPage.total) {
      trackPage = await spotify.getPlaylistTracks(playlistId, {
        offset: tracks.length,
        market: country
      });
      tracks = tracks.concat(trackPage.items);
    }

    tracks = tracks.filter((track) => track.track.is_playable);
    spotifyPlaylist.tracks = tracks;

    this.setState({
      spotifyPlaylist
    });
  }

  async componentDidMount() {
    try {
      // updating the playlist state and user state can occur in parallel
      await Promise.all([this.updatePlaylistState(), this.updateUserState(), this.getSpotifyPlaylist()]);
      this.setState({ loading: false });
    } catch (err) {
      this.setState({
        loading: false,
        error: err
      });
    }
  }
}

export default Playlist;
