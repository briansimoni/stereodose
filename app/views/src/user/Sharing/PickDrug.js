import React from "react";

export default function PickDrug(props) {
  const { onSelectDrug, categories, playlist } = props;
  const drugs = Object.keys(categories);
  return (
    <div>
      <h4>Choose Drug for {playlist.name}</h4>
      <ul>
        {drugs.map((drug, index) =>
          <li key={index} onClick={() => { onSelectDrug(drug) }}>
            {drug}
          </li>
        )}
      </ul>
    </div>
  )
}