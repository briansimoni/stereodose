import React from "react";
import { Link } from "react-router-dom";
import "./Screens.css";
import Pagination from "./Pagination";

// Playlists renders the playlists that correspond to a particular drug + mood
class Playlists extends React.Component {

  resultsPerPage = 15;

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
      return (
        <div className="row justify-content-md-center">
          <div className="spinner-grow text-success text-center" role="status">
            <span className="sr-only">Loading...</span>
          </div>
        </div>
      )
    }

    if (err) {
      return <h3>Error: {err}</h3>
    }

    if (playlists) {
      const match = this.props.match;
      // reduce a large array into multidimensinal array
      // where we have m x 3 matrix (m rows of 3 columns)
      // this makes it way easier to render with react
      // With Bootstrap remember total row width is 12 columns.
      // So columns of length 4 mean you get 3 columns per row
      // A slice of total playlists is used to allow pagination to work correctly
      // Only the first 15 (or whatever the desired number of results per page)
      // is taken as a slice
      const playlistsSlice = playlists.slice(0, this.resultsPerPage);
      const rows = playlistsSlice.reduce((accumulator, currentPlaylist, index) => {
        if (index % 3 === 0) {
          return accumulator.concat([playlists.slice(index, index + 3)])
        }
        return accumulator;
      }, [])

      return (
        <div className="playlists">
          <h2 id="choose-a-playlist">Choose A Playlist</h2>

          {rows.map((row, index) => {
            return (
              <div className="row" key={index}>
                {row.map((playlist) => {
                  const thumbnailImageURL = playlist.bucketThumbnailURL ? playlist.bucketThumbnailURL : "https://via.placeholder.com/250x200";
                  return (
                    <div className="col-md-4" key={playlist.spotifyID}>
                      <Link to={`${match.url}/${playlist.spotifyID}`}>
                        <img src={thumbnailImageURL} alt="playlist-artwork" />
                      </Link>

                      <Link to={`${match.url}/${playlist.spotifyID}`}>
                        <h4>{playlist.name}</h4>
                      </Link>
                    </div>
                  )
                })}
              </div>
            )
          })}

          <Pagination resultsPerPage={this.resultsPerPage} match={match} playlists={playlists} />


        </div>

      );
    }
  }

  // when this component loads, grab the query params for pagination
  async componentDidMount() {
    try {
      const playlists = await this.fetchPlaylists();
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

  // updating the page queryParam seems to trigger an update
  async componentDidUpdate() {
    try {
      const oldPlaylists = this.state.playlists;
      if (!oldPlaylists) {
        return;
      }
      const newPlaylists = await this.fetchPlaylists();
      // maybe want to deep equal instead of JSON.stringify
      if (JSON.stringify(newPlaylists) === JSON.stringify(oldPlaylists)) {
        return;
      }
      this.setState({ playlists: newPlaylists });
    } catch (err) {
      this.setState({
        loading: false,
        error: err.message
      });
    }
  }

  fetchPlaylists = async () => {
    let drug = this.props.match.params.drug;
    let subcategory = this.props.match.params.subcategory;

    const params = new URLSearchParams(window.location.search);
    const page = params.get('page') !== null ? parseInt(params.get('page')) : 1;
    const offset = page === 1 ? 0 : (page * this.resultsPerPage) - this.resultsPerPage;

    const response = await fetch(`/api/playlists/?category=${drug}&subcategory=${subcategory}&limit=20&offset=${offset}`, { credentials: "same-origin" });
    if (response.status !== 200) {
      throw new Error(`Error fetching playlists ${response.status}, ${response.statusText}`);
    }
    const playlists = await response.json();
    if (playlists.length === 0) {
      throw new Error(`No playlists found for drug: ${drug}, mood: ${subcategory}`);
    }
    return playlists;
  }

}

export default Playlists