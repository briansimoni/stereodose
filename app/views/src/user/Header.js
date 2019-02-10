import React from 'react';
import { Link } from "react-router-dom";
import stereodoseLogo from "../images/logo.png";
import UserStatusIndicator from "./StatusIndicator";

export default function Header(props) {
  return (
    <header>
      <nav className="navbar navbar-expand-lg navbar-dark">

      <Link className="navbar-brand" to="/"><img src={stereodoseLogo} id="logo" alt="logo" /></Link>

        <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
          <span className="navbar-toggler-icon"></span>
        </button>

        <div className="collapse navbar-collapse" id="navbarSupportedContent">
          <ul className="navbar-nav ml-auto float-right">
            <li className="nav-item">
              <Link className="nav-link" to="/">Home</Link>
            </li>
            <UserStatusIndicator isUserLoggedIn={props.isUserLoggedIn}/>
          </ul>
        </div>

      </nav>
    </header>
  );
};