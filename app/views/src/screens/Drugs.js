import React from 'react';
import {Link} from 'react-router-dom';
import "./Screens.css";

// Drugs renders all the drug choices
// Weed, Ecstacy, Shrooms, LSD
class Drugs extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      loading: true
    }
  }

  render() {
    if (this.state.loading) {
      return <p>Loading...</p>
    }

    if (this.state.categories !== null) {
      let drugNames = Object.keys(this.state.categories);

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
      )
    }

    if (this.state.error) {
      return <p>{this.state.error}</p>
    }
  }

  componentDidMount() {
    fetch("/api/categories/", { credentials: "same-origin" })
      .then((response) => {
        return response.json();
      })
      .then((json) => {
        this.setState({ loading: false, categories: json });
      })
      .catch((err) => {
        this.setState({ loading: false, error: err });
      })
  }
}

export default Drugs;