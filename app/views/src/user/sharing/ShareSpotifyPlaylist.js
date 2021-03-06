import React from 'react';
import { Fragment } from 'react';
import PickPlaylist from './PickPlaylist';
import PickDrug from './PickDrug';
import PickMood from './PickMood';
import PickImage from './PickImage';
import SharePlaylistButton from './SharePlaylistButton';

// ShareSpotifyPlaylist is the component that allows users to share Spotify Playlists to Stereodose
// It will maintain the state of a particular playlist as a user selects different options before uploading
// The parent component (ShareController) is responsible for passing required information
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
      inFlight: false
    };
  }

  render() {
    const { playlists, categories } = this.props;
    const { selectedPlaylist, selectedDrug, selectedMood, imageBlob } = this.state;
    if (!playlists || !categories) {
      return <div />;
    }

    if (selectedMood) {
      let buttonDisabled = true;
      if (imageBlob) {
        buttonDisabled = false;
      }

      return (
        <Fragment>
          <div>
            <h2 id="content-title">Upload An Image</h2>
          </div>

          <div>
            <PickImage onBlobCreated={this.onBlobCreated} />
          </div>

          <div>
            <h4>Playlist: {this.state.selectedPlaylist.name}</h4>
            <h4>Drug: {this.state.selectedDrug}</h4>
            <h4>Mood: {this.state.selectedMood}</h4>
          </div>

          <div className="text-center">
            <SharePlaylistButton
              disabled={buttonDisabled}
              onShareStereodose={this.shareToStereodose}
              inFlight={this.state.inFlight}
            />
          </div>

          <div className="cancel text-center">
            <button className="btn btn-danger" onClick={this.cancel}>
              Cancel
            </button>
          </div>
        </Fragment>
      );
    }

    if (selectedPlaylist && selectedDrug) {
      return (
        <Fragment>
          <PickMood
            onSelectMood={this.onSelectMood}
            categories={categories}
            playlist={selectedPlaylist}
            drug={selectedDrug}
          />
          <div className="cancel text-center">
            <button className="btn btn-danger" onClick={this.cancel}>
              Cancel
            </button>
          </div>
        </Fragment>
      );
    }

    if (selectedPlaylist) {
      return (
        <Fragment>
          <PickDrug onSelectDrug={this.onSelectDrug} categories={categories} playlist={selectedPlaylist} />
          <div className="cancel text-center">
            <button className="btn btn-danger" onClick={this.cancel}>
              Cancel
            </button>
          </div>
        </Fragment>
      );
    }

    if (playlists) {
      return <PickPlaylist onSelectPlaylist={this.onSelectPlaylist} playlists={playlists} />;
    }
  }

  onBlobCreated = blob => {
    this.setState({ imageBlob: blob });
  };

  onSelectPlaylist = playlist => {
    this.setState({ selectedPlaylist: playlist });
  };

  onSelectDrug = drug => {
    this.setState({ selectedDrug: drug });
  };

  onSelectMood = mood => {
    this.setState({ selectedMood: mood });
  };

  uploadImage = async blob => {
    const data = new FormData();
    data.append('playlist-image', blob);
    data.append('filename', 'playlist-image');
    const response = await fetch(`/api/playlists/image`, {
      method: 'POST',
      body: data
    });
    if (response.status !== 201) {
      const errorMessage = await response.text();
      throw new Error(`Problem uploading image, ${errorMessage} ${response.status}: ${response.statusText}`);
    }
    const json = await response.json();
    return json;
  };

  // shareToStereodose kicks off two requests
  // First it tries to upload the image. Once it is successful, it will take the permalink of the image
  // in the returned JSON and use that in the subsequent request to create the playlist in the database
  // this works well at the consequence of having a bunch of trash images in S3.
  // If the image upload is successful, but the playlist upload is not, there will just be an image
  // out there with nothing pointing to it.
  shareToStereodose = async () => {
    // disable button while request is in flight
    if (this.state.inFlight) {
      return;
    }
    this.setState({ inFlight: true });
    let imageURL, thumbnailURL;
    try {
      const uploadResult = await this.uploadImage(this.state.imageBlob);
      imageURL = uploadResult.imageURL;
      thumbnailURL = uploadResult.thumbnailURL;
    } catch (error) {
      this.setState({ inFlight: false });
      alert(error.message);
      return;
    }

    const { selectedPlaylist, selectedMood, selectedDrug } = this.state;
    const resp = await fetch(`/api/playlists/`, {
      method: 'POST',
      credentials: 'same-origin',
      body: JSON.stringify({
        SpotifyID: selectedPlaylist.id,
        Category: selectedDrug,
        Subcategory: selectedMood,
        ImageURL: imageURL,
        ThumbnailURL: thumbnailURL
      })
    });
    if (resp.status !== 201) {
      const errorMessage = await resp.text();
      alert(`Problem sharing playlist: ${errorMessage}`);
      this.setState({
        selectedPlaylist: null,
        selectedDrug: null,
        selectedMood: null,
        imagePermaLink: null,
        inFlight: false
      });
      this.props.onUpdate();
      return;
    }

    alert('Share Successful!');
    this.setState({
      selectedPlaylist: null,
      selectedDrug: null,
      selectedMood: null,
      imagePermaLink: null,
      inFlight: false
    });
    this.props.onUpdate();
  };

  cancel = () => {
    this.setState({
      selectedPlaylist: null,
      selectedDrug: null,
      selectedMood: null,
      imageBlob: null
    });
  };
}

export default ShareSpotifyPlaylist;
