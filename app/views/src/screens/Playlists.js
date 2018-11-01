import React from "react";
import { Link } from "react-router-dom";
import "./Playlists.css";

// Playlists renders the playlists that correspond to a particular drug + mood
class Playlists extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      error: null,
      playlists: null
    }
  }

  render() {
    let loading = this.state.loading;
    let err = this.state.error;
    let playlists = this.state.playlists;
    if (loading) {
      return <h3>Loading</h3>
    }

    if (err) {
      return <h3>Error: {err}</h3>
    }

    if (playlists) {
      let match = this.props.match;
      return (
        <div className="row">
          <div className="col">
            <h2 id="choose-a-playlist">Choose A Playlist</h2>
            <ul className="playlists">
              {playlists.map((playlist) => {
                return (
                  <Link key={playlist.spotifyID} to={`${match.url}/${playlist.spotifyID}`}>
                    <li><h4>{playlist.name}</h4></li>
                  </Link>
                )
              })}
            </ul>
          </div>
        </div>

      );
    }
  }

  // TODO: if drug or subcategory is not found, 404
  async componentDidMount() {
    let drug = this.props.match.params.drug;
    let subcategory = this.props.match.params.subcategory;

    try {
      let response = await fetch(`/api/playlists/?category=${drug}&subcategory=${subcategory}`, { credentials: "same-origin" });
      if (response.status !== 200) {
        throw new Error(`Error fetching playlists ${response.status}, ${response.statusText}`);
      }
      let playlists = await response.json();
      if (playlists.length === 0) {
        throw new Error(`No playlists found for drug: ${drug}, mood: ${subcategory}`);
      }
      this.setState({
        loading: false,
        playlists: playlists
      });
    } catch (err) {
      this.setState({
        loading: false,
        error: err.message
      });
    }
  }

}

export default Playlists