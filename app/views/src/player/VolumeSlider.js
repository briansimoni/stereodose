import React from "react";
import Slider, { Rail, Handles, Tracks } from 'react-compound-slider';
import Octicon from "react-octicon";

const railStyle = {
  position: 'absolute',
  width: '100%',
  height: 10,
  marginTop: 35,
  borderRadius: 5,
  backgroundColor: '#8B9CB6',
}

const defaultValues = [100]

class VolumeSlider extends React.Component {

  state = {
    values: defaultValues,
    update: defaultValues,
  }

  onUpdate = update => {
    // nothing right now
  }

  onChange = values => {
    const percent = Math.round(values[0]);
    this.props.onChangeVolume(percent);
  }

  render() {
    const values = this.state.values;
    const className = this.props.className;

    return (
      <Slider
        className={className}
        domain={[0, 100]}
        values={values}
        onChange={this.onChange}
        disabled={this.props.disabled}
      >
      <Octicon className="unmute" name="unmute"/>
        <Rail>
          {({ getRailProps }) => (  // adding the rail props sets up events on the rail
            <div style={railStyle} {...getRailProps()} />
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


function Handle({ // your handle component
  handle: { id, value, percent },
  getHandleProps
}) {
  return (
    <div
      style={{
        left: `${percent}%`,
        position: 'absolute',
        marginLeft: -5.5,
        marginTop: 33,
        zIndex: 2,
        width: 15,
        height: 15,
        border: 0,
        textAlign: 'center',
        cursor: 'pointer',
        borderRadius: '50%',
        backgroundColor: 'green',
        color: '#333',
      }}
      {...getHandleProps(id)}
    >
    </div>
  )
}


function Track({ source, target, getTrackProps }) { // your own track component
  return (
    <div
      style={{
        position: 'absolute',
        height: 10,
        zIndex: 1,
        marginTop: 35,
        backgroundColor: '#145814',
        borderRadius: 5,
        cursor: 'pointer',
        left: `${source.percent}%`,
        width: `${target.percent - source.percent}%`,
      }}
      {...getTrackProps()} // this will set up events if you want it to be clickeable (optional)
    />
  )
}

export default VolumeSlider;