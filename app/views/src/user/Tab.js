import React, { Component } from 'react';

// from https://alligator.io/react/tabs-component/
// required props:
// 1. activeTab : bool
// 2. label : string
// 3. onClick : func
class Tab extends Component {

  onClick = () => {
    const { label, onClick } = this.props;
    onClick(label);
  }

  render() {
    const {
      onClick,
      props: {
        activeTab,
        label,
      },
    } = this;

    let className = 'nav-link';

    if (activeTab === label) {
      className += ' nav-link active';
    }

    return (
      <li
        className="nav-item"
        onClick={onClick}
      >
        <span className={className}>{label}</span>
      </li>
    );
  }
}

export default Tab;