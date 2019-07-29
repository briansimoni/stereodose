import React from 'react';
import Slider, { Rail, Handles, Tracks } from 'react-compound-slider';
import Octicon from 'react-octicon';

const railStyle = {
  position: 'absolute',
  width: '100%',
  height: 10,
  marginTop: 35,
  borderRadius: 5,
  backgroundColor: '#8B9CB6'
};

const defaultValues = [50];

class VolumeSlider extends React.Component {
  state = {
    values: defaultValues,
    update: defaultValues
  };

  onUpdate = update => {
    // nothing right now
  };

  onChange = values => {
    // values is a set of percentages
    // the web SDK takes a value from 0 - 1
    this.props.onChangeVolume(values[0] / 100);
  };

  render() {
    const values = this.state.values;
    const { className, disabled } = this.props;

    return (
      <Slider className={className} domain={[0, 100]} values={values} onChange={this.onChange} disabled={disabled}>
        <Octicon className="unmute" name="unmute" />
        <Rail>
          {(
            { getRailProps } // adding the rail props sets up events on the rail
          ) => <div style={railStyle} {...getRailProps()} />}
        </Rail>
        <Handles>
          {({ handles, getHandleProps }) => (
            <div className="slider-handles">
              {handles.map(handle => (
                <Handle key={handle.id} handle={handle} getHandleProps={getHandleProps} disabled={disabled} />
              ))}
            </div>
          )}
        </Handles>
        <Tracks right={false}>
          {({ tracks, getTrackProps }) => (
            <div className="slider-tracks">
              {tracks.map(({ id, source, target }) => (
                <Track key={id} source={source} target={target} getTrackProps={getTrackProps} disabled={disabled} />
              ))}
            </div>
          )}
        </Tracks>
      </Slider>
    );
  }
}

function Handle({
  // your handle component
  handle: { id, value, percent },
  getHandleProps,
  disabled
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
        cursor: disabled ? 'not-allowed' : 'pointer',
        borderRadius: '50%',
        backgroundColor: 'green',
        color: '#333'
      }}
      {...getHandleProps(id)}
    />
  );
}

function Track({ source, target, getTrackProps, disabled }) {
  // your own track component
  return (
    <div
      style={{
        position: 'absolute',
        height: 10,
        zIndex: 1,
        marginTop: 35,
        backgroundColor: '#145814',
        borderRadius: 5,
        cursor: disabled ? 'not-allowed' : 'pointer',
        left: `${source.percent}%`,
        width: `${target.percent - source.percent}%`
      }}
      {...getTrackProps()} // this will set up events if you want it to be clickeable (optional)
    />
  );
}

export default VolumeSlider;
