
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
  fetchingData = true;

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
    if (this.fetchingData) {
      console.log('fetching data');
      return;
    }


    const playerState = await this.props.app.player.getCurrentState();
    const songDuration = playerState.track_window.current_track.duration_ms;
    const position = playerState.position;
    const positionInSeconds = position / 1000;
    const segments = this.analysis.segments;

    // important to set the position everytime.
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
    // if (position >= this.nextSegment.start) {
    //   console.log('bigger');
    //   this.currentSegmentIndex++;
    //   this.currentSegment = segments[this.currentSegmentIndex]; // n
    //   this.nextSegment = segments[this.currentSegmentIndex + 1]; // n + 1
    // }



    // const segmentIndex = Math.floor((playerState.position * segments.length) / songDuration);
    // const currentSegment = segments[segmentIndex];
    // this.currentSegment = currentSegment;
    // if (segments[segmentIndex +1] !== undefined) {
    //   this.nextSegment = segments[segmentIndex+1];
    // }

    const sections = this.analysis.sections;
    const sectionIndex = Math.floor((playerState.position * sections.length) / songDuration)
    const currentSection = sections[sectionIndex]
    this.currentSection = currentSection;
  }

  async componentDidMount() {
    const playerState = await this.props.app.player.getCurrentState();
    this.fetchingData = true;
    this.getTrackInfo(playerState.track_window.current_track.id);
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