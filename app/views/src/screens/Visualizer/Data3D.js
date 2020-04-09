import React from 'react';
import * as Three from 'three';
import Visualizer from './Visualizer';
import { GLTFLoader, GLTF } from 'three/examples/jsm/loaders/GLTFLoader';


export default class Data3D extends Visualizer {
  /**
  *
  * @param {object} props
  * @param {App} props.app - a reference to the App component
  * @param {function} props.toggleVisualizer
  */
  constructor(props) {
    super(props);
    this.ref = React.createRef();

    this.loadJosiahsDoughnut().then((nut) => {
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

      console.log(nut);
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
  loadJosiahsDoughnut() {
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
    const nut = await this.loadJosiahsDoughnut();

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
    requestAnimationFrame(this.animate);
    if (this.nut) {
      console.log(this.nut);
      // this.nut.rotation.x += 0.01;
      // this.nut.rotation.y += 0.01;
    }
    this.nut.scene.children[0].rotation.x += 0.01;
    this.nut.scene.children[0].rotation.y += 0.02;
    this.nut.scene.children[0].rotation.z += 0.03;
    this.renderer.render(this.scene, this.camera);
  };
}
