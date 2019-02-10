import React from "react";
import { Fragment } from "react";
import PickPlaylist from "./PickPlaylist";
import PickDrug from "./PickDrug";
import PickMood from "./PickMood";
import PickImage from "./PickImage";
import SharePlaylistButton from "./SharePlaylistButton";

// ShareSpotifyPlaylist is the component that allows users to share Spotify Playlists to Stereodose
// It will maintain the state of a particular playlist as a user selects different options before uploading
// The parent component (Profile) is responsible for passing required information
// e.g. access tokens, the set of playlists, moods, categories, update callbacks, etc...
class ShareSpotifyPlaylist extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      selectedPlaylist: null,
      selectedDrug: null,
      selectedMood: null,
      imageBlob: null,
      // inFlight is a boolean to usesd to disable stuff while requests are pending
      inFlight: false,
    }

    this.onSelectPlaylist = this.onSelectPlaylist.bind(this);
    this.onSelectDrug = this.onSelectDrug.bind(this);
    this.onSelectMood = this.onSelectMood.bind(this);
    this.uploadImage = this.uploadImage.bind(this);
    this.shareToStereodose = this.shareToStereodose.bind(this);
  }

  render() {
    const { playlists, categories } = this.props;
    const { selectedPlaylist, selectedDrug, selectedMood, imageBlob } = this.state;
    if (!playlists) {
      return <div></div>
    }

    if (selectedMood) {
      let buttonDisabled = true;
      if (imageBlob) {
        buttonDisabled = false;
      }

      return (
        <Fragment>
          <PickImage onBlobCreated={this.onBlobCreated} />
          <SharePlaylistButton
            disabled={buttonDisabled}
            onShareStereodose={this.shareToStereodose}
            inFlight={this.state.inFlight} />
        </Fragment>
      )
    }

    if (selectedPlaylist && selectedDrug) {
      return (
        <PickMood
          onSelectMood={this.onSelectMood}
          categories={categories}
          playlist={selectedPlaylist}
          drug={selectedDrug}
        />
      );
    }

    if (selectedPlaylist) {
      return (
        <PickDrug
          onSelectDrug={this.onSelectDrug}
          categories={categories}
          playlist={selectedPlaylist}
        />
      );
    }



    if (playlists) {
      return <PickPlaylist onSelectPlaylist={this.onSelectPlaylist} playlists={playlists} />
    }
  }

  onBlobCreated = blob => {
    this.setState({imageBlob: blob});
  }

  onSelectPlaylist(playlist) {
    this.setState({ selectedPlaylist: playlist });
  }

  onSelectDrug(drug) {
    this.setState({ selectedDrug: drug });
  }

  onSelectMood(mood) {
    this.setState({ selectedMood: mood })
  }

  async uploadImage(blob) {
      const data = new FormData();
      data.append('playlist-image', blob);
      data.append('filename', 'playlist-image');
      let response = await fetch(`/api/playlists/${this.state.selectedPlaylist.id}/image`, {
        method: "POST",
        body: data
      });
      if (response.status !== 201) {
        const errorMessage = await response.text();
        throw new Error(`Problem uploading image, ${errorMessage} ${response.status}: ${response.statusText}`);
      }
      const json = await response.json();
      return json.imageURL;
  }

  async shareToStereodose() {
    // disable button while request is in flight
    if (this.state.inFlight) {
      return;
    }
    this.setState({inFlight: true});
    const imageURL = await this.uploadImage(this.state.imageBlob);

    const { selectedPlaylist, selectedMood, selectedDrug } = this.state;
    let resp = await fetch(`/api/playlists/`, {
      method: "POST",
      credentials: "same-origin",
      body: JSON.stringify({
        SpotifyID: selectedPlaylist.id,
        Category: selectedDrug,
        Subcategory: selectedMood,
        ImageURL: imageURL
      })
    });
    if (resp.status !== 201) {
      alert("error! " + resp.status + " " + resp.statusText);
    }
    alert("Share Successful!");
    this.setState({
      selectedPlaylist: null,
      selectedDrug: null,
      selectedMood: null,
      imagePermaLink: null,
      inFlight: false
    });
    this.props.onUpdate();
  }

}

export default ShareSpotifyPlaylist;