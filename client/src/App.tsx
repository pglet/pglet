import React from 'react';
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import PageLanding from './controls/PageLanding'
import AccountLanding from './controls/AccountLanding';

export const App: React.FunctionComponent = () => {
  return (
    <Router>
      <Switch>
        <Route path="/p/:accountName/:pageName" children={<PageLanding />} />
        <Route path="/p/:accountName" children={<AccountLanding />} />
      </Switch>
    </Router>
  );
};
