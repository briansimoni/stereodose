import React from 'react';
import * as Three from 'three';


export default class Visualizer extends React.Component {
  constructor(props) {
    super(props);
    this.ref = React.createRef();

    this.scene = new Three.Scene();
    this.camera = new Three.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);

    const renderer = new Three.WebGLRenderer();
    Three.Vector3()
    renderer.setSize(window.innerWidth, window.innerHeight);
    this.renderer = renderer;

    this.geometry = new Three.BoxGeometry(1, 1, 1);
    this.material = new Three.MeshBasicMaterial({ color: 0x00ff00 });
    this.cube = new Three.Mesh(this.geometry, this.material);
    this.scene.add(this.cube);

    this.camera.position.z = 5;
    this.animate();

    window.addEventListener('resize', this.resize, true);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resize, true);
    cancelAnimationFrame(this.animationFrameId);
  }

  resize = event => {
    this.forceUpdate();
    this.camera.aspect = window.innerWidth / window.innerHeight;
    this.renderer.setSize(window.innerWidth, window.innerHeight);
  };

  render() {
    this.camera.aspect = window.innerWidth / window.innerHeight;

    this.renderer.setSize(window.innerWidth, window.innerHeight);

    return (
      <div
        id="visualizer"
        ref={nodeElement => {
          nodeElement && nodeElement.appendChild(this.renderer.domElement);
        }}
      >
        <button onClick={this.props.toggleVisualizer} id="close-visualzier-button" className="btn btn-danger">CLOSE</button>
      </div>
    );
  }

  animate = async () => {
    const playerState = await this.props.app.player.getCurrentState();
    const analysis = this.props.analysis;

    const position = playerState.position;
    const segmentLength = analysis.segments.length;
    const songDuration = playerState.track_window.current_track.duration_ms;
    const segmentIndex = Math.floor((segmentLength / songDuration) * position);
    if (segmentIndex !== this.segmentIndex) {
      console.log(segmentIndex);
    }
    this.segmentIndex = segmentIndex;

    const segment = analysis.segments[segmentIndex];

    const animationFrameId = requestAnimationFrame(this.animate);
    this.animationFrameId = animationFrameId;
    // this.cube.translateY(0.01);
    // console.log(this.cube.scale);
    if (this.currentSegmentLoudness === Math.abs(segment.loudness_max)) {
      this.cube.scale.y -= 0.01;
    } else {
      this.cube.scale.y = Math.abs(segment.loudness_max / 10);
    }
    this.currentSegmentLoudness = Math.abs(segment.loudness_max);

    // console.log(this.cube.geometry);
    // this.cube.rotation.x += 0.01;
    // this.cube.rotation.y += 0.01;
    // this.cube.rotation.y = Math.abs(segment.loudness_max + Math.random());
    this.renderer.render(this.scene, this.camera);
  };
}
