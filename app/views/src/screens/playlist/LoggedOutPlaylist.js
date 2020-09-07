import React from 'react';
import LoggedOutTrack from './LoggedOutTrack';
import Comments from './Comments';
import Likes from './Likes';
import Helmet from 'react-helmet';
import { captializeFirstLetter } from '../../util';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye } from '@fortawesome/free-solid-svg-icons';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';


class LoggedOutPlaylist extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: true,
      showComments: false,
      playlist: null,
      error: null
    };
  }

  render() {
    let { loading, showComments, playlist, error } = this.state;
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
            <span>
              <FontAwesomeIcon onClick={this.toggleVisualizer} icon={faEye} />
              Visualizer - Alpha
            </span>
          </span>

          {/* Conditionally render either the comments or playlist tracks */}
          {!showComments ? (
            <ul className="list-group playlist">
              {playlist.tracks &&
                playlist.tracks.map((track) => {
                  return (
                    <li className="list-group-item" key={track.spotifyID}>
                      <LoggedOutTrack
                        track={track}
                        playlist={playlist}
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

  toggleComments = () => {
    this.setState({ showComments: !this.state.showComments });
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

  async componentDidMount() {
    try {
      // updating the playlist state and user state can occur in parallel
      await this.updatePlaylistState();
    } catch (err) {
      this.setState({
        loading: false,
        error: err
      });
    }
  }
}

export default LoggedOutPlaylist;
