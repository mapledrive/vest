import React, { Component } from "react";
import "./sidebar.css";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import PropTypes from 'prop-types';

class Example extends Component {
  constructor(props) {
    super(props);
    this.state = {
      startDate: new Date(2018,6,10),
	  endDate: new Date()
    };
    this.handleChangeStart = this.handleChangeStart.bind(this);
	this.handleChangeEnd = this.handleChangeEnd.bind(this);
  }
 
  handleChangeStart(date) {
    this.setState({ startDate: date });
	this.props.getRepos(this.state);
	this.props.getStats(this.state);
  }
 
   handleChangeEnd(date) {
    this.setState({ endDate: date });
	this.props.getRepos(this.state);
	this.props.getStats(this.state);
  }
  
	
  render() {
    return (
	<div>
        <DatePicker
			customInput={<ExampleCustomInput />}
			
			selected={this.state.startDate}
			selectsStart
			startDate={this.state.startDate}
			endDate={this.state.endDate}
			onChange={this.handleChangeStart}
		/>
		<span className="date-separator"> - </span>
		<DatePicker
			customInput={<ExampleCustomInput />}
			
			selected={this.state.endDate}
			selectsEnd
			startDate={this.state.startDate}
			endDate={this.state.endDate}
			onChange={this.handleChangeEnd}
		/>
	</div>
    );
  }
}

class ExampleCustomInput extends React.Component {
  render () {
    return (
      <button className="btn btn-info pattern-filter" onClick={this.props.onClick}>{this.props.value}</button>
    )
  }
}

ExampleCustomInput.propTypes = {
  onClick: PropTypes.func,
  value: PropTypes.string
};


export default Example;