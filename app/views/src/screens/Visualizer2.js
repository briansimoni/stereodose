import React from 'react';
import Spotify from 'spotify-web-api-js';

export default class Visualizer2 extends React.Component {
  analysis
  features
  playerState
  currentTrackId // string spotify id
  intervalTypes = ['tatums', 'segments', 'beats', 'bars', 'sections']
  /**
   * 
   * @param {object} props
   * @param {App} props.app - a reference to the App component
   * @param {function} props.toggleVisualizer
   */
  constructor(props) {
    super(props);
    this.ref = React.createRef();
    this.canvas = document.createElement('canvas')
    this.ctx = this.canvas.getContext('2d')
    this.canvas.style.width = '500px';
    this.canvas.style.height = '500px';
    requestAnimationFrame(this.animate);
  }

  render() {
    return (
      <div
        id="visualizer"
        ref={nodeElement => {
          nodeElement && nodeElement.appendChild(this.canvas);
        }}
      >
        <button onClick={this.props.toggleVisualizer} id="close-visualzier-button" className="btn btn-danger">CLOSE</button>
      </div>
    )
  }

  // whenever the component mounts or the playback changes, we make sure to do the calls
  // to update the audio analysis. These are large network calls, so try to not make them too often
  async componentDidMount() {
    const playerState = await this.props.app.player.getCurrentState();
    await this.getTrackInfo(playerState.track_window.current_track.id);
    this.props.app.player.addListener('player_state_changed', async (playerState) => {
      if (this.currentTrack !== playerState.track_window.current_track.id) {
        this.playerState = playerState;
        this.currentTrackId = playerState.track_window.current_track.id;
        await this.getTrackInfo(this.currentTrackId);
      }
    })
  }

  componentWillUnmount() {

  }

  animate = () => {
    console.log('animating');
    this.ctx.save()
    this.ctx.fillStyle = this.fill
    this.ctx.fillRect(0, 0, 500, 500)
    this.ctx.restore()
    requestAnimationFrame(this.animate);
  }

  getTrackInfo = async (trackId) => {
    const accessToken = await this.props.app.getAccessToken();
    const spotify = new Spotify();
    spotify.setAccessToken(accessToken);
    try {
      const [analysis, features] = await Promise.all([
        spotify.getAudioAnalysisForTrack(trackId),
        spotify.getAudioFeaturesForTrack(trackId)
      ]);
      this.analysis = analysis;
      this.features = features;
    } catch (err) {
      console.error(err);
    }
  }

  setActiveIntervals = () => {

  }
}