import React from "react";

export default function PickMood(props) {
  const { playlist, drug, categories, onSelectMood } = props;
  return (
    <div>
      <h4>Choose Mood for {playlist.name}</h4>
      <ul>
        {categories[drug].map((mood, index) =>
          <li key={index} onClick={() => { onSelectMood(mood) }}>
            {mood}
          </li>
        )}
      </ul>
    </div>
  )
}