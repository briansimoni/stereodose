import React from 'react';
import ReactCrop from 'react-image-crop';
import 'react-image-crop/dist/ReactCrop.css';

// PickImage presents a way to upload and crop an image to predetermined aspect ratio
// When the user crops the image, it sets a blob value in the parent component
// Using a lot of code from the example referenced on https://github.com/DominicTobias/react-image-crop
// https://codesandbox.io/s/72py4jlll6
export default class PickImage extends React.Component {
  state = {
    fileInputRef: React.createRef(),
    src: null,
    crop: {
      aspect: 1
    }
  };

  onSelectFile = e => {
    if (e.target.files && e.target.files.length > 0) {
      const reader = new FileReader();
      reader.addEventListener('load', () => this.setState({ src: reader.result }));
      reader.readAsDataURL(e.target.files[0]);
    }
  };

  onCropChange = crop => {
    this.setState({ crop });
  };

  onCropComplete = (crop, pixelCrop) => {
    this.makeClientCrop(crop, pixelCrop);
  };

  async makeClientCrop(crop, pixelCrop) {
    if (this.imageRef && crop.width && crop.height) {
      const croppedImageBlob = await this.getCroppedImg(this.imageRef, pixelCrop, 'newFile.jpeg');
      this.setState({ croppedImageBlob });
      const onBlobCreated = this.props.onBlobCreated;
      onBlobCreated(croppedImageBlob);
    }
  }

  onImageLoaded = (image, pixelCrop) => {
    this.imageRef = image;
  };

  getCroppedImg(image, pixelCrop, fileName) {
    const canvas = document.createElement('canvas');
    canvas.width = pixelCrop.width;
    canvas.height = pixelCrop.height;
    const ctx = canvas.getContext('2d');

    ctx.drawImage(
      image,
      pixelCrop.x,
      pixelCrop.y,
      pixelCrop.width,
      pixelCrop.height,
      0,
      0,
      pixelCrop.width,
      pixelCrop.height
    );

    return new Promise((resolve, reject) => {
      canvas.toBlob(blob => {
        if (!blob) {
          //reject(new Error('Canvas is empty'));
          console.error('Canvas is empty');
          return;
        }
        resolve(blob);
      }, 'image/jpeg');
    });
  }

  render() {
    let fileInputRef = React.createRef();
    const { src, crop } = this.state;
    return (
      <div>
        <form>
          <div className="input-group mb-3">
            <div className="custom-file">
              <input
                className="custom-file-input"
                ref={fileInputRef}
                type="file"
                id="file-upload"
                name="playlist-image"
                accept=".jpg,.jpeg,.png"
                onChange={this.onSelectFile}
              />
              <label className="custom-file-label" htmlFor="file-upload" aria-describedby="file-upload">
                Choose file
              </label>
            </div>
          </div>
        </form>
        {!src && <img alt="playlist placeholder" src="https://via.placeholder.com/350" />}
        {src && (
          <ReactCrop
            src={src}
            crop={crop}
            onImageLoaded={this.onImageLoaded}
            onComplete={this.onCropComplete}
            onChange={this.onCropChange}
          />
        )}
      </div>
    );
  }
}
