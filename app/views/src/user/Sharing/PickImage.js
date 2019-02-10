import React from "react";

export default function PickImage(props) {
  let fileInputRef = React.createRef();
  const { onUploadImage } = props;
  return (
    <form>
      <input
        ref={fileInputRef}
        type="file"
        id="file-upload"
        name="playlist-image"
        accept=".jpg,.jpeg,.png"
        onChange={() => {onUploadImage(fileInputRef)}} />
    </form>
  )
}