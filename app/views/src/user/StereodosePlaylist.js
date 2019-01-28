import React from "react";
// react-octicons increases build artifact by over 20 kb
import {TrashcanIcon} from 'react-octicons'

class StereodosePlaylist extends React.Component {
  constructor(props) {
    super(props);
    this.deleteFromStereodose = this.deleteFromStereodose.bind(this);
  }
  render() {
    return (
      <tr>
        <td>{this.props.playlist.name}</td>
        <td>{this.props.playlist.category}</td>
        <td>{this.props.playlist.subCategory}</td>
        <td>
          <button type="button" className="btn btn-danger" onClick={this.deleteFromStereodose}>
            <TrashcanIcon/>
          </button>
        </td>
      </tr>
    )
  }

  async deleteFromStereodose() {
    let id = this.props.playlist.spotifyID;
    let resp = await fetch(`/api/playlists/${id}`, {
      method: "DELETE",
      credentials: "same-origin"
    });
    if (resp.status !== 200) {
      alert(resp.status + " " + resp.statusText);
    }
    this.props.onUpdate();
  }
}

export default StereodosePlaylist;