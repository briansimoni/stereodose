import React from "react";
import Spotify from "spotify-web-api-js";
import { Link } from "react-router-dom";
import ShareSpotifyPlaylist from "./sharing/ShareSpotifyPlaylist";
import StereodosePlaylist from "./StereodosePlaylist";
import "./Profile.css";
import profilePlaceholder from "../images/profile-placeholder.jpeg";

class UserProfile extends React.Component {

  state = {
    spotifyPlaylists: null,
    stereodosePlaylists: null,
    user: null
  }

  render() {
    const { spotifyPlaylists, stereodosePlaylists, user } = this.state;
    const categories = this.props.app.state.categories;

    if (spotifyPlaylists && stereodosePlaylists && categories) {
      return (
        <div className="container">
          <div className="row">
            <div className="col">
              {this.props.location.pathname === "/profile/shared" &&
                <div label="Playlists Shared to Stereodose">
                  <h2 id="content-title">Playlists Shared to Stereodose</h2>
                  <table className="table">
                    <tbody>
                      <tr>
                        <th>Playlist Name</th>
                        <th>Drug</th>
                        <th>Mood</th>
                        <th>Delete?</th>
                      </tr>
                      {stereodosePlaylists.map((playlist) => {
                        return <StereodosePlaylist
                          key={playlist.spotifyID}
                          playlist={playlist}
                          onUpdate={() => { this.checkPlaylists() }}
                        />
                      })}
                    </tbody>
                  </table>
                </div>
              }

              {this.props.location.pathname === "/profile/available" &&
                <div label="Playlists Available">
                  <ShareSpotifyPlaylist
                    playlists={spotifyPlaylists}
                    categories={categories}
                    onUpdate={() => { this.checkPlaylists() }}
                  />
                </div>
              }

              {this.props.location.pathname === "/profile" && user &&

                <div className="text-center profile-main">

                  <div className="row">
                    <div className="col">

                      {/*hotfix*/}
                      {user.images && user.images.length > 0 &&
                        <img src={user.images[user.images.length - 1].url} alt="profile" className="img-thumbnail rounded-circle" />
                      }

                      {(!user.images || !user.images.length > 0) &&
                        <img src={profilePlaceholder} alt="profile" className="img-thumbnail rounded-circle" />
                      }

                      <br />
                      {user.displayName}
                    </div>
                  </div>

                  <div className="row">
                    <div className="col-md-6"><h2><Link className="nav-link" to="/profile/shared">Playlists Shared</Link></h2></div>
                    <div className="col-md-6"><h2><Link to="/profile/available">Playlists Available</Link></h2></div>
                  </div>

                  <div className="row">
                    <div className="col-md-4">
                      <h3>Likes: {user.likes.length}</h3>
                      <ul>
                        {user.likes.map((like) =>
                          <li key={like.ID}><Link to={like.permalink}>{like.playlistName}</Link></li>
                        )}
                      </ul>
                    </div>

                    <div className="col-md-4">
                      <h3>Comments: {user.comments.length}</h3>
                      <ul>
                        {user.comments.map((comment) =>
                            <li key={comment.ID}><Link to={comment.permalink}>{`${comment.content.slice(0,15)}...`}</Link></li>
                          )}
                      </ul>
                    </div>

                    <div className="col-md-4">
                      <h3>Shared: {stereodosePlaylists.length}</h3>
                      <ul>
                        {stereodosePlaylists.map((playlist) =>
                          <li key={playlist.spotifyID}><Link to={playlist.permalink}>{playlist.name}</Link></li>
                        )}
                      </ul>
                    </div>
                  </div>
                </div>

              }

            </div>
          </div>
        </div>
      )
    }
    return (
      <div>...loading</div>
    );
  }

  async componentDidMount() {
    try {
      await this.checkPlaylists();
      await this.fetchUserData();
    } catch (err) {
      alert(err.message);
    }
  }

  checkPlaylists = async () => {
    const SDK = new Spotify();
    const token = await this.props.app.getAccessToken();
    SDK.setAccessToken(token);

    let spotifyPlaylists = [];
    let allPlaylistsLoaded = false;
    let offset = 0;
    while (!allPlaylistsLoaded) {
      const userPlaylists = await SDK.getUserPlaylists({
        limit: 50,
        offset: offset
      });
      spotifyPlaylists = spotifyPlaylists.concat(userPlaylists.items);
      if (userPlaylists.items.length < 50) {
        allPlaylistsLoaded = true;
      }
      offset = offset + 50;
    }

    const response = await fetch("/api/playlists/me", { credentials: "same-origin" });
    if (response.status !== 200) {
      throw new Error(`${response.status} Unable to fetch user profile`);
    }
    const stereodosePlaylists = await response.json();

    const diffedSpotifyPlaylists = [];
    const diffedStereodosePlaylists = [];

    // const spotifyPlaylists = userPlaylists.items;
    for (let i = 0; i < spotifyPlaylists.length; i++) {
      let match = false;
      for (let j = 0; j < stereodosePlaylists.length; j++) {
        if (spotifyPlaylists[i].id === stereodosePlaylists[j].spotifyID) {
          diffedStereodosePlaylists.push(stereodosePlaylists[j]);
          match = true;
          break;
        }
      }
      if (match === false) {
        diffedSpotifyPlaylists.push(spotifyPlaylists[i]);
      }
    }

    this.setState({
      spotifyPlaylists: diffedSpotifyPlaylists,
      stereodosePlaylists: diffedStereodosePlaylists
    });
  }

  fetchUserData = async () => {
    const response = await fetch("/api/users/me", { credentials: "same-origin" });
    if (response.status !== 200) {
      throw new Error(`${response.status} Unable to fetch user profile`);
    }
    const user = await response.json();
    this.setState({ user: user });
  }
}

export default UserProfile