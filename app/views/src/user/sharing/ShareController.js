import React from 'react';
import Spotify from 'spotify-web-api-js';
import ShareSpotifyPlaylist from './ShareSpotifyPlaylist';
import '../Profile.css';

// ShareController coordinates and holds data to most of the things
// related to sharing a playlist
class ShareController extends React.Component {
  state = {
    spotifyPlaylists: null,
    stereodosePlaylists: null,
    user: null
  };

  render() {
    const { spotifyPlaylists, stereodosePlaylists } = this.state;
    const categories = this.props.app.state.categories;

    if (spotifyPlaylists && stereodosePlaylists && categories) {
      return (
        <div className="container profile">

          {this.props.location.pathname === '/profile/available' && (
            <div label="Playlists Available" className="container">
              <div className="row justify-content-md-center">
                <div className="col col-md-auto">
                  <ShareSpotifyPlaylist
                    playlists={spotifyPlaylists}
                    categories={categories}
                    onUpdate={() => {
                      this.checkPlaylists();
                    }}
                  />
                </div>
              </div>
            </div>
          )}

        </div>
      );
    }
    return (
      <div className="row justify-content-md-center">
        <div className="spinner-grow text-center" role="status">
          <span className="sr-only">Loading...</span>
        </div>
      </div>
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

    const spotifyPlaylistIds = spotifyPlaylists.map((playlist) => playlist.id);
    const idQueryString = spotifyPlaylistIds.join(' ');

    const response = await fetch(`/api/playlists/?spotify-ids=${idQueryString}`, { credentials: 'same-origin' });
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
  };

  fetchUserData = async () => {
    const response = await fetch('/api/users/me', { credentials: 'same-origin' });
    if (response.status !== 200) {
      throw new Error(`${response.status} Unable to fetch user profile`);
    }
    const user = await response.json();
    this.setState({ user: user });
  };
}

export default ShareController;
