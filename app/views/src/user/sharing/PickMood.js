import React from "react";

export default function PickMood(props) {
  const { playlist, drug, categories, onSelectMood } = props;
  return (
    <div>
      <h2 id="content-title">Choose Mood for {playlist.name}</h2>
      <div className="list-group">
        {categories[drug].map((mood, index) =>
          <button
            type="button"
            className="list-group-item list-group-item-action"
            key={index}
            onClick={() => { onSelectMood(mood) }}>
            {mood}
          </button>
        )}
      </div>
    </div>
  )
}