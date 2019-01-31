import React from 'react';
import { connect } from "react-redux";


const addRepos = repos => ({ type: "ADD_REPOS", repos });
const clearRepos = () => ({ type: "CLEAR_REPOS" });

const getRepos = username => async dispatch => {
  try {
	var rangestart;
	if (username === "5") {rangestart = '2016-01-01'};
	if (username === "10") {rangestart = '2015-01-01'};
	if (username === "15") {rangestart = '2014-01-01'};
	if (username === "20") {rangestart = '2013-01-01'};
	if (username === "25") {rangestart = '2012-01-01'};
	var url = `http://ec2-54-84-182-191.compute-1.amazonaws.com:8080/vest/v1/charts?code=gc&rs=${rangestart}&re=2019-01-26`;
    const response = await fetch(url);
    const responseBody = await response.json();
    dispatch(addRepos(responseBody));
  } catch (error) {
    console.log(error);
    dispatch(clearRepos());
  }
};


class Control extends React.Component {
	constructor(props) {
		super(props);
		this.handleChange = this.handleChange.bind(this);
		this.state = {isFilterOn: false, isShareOn: false};
		this.handleFilter = this.handleFilter.bind(this);
		this.handleShare = this.handleShare.bind(this);
	}

	handleChange(event) {
		this.props.getRepos(event.target.value);
	}
	
	handleFilter() {
		this.setState(state => ({ isFilterOn: !state.isFilterOn }));
	}
	
	handleShare() {
		this.setState(state => ({ isShareOn: !state.isShareOn }));
	}
	
	render(){
		return (
			<div className="col-12 mini-padding">
			  <div className="card">
					<div className="card-body">
				<span className="btn controller float-right"><i className="fa fa-undo" /></span>
				
				  <div className="dropdown controller float-right mr-2">
					<button type="button" onClick={this.handleFilter} className="btn controller float-right dropdown-toggle" data-toggle="dropdown">
					  Filter
					</button>
					<div className={"dropdown-menu " + (this.state.isFilterOn ? 'show' : '')} >
					  <a className="dropdown-item" href="i.htm"><i className="fa fa-arrow-circle-up mr-2"></i>Bullish</a>
					  <a className="dropdown-item" href="i.htm"><i className="fa fa-arrow-circle-down mr-2"></i>Bearish</a>
					  <a className="dropdown-item" href="i.htm"><i className="fa fa-arrow-circle-right mr-2"></i>Neutral</a>
					</div>
				  </div>
				  <select onChange={this.handleChange} className="custom-select year float-right mr-2" id="selectRange">
					<option value="5">  5 years</option>
					<option value="10">10 years</option>
					<option value="15">15 years</option>
					<option value="20">20 years</option>
					<option value="25">25 years</option>
				  </select>
				  
				  <div className="btn-group float-right mr-2">
					<button type="button" className="btn controller"><i className="fa fa-line-chart" /></button>
					<button type="button" className="btn controller"><i className="fa fa-calendar-check-o" /></button>
				  </div>
				  
				  <div className="btn-group float-right mr-2">
					<button type="button" className="btn controller"><i className="fa fa-chevron-left"></i></button>
					<button type="button" className="btn controller"><i className="fa fa-chevron-right"></i></button>
				  </div>
					
				  <div className="dropdown controller float-right mr-2">
					<button type="button" onClick={this.handleShare} className="btn controller float-right dropdown-toggle" data-toggle="dropdown">
					  Share
					</button>
					<div className={"dropdown-menu " + (this.state.isShareOn ? 'show' : '')} >
					  <a className="dropdown-item" href="i.htm"><i className="fa fa-files-o mr-2"></i>Copy Link</a>
					  <a className="dropdown-item" href="i.htm"><i className="fa fa-file-image-o mr-2"></i>Screenshot</a>
					  <a className="dropdown-item" href="i.htm"><i className="fa fa-file-code-o mr-2"></i>Embed code</a>
					</div>
				  </div>
				  <span className="btn controller float-right mr-2"><i className="fa fa-question-circle" /></span>
			
			
					</div>
				</div>
			</div>
			
			
			
		);
	}
}

const mapDispatchToProps = { getRepos };
const ControlContainer = connect(null, mapDispatchToProps)(Control);

export default ControlContainer;