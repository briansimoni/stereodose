import React from 'react';
import profilePlaceholder from '../images/profile-placeholder.jpeg';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHeart } from '@fortawesome/free-solid-svg-icons';
import { faComments } from '@fortawesome/free-solid-svg-icons';
import { faMusic } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';

class Profile extends React.Component {
  state = {
    spotifyPlaylists: null,
    stereodosePlaylists: null,
    user: null
  };

  render() {
    if (!this.state.user) {
      return (
        <div className="row justify-content-md-center">
          <div className="spinner-grow text-center" role="status">
            <span className="sr-only">Loading...</span>
          </div>
        </div>
      );
    }

    return (
      <div>
        <div className="row justify-content-center">
          <div className="col col-auto">
            {/*this strange-looking code is here to compensate in the edge cases where a user image may not exist or the API has bad data*/}
            {this.state.user.images && this.state.user.images.length > 0 && (
              <img id="profile-image" src={this.state.user.images[this.state.user.images.length - 1].url} alt="profile" />
            )}

            {(!this.state.user.images || !this.state.user.images.length > 0) && <img src={profilePlaceholder} alt="profile" />}
          </div>
        </div>

        <div className="row justify-content-center">
            <h5>{this.state.user.displayName}</h5>
        </div>

        <div className="row">
          <div className="col">
            <div className="accordion" id="profile-accordion">
              <div className="card">
                <div className="card-header" id="headingOne">
                  <h2 className="mb-0">
                    <button
                      className="btn btn-link btn-block text-left"
                      type="button"
                      data-toggle="collapse"
                      data-target="#collapseOne"
                      aria-expanded="false"
                      aria-controls="collapseOne"
                    >
                      <FontAwesomeIcon icon={faHeart} /> Likes ({this.state.user.likes.length})
                    </button>
                  </h2>
                </div>

                <div
                  id="collapseOne"
                  className="collapse"
                  aria-labelledby="headingOne"
                  data-parent="#profile-accordion"
                >
                  <div className="card-body">
                    {this.state.user.likes.map((like, index) => {
                      return <Like key={index} like={like} />;
                    })}
                  </div>
                </div>
              </div>
              <div className="card">
                <div className="card-header" id="headingTwo">
                  <h2 className="mb-0">
                    <button
                      className="btn btn-link btn-block text-left collapsed"
                      type="button"
                      data-toggle="collapse"
                      data-target="#collapseTwo"
                      aria-expanded="false"
                      aria-controls="collapseTwo"
                    >
                      <FontAwesomeIcon icon={faComments} /> Comments ({this.state.user.comments.length})
                    </button>
                  </h2>
                </div>
                <div
                  id="collapseTwo"
                  className="collapse"
                  aria-labelledby="headingTwo"
                  data-parent="#profile-accordion"
                >
                  <div className="card-body">
                    {this.state.user.comments.map((comment, index) => (
                      <Comment key={index} comment={comment} />
                    ))}
                  </div>
                </div>
              </div>
              <div className="card">
                <div className="card-header" id="headingThree">
                  <h2 className="mb-0">
                    <button
                      className="btn btn-link btn-block text-left collapsed"
                      type="button"
                      data-toggle="collapse"
                      data-target="#collapseThree"
                      aria-expanded="false"
                      aria-controls="collapseThree"
                    >
                      <FontAwesomeIcon icon={faMusic} /> Playlists ({this.state.user.playlists.length})
                    </button>
                  </h2>
                </div>
                <div
                  id="collapseThree"
                  className="collapse"
                  aria-labelledby="headingThree"
                  data-parent="#profile-accordion"
                >
                  <div className="card-body">
                    {this.state.user.playlists.map((playlist, index) => (
                      <Playlist key={index} playlist={playlist} />
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  async componentDidMount() {
    try {
      await this.fetchUserData();
    } catch (err) {
      alert(err.message);
    }
  }

  fetchUserData = async () => {
    const response = await fetch('/api/users/me', { credentials: 'same-origin' });
    if (response.status !== 200) {
      throw new Error(`${response.status} Unable to fetch user profile`);
    }
    const user = await response.json();
    this.setState({ user: user });
  };
}

function Like(props) {
  const like = props.like;

  return (
    <div className="card mb-3">
      <div className="row no-gutters">
        <div className="col-3 col-sm-2">
          <img src={like.playlist.bucketThumbnailURL} className="card-img" alt="playlist-artwork"></img>
        </div>
        <div className="col-9 col-sm-6">
          <div className="card-body">
            <Link to={like.playlist.permalink}>
              <h6 className="card-title">{like.playlist.name}</h6>
            </Link>
            <p className="card-text">
              <small>
                {like.playlist.category} - {like.playlist.subCategory}
              </small>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

function Comment(props) {
  const comment = props.comment;

  return (
    <div className="card mb-3">
      <div className="row no-gutters">
        <div className="col-3 col-sm-2">
          <img src={comment.playlist.bucketThumbnailURL} className="card-img" alt="playlist-artwork"></img>
        </div>
        <div className="col-9 col-sm-6">
          <div className="card-body">
            <Link to={comment.playlist.permalink}>
              <h6 className="card-title">{comment.playlist.name}</h6>
            </Link>
            <p className="card-text">
              <small>
                {comment.playlist.category} - {comment.playlist.subCategory}
              </small>
            </p>
            <p className="card-text">
              <small className="text-muted">{comment.content}</small>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

function Playlist(props) {
  const playlist = props.playlist;

  return (
    <div className="card mb-3">
      <div className="row no-gutters">
        <div className="col-3 col-sm-2">
          <img src={playlist.bucketThumbnailURL} className="card-img" alt="playlist-artwork"></img>
        </div>
        <div className="col-9 col-sm-6">
          <div className="card-body">
            <Link to={playlist.permalink}>
              <h6 className="card-title">{playlist.name}</h6>
            </Link>
            <p className="card-text">
              <small>
                {playlist.category} - {playlist.subCategory}
              </small>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Profile;