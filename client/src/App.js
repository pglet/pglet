import React from 'react';
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
//import logo from './logo.svg';
import './pglet.scss';
import PageLanding from './components/PageLanding'
import AccountLanding from './components/AccountLanding';

const App = () => {
  return (
    <div className="container-fluid">
      <Router>
        <Switch>
          <Route path="/p/:accountName/:pageName" children={<PageLanding />} />
          <Route path="/p/:accountName" children={<AccountLanding />} />
        </Switch>
      </Router>
    </div>
  );
}

export default App;
