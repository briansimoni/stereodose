
import React from 'react';
import Spotify from 'spotify-web-api-js';
import './Visualizer.css'

// props need to include an App instance
export default class Visualizer extends React.Component {
  analysis = null
  features = null

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

  /**
   * setActiveIntervals should be called in the animate function
   * Thus we guarantee that every frame has the correct data available to it
   */
  setActiveIntervals = async () => {
    if (!this.analysis || !this.features) {
      await this.getTrackInfo();
    }
    console.log(this.analysis);

    const playerState = await this.props.app.player.getCurrentState();
    const songDuration = playerState.track_window.current_track.duration_ms;

    const segments = this.analysis.segments;
    const segmentIndex = Math.floor((playerState.position * segments.length) / songDuration);
    const currentSegment = segments[segmentIndex];
    this.currentSegment = currentSegment;

    const sections = this.analysis.sections;
    const sectionIndex = Math.floor((playerState.position * sections.length) / songDuration)
    const currentSection = sections[sectionIndex]
    this.currentSection = currentSection;
  }

  async componentDidMount() {
    const playerState = await this.props.app.player.getCurrentState();
    await this.getTrackInfo(playerState.track_window.current_track.id);
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