import React from 'react';
import { Link } from 'react-router-dom';
import "./Screens.css";

// Drugs renders all the drug choices
// Weed, Ecstacy, Shrooms, LSD
export default function Drugs(props) {

  const categories = props.app.state.categories;
  if (categories) {
    let drugNames = Object.keys(categories);

    return (
      <div className="row">
        <div className="col">
          <h2 className="drug-header">Choose Your Drug</h2>
          <ul className="drugs">
            {
              drugNames.map((drug, index) =>
                <li key={index}>
                  <h3><Link to={`/${drug}`}>{drug}</Link></h3>
                </li>
              )
            }
          </ul>
        </div>
      </div>
    );
  }

  return (
    <div>loading...</div>
  );
}