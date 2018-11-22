import React from "react";
import { Link } from "react-router-dom";
import "./Screens.css";

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
      const match = this.props.match;
      // reduce a large array into multidimensinal array
      // where we have m x 3 matrix (m rows of 3 columns)
      // this makes it way easier to render with react
      // With Bootstrap remember total row width is 12 columns.
      // So columns of length 4 mean you get 3 columns per row
      const rows = playlists.reduce( (accumulator, currentPlaylist, index) =>{
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
                      {row.map( (playlist) => {
                        const bucketImageURL = playlist.bucketImageURL ?  playlist.bucketImageURL: "https://via.placeholder.com/250x200";
                        return (
                          <div className="col-md-4" key={playlist.spotifyID}>
                            <img src={bucketImageURL} alt="playlist-artwork"/> 
                              <Link to={`${match.url}/${playlist.spotifyID}`}>
                                <h4>{playlist.name}</h4>
                              </Link>
                          </div>
                        )
                      })}
                    </div>
                  )
              })}
            </div>

      );
    }
  }

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