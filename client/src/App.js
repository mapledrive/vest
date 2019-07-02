import React, { Component } from "react";
import Header from "./components/Header";
import CumulativeContainer from "./containers/CumulativeContainer";
import SidebarContainer from "./containers/SidebarContainer";
import PatternContainer from "./containers/PatternContainer";
import TableContainer from "./containers/TableContainer";
import ChartContainer from "./containers/ChartContainer";
import ControlContainer from "./containers/ControlContainer";
import Search from "./components/Search";
import Title from "./components/Title";
import Empty from "./components/Empty";
import { Route, Switch } from "react-router-dom";

class App extends Component {
  render() {
    return (
      <div>
        <Header />
        <Switch>
          <div className="container-fluid">
            <div className="row">
              <div className="col-12 col-lg-9">
                <div className="row">
                  <Title />
                  <Search />
                  <ControlContainer />
                  <div className="col-12 mini-padding">
                    <div className="card">
                      <div className="card-body">
                        <Route exact path="/" component={ChartContainer} />
                        <Route path="/crypto" component={Empty} />
                        <Route path="/commodities" component={Empty} />
                        <Route path="/american" component={Empty} />
                        <Route path="/european" component={Empty} />
                        <Route path="/international" component={Empty} />
                        <Route path="/stocks" component={Empty} />
                        <Route path="/ussectors" component={Empty} />
                      </div>
                    </div>
                  </div>
                  <CumulativeContainer />
                  <PatternContainer />
                  <TableContainer />
                </div>
              </div>
              <SidebarContainer />
            </div>
          </div>
        </Switch>
      </div>
    );
  }
}

export default App;
