import React from "react";
import Spotify from "spotify-web-api-js";
import ShareSpotifyPlaylist from "./Sharing/ShareSpotifyPlaylist";
import StereodosePlaylist from "./StereodosePlaylist";
import Tabs from "./Tabs";
import "./Profile.css";

class UserProfile extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      spotifyPlaylists: null,
      stereodosePlaylists: null,
      categories: null,
      loading: true
    }
  }

  render() {
    let { spotifyPlaylists, stereodosePlaylists, categories, loading } = this.state;
    if (spotifyPlaylists !== null && stereodosePlaylists !== null && !loading) {
      return (
        <div>
          {/* each child of <Tabs> needs to be a <div> with a label attribute*/}
          <Tabs>
            <div label="Playlists Shared to Stereodose">
              <h2 id="tab-content-title">Playlists Shared to Stereodose</h2>
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

            <div label="Playlists Available">
              <ShareSpotifyPlaylist
                playlists={spotifyPlaylists}
                categories={categories}
                onUpdate={() => { this.checkPlaylists() }}
              />
            </div>
          </Tabs>
        </div>
      )
    }
    return (
      <div>
        <h2>Loading...</h2>
      </div>
    );
  }

  async componentDidMount() {
    let resp = await fetch("/api/categories/", { credentials: "same-origin" });
    let categories = await resp.json();
    let state = this.state;
    state.categories = categories;
    this.setState(state);
    this.checkPlaylists();
  }

  checkPlaylists = async() => {
    let SDK = new Spotify();
    // TODO: catch errors here
    let token = await this.props.getAccessToken();
    SDK.setAccessToken(token);
    let userPlaylists = await SDK.getUserPlaylists();

    let response = await fetch("/api/playlists/me", { credentials: "same-origin" });
    let stereodosePlaylists = await response.json();

    let diffedSpotifyPlaylists = [];
    let diffedStereodosePlaylists = [];

    // so old school
    let spotifyPlaylists = userPlaylists.items;
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

    let state = this.state;
    state.spotifyPlaylists = diffedSpotifyPlaylists;
    state.stereodosePlaylists = diffedStereodosePlaylists;
    state.loading = false;
    this.setState(state);
  }
}

export default UserProfile