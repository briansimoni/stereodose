
import React from 'react';
import Spotify from 'spotify-web-api-js';
import './Visualizer.css'

// props need to include an App instance
export default class Visualizer extends React.Component {
  analysis = undefined;
  features = undefined;
  currentSegment = null;
  nextSegment = null;
  currentSection = null;
  fetchingData = false;
  componentID = Math.random()

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
      this.fetchingData = false;
    } catch (err) {
      console.error(err);
    }
  };

  /**
   * setActiveIntervals should be called in the animate function
   * Thus we guarantee that every frame has the correct data available to it
   */
  setActiveIntervals = async () => {
    if (!this.analysis) {
      console.log('fetching data', this.componentID);
      return;
    }


    const playerState = await this.props.app.player.getCurrentState();
    const songDuration = playerState.track_window.current_track.duration_ms;
    const position = playerState.position;
    const positionInSeconds = position / 1000;
    const segments = this.analysis.segments;

    // important to set the position every time.
    this.position = position;

    // these loops and conditions helps to synchronize segments to the position
    // this problem is hard to solve with proportions for some reason
    if (this.currentSegment === null) {
      // start at 1 because the first element's start is undefined
      for(let i = 1; i < segments.length; i++) {
        if (positionInSeconds > segments[i].start && positionInSeconds < segments[i+1].start) {
          this.currentSegmentIndex = i;
          break;
        }
      }
      this.currentSegment = segments[this.currentSegmentIndex]; // n
      this.nextSegment = segments[this.currentSegmentIndex + 1]; // n + 1
    }

    if (positionInSeconds >= this.nextSegment.start) {
      this.currentSegmentIndex++;
      this.currentSegment = segments[this.currentSegmentIndex]; // n
      this.nextSegment = segments[this.currentSegmentIndex + 1]; // n + 1

    }


    // sections is calculated with a proportion and may not be at a desirable level of synchronization
    const sections = this.analysis.sections;
    const sectionIndex = Math.floor((playerState.position * sections.length) / songDuration)
    const currentSection = sections[sectionIndex]
    this.currentSection = currentSection;
  }

  /**
   * For this function to work properly, the active intervals need to have been set.
   * It will take the current segment and the next segment and calculate the slope-intercept (y=mx+b)
   * equation between the volume in decibels between the two segments. It will plug in
   * the current position of the song into the equation to get the approximate volume at any given time
   */
  getVolume = () => {
    if (this.currentSegment) {
      const x = this.position;
      let { loudness_start, start, duration } = this.currentSegment;
      let loudness_end = this.nextSegment.loudness_start;
      loudness_start = Math.abs(loudness_start);
      loudness_end = Math.abs(loudness_end);

      const x1 = start * 1000; // convert seconds to ms
      const y1 = loudness_start;
      const x2 = (start + duration) * 1000 // convert seconds to ms
      const y2 = loudness_end;

      const m = (y2 - y1) / (x2 - x1);
      const b = y1 - (m * x1);

      return 100 - (m * x + b);
    }
  }

  async componentDidMount() {
    const playerState = await this.props.app.player.getCurrentState();
    if (!this.fetchingData) {
      console.log('getting the info');
      this.getTrackInfo(playerState.track_window.current_track.id);
    }
    this.fetchingData = true;
    this.props.app.player.addListener('player_state_changed', this.playerStateChangeCallback);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeCanvas);
    this.props.app.player.removeListener('player_state_changed', this.playerStateChangeCallback);
  }

  playerStateChangeCallback = async playerState => {
    if (this.currentTrack !== playerState.track_window.current_track.id) {
      this.playerState = playerState;
      this.currentTrackId = playerState.track_window.current_track.id;
      await this.getTrackInfo(this.currentTrackId);
    }
  };
}