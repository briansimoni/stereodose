import React from 'react';
import Visualizer from './Visualizer';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMusic } from '@fortawesome/free-solid-svg-icons';

export default class Data2D extends Visualizer{

  constructor(props) {
    super(props);
    this.ref = React.createRef();
    this.canvas = document.createElement('canvas');
    this.ctx = this.canvas.getContext('2d');
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight - 125;

    this.ctx.beginPath();

    window.addEventListener('resize', this.resizeCanvas);
    requestAnimationFrame(this.animate);

    this.x = 0;

    this.volumePoint = new Point(0, window.innerHeight / 2 - 100);
    this.pitchPoint = new Point(0, window.innerHeight / 2 - 50);

    this.sweeping = false;
  }

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

        <div id="y-axis" />

        <div id="pitch-Ab" />
        <div id="pitch-G" />
        <div id="pitch-Gb" />
        <div id="pitch-F" />
        <div id="pitch-E" />
        <div id="pitch-Eb" />
        <div id="pitch-D" />
        <div id="pitch-Db" />
        <div id="pitch-C" />
        <div id="pitch-B" />
        <div id="pitch-Bb" />
        <div id="pitch-A" />

        <FontAwesomeIcon id="beat-icon" className="active" icon={faMusic} />

      </div>
    );
  }

  resizeCanvas = () => {
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight - 125;
  };

  animate = async () => {
    requestAnimationFrame(this.animate);
    if (!this.analysis) {
      return;
    }
    await this.setActiveIntervals();
    
    const { currentSegment } = this;
    const { currentSection } = this;

    const playerState = await this.props.app.player.getCurrentState();
    const songDuration = playerState.track_window.current_track.duration_ms;

    const mode = this.getMode(currentSection.mode);

    const beats = this.analysis.beats;
    const beatIndex = Math.floor((playerState.position * beats.length) / songDuration);
    const currentBeat = beats[beatIndex];
    document.getElementById('beat-icon').style['font-size'] = currentBeat.confidence * 100

    const modePlaceHolderButton = document.getElementById('temp-mode-indicator');
    if (modePlaceHolderButton) {
      modePlaceHolderButton.innerHTML = mode;
    }

    if (this.x >= window.innerWidth) {
      this.x = -1;
      this.volumePoint = new Point(0, 0);
      this.pitchPoint = new Point(0, 0);
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
    if (pitchPlaceHolderButton) {
      pitchPlaceHolderButton.innerHTML = pitch;
    }

    // const y = window.innerHeight - 20;
    this.ctx.beginPath();
    this.ctx.moveTo(this.pitchPoint.x, this.pitchPoint.y);
    // const pitchMax = Math.max(...segment.pitches)
    // const pitchIndex = segment.pitches.indexOf(pitchMax);
    // const amplitude = segment.pitches[pitchIndex];
    const amplitude = this.getPitchAmplitude(segment.pitches)

    this.ctx.lineTo(this.x, amplitude);
    this.ctx.strokeStyle = 'green';
    this.pitchPoint.x = this.x;
    this.pitchPoint.y = amplitude;
    this.ctx.closePath();
    this.ctx.stroke();

  };


  /**
 * @param {Array} pitches an array of pitches from a segment object
 * this will return a number in pixels of the y coordinate of what the pitch
 * line should be
 */
  getPitchAmplitude = (pitches) => {
    const pitch = Math.max(...pitches);
    const pitchIndex = pitches.indexOf(pitch)
    // height of the screen minus the height of the player controls
    // divided by 2 (half of that height) with a 5 % padding on both
    // top and bottom (90%) divided by 12 distinct pitches
    const multiple = (((window.innerHeight - 125) / 2) * .9) / 12;
    const p = {
      0: multiple * 9, // "C",
      1: multiple * 8, // "Db",
      2: multiple * 7, // "D",
      3: multiple * 6, // "Eb",
      4: multiple * 5, // "E",
      5: multiple * 4, // "F",
      6: multiple * 3, // "Gb",
      7: multiple * 2, // "G",
      8: multiple * 1, // "Ab",
      9: multiple * 12, // "A",
      10: multiple * 11, // "Bb",
      11: multiple * 10, // "B",
    } 
    return p[pitchIndex];
  }

  
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
