import React from 'react';
import "./App.css";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import { Page } from './controls/Page'
import { AccountLanding } from './controls/AccountLanding';

export const App: React.FunctionComponent = () => {
  return (
    <Router>
      <Switch>
        <Route path="/:accountName/:pageName" children={<Page />} />
        <Route path="/:accountName" children={<AccountLanding />} />
        <Route path="/" children={<Page />} />
      </Switch>
    </Router>
  );
};
