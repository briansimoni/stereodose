import React from "react";
import Spotify from "spotify-web-api-js";
import ShareSpotifyPlaylist from "./sharing/ShareSpotifyPlaylist";
import StereodosePlaylist from "./StereodosePlaylist";
import "./Profile.css";

class UserProfile extends React.Component {

  state = {
    spotifyPlaylists: null,
    stereodosePlaylists: null,
  }

  render() {
    const { spotifyPlaylists, stereodosePlaylists } = this.state;
    const categories = this.props.categories;

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

              {this.props.location.pathname === "/profile" &&
                <div>TODO: add some profile data stuff here</div>
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
      this.checkPlaylists();
    } catch (err) {
      alert(err.message);
    }
  }

  checkPlaylists = async () => {
    const SDK = new Spotify();
    const token = await this.props.getAccessToken();
    SDK.setAccessToken(token);
    const userPlaylists = await SDK.getUserPlaylists();

    const response = await fetch("/api/playlists/me", { credentials: "same-origin" });
    if (response.status !== 200) {
      throw new Error(`${response.status} Unable to fetch user profile`);
    }
    const stereodosePlaylists = await response.json();

    const diffedSpotifyPlaylists = [];
    const diffedStereodosePlaylists = [];

    const spotifyPlaylists = userPlaylists.items;
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
}

export default UserProfile