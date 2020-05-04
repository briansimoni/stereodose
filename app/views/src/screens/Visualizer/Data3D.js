import React from 'react';
import * as Three from 'three';
import Visualizer from './Visualizer';
import { GLTFLoader, GLTF } from 'three/examples/jsm/loaders/GLTFLoader';


export default class Data3D extends Visualizer {


  segmentChanged = true // this is a flag that can be used later to check if the currentSegment changed or not

  /**
  *
  * @param {object} props
  * @param {App} props.app - a reference to the App component
  * @param {function} props.toggleVisualizer
  */
  constructor(props) {
    super(props);
    this.ref = React.createRef();

    this.loadJosiahsDonut().then((nut) => {
      this.scene = new Three.Scene();
      this.camera = new Three.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);

      const renderer = new Three.WebGLRenderer();
      Three.Vector3()
      renderer.setSize(window.innerWidth, window.innerHeight);
      this.renderer = renderer;

      this.geometry = new Three.BoxGeometry(1, 1, 1);
      this.material = new Three.MeshBasicMaterial({ color: 0x00FFFF });
      this.cube = new Three.Mesh(this.geometry, this.material);
      // nut.color(blue);
      // let bigNut = new Three.Object3D();
      // bigNut = nut;

      this.scene.add(nut.scene);
      // this.scene.add(this.cube);
      this.camera.position.z = 5;
      nut.scene.children[0].material = new Three.MeshBasicMaterial({ color: 0x00FFFF });
      this.nut = nut;

      this.animate();

      
      // nut.material.color.setHex( 0x00FFFF );
      // renderer.material.color = 'blue';
      window.addEventListener('resize', this.resize, true);
      this.forceUpdate();
    })
  }

  /**
   * returns a Promise that resolves to a three js object
   * @returns {Promise<GLTF>}
   */
  loadJosiahsDonut() {
    return new Promise((resolve, reject) => {
      const loader = new GLTFLoader();
      loader.load('/public/three/blobvis.glb', (gltf, thing) => {
        resolve(gltf);
        return;
      })
    })
  }

  playerStateChangeCallback = async playerState => {
    if (this.currentTrack !== playerState.track_window.current_track.id) {
      this.playerState = playerState;
      this.currentTrackId = playerState.track_window.current_track.id;
      await this.getTrackInfo(this.currentTrackId);
    }
  };

  async componentWillUnmount() {
    super.componentWillUnmount();
    const nut = await this.loadJosiahsDonut();

    this.scene.add(nut);
    window.removeEventListener('resize', this.resize, true);
    cancelAnimationFrame(this.animationFrameId);
  }

  resize = event => {
    this.forceUpdate();
    this.camera.aspect = window.innerWidth / window.innerHeight;
    this.renderer.setSize(window.innerWidth, window.innerHeight);
  };

  render() {
    if (!this.scene) {
      return <div></div>
    }
    this.camera.aspect = window.innerWidth / window.innerHeight;

    this.renderer.setSize(window.innerWidth, window.innerHeight);

    return (
      <div
        id="visualizer"
        ref={nodeElement => {
          nodeElement && nodeElement.appendChild(this.renderer.domElement);
        }}
      >
        <button onClick={this.props.toggleVisualizer} id="close-visualizer-button" className="btn btn-danger">CLOSE</button>
      </div>
    );
  }

  animate = async () => {
    await this.setActiveIntervals();
    let donut = null;
    if (!this.nut) {
      return;
    }
    donut = this.nut.scene.children[0];
    // if (donut.rotation.x < 0.5) {
    //   donut.rotation.x += 0.01;
    // }
    donut.rotation.x += 0.001;
    donut.rotation.y += 0.002;
    donut.rotation.z += 0.003;

    // scale the donut to the volume if the volume data is available.
    if (this.currentSegment) {
      const x = this.position;
      let { loudness_start, start, duration } = this.currentSegment;
      let loudness_end = this.nextSegment.loudness_start;
      loudness_start = Math.abs(loudness_start);
      loudness_end = Math.abs(loudness_end);

      if (loudness_start !== this.placeholder) {
        this.segmentChanged = true;
      }

      if (this.segmentChanged) {
        this.placeholder = loudness_start; // used to detect if the segment changed
        this.renderVolume = loudness_start; // mutated and used to render the donut scale
        this.segmentChanged = false;
      }

      const x1 = start * 1000; // convert seconds to ms
      const y1 = loudness_start;
      const x2 = (start + duration) * 1000 // convert seconds to ms
      const y2 = loudness_end;


      const m = (y2 - y1) / (x2 - x1);
      const b = y1 - (m * x1);

      // console.log({
      //   y: (m *x + b),
      //   m,
      //   x,
      //   b,
      //   x1,
      //   y1,
      //   x2,
      //   y2,
      //   duration,
      //   renderVolume: this.renderVolume
      // })

      // console.log(x);

      //y=mx+b
      // (y2-y1)/(x2-x1) = m
      //delta y/delta x = m
      // console.log((m *x) + b);
      // const volume = (m * x + b) // = y
      const volume = 100 - (m * x + b);
      console.log(volume);
      donut.scale.x = Math.abs(volume) / 35;
      donut.scale.y = Math.abs(volume) / 35;
      donut.scale.z = Math.abs(volume) / 35;
      // console.log('new scale x', donut.scale.x);

    }

    this.renderer.render(this.scene, this.camera);
    requestAnimationFrame(this.animate);
  }

}
