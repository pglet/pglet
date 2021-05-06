import React from 'react';
import "./App.css";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import { PageLanding } from './controls/PageLanding'
import { AccountLanding } from './controls/AccountLanding';

export const App: React.FunctionComponent = () => {
  return (
    <Router>
      <Switch>
        <Route path="/:accountName/:pageName" children={<PageLanding />} />
        <Route path="/:accountName" children={<AccountLanding />} />
        <Route path="/" children={<PageLanding />} />
      </Switch>
    </Router>
  );
};
