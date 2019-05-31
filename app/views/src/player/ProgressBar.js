import React from "react";
import { Slider, Handles, Tracks, Rail } from "react-compound-slider";
import "./ProgressBar.css";

// ProgressBar represents how far into the song you are
// It displays visually like a loading bar
export default class ProgressBar extends React.Component {

  onSlideEnd = async (values) => {
    await this.props.onSeek(values, this.props.duration);
  }

  render() {
    const progress = this.props.position / this.props.duration;
    const percentage = Math.round(progress * 1000) / 10;

    return (
      <Slider
        onSlideEnd={this.onSlideEnd}
        className="progress-bar-slider"
        domain={[0, 100]}
        values={[percentage]}
      >
        <Rail>
          {({ getRailProps }) => (  // adding the rail props sets up events on the rail
            <div className="progress-bar-rail" {...getRailProps()} />
          )}
        </Rail>
        <Handles>
          {({ handles, getHandleProps }) => (
            <div className="slider-handles">
              {handles.map(handle => (
                <Handle
                  key={handle.id}
                  handle={handle}
                  getHandleProps={getHandleProps}
                />
              ))}
            </div>
          )}
        </Handles>
        <Tracks right={false}>
          {({ tracks, getTrackProps }) => (
            <div className="slider-tracks">
              {tracks.map(({ id, source, target }) => (
                <Track
                  key={id}
                  source={source}
                  target={target}
                  getTrackProps={getTrackProps}
                />
              ))}
            </div>
          )}
        </Tracks>
      </Slider>
    )
  }
}

// Handle marks where the song currently is
// Visually, it's what you can click and drag to change where the song is playing from
function Handle(props) {
  const { id, percent } = props.handle; // handle also has 'value' prop
  const { getHandleProps } = props;
  return (
    <div className="progress-bar-handle"
      style={{ left: `${percent}%` }}
      {...getHandleProps(id)}
    >
    </div>
  )
}

// Track shows the actual progression
// - a visual representation of elapsed time
// - the bar to the left of the handle
function Track(props) {
  const { source, target, getTrackProps } = props;
  return (
    <div
      className="progress-bar-track"
      style={{
        left: `${source.percent}%`,
        width: `${target.percent - source.percent}%`
      }}
      {...getTrackProps()} // this will set up events if you want it to be clickeable (optional)
    />
  )
}