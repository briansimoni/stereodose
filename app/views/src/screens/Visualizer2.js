import React from 'react';
import Spotify from 'spotify-web-api-js';

export default class Visualizer2 extends React.Component {
  analysis;
  features;
  playerState;
  currentTrackId; // string spotify id
  intervalTypes = ['tatums', 'segments', 'beats', 'bars', 'sections'];
  /**
   *
   * @param {object} props
   * @param {App} props.app - a reference to the App component
   * @param {function} props.toggleVisualizer
   */
  constructor(props) {
    super(props);
    this.ref = React.createRef();
    this.canvas = document.createElement('canvas');
    this.ctx = this.canvas.getContext('2d');
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight;

    this.ctx.beginPath();

    window.addEventListener('resize', this.resizeCanvas);
    requestAnimationFrame(this.animate);

    this.x = 0;

    this.volumePoint = new Point(0, window.innerHeight / 2 - 100);
    this.pitchPoint = new Point(0, window.innerHeight / 2 - 50);

    this.sweeping = false;
  }

  resizeCanvas = () => {
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight;
  };

  render() {
    return (
      <div
        id="visualizer"
        ref={nodeElement => {
          nodeElement && nodeElement.appendChild(this.canvas);
        }}
      >
        <button onClick={this.props.toggleVisualizer} id="close-visualizer-button" className="btn btn-danger">
          CLOSE
        </button>

        <button id="temp-pitch-indicator" className="btn btn-success">
          pitch placeholder
        </button>

        <button id="temp-mode-indicator" className="btn btn-success">
          mode placeholder
        </button>
      </div>
    );
  }

  // whenever the component mounts or the playback changes, we make sure to do the calls
  // to update the audio analysis. These are large network calls, so try to not make them too often
  async componentDidMount() {
    const playerState = await this.props.app.player.getCurrentState();
    await this.getTrackInfo(playerState.track_window.current_track.id);
    this.props.app.player.addListener('player_state_changed', this.playerStateChangeCallback);
  }

  playerStateChangeCallback = async playerState => {
    if (this.currentTrack !== playerState.track_window.current_track.id) {
      this.playerState = playerState;
      this.currentTrackId = playerState.track_window.current_track.id;
      await this.getTrackInfo(this.currentTrackId);
    }
  };

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeCanvas);
    this.props.app.player.removeListener('player_state_changed', this.playerStateChangeCallback);
  }

  animate = async () => {
    requestAnimationFrame(this.animate);
    if (!this.analysis) {
      return;
    }
    const segments = this.analysis.segments;
    const playerState = await this.props.app.player.getCurrentState();
    const songDuration = playerState.track_window.current_track.duration_ms;
    const segmentIndex = Math.floor((playerState.position * segments.length) / songDuration);
    const currentSegment = segments[segmentIndex];


    const sections = this.analysis.sections;
    const sectionIndex = Math.floor((playerState.position * sections.length) / songDuration)
    const currentSection = sections[sectionIndex]
    const mode = this.getMode(currentSection.mode);

    const modePlaceHolderButton = document.getElementById('temp-mode-indicator');
    if(modePlaceHolderButton) {
      modePlaceHolderButton.innerHTML = mode;
    }

    console.log(segmentIndex, currentSegment.start * 1000, playerState.position);
    if (this.x >= window.innerWidth) {
      this.x = -1;
      this.volumePoint = new Point(0,0);
      this.pitchPoint = new Point(0,0);
    }

    // TODO: remove this hot fix
    if (!currentSegment) {
      return;
    }

    this.sweep();
    this.animateVolume(currentSegment);
    this.animatePitch(currentSegment);
    this.x += 1;
  };

  animateVolume = (segment) => {
    const y = window.innerHeight / 2 - 100;
    this.ctx.beginPath();
    this.ctx.moveTo(this.volumePoint.x, this.volumePoint.y);
    const amplitude = Math.abs(segment.loudness_start) * 5;
    this.ctx.lineTo(this.x, y + amplitude);
    this.ctx.strokeStyle = 'blue';
    this.volumePoint.x = this.x;
    this.volumePoint.y = y + amplitude;
    this.ctx.closePath();
    this.ctx.stroke();
  };

  animatePitch = (segment) => {
    const pitchPlaceHolderButton = document.getElementById('temp-pitch-indicator');
    const pitch = this.getPitch(segment.pitches);
    if(pitchPlaceHolderButton) {
      pitchPlaceHolderButton.innerHTML = pitch;
    }
    console.log(pitch);
    const y = window.innerHeight / 2 - 50;
    this.ctx.beginPath();
    this.ctx.moveTo(this.pitchPoint.x, this.pitchPoint.y);
    const amplitude = Math.abs(segment.pitches[0]) * 100;
    this.ctx.lineTo(this.x, y + amplitude);
    this.ctx.strokeStyle = 'red';
    this.pitchPoint.x = this.x;
    this.pitchPoint.y = y + amplitude;
    this.ctx.closePath();
    this.ctx.stroke();
  };

  /**
   * @param {Array} pitches an array of pitches from a segment object
   */
  getPitch = (pitches) => {
    const pitch = Math.max(...pitches);
    const pitchIndex = pitches.indexOf(pitch)
    const p = {
      0: "C",
      1: "Db",
      2: "D",
      3: "Eb",
      4: "E",
      5: "F",
      6: "Gb",
      7: "G",
      8: "Ab",
      9: "A",
      10: "Bb",
      11: "B",
    }
    return p[pitchIndex];
  }

  sweep = () => {
    this.ctx.clearRect(this.x + 1, 0, 20, window.innerHeight);
  };

  getTrackInfo = async trackId => {
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
  };

  setActiveIntervals = () => {};

  /**
   * @param {number} mode
   */
  getMode = (mode) => {
    const m = {
      0: "minor",
      1: "major",
    }
    return m[mode];
  }
}

class Point {
  /**
   * @param {number} x
   * @param {number} y
   */
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }
}
