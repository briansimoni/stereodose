import React from "react";

export default function PickDrug(props) {
  const { onSelectDrug, categories, playlist } = props;
  const drugs = Object.keys(categories);
  return (
    <div>
      <h2 id="tab-content-title">Choose Drug for {playlist.name}</h2>
      <div className="list-group">
        {drugs.map((drug, index) =>
          <button
            type="button"
            className="list-group-item list-group-item-action"
            key={index}
            onClick={() => { onSelectDrug(drug) }}>
            {drug}
          </button>
        )}
      </div>
    </div>
  )
}