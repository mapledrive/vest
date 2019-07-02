import React, {Component} from "react";
import { Link } from 'react-router-dom';
import Logo from "../components/Logo";

class Header extends Component {
	constructor(props) {
    super(props);
    this.state = {isToggleOn: true};
    this.handleClick = this.handleClick.bind(this);
  }
  
  handleClick() {
    this.setState(state => ({ isToggleOn: !state.isToggleOn }));
  }
	render(){
		return (
    <nav className="navbar navbar-expand-xl bg-light navbar-light">
	  <div className="navbar-brand">
          <Logo />
        </div>
      <button
        className="navbar-toggler"
        type="button"
		onClick={this.handleClick}
        data-toggle="collapse"
        data-target="#collapsibleNavbar"
      >
        <span className="navbar-toggler-icon" />
      </button>
	  <div className={"navbar-collapse " + (this.state.isToggleOn ? 'collapse' : '')} id="collapsibleNavbar">
        <ul className="navbar-nav">
          <li className="nav-item">
              <Link className="nav-link" to='/'>Gold</Link> 
          </li>
          <li className="nav-item">
              <Link className="nav-link" to='/crypto'>Crypto</Link>
          </li>
		  <li className="nav-item">
              <Link className="nav-link" to='/commodities'>Commodities</Link>
          </li>
		  <li className="nav-item">
              <Link className="nav-link" to='/american'>American</Link>
          </li>
		  <li className="nav-item">
              <Link className="nav-link" to='/european'>European</Link>
          </li>
		  <li className="nav-item">
              <Link className="nav-link" to='/international'>International</Link>
          </li>
		  <li className="nav-item">
              <Link className="nav-link" to='/stocks'>Stocks</Link>
          </li>
		  <li className="nav-item">
              <Link className="nav-link" to='/ussectors'>US Sectors</Link>
          </li>
        </ul>
        <ul className="navbar-nav ml-auto d-none d-xl-block">
          <button type="button" className="btn btn-outline-info mr-2">
            <i className="fa fa-heart-o mr-1" /> Signup
          </button>
          <button type="button" className="btn btn-outline-info">
            <i className="fa fa-sign-in mr-1" /> Login
          </button>
        </ul>
      </div>
    </nav>
  )
	}
}

export default Header;