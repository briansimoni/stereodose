import React from "react";
import { Fragment } from "react";
import Slider, { Rail, Handles, Tracks } from 'react-compound-slider'

const sliderStyle = {
  position: 'relative',
  height: '50px',
  display: 'inline-block',
  'marginLeft': '200px',
  top: '-55px'
}

const domain = [0, 100]
const defaultValues = [0]

class VolumeSlider extends React.Component {

  state = {
    values: defaultValues,
    update: defaultValues,
  }

  onUpdate = update => {
    // nothing right now
  }

  onChange = values => {
    const percent = values[0];
    this.props.onChangeVolume(percent);
  }

  render() {
    // const { state: { values }} = this

    // const state = this.state;
    const values = this.state.values;

    return (
      <div style={{position: 'absolute', width: '100%' }}>
        <Slider
          vertical
          mode={2}
          step={5}
          domain={domain}
          rootStyle={sliderStyle}
          onUpdate={this.onUpdate}
          onChange={this.onChange}
          values={values}
          reversed
          disabled={this.props.disabled}
        >
          <Rail>
            {({ getRailProps }) => <SliderRail getRailProps={getRailProps} />}
          </Rail>
          <Handles>
            {({ handles, getHandleProps }) => (
              <div className="slider-handles">
                {handles.map(handle => (
                  <Handle
                    key={handle.id}
                    handle={handle}
                    domain={domain}
                    getHandleProps={getHandleProps}
                  />
                ))}
              </div>
            )}
          </Handles>
          <Tracks left={false}>
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
          {/* <Ticks count={5}>
            {({ ticks }) => (
              <div className="slider-ticks">
                {ticks.map(tick => (
                  <Tick key={tick.id} tick={tick} />
                ))}
              </div>
            )}
          </Ticks> */}
        </Slider>
      </div>
    )
  }

}


const railOuterStyle = {
  position: 'absolute',
  height: '100%',
  width: 21,
  transform: 'translate(-50%, 0%)',
  borderRadius: 7,
  cursor: 'pointer',
  // border: '1px solid white',
}

const railInnerStyle = {
  position: 'absolute',
  height: '100%',
  width: 7,
  transform: 'translate(-50%, 0%)',
  borderRadius: 7,
  pointerEvents: 'none',
  backgroundColor: 'rgb(155,155,155)',
}

function SliderRail({ getRailProps }) {
  return (
    <Fragment>
      <div style={railOuterStyle} {...getRailProps()} />
      <div style={railInnerStyle} />
    </Fragment>
  )
}



function Handle({
  domain: [min, max],
  handle: { id, value, percent },
  getHandleProps,
}) {
  return (
    <Fragment>
      <div
        style={{
          top: `${percent}%`,
          position: 'absolute',
          transform: 'translate(-50%, -50%)',
          WebkitTapHighlightColor: 'rgba(0,0,0,0)',
          zIndex: 5,
          width: 21,
          height: 14,
          cursor: 'pointer',
          // border: '1px solid white',
          backgroundColor: 'none',
        }}
        {...getHandleProps(id)}
      />
      <div
        role="slider"
        aria-valuemin={min}
        aria-valuemax={max}
        aria-valuenow={value}
        style={{
          top: `${percent}%`,
          position: 'absolute',
          transform: 'translate(-50%, -50%)',
          zIndex: 2,
          width: 12,
          height: 12,
          borderRadius: '50%',
          boxShadow: '1px 1px 1px 1px rgba(0, 0, 0, 0.3)',
          backgroundColor: 'green',
        }}
      />
    </Fragment>
  )
}


function Track({ source, target, getTrackProps }) {
  return (
    <div
      style={{
        position: 'absolute',
        zIndex: 1,
        backgroundColor: 'green',
        borderRadius: 7,
        cursor: 'pointer',
        width: 7,
        transform: 'translate(-50%, 0%)',
        top: `${source.percent}%`,
        height: `${target.percent - source.percent}%`,
      }}
      {...getTrackProps()}
    />
  )
}

export default VolumeSlider;