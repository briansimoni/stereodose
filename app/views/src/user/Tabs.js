import React, { Component } from 'react';

import Tab from './Tab';

// from https://alligator.io/react/tabs-component/
class Tabs extends Component {

  constructor(props) {
    super(props);

    this.state = {
      activeTab: this.props.children[0].props.label,
    };
  }

  onClickTabItem = (tab) => {
    this.setState({ activeTab: tab });
  }

  render() {
    const {
      onClickTabItem,
      props: {
        children,
      },
      state: {
        activeTab,
      }
    } = this;

    return (
      <div className="tabs container">

        <div className="row">
          <div className="col">
            <ul className="nav nav-pills justify-content-center">
              {children.map((child) => {
                const { label } = child.props;

                return (
                  <Tab
                    activeTab={activeTab}
                    key={label}
                    label={label}
                    onClick={onClickTabItem}
                  />
                );
              })}
            </ul>
          </div>
        </div>

        <div className="row">
          <div className="col">
            <div className="tab-content">
              {children.map((child) => {
                if (child.props.label !== activeTab) return undefined;
                return child.props.children;
              })}
            </div>
          </div>
        </div>

      </div>
    );
  }
}

export default Tabs;