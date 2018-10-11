import React from 'react';
import { Link } from "react-router-dom";
import stereodoseLogo from "../images/logo.png";

export default (props) => {
  return (
    <header>
      <nav className="navbar">
	  	<Link className="navbar-brand" to="/"><img src={stereodoseLogo} id="logo" alt="logo"/></Link>
		{props.children}
      </nav>
    </header>
  );
};