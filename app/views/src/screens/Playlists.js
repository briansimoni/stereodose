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

  componentDidMount() {
    let drug = this.props.match.params.drug;
    let subcategory = this.props.match.params.subcategory;

    fetch(`/api/playlists/?category=${drug}&subcategory=${subcategory}`, { credentials: "same-origin" })
      .then((response) => {
        if (response.status === 200) {
          return response.json();
        }
        return Promise.reject(`${response.status} ${response.statusText}`)
      })
      .then((json) => {
        this.setState({
          loading: false,
          playlists: json
        })
      })
      .catch((err) => {
        this.setState({
          loading: false,
          error: err
        })
      });
  }
}

export default Playlists