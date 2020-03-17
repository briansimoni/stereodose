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
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight;

    this.ctx.beginPath();

    window.addEventListener('resize', this.resizeCanvas);
    requestAnimationFrame(this.animate);

    this.x = 0;
    this.y = 0;
    this.sweeping = false;
  }

  resizeCanvas = () => {
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight;
  }

  render() {
    return (
      <div
        id="visualizer"
        ref={nodeElement => {
          nodeElement && nodeElement.appendChild(this.canvas);
        }}
      >
        <button onClick={this.props.toggleVisualizer} id="close-visualizer-button" className="btn btn-danger">CLOSE</button>
      </div>
    )
  }

  // whenever the component mounts or the playback changes, we make sure to do the calls
  // to update the audio analysis. These are large network calls, so try to not make them too often
  async componentDidMount() {
    const playerState = await this.props.app.player.getCurrentState();
    await this.getTrackInfo(playerState.track_window.current_track.id);
    this.props.app.player.addListener('player_state_changed', this.playerStateChangeCallback)
  }

  playerStateChangeCallback = async (playerState) => {
    if (this.currentTrack !== playerState.track_window.current_track.id) {
      this.playerState = playerState;
      this.currentTrackId = playerState.track_window.current_track.id;
      await this.getTrackInfo(this.currentTrackId);
    }
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeCanvas);
    this.props.app.player.removeListener('player_state_changed', this.playerStateChangeCallback)
  }

  animate = async () => {
    requestAnimationFrame(this.animate);
    if (!this.analysis) {
      return;
    }
    const segments = this.analysis.segments;
    const playerState = await this.props.app.player.getCurrentState();
    const songDuration = playerState.track_window.current_track.duration_ms;
    const segmentIndex = Math.floor((segments.length / songDuration) * playerState.position);
    const currentSegment = segments[segmentIndex];
    if (this.x >= window.innerWidth) {
      this.ctx.closePath();
      this.x = -1;
      this.y += 10;
      this.ctx.beginPath();
    }


    const y = (window.innerHeight / 2) -100;
    const amplitude = (Math.abs(currentSegment.loudness_start)) * 5
    console.log(amplitude, currentSegment.loudness_start);
    this.sweep();
    this.ctx.lineTo(this.x, y + amplitude);
    this.ctx.strokeStyle = 'blue';
    this.ctx.stroke()
    this.x += 1;
  }

  sweep = () => {
    this.ctx.clearRect(this.x + 1, 0, 20, window.innerHeight);
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